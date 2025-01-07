package checker

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptrace"
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
	latencyMetrics sync.Map
	ipInitialized  bool
	ipCheckTimeout int
	genMethodURL   string
	checkMethod    string
	lastDuration   time.Duration
}

func NewProxyChecker(proxies []*models.ProxyConfig, startPort int, ipCheckURL string, ipCheckTimeout int, genMethodURL string, checkMethod string) *ProxyChecker {
	return &ProxyChecker{
		proxies:   proxies,
		startPort: startPort,
		ipCheck:   ipCheckURL,
		httpClient: &http.Client{
			Timeout: time.Second * time.Duration(ipCheckTimeout),
		},
		ipCheckTimeout: ipCheckTimeout,
		genMethodURL:   genMethodURL,
		checkMethod:    checkMethod,
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

	setFailedLatency := func() {
		metrics.RecordProxyLatency(
			proxy.Protocol,
			fmt.Sprintf("%s:%d", proxy.Server, proxy.Port),
			proxy.Name,
			0,
		)
		pc.latencyMetrics.Store(metricKey, 0)
	}

	proxyURL := fmt.Sprintf("socks5://127.0.0.1:%d", pc.startPort+proxy.Index)
	proxyURLParsed, err := url.Parse(proxyURL)
	if err != nil {
		log.Printf("Error parsing proxy URL %s: %v", proxyURL, err)
		setFailedStatus()
		setFailedLatency()
		return
	}

	var checkSuccess bool
	var checkErr error
	var logMessage string

	if pc.checkMethod == "ip" {
		checkSuccess, logMessage, checkErr = pc.checkByIP(proxyURLParsed)
	} else if pc.checkMethod == "gen" {
		checkSuccess, checkErr = pc.checkByGen(proxyURLParsed)
		if checkSuccess {
			logMessage = "Status: 204"
		} else {
			logMessage = "Check failed"
		}
	}

	if checkErr != nil {
		log.Printf("%s | Error | %v", proxy.Name, checkErr)
		setFailedStatus()
		setFailedLatency()
		return
	}

	if !checkSuccess {
		log.Printf("%s | Failed | %s | Latency: %s", proxy.Name, logMessage, pc.lastDuration)
		setFailedStatus()
		setFailedLatency()
	} else {
		log.Printf("%s | Success | %s | Latency: %s", proxy.Name, logMessage, pc.lastDuration)
		metrics.RecordProxyStatus(
			proxy.Protocol,
			fmt.Sprintf("%s:%d", proxy.Server, proxy.Port),
			proxy.Name,
			1,
		)
		metrics.RecordProxyLatency(
			proxy.Protocol,
			fmt.Sprintf("%s:%d", proxy.Server, proxy.Port),
			proxy.Name,
			float64(pc.lastDuration.Milliseconds()),
		)

		pc.latencyMetrics.Store(metricKey, pc.lastDuration)
		pc.currentMetrics.Store(metricKey, true)
	}
}

func (pc *ProxyChecker) checkByIP(proxyURL *url.URL) (bool, string, error) {
	var start time.Time

	wrappedTrace := &httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) {
			start = time.Now()
		},
		GotFirstResponseByte: func() {
			pc.lastDuration = time.Since(start)
		},
	}

	req, err := http.NewRequest("GET", pc.ipCheck, nil)
	if err != nil {
		return false, "", err
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), wrappedTrace))

	client := &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyURL(proxyURL),
			DisableKeepAlives: true,
		},
		Timeout: time.Second * time.Duration(pc.ipCheckTimeout),
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", err
	}

	proxyIP := string(body)
	logMessage := fmt.Sprintf("Source IP: %s | Proxy IP: %s", pc.currentIP, proxyIP)
	return proxyIP != pc.currentIP, logMessage, nil
}

func (pc *ProxyChecker) checkByGen(proxyURL *url.URL) (bool, error) {
	var start time.Time

	wrappedTrace := &httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) {
			start = time.Now()
		},
		GotFirstResponseByte: func() {
			pc.lastDuration = time.Since(start)
		},
	}

	req, err := http.NewRequest("GET", pc.genMethodURL, nil)
	if err != nil {
		return false, err
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), wrappedTrace))

	client := &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyURL(proxyURL),
			DisableKeepAlives: true,
		},
		Timeout: time.Second * time.Duration(pc.ipCheckTimeout),
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusNoContent, nil
}

func (pc *ProxyChecker) ClearMetrics() {
	pc.currentMetrics.Range(func(key, _ interface{}) bool {
		metricKey := key.(string)
		parts := strings.Split(metricKey, "|")
		if len(parts) == 3 {
			metrics.DeleteProxyStatus(parts[0], parts[1], parts[2])
			metrics.DeleteProxyLatency(parts[0], parts[1], parts[2])
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

func (pc *ProxyChecker) GetProxyStatus(name string) (bool, time.Duration, error) {
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
		return false, 0, fmt.Errorf("proxy not found")
	}

	status, ok := pc.currentMetrics.Load(metricKey)
	if !ok {
		return false, 0, fmt.Errorf("metric not found")
	}

	latency, _ := pc.latencyMetrics.Load(metricKey)
	if latency == nil {
		latency = time.Duration(0)
	}

	return status.(bool), latency.(time.Duration), nil
}

func (pc *ProxyChecker) GetProxies() []*models.ProxyConfig {
	return pc.proxies
}
