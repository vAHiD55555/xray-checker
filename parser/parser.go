package parser

import (
	"encoding/base64"
	"encoding/json"
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
	configs, err := ParseSubscription(config.CLIConfig.SubscriptionURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing subscription: %v", err)
	}
	proxyConfigs := &configs

	xray.PrepareProxyConfigs(*proxyConfigs)
	if err := xray.GenerateAndSaveConfig(*proxyConfigs, config.CLIConfig.StartPort, configFile, config.CLIConfig.XrayLogLevel); err != nil {
		return nil, fmt.Errorf("error generating Xray config: %v", err)
	}

	return proxyConfigs, nil
}

func ParseSubscriptionURL(subscriptionURL string) ([]string, error) {
	if !strings.HasSuffix(subscriptionURL, "/info") {
		subscriptionURL = subscriptionURL + "/info"
	}

	resp, err := http.Get(subscriptionURL)
	if err != nil {
		return nil, fmt.Errorf("error getting subscription: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var subResp models.SubscriptionResponse
	if err := json.Unmarshal(body, &subResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return subResp.Links, nil
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
