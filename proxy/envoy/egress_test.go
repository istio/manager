package envoy

import (
	"testing"

	"istio.io/manager/test/mock"
)

const (
	egressEnvoyConfig = "testdata/egress-envoy.json"
)

func testEgressConfig(c *EgressConfig, envoyConfig string, t *testing.T) {
	config := generateEgress(c)
	if config == nil {
		t.Fatal("Failed to generate config")
	}

	if err := config.WriteFile(envoyConfig); err != nil {
		t.Fatalf(err.Error())
	}

	compareJSON(envoyConfig, t)
}

func TestEgressRoutes(t *testing.T) {
	r := mock.Discovery
	testEgressConfig(&EgressConfig{
		Services: r,
		Mesh:     DefaultMeshConfig,
	}, egressEnvoyConfig, t)
}
