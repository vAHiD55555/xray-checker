package checker

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"xray-checker/metrics"
	"xray-checker/models"
)

type ProxyChecker struct {
	proxies        []*models.ProxyConfig
	startPort      int
	ipCheck        string
	currentIP      string
	httpClient     *http.Client
	currentMetrics sync.Map
	ipInitialized  bool
	ipCheckTimeout int
}

func NewProxyChecker(proxies []*models.ProxyConfig, startPort int, ipCheckURL string, ipCheckTimeout int) *ProxyChecker {
	return &ProxyChecker{
		proxies:   proxies,
		startPort: startPort,
		ipCheck:   ipCheckURL,
		httpClient: &http.Client{
			Timeout: time.Second * time.Duration(ipCheckTimeout),
		},
		ipCheckTimeout: ipCheckTimeout,
	}
}

func (pc *ProxyChecker) GetCurrentIP() (string, error) {
	if pc.ipInitialized && pc.currentIP != "" {
		return pc.currentIP, nil
	}

	resp, err := pc.httpClient.Get(pc.ipCheck)
	if err != nil {
		return "", fmt.Errorf("error getting current IP: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	pc.currentIP = string(body)
	pc.ipInitialized = true
	return pc.currentIP, nil
}

func (pc *ProxyChecker) CheckProxy(proxy *models.ProxyConfig) {
	metricKey := fmt.Sprintf("%s|%s:%d|%s",
		proxy.Protocol,
		proxy.Server,
		proxy.Port,
		proxy.Name,
	)

	setFailedStatus := func() {
		metrics.RecordProxyStatus(
			proxy.Protocol,
			fmt.Sprintf("%s:%d", proxy.Server, proxy.Port),
			proxy.Name,
			0,
		)
		pc.currentMetrics.Store(metricKey, false)
	}

	proxyURL := fmt.Sprintf("socks5://127.0.0.1:%d", pc.startPort+proxy.Index)
	proxyURLParsed, err := url.Parse(proxyURL)
	if err != nil {
		log.Printf("Error parsing proxy URL %s: %v", proxyURL, err)
		setFailedStatus()
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURLParsed),
		},
		Timeout: time.Second * time.Duration(pc.ipCheckTimeout),
	}

	resp, err := client.Get(pc.ipCheck)
	if err != nil {
		log.Printf("%s | Error | %v", proxy.Name, err)
		setFailedStatus()
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response from proxy %s: %v", proxy.Name, err)
		setFailedStatus()
		return
	}

	proxyIP := string(body)
	if proxyIP == pc.currentIP {
		log.Printf("%s | Failed | Source IP: %s | Proxy IP: %s", proxy.Name, pc.currentIP, proxyIP)
		setFailedStatus()
	} else {
		log.Printf("%s | Success | Source IP: %s | Proxy IP: %s", proxy.Name, pc.currentIP, proxyIP)
		metrics.RecordProxyStatus(
			proxy.Protocol,
			fmt.Sprintf("%s:%d", proxy.Server, proxy.Port),
			proxy.Name,
			1,
		)
		pc.currentMetrics.Store(metricKey, true)
	}
}

func (pc *ProxyChecker) ClearMetrics() {
	pc.currentMetrics.Range(func(key, _ interface{}) bool {
		metricKey := key.(string)
		parts := strings.Split(metricKey, "|")
		if len(parts) == 3 {
			metrics.DeleteProxyStatus(parts[0], parts[1], parts[2])
		}
		pc.currentMetrics.Delete(key)
		return true
	})
}

func (pc *ProxyChecker) UpdateProxies(newProxies []*models.ProxyConfig) {
	pc.ClearMetrics()
	pc.proxies = newProxies
}

func (pc *ProxyChecker) CheckAllProxies() {
	if _, err := pc.GetCurrentIP(); err != nil {
		log.Printf("Error getting current IP: %v", err)
		return
	}

	var wg sync.WaitGroup
	for _, proxy := range pc.proxies {
		wg.Add(1)
		go func(p *models.ProxyConfig) {
			defer wg.Done()
			pc.CheckProxy(p)
		}(proxy)
	}
	wg.Wait()
}

func (pc *ProxyChecker) GetProxyStatus(name string) (bool, error) {
	var metricKey string
	for _, proxy := range pc.proxies {
		if proxy.Name == name {
			metricKey = fmt.Sprintf("%s|%s:%d|%s",
				proxy.Protocol,
				proxy.Server,
				proxy.Port,
				proxy.Name,
			)
			break
		}
	}

	if metricKey == "" {
		return false, fmt.Errorf("proxy not found")
	}

	if value, ok := pc.currentMetrics.Load(metricKey); ok {
		return value.(bool), nil
	}

	return false, fmt.Errorf("metric not found")
}

func (pc *ProxyChecker) GetProxies() []*models.ProxyConfig {
	return pc.proxies
}
