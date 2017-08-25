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
	"context"
	"fmt"
	"os"

	"github.com/golang/glog"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"

	"istio.io/pilot/cmd"
	"istio.io/pilot/proxy"
	"istio.io/pilot/proxy/envoy"
	"istio.io/pilot/tools/version"
)

// ConsulArgs store the args related to Consul configuration
type ConsulArgs struct {
	config    string
	serverURL string
}

var (
	configpath      string
	meshconfig      string
	role            proxy.Node
	serviceregistry proxy.ServiceRegistry
	consulargs      ConsulArgs

	rootCmd = &cobra.Command{
		Use:   "agent",
		Short: "Istio Pilot agent",
		Long:  "Istio Pilot provides management plane functionality to the Istio service mesh and Istio Mixer.",
	}

	proxyCmd = &cobra.Command{
		Use:   "proxy",
		Short: "Envoy proxy agent",
		RunE: func(c *cobra.Command, args []string) error {
			// set values from environment variables
			if serviceregistry == proxy.KubernetesRegistry {
				if role.IPAddress == "" {
					role.IPAddress = os.Getenv("INSTANCE_IP")
				}
				if role.ID == "" {
					role.ID = os.Getenv("POD_NAME") + "." + os.Getenv("POD_NAMESPACE")
				}
				if role.Domain == "" {
					role.Domain = os.Getenv("POD_NAMESPACE") + ".svc.cluster.local"
				}

				role.Type = proxy.Sidecar
				if len(args) > 0 {
					role.Type = proxy.NodeType(args[0])
				}

				// receive mesh configuration
				mesh, err := cmd.ReadMeshConfig(meshconfig)
				if err != nil {
					return multierror.Prefix(err, "failed to read mesh configuration.")
				}

				glog.V(2).Infof("version %s", version.Line())
				glog.V(2).Infof("mesh configuration %#v", mesh)

				watcher, err := envoy.NewWatcher(mesh, role, configpath)
				if err != nil {
					return err
				}
				ctx, cancel := context.WithCancel(context.Background())
				go watcher.Run(ctx)

				stop := make(chan struct{})
				cmd.WaitSignal(stop)
				<-stop
				cancel()
			} else if serviceregistry == proxy.ConsulRegistry {
				mesh := proxy.DefaultMeshConfig()
				ipAddr := "127.0.0.1"
				available := envoy.WaitForPrivateNetwork()
				if available {
					ipAddr = envoy.GetPrivateIP().String()
					glog.V(2).Infof("obtained private IP %v", ipAddr)
				}

				role.IPAddress = ipAddr
				role.Type = proxy.Sidecar
				if len(args) > 0 {
					role.Type = proxy.NodeType(args[0])
				}

				watcher, err := envoy.NewWatcher(&mesh, role, configpath)
				if err != nil {
					return err
				}
				ctx, cancel := context.WithCancel(context.Background())
				go watcher.Run(ctx)
				stop := make(chan struct{})
				cmd.WaitSignal(stop)
				<-stop
				cancel()
			}
			return nil
		},
	}
)

func init() {
	proxyCmd.PersistentFlags().StringVar((*string)(&serviceregistry), "serviceregistry", string(proxy.KubernetesRegistry),
		fmt.Sprintf("Select the platform for service registry, options are {%s, %s}",
			string(proxy.KubernetesRegistry), string(proxy.ConsulRegistry)))
	proxyCmd.PersistentFlags().StringVar(&meshconfig, "meshconfig", "/etc/istio/config/mesh",
		"File name for Istio mesh configuration")
	proxyCmd.PersistentFlags().StringVar(&configpath, "configpath", "/etc/istio/proxy",
		"Path to generated proxy configuration directory")
	proxyCmd.PersistentFlags().StringVar(&role.IPAddress, "ip", "",
		"Proxy IP address. If not provided uses ${INSTANCE_IP} environment variable.")
	proxyCmd.PersistentFlags().StringVar(&role.ID, "id", "",
		"Proxy unique ID. If not provided uses ${POD_NAME}.${POD_NAMESPACE} from environment variables")
	proxyCmd.PersistentFlags().StringVar(&role.Domain, "domain", "",
		"DNS domain suffix. If not provided uses ${POD_NAMESPACE}.svc.cluster.local")
	proxyCmd.PersistentFlags().StringVar(&consulargs.config, "consulconfig", "",
		"Consul config file for sidecar")

	proxyCmd.PersistentFlags().StringVar(&consulargs.serverURL, "consulserverURL", "",
		"URL for the Consul registry server")

	cmd.AddFlags(rootCmd)

	rootCmd.AddCommand(proxyCmd)
	rootCmd.AddCommand(cmd.VersionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		glog.Error(err)
		os.Exit(-1)
	}
}
