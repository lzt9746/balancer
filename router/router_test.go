package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
	"testing"
)

func TestNewRouter(t *testing.T) {
	tests := []struct {
		name string
		want *Router
	}{
		{
			name: "TestNewRouter",
			want: &Router{
				Mapping:            make(map[int]*mux.Router),
				SchemaMapping:      make(map[int]string),
				CertificateMapping: make(map[int]*Certificate),
				ServerMapping:      make(map[int]*http.Server),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRouter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouter_AddCertificateMapping(t *testing.T) {
	type fields struct {
		Mapping            map[int]*mux.Router
		SchemaMapping      map[int]string
		CertificateMapping map[int]*Certificate
		ServerMapping      map[int]*http.Server
	}
	type args struct {
		port int
		crt  string
		key  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "TestAddCertificateMapping",
			fields: fields{
				Mapping:            make(map[int]*mux.Router),
				SchemaMapping:      make(map[int]string),
				CertificateMapping: make(map[int]*Certificate),
				ServerMapping:      make(map[int]*http.Server),
			},
			args: args{
				port: 8080,
				crt:  "test.crt",
				key:  "test.key",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Router{
				Mapping:            tt.fields.Mapping,
				SchemaMapping:      tt.fields.SchemaMapping,
				CertificateMapping: tt.fields.CertificateMapping,
				ServerMapping:      tt.fields.ServerMapping,
			}
			r.AddCertificateMapping(tt.args.port, tt.args.crt, tt.args.key)
		})
	}
}

func TestRouter_AddRouter(t *testing.T) {
	type fields struct {
		Mapping            map[int]*mux.Router
		SchemaMapping      map[int]string
		CertificateMapping map[int]*Certificate
		ServerMapping      map[int]*http.Server
	}
	type args struct {
		port int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Router{
				Mapping:            tt.fields.Mapping,
				SchemaMapping:      tt.fields.SchemaMapping,
				CertificateMapping: tt.fields.CertificateMapping,
				ServerMapping:      tt.fields.ServerMapping,
			}
			r.AddRouter(tt.args.port)
		})
	}
}

func TestRouter_AddSchemaMapping(t *testing.T) {
	type fields struct {
		Mapping            map[int]*mux.Router
		SchemaMapping      map[int]string
		CertificateMapping map[int]*Certificate
		ServerMapping      map[int]*http.Server
	}
	type args struct {
		port   int
		schema string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Router{
				Mapping:            tt.fields.Mapping,
				SchemaMapping:      tt.fields.SchemaMapping,
				CertificateMapping: tt.fields.CertificateMapping,
				ServerMapping:      tt.fields.ServerMapping,
			}
			r.AddSchemaMapping(tt.args.port, tt.args.schema)
		})
	}
}

func TestRouter_checkSchema(t *testing.T) {
	type fields struct {
		Mapping            map[int]*mux.Router
		SchemaMapping      map[int]string
		CertificateMapping map[int]*Certificate
		ServerMapping      map[int]*http.Server
	}
	type args struct {
		port   int
		schema string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Router{
				Mapping:            tt.fields.Mapping,
				SchemaMapping:      tt.fields.SchemaMapping,
				CertificateMapping: tt.fields.CertificateMapping,
				ServerMapping:      tt.fields.ServerMapping,
			}
			if got := r.checkSchema(tt.args.port, tt.args.schema); got != tt.want {
				t.Errorf("checkSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}
