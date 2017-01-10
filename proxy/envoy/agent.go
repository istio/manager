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

package envoy

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/golang/glog"
)

// Agent manages a proxy instance
type Agent interface {
	Reload(config *Config) error
}

type agent struct {
	// Envoy binary path
	binary string
	// Map of known running Envoy processes and their restart epochs.
	epoch  int
	cmdMap map[*exec.Cmd]int
	mutex  sync.Mutex
}

const (
	EnvoyConfigPath   = "/etc/envoy/"
	EnvoyFileTemplate = "envoy-rev%d.json"
)

// NewAgent creates a new instance.
func NewAgent(binary string) Agent {
	return &agent{
		binary: binary,
		cmdMap: make(map[*exec.Cmd]int),
	}
}

func configFile(epoch int) string {
	return fmt.Sprintf(EnvoyConfigPath+EnvoyFileTemplate, epoch)
}

// Reload Envoy with a hot restart. Envoy hot restarts are performed by launching a new Envoy process with an
// incremented restart epoch. To successfully launch a new Envoy process that will replace the running Envoy processes,
// the restart epoch of the new process must be exactly 1 greater than the highest restart epoch of the currently
// running Envoy processes. To ensure that we launch the new Envoy process with the correct restart epoch, we keep track
// of all running Envoy processes and their restart epochs.
//
// Envoy hot restart documentation: https://lyft.github.io/envoy/docs/intro/arch_overview/hot_restart.html
func (s *agent) Reload(config *Config) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Write config file
	fname := configFile(s.epoch)
	if err := config.WriteFile(fname); err != nil {
		return err
	}

	// Spin up a new Envoy process.
	cmd := exec.Command(s.binary, "-c", fname, "--restart-epoch", fmt.Sprint(s.epoch))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	s.cmdMap[cmd] = s.epoch
	s.epoch++

	// Start tracking the process.
	go s.waitForExit(cmd)

	return nil
}

// waitForExit waits until the command exits and removes it from the set of known running Envoy processes.
func (s *agent) waitForExit(cmd *exec.Cmd) {
	if err := cmd.Wait(); err != nil {
		glog.V(2).Infof("Envoy terminated: %v", err.Error())
	} else {
		glog.V(2).Infof("Envoy process exited")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// delete config file
	epoch := s.cmdMap[cmd]
	path := configFile(epoch)
	if err := os.Remove(path); err != nil {
		glog.Warningf("Failed to delete config file %s, %v", path, err)
	}
	delete(s.cmdMap, cmd)
}
