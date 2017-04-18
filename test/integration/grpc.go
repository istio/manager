// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"

	proxyconfig "istio.io/api/proxy/v1/config"
)

type grpc struct {
	*infra
	logs *accessLogs
}

func (t *grpc) String() string {
	return "HTTP/2 reachability"
}

func (t *grpc) setup() error {
	t.logs = makeAccessLogs()
	return nil
}

func (t *grpc) teardown() {
}

func (t *grpc) run() error {
	if err := t.makeRequests(); err != nil {
		return err
	}
	if err := t.logs.check(t.infra); err != nil {
		return err
	}
	return nil
}

func (t *grpc) makeRequests() error {
	testPods := []string{"a", "b"}
	if t.Auth == proxyconfig.ProxyMeshConfig_NONE {
		// t is not behind proxy, so it cannot talk in Istio auth.
		testPods = append(testPods, "t")
	}
	funcs := make(map[string]func() status)
	for _, src := range testPods {
		for _, dst := range testPods {
			for _, port := range []string{":70", ":7070"} {
				for _, domain := range []string{"", "." + t.Namespace} {
					name := fmt.Sprintf("grpc connection from %s to %s%s%s", src, dst, domain, port)
					funcs[name] = (func(src, dst, port, domain string) func() status {
						url := fmt.Sprintf("grpc://%s%s%s", dst, domain, port)
						return func() status {
							resp := t.clientRequest(src, url)
							if len(resp.id) > 0 {
								id := resp.id[0]
								if src != "t" {
									t.logs.add(src, id, name)
								}
								if dst != "t" {
									t.logs.add(dst, id, name)
								}
								// mixer filter is invoked on the server side, that is when dst is not "t"
								if t.Mixer && dst != "t" {
									t.logs.add("mixer", id, name)
								}
								return success
							}
							if src == "t" && dst == "t" {
								// Expected no match for t->t
								return success
							}
							return again
						}
					})(src, dst, port, domain)
				}
			}
		}
	}
	return parallel(funcs)
}
