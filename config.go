// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

var (
	ascii = `
___ _ _  _ _   _ ___  ____ _    ____ _  _ ____ ____ ____ 
 |  | |\ |  \_/  |__] |__| |    |__| |\ | |    |___ |__/ 
 |  | | \|   |   |__] |  | |___ |  | | \| |___ |___ |  \                                        
`
)

// Config configuration details of balancer
type Config struct {
	Location            []*Location `yaml:"location"`
	HealthCheck         bool        `yaml:"tcp_health_check"`
	HealthCheckInterval uint        `yaml:"health_check_interval"`
	MaxAllowed          uint        `yaml:"max_allowed"`
}

type Header struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// Location routing details of balancer
type Location struct {
	Name              string    `yaml:"name"`
	Listen            int       `yaml:"listen"`
	Schema            string    `yaml:"schema"`
	SSLCertificate    string    `yaml:"ssl_certificate"`
	SSLCertificateKey string    `yaml:"ssl_certificate_key"`
	Pattern           string    `yaml:"pattern"`
	ProxyPass         []string  `yaml:"proxy_pass"`
	BalanceMode       string    `yaml:"balance_mode"`
	SetHeader         []*Header `yaml:"set_header"`
}

// ReadConfig read configuration from `fileName` file
func ReadConfig(fileName string) (*Config, error) {
	in, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(in, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Print print config details
func (c *Config) Print() {
	fmt.Printf("%s\nHealth Check: %v\nLocation:\n",
		ascii, c.HealthCheck)
	for _, l := range c.Location {
		fmt.Printf("\tRoute: %s\n\tName: %s\n\tSchema:%s\n\tListen:%d\n\tProxy Pass: %s\n\tMode: %s\n\n",
			l.Pattern, l.Name, l.Schema, l.Listen, l.ProxyPass, l.BalanceMode)
	}
}

// Validation verify the configuration details of the balancer
func (c *Config) Validation() error {
	for _, l := range c.Location {
		if l.Schema != "http" && l.Schema != "https" {
			return fmt.Errorf("the schema \"%s\" not supported", l.Schema)
		}
		if l.Schema == "https" && (len(l.SSLCertificate) == 0 || len(l.SSLCertificateKey) == 0) {
			return errors.New("the https proxy requires ssl_certificate_key and ssl_certificate")
		}
	}

	if len(c.Location) == 0 {
		return errors.New("the details of location cannot be null")
	}

	if c.HealthCheckInterval < 1 {
		return errors.New("health_check_interval must be greater than 0")
	}
	return nil
}
