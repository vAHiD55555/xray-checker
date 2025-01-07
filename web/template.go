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
	Version            string
	Commit             string
	Port               string
	CheckInterval      int
	IPCheckUrl         string
	SimulateLatency    bool
	CheckMethod        string
	StatusCheckUrl     string
	Timeout            int
	SubscriptionUpdate bool
	StartPort          int
	Endpoints          []EndpointInfo
}

func RenderIndex(w io.Writer, data PageData) error {
	return indexTmpl.Execute(w, data)
}
