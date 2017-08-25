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

package kube

import (
	"testing"
)

func TestStr2NamedPort(t *testing.T) {
	var tests = []struct {
		input  string    // input
		expVal NamedPort // output
		expErr bool      // error
	}{
		// Good cases:
		{"http:5555", NamedPort{5555, "http"}, false},
		{"80", NamedPort{80, "http"}, false},
		{"443", NamedPort{443, "https"}, false},
		{"1234", NamedPort{1234, "1234"}, false},
		// Error case:
		{"", NamedPort{0, ""}, true},
		{"foo:bar", NamedPort{0, "foo"}, true},
	}
	for _, tst := range tests {
		actVal, actErr := Str2NamedPort(tst.input)
		if tst.expVal != actVal {
			t.Errorf("Got '%+v', expecting '%+v' for Str2NamedPort('%s')", actVal, tst.expVal, tst.input)
		}
		if tst.expErr {
			if actErr == nil {
				t.Errorf("Got no error when expecting an error for for Str2NamedPort('%s')", tst.input)
			}
		} else {
			if actErr != nil {
				t.Errorf("Got unexpected error '%+v' when expecting none for Str2NamedPort('%s')", actErr, tst.input)
			}
		}
	}
}

func TestSplitEqual(t *testing.T) {
	var tests = []struct {
		input string // input
		expK  string // output1
		expV  string // output2
	}{
		{"foo=bar", "foo", "bar"},
		{"foo=bar=blah", "foo", "bar=blah"},
		{"foo", "foo", ""},
	}
	for _, tst := range tests {
		actK, actV := splitEqual(tst.input)
		if tst.expK != actK {
			t.Errorf("Got key '%+v', expecting '%+v' for splitEqual('%s')", actK, tst.expK, tst.input)
		}
		if tst.expV != actV {
			t.Errorf("Got value '%+v', expecting '%+v' for splitEqual('%s')", actV, tst.expV, tst.input)
		}
	}
}

/*
func TestSamePorts(t *testing.T) {
	var tests = []struct {
		input1 []v1.Endpoints,
		input2 map[int32]bool,
		expected  bool // result
	}{
	}
	for _, tst := range tests {
		actual := samePorts(tst.input1, tst.input2)
		if tst.actual != expected {
			t.Errorf("Got unexpected samePorts(%+v, %+v) = %v", tst.input1, tst.input2, actual)
		}
		}
}
*/
