package web

import (
	"embed"
	"html/template"
	"io"
)

//go:embed templates/*.html
var content embed.FS

var indexTmpl *template.Template

func init() {
	var err error
	indexTmpl, err = template.ParseFS(content, "templates/*.html")
	if err != nil {
		panic(err)
	}
}

type PageData struct {
	Version         string
	Commit          string
	Port            string
	CheckInterval   int
	IPCheckService  string
	ResponseWithLatency bool
	CheckMethod     string
	GenMethodURL    string
	IpCheckTimeout  int
	RecheckSubscription bool
	StartPort       int
	Endpoints       []EndpointInfo
}

func RenderIndex(w io.Writer, data PageData) error {
	return indexTmpl.Execute(w, data)
}
