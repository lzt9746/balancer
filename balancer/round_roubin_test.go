// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRoundRobin_Add .
func TestRoundRobin_Add(t *testing.T) {
	cases := []struct {
		name   string
		lb     Balancer
		args   string
		expect Balancer
	}{
		{
			"test-1",
			NewRoundRobin([]string{"http://127.0.0.1:1011",
				"http://127.0.0.1:1012", "http://127.0.0.1:1013"}, nil),
			"http://127.0.0.1:1013",
			&RoundRobin{hosts: []string{"http://127.0.0.1:1011",
				"http://127.0.0.1:1012", "http://127.0.0.1:1013"}, i: 0},
		},
		{
			"test-2",
			NewRoundRobin([]string{"http://127.0.0.1:1011",
				"http://127.0.0.1:1012"}, nil),
			"http://127.0.0.1:1012",
			&RoundRobin{hosts: []string{"http://127.0.0.1:1011",
				"http://127.0.0.1:1012"}, i: 0},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.lb.Add(c.args)
			assert.Equal(t, c.expect, c.lb)
		})
	}
}

// TestRoundRobin_Remove .
func TestRoundRobin_Remove(t *testing.T) {
	cases := []struct {
		name   string
		lb     Balancer
		args   string
		expect Balancer
	}{
		{
			"test-1",
			NewRoundRobin([]string{"http://127.0.0.1:1011", "http://127.0.0.1:1012"}, nil),
			"http://127.0.0.1:1013",
			&RoundRobin{hosts: []string{"http://127.0.0.1:1011", "http://127.0.0.1:1012"}},
		},
		{
			"test-2",
			NewRoundRobin([]string{"http://127.0.0.1:1011", "http://127.0.0.1:1012"}, nil),
			"http://127.0.0.1:1012",
			&RoundRobin{hosts: []string{"http://127.0.0.1:1011"}},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.lb.Remove(c.args)
			assert.Equal(t, c.expect, c.lb)
		})
	}
}

// TestRoundRobin_Balance .
func TestRoundRobin_Balance(t *testing.T) {
	type expect struct {
		reply string
		err   error
	}
	cases := []struct {
		name   string
		lb     Balancer
		args   string
		expect expect
	}{
		{
			"test-1",
			NewRoundRobin([]string{"http://127.0.0.1:1011"}, nil),
			"",
			expect{
				"http://127.0.0.1:1011",
				nil,
			},
		},
		{
			"test-2",
			NewRoundRobin([]string{}, nil),
			"",
			expect{
				"",
				NoHostError,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reply, err := c.lb.Balance(c.args)
			assert.Equal(t, c.expect.reply, reply)
			assert.Equal(t, c.expect.err, err)
		})
	}
}
