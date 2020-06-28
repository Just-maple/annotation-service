package http_gin

import "text/template"

var ginSvcTemplate = template.Must(template.New("svc").Parse(`// Code generated by service2api. DO NOT EDIT.
package {{.Package}}

import (
	svrlessgin "github.com/Just-maple/serverless-gin"
	"github.com/gin-gonic/gin"
)

func Register{{ .GroupName }}Group(svc {{ .ServiceName }}, router gin.IRouter, svcH svrlessgin.GinSvcHandler) {
	router = router.Group({{ .GroupRoute }})
	
	{{ range .Apis }}
	// {{ .Title }}
	router.{{ .HttpMethod }}({{ .Route }}, svcH(svc.{{ .Handler }})){{ end }}
}
`))

type GinRoute struct {
	Package     string
	GroupName   string
	ServiceName string
	GroupRoute  string
	Apis        []GinApi
}

type GinApi struct {
	HttpMethod string
	Route      string
	Handler    string
	Title      string
}
