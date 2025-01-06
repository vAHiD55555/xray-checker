package models

import "fmt"

type ProxyConfig struct {
	Protocol    string
	Server      string
	Port        int
	Name        string
	Security    string
	Type        string
	UUID        string
	Flow        string
	HeaderType  string
	Path        string
	Host        string
	SNI         string
	Fingerprint string
	PublicKey   string
	ShortID     string
	Password    string
	Method      string
	Index       int
}

func (pc *ProxyConfig) Validate() error {
	if pc.Protocol == "" {
		return fmt.Errorf("protocol is required")
	}
	if pc.Server == "" {
		return fmt.Errorf("server is required")
	}
	if pc.Port <= 0 || pc.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", pc.Port)
	}
	return nil
}

func (pc *ProxyConfig) GetEndpointPath() string {
	return fmt.Sprintf("%s-%s-%d", pc.Protocol, pc.Server, pc.Port)
}