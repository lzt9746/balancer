package router

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Router struct {
	Mapping            map[int]*mux.Router
	SchemaMapping      map[int]string
	CertificateMapping map[int]*Certificate
	ServerMapping      map[int]*http.Server
}

type Certificate struct {
	SSLCertificate    string `yaml:"ssl_certificate"`
	SSLCertificateKey string `yaml:"ssl_certificate_key"`
}

func NewRouter() *Router {
	return &Router{
		Mapping:            make(map[int]*mux.Router),
		SchemaMapping:      make(map[int]string),
		CertificateMapping: make(map[int]*Certificate),
		ServerMapping:      make(map[int]*http.Server),
	}
}

func (r *Router) AddRouter(port int) {
	if r.Mapping[port] == nil {
		r.Mapping[port] = mux.NewRouter()
	}
}

func (r *Router) AddSchemaMapping(port int, schema string) {
	if r.checkSchema(port, schema) {
		r.SchemaMapping[port] = schema
	} else {
		panic("schema conflict on port {port}")
	}
	if schema != "http" && schema != "https" {
		panic("schema error")
	}
	if schema == "https" {

	}
}

func (r *Router) AddCertificateMapping(port int, crt, key string) {
	r.CertificateMapping[port] = &Certificate{
		SSLCertificate:    crt,
		SSLCertificateKey: key,
	}
}

func (r *Router) checkSchema(port int, schema string) bool {
	if v, ok := r.SchemaMapping[port]; ok {
		if v == schema {
			return true
		}
		return false
	}
	return true
}
