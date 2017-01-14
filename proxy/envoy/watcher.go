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

package envoy

import (
	"reflect"
	"time"

	"github.com/golang/glog"
	"istio.io/manager/model"
)

// Watcher observes service registry and triggers a reload on a change
type Watcher interface {
}

type watcher struct {
	agent     Agent
	discovery model.ServiceDiscovery
	mesh      *MeshConfig
}

func NewWatcher(discovery model.ServiceDiscovery, ctl model.Controller, mesh *MeshConfig, binary string) Watcher {
	out := &watcher{
		agent:     NewAgent(binary),
		discovery: discovery,
		mesh:      mesh,
	}
	ctl.AppendServiceHandler(out.notify)
	return out
}

func (w *watcher) notify(svc *model.Service, ev model.Event) {
	config, err := Generate(nil, w.discovery.Services(), w.mesh)
	if err != nil {
		glog.Warningf("Failed to generate Envoy configuration: %v", err)
		return
	}

	current := w.agent.ActiveConfig()

	if reflect.DeepEqual(config, current) {
		glog.V(2).Info("Configuration is identical, skipping reload")
		return
	}

	// TODO: add retry logic
	if err := w.agent.Reload(config); err != nil {
		glog.Warningf("Envoy reload error: %v", err)
		return
	}

	// Add a short delay to de-risk a potential race condition in envoy hot reload code.
	// The condition occurs when the active Envoy instance terminates in the middle of
	// the Reload() function.
	time.Sleep(256 * time.Millisecond)
}
