// Copyright 2017 Google Inc.
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
	"fmt"
	"regexp"

	"github.com/golang/protobuf/proto"

	multierror "github.com/hashicorp/go-multierror"
)

const (
	dns1123LabelMaxLength int    = 63
	dns1123LabelFmt       string = "[a-z0-9]([-a-z0-9]*[a-z0-9])?"
)

var (
	dns1123LabelRex = regexp.MustCompile("^" + dns1123LabelFmt + "$")
	kindRegexp      = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]*$")
)

// IsDNS1123Label tests for a string that conforms to the definition of a label in
// DNS (RFC 1123).
func IsDNS1123Label(value string) bool {
	return len(value) <= dns1123LabelMaxLength && dns1123LabelRex.MatchString(value)
}

// Validate confirms that the names in the configuration key are appropriate
func (k *ConfigKey) Validate() error {
	var errs error
	if !kindRegexp.MatchString(k.Kind) {
		errs = multierror.Append(errs, fmt.Errorf("Invalid kind: %q", k.Kind))
	}
	if !IsDNS1123Label(k.Name) {
		errs = multierror.Append(errs, fmt.Errorf("Invalid name: %q", k.Name))
	}
	if !IsDNS1123Label(k.Namespace) {
		errs = multierror.Append(errs, fmt.Errorf("Invalid namespace: %q", k.Namespace))
	}
	return errs
}

// Validate checks that each name conforms to the spec and has a ProtoMessage
func (km KindMap) Validate() error {
	var errs error
	for k, v := range km {
		if !kindRegexp.MatchString(k) {
			errs = multierror.Append(errs, fmt.Errorf("Invalid kind: %q", k))
		}
		if proto.MessageType(v.MessageName) == nil {
			errs = multierror.Append(errs, fmt.Errorf("Cannot find proto message type: %q", v.MessageName))
		}
	}
	return errs
}

// ValidateKey ensures that the key is well-defined and kind is well-defined
func (km KindMap) ValidateKey(k *ConfigKey) error {
	if err := k.Validate(); err != nil {
		return err
	}
	if _, ok := km[k.Kind]; !ok {
		return fmt.Errorf("Kind %q is not defined", k.Kind)
	}
	return nil
}

// ValidateConfig ensures that the config object is well-defined
func (km KindMap) ValidateConfig(obj *Config) error {
	if obj == nil {
		return fmt.Errorf("Invalid nil configuration object")
	}

	if err := obj.ConfigKey.Validate(); err != nil {
		return err
	}
	t, ok := km[obj.Kind]
	if !ok {
		return fmt.Errorf("Undeclared kind: %q", obj.Kind)
	}

	// Validate spec field
	if obj.Spec == nil {
		return fmt.Errorf("Want a proto message, received empty content")
	}
	v, ok := obj.Spec.(proto.Message)
	if !ok {
		return fmt.Errorf("Cannot cast spec to a proto message")
	}
	if proto.MessageName(v) != t.MessageName {
		return fmt.Errorf("Mismatched spec message type %q and kind %q",
			proto.MessageName(v), t.MessageName)
	}
	if err := t.Validate(v); err != nil {
		return err
	}

	// Validate status field
	if obj.Status != nil {
		v, ok := obj.Status.(proto.Message)
		if !ok {
			return fmt.Errorf("Cannot cast status to a proto message")
		}
		if proto.MessageName(v) != t.StatusMessageName {
			return fmt.Errorf("Mismatched status message type %q and kind %q",
				proto.MessageName(v), t.StatusMessageName)
		}
	}

	return nil
}

// Validate ensures that the service object is well-defined
func (s *Service) Validate() error {
	var errs error
	if !IsDNS1123Label(s.Name) {
		errs = multierror.Append(errs, fmt.Errorf("Invalid name: %q", s.Name))
	}
	if s.Namespace != "" && !IsDNS1123Label(s.Namespace) {
		errs = multierror.Append(errs, fmt.Errorf("Invalid namespace: %q", s.Namespace))
	}
	for _, tag := range s.Tags {
		if !IsDNS1123Label(tag) {
			errs = multierror.Append(errs, fmt.Errorf("Invalid service tag: %q", tag))
		}
	}
	// Require at least one port
	if len(s.Ports) == 0 {
		errs = multierror.Append(errs, fmt.Errorf("Service must have at least one declared port"))
	}

	// Port names can be empty if there exists only one port
	for _, port := range s.Ports {
		if port.Name == "" {
			if len(s.Ports) > 1 {
				errs = multierror.Append(errs,
					fmt.Errorf("Empty port names are not allowed for services with multiple ports"))
			}
		} else if !IsDNS1123Label(port.Name) {
			errs = multierror.Append(errs, fmt.Errorf("Invalid name: %q", port.Name))
		}
		if port.Port < 0 {
			errs = multierror.Append(errs, fmt.Errorf("Invalid port value %d for %q", port.Port, port.Name))
		}
	}
	return errs
}
