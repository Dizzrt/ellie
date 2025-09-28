package main

import (
	"bytes"
	_ "embed"
	"html/template"
	"strings"
)

//go:embed http.tpl
var httpTmpl string

type serviceDesc struct {
	ServiceType string
	ServiceName string
	Metadata    string
	Methods     []*methodDesc
	MethodSets  map[string]*methodDesc
}

type methodDesc struct {
	Name         string
	OriginalName string
	Num          int
	Request      string
	Response     string
	Comment      string
	Path         string
	Method       string
	HasBody      bool
	Body         string
	HasVars      bool
}

func (sd *serviceDesc) excute() string {
	sd.MethodSets = make(map[string]*methodDesc)
	for _, m := range sd.Methods {
		sd.MethodSets[m.Name] = m
	}

	buf := new(bytes.Buffer)
	tmpl, err := template.New("http").Parse(strings.TrimSpace(httpTmpl))
	if err != nil {
		panic(err)
	}

	if err := tmpl.Execute(buf, sd); err != nil {
		panic(err)
	}

	return strings.Trim(buf.String(), "\n")
}
