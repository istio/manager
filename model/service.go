// Copyright 2016 Google Inc.
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

package model

import (
	"bytes"
	"sort"
	"strings"
)

// ServiceDiscovery enumerates Istio service instances
type ServiceDiscovery interface {
	// Services list all services and their tags
	Services() []*Service
	// Endpoints retrieves service instances for a service and a tag.
	// An empty tag value is allowed only for a tag-less service.
	Endpoints(s *Service, tag string) []*ServiceInstance
}

// Service describes an Istio service
type Service struct {
	// Name of the service
	Name string `json:"name"`
	// Namespace of the service name, optional
	Namespace string `json:"namespace,omitempty"`
	// Tags is a set of declared tags for the service.
	// An empty set is allowed but tag values must be non-empty strings.
	Tags []string `json:"tags,omitempty"`
}

// Endpoint defines a network endpoint
type Endpoint struct {
	// Address of the endpoint, typically an IP address
	Address string `json:"ip_address,omitempty"`
	// Port on the host address
	Port int `json:"port"`
	// Name of the port classifies ports for a single service
	Name string `json:"name,omitempty"`
	// Protocol for the port: TCP, UDP (default is TCP)
	Protocol Protocol `json:"protocol,omitempty"`
}

// Protocol defines network protocols for ports
type Protocol string

const (
	ProtocolHTTP Protocol = "HTTP"
	ProtocolTCP  Protocol = "TCP"
	ProtocolUDP  Protocol = "UDP"
)

// ServiceInstance binds an endpoint to a service and a tag.
// If the service has no tags, the tag value is an empty string;
// otherwise, the tag value is an element in the set of service tags.
type ServiceInstance struct {
	Endpoint Endpoint `json:"endpoint,omitempty"`
	Service  *Service `json:"service,omitempty"`
	Tag      string   `json:"tag,omitempty"`
}

func (s *Service) String() string {
	// example: name.namespace:my-v1,prod
	var buffer bytes.Buffer
	buffer.WriteString(s.Name)
	if len(s.Namespace) > 0 {
		buffer.WriteString(".")
		buffer.WriteString(s.Namespace)
	}
	n := len(s.Tags)
	if n > 0 {
		buffer.WriteString(":")
		tags := make([]string, n)
		copy(tags, s.Tags)
		sort.Strings(tags)
		for i := 0; i < n; i++ {
			if i > 0 {
				buffer.WriteString(",")
			}
			buffer.WriteString(tags[i])
		}
	}
	return buffer.String()
}

// ParseServiceString is the inverse of the Service.String() method
func ParseServiceString(s string) *Service {
	var tags []string
	sep := strings.Index(s, ":")
	if sep < 0 {
		sep = len(s)
	} else {
		tags = strings.Split(s[sep+1:], ",")
	}

	var name, namespace string
	dot := strings.Index(s[:sep], ".")
	if dot < 0 {
		name = s[:sep]
	} else {
		name = s[:dot]
		namespace = s[dot+1 : sep]
	}

	return &Service{
		Name:      name,
		Namespace: namespace,
		Tags:      tags,
	}
}
