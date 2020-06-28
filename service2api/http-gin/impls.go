package gin

import "text/template"

var implTemplate = template.Must(template.New("impl").Parse(`package svc_{{ .Service }}

var _ {{ .Interface }} = &Service{}

//@autowire({{ .Interface }},set=service)
type Service struct {
}
`))

type Impl struct {
	Service   string
	Interface string
}
