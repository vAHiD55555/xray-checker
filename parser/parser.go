package parser

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"xray-checker/config"
	"xray-checker/models"
	"xray-checker/xray"
)

func InitializeConfiguration(configFile string) (*[]*models.ProxyConfig, error) {
	configs, err := ParseSubscription(config.CLIConfig.Subscription.URL)
	if err != nil {
		return nil, fmt.Errorf("error parsing subscription: %v", err)
	}
	proxyConfigs := &configs

	xray.PrepareProxyConfigs(*proxyConfigs)
	if err := xray.GenerateAndSaveConfig(*proxyConfigs, config.CLIConfig.Xray.StartPort, configFile, config.CLIConfig.Xray.LogLevel); err != nil {
		return nil, fmt.Errorf("error generating Xray config: %v", err)
	}

	return proxyConfigs, nil
}

func ParseSubscriptionURL(subscriptionURL string) ([]string, error) {
	if _, err := url.Parse(subscriptionURL); err != nil {
		return nil, fmt.Errorf("error parsing URL: %v", err)
	}

	req, err := http.NewRequest("GET", subscriptionURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("User-Agent", "Xray-Checker")
	req.Header.Set("Accept", "*/*")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting subscription: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		links := strings.Split(string(body), "\n")
		_, err = ParseProxyURL(links[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse config: %v", err)
		}
		return filterEmptyLinks(links), nil
	}

	links := strings.Split(string(decoded), "\n")
	return filterEmptyLinks(links), nil
}

func filterEmptyLinks(links []string) []string {
	var filtered []string
	for _, link := range links {
		if strings.TrimSpace(link) != "" {
			filtered = append(filtered, link)
		}
	}
	return filtered
}

func ParseProxyURL(proxyURL string) (*models.ProxyConfig, error) {
	u, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing proxy URL: %v", err)
	}

	config := &models.ProxyConfig{}
	config.Protocol = u.Scheme

	switch config.Protocol {
	case "vless":
		return ParseVLESSConfig(u)
	case "trojan":
		return ParseTrojanConfig(u)
	case "ss":
		return ParseShadowsocksConfig(u)
	default:
		return nil, fmt.Errorf("unsupported protocol: %s", config.Protocol)
	}
}

func ParseVLESSConfig(u *url.URL) (*models.ProxyConfig, error) {
	config := &models.ProxyConfig{
		Protocol: "vless",
		Name:     strings.TrimPrefix(u.Fragment, ""),
	}

	config.UUID = u.User.Username()

	hostParts := strings.Split(u.Host, ":")
	if len(hostParts) != 2 {
		return nil, fmt.Errorf("invalid server address format: %s", u.Host)
	}
	config.Server = hostParts[0]
	fmt.Sscanf(hostParts[1], "%d", &config.Port)
	if config.Port == 0 || config.Port == 1 {
		return nil, fmt.Errorf("skipping port: %d", config.Port)
	}

	query := u.Query()

	config.Security = query.Get("security")
	config.Type = query.Get("type")
	config.HeaderType = query.Get("headerType")
	config.Flow = query.Get("flow")
	config.Path = query.Get("path")
	config.Host = query.Get("host")
	config.SNI = query.Get("sni")
	config.Fingerprint = query.Get("fp")
	config.PublicKey = query.Get("pbk")
	config.ShortID = query.Get("sid")

	return config, nil
}

func ParseTrojanConfig(u *url.URL) (*models.ProxyConfig, error) {
	config := &models.ProxyConfig{
		Protocol: "trojan",
		Name:     strings.TrimPrefix(u.Fragment, ""),
	}

	config.Password = u.User.Username()

	hostParts := strings.Split(u.Host, ":")
	if len(hostParts) != 2 {
		return nil, fmt.Errorf("invalid server address format: %s", u.Host)
	}
	config.Server = hostParts[0]
	fmt.Sscanf(hostParts[1], "%d", &config.Port)
	if config.Port == 0 || config.Port == 1 {
		return nil, fmt.Errorf("skipping port: %d", config.Port)
	}

	query := u.Query()

	config.Security = query.Get("security")
	config.Type = query.Get("type")
	config.HeaderType = query.Get("headerType")
	config.Path = query.Get("path")
	config.Host = query.Get("host")
	config.SNI = query.Get("sni")
	config.Fingerprint = query.Get("fp")

	return config, nil
}

func ParseShadowsocksConfig(u *url.URL) (*models.ProxyConfig, error) {
	config := &models.ProxyConfig{
		Protocol: "shadowsocks",
		Name:     strings.TrimPrefix(u.Fragment, ""),
	}

	methodPass, err := base64.StdEncoding.DecodeString(u.User.String())
	if err != nil {
		return nil, fmt.Errorf("error decoding method and password: %v", err)
	}

	parts := strings.SplitN(string(methodPass), ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid method:password format")
	}
	config.Method = parts[0]
	config.Password = parts[1]

	hostParts := strings.Split(u.Host, ":")
	if len(hostParts) != 2 {
		return nil, fmt.Errorf("invalid server address format: %s", u.Host)
	}
	config.Server = hostParts[0]
	fmt.Sscanf(hostParts[1], "%d", &config.Port)
	if config.Port == 0 || config.Port == 1 {
		return nil, fmt.Errorf("skipping port: %d", config.Port)
	}

	return config, nil
}

func ParseSubscription(subscriptionURL string) ([]*models.ProxyConfig, error) {
	links, err := ParseSubscriptionURL(subscriptionURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing subscription URL: %v", err)
	}

	var configs []*models.ProxyConfig
	for _, link := range links {
		config, err := ParseProxyURL(link)
		if err != nil {
			log.Printf("Warning: error parsing proxy URL %s: %v", link, err)
			continue
		}
		configs = append(configs, config)
	}

	return configs, nil
}
