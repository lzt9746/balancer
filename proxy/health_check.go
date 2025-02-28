// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proxy

import (
	"github.com/zehuamama/balancer/utils"
	"log"
	"time"
)

// ReadAlive reads the alive status of the site
func (h *HTTPProxy) ReadAlive(url string) bool {
	h.RLock()
	defer h.RUnlock()
	return h.alive[url]
}

// SetAlive sets the alive status to the site
func (h *HTTPProxy) SetAlive(url string, alive bool) {
	h.Lock()
	defer h.Unlock()
	h.alive[url] = alive
}

// HealthCheck enable a health check goroutine for each agent
func (h *HTTPProxy) HealthCheck(interval uint) {
	for host := range h.hostMap {
		go h.healthCheck(host, interval)
	}
}

func (h *HTTPProxy) healthCheck(host string, interval uint) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for range ticker.C {
		if !utils.IsBackendAlive(host) && h.ReadAlive(host) {
			log.Printf("Router \"%s\" unreachable, remove %s from load balancer.", h.name, host)

			h.SetAlive(host, false)
			h.lb.Remove(host)
		} else if utils.IsBackendAlive(host) && !h.ReadAlive(host) {
			log.Printf("Router \"%s\" reachable, add %s to load balancer.", h.name, host)

			h.SetAlive(host, true)
			h.lb.Add(host)
		}
	}

}
