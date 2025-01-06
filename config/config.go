package config

import (
	"fmt"

	"github.com/alecthomas/kong"
)

var CLIConfig CLI

func Parse(version, commit string) {
	ctx := kong.Parse(&CLIConfig,
		kong.Name("xray-checker"),
		kong.Description("Xray Checker: A Prometheus exporter for monitoring Xray proxies"),
		kong.Vars{
			"version": version,
			"commit":  commit,
		},
	)
	_ = ctx
}

type CLI struct {
	SubscriptionURL  string      `name:"subscription-url" help:"URL of the subscription" required:"true" env:"SUBSCRIPTION_URL"`
	RecheckSubscription bool     `name:"recheck-subscription" help:"Whether to recheck the subscription" default:"true" env:"RECHECK_SUBSCRIPTION"`
	CheckInterval    int         `name:"check-interval" help:"Interval for proxy checks in seconds" default:"300" env:"CHECK_INTERVAL"`
	IPCheckService   string      `name:"ip-check-service" help:"Service URL for IP checking" default:"https://api.ipify.org?format=text" env:"IP_CHECK_SERVICE"`
	GenMethodURL     string      `name:"gen-method-url" help:"Response status generator, used by check-method=gen" default:"http://cp.cloudflare.com/generate_204" env:"GEN_METHOD_URL"`
	ResponseWithLatency bool     `name:"response-with-latency" help:"Whether to add latency to the response" default:"false" env:"RESPONSE_WITH_LATENCY"`
	CheckMethod      string      `name:"check-method" help:"Method for checking proxy, ip or gen" default:"ip" env:"CHECK_METHOD"`
	IpCheckTimeout   int         `name:"ip-check-timeout" help:"Timeout for IP checking in seconds" default:"30" env:"IP_CHECK_TIMEOUT"`
	StartPort        int         `name:"start-port" help:"Start port for proxy configuration" default:"10000" env:"START_PORT"`
	XrayLogLevel     string      `name:"xray-log-level" help:"Xray log level (debug|info|warning|error|none)" default:"none" env:"XRAY_LOG_LEVEL"`
	Port             string      `name:"metrics-port" help:"Port to listen on" default:"2112" env:"METRICS_PORT"`
	ProtectedMetrics bool        `name:"metrics-protected" help:"Whether metrics are protected by basic auth" default:"false" env:"METRICS_PROTECTED"`
	MetricsUsername  string      `name:"metrics-username" help:"Username for metrics if protected by basic auth" default:"metricsUser" env:"METRICS_USERNAME"`
	MetricsPassword  string      `name:"metrics-password" help:"Password for metrics if protected by basic auth" default:"MetricsVeryHardPassword" env:"METRICS_PASSWORD"`
	Version          VersionFlag `name:"version" help:"Print version information and quit"`
}

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println("Xray Checker: A Prometheus exporter for monitoring Xray proxies")
	fmt.Printf("Version:\t %s\n", vars["version"])
	fmt.Printf("Commit:\t %s\n", vars["commit"])
	fmt.Printf("GitHub: https://github.com/kutovoys/xray-checker\n")
	app.Exit(0)
	return nil
}
