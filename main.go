// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"github.com/zehuamama/balancer/proxy"
	"github.com/zehuamama/balancer/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func main() {
	config, err := ReadConfig("config.yaml")
	if err != nil {
		log.Fatalf("read config error: %s", err)
	}

	err = config.Validation()
	if err != nil {
		log.Fatalf("verify config error: %s", err)
	}

	//router := mux.NewRouter()
	rt := router.NewRouter()
	for _, l := range config.Location {
		rt.AddRouter(l.Listen)
		rt.AddSchemaMapping(l.Listen, l.Schema)
		if l.Schema == "https" {
			rt.AddCertificateMapping(l.Listen, l.SSLCertificate, l.SSLCertificateKey)
		}

		r := rt.Mapping[l.Listen]
		httpProxy, err := proxy.NewHTTPProxy(l.ProxyPass, l.BalanceMode, l.Name)
		if err != nil {
			log.Fatalf("create proxy error: %s", err)
		}
		// start health check
		if config.HealthCheck {
			httpProxy.HealthCheck(config.HealthCheckInterval)
		}
		r.Handle(l.Pattern, httpProxy)
		r.Use(setHeaderMiddleware(l.SetHeader))
		if config.MaxAllowed > 0 {
			r.Use(maxAllowedMiddleware(config.MaxAllowed))
		}
	}
	for p, r := range rt.Mapping {
		svr := http.Server{
			Addr:    ":" + strconv.Itoa(p),
			Handler: r,
		}
		rt.ServerMapping[p] = &svr
		if rt.SchemaMapping[p] == "http" {
			go func() {
				err := svr.ListenAndServe()
				if err != nil {
					log.Fatalf("listen and serve error: %s", err)
				}
			}()
		} else if rt.SchemaMapping[p] == "https" {
			go func() {
				err := svr.ListenAndServeTLS(rt.CertificateMapping[p].SSLCertificate, rt.CertificateMapping[p].SSLCertificateKey)
				if err != nil {
					log.Fatalf("listen and serve error: %s", err)
				}
			}()
		}
	}
	// print config detail
	config.Print()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	for _, svr := range rt.ServerMapping {
		var wait time.Duration
		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()
		svr.Shutdown(ctx)
	}

	log.Println("shutting down")
	os.Exit(0)

}
