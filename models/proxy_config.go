package models

import "fmt"

type ProxyConfig struct {
	Protocol      string            // vless, vmess, trojan, shadowsocks
	Server        string            // server address
	Port          int               // server port
	Name          string            // display name
	Security      string            // tls, reality, none
	Type          string            // tcp, ws, grpc, quic
	UUID          string            // for vless/vmess
	Flow          string            // for vless
	HeaderType    string            // http, none
	Path          string            // for websocket/grpc
	Host          string            // for websocket
	SNI           string            // for tls/reality
	Fingerprint   string            // for tls/reality
	PublicKey     string            // for reality
	ShortID       string            // for reality
	Password      string            // for trojan/shadowsocks
	Method        string            // for shadowsocks
	Level         int               // user level
	AlterId       int               // for vmess
	VMessAid      int               // alternative to AlterId for vmess
	MultiMode     bool              // for grpc
	ServiceName   string            // for grpc
	IdleTimeout   int               // for grpc
	WindowsSize   int               // for grpc
	AllowInsecure bool              // for tls
	ALPN          []string          // for tls
	Index         int               // internal index for port allocation
	Settings      map[string]string // additional settings
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

	// Protocol-specific validation
	switch pc.Protocol {
	case "vless", "vmess":
		if pc.UUID == "" {
			return fmt.Errorf("UUID is required for %s", pc.Protocol)
		}
	case "trojan":
		if pc.Password == "" {
			return fmt.Errorf("password is required for Trojan")
		}
	case "shadowsocks":
		if pc.Password == "" || pc.Method == "" {
			return fmt.Errorf("password and method are required for Shadowsocks")
		}
	default:
		return fmt.Errorf("unsupported protocol: %s", pc.Protocol)
	}

	return nil
}

// GetEndpointPath returns the unique identifier for this proxy
func (pc *ProxyConfig) GetEndpointPath() string {
	return fmt.Sprintf("%s-%s-%d", pc.Protocol, pc.Server, pc.Port)
}

// GetTransportType returns the effective transport type
func (pc *ProxyConfig) GetTransportType() string {
	if pc.Type == "" {
		return "tcp"
	}
	return pc.Type
}

// GetSecurityType returns the effective security type
func (pc *ProxyConfig) GetSecurityType() string {
	if pc.Security == "" {
		return "none"
	}
	return pc.Security
}

// GetAlterId returns the effective alterId for VMess
func (pc *ProxyConfig) GetAlterId() int {
	if pc.AlterId == 0 {
		return pc.VMessAid
	}
	return pc.AlterId
}

// GetVMessSecurity returns the VMess security method
func (pc *ProxyConfig) GetVMessSecurity() string {
	if pc.Security == "" {
		return "auto"
	}
	return pc.Security
}

// GetUserLevel returns the user level
func (pc *ProxyConfig) GetUserLevel() int {
	if pc.Level == 0 {
		return 0
	}
	return pc.Level
}

// HasGRPCSettings returns true if gRPC settings are configured
func (pc *ProxyConfig) HasGRPCSettings() bool {
	return pc.Type == "grpc" && pc.ServiceName != ""
}

// GetServiceName returns the gRPC service name or default value
func (pc *ProxyConfig) GetServiceName() string {
	if pc.ServiceName == "" {
		return "GunService"
	}
	return pc.ServiceName
}

// GetALPNSettings returns the ALPN settings or default values
func (pc *ProxyConfig) GetALPNSettings() []string {
	if len(pc.ALPN) == 0 {
		return []string{"h2", "http/1.1"}
	}
	return pc.ALPN
}
