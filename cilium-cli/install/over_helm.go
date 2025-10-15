// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

// Copyright The Helm Authors.

package install

import (
	"fmt"

	"github.com/cilium/cilium/cilium-cli/internal/helm"
	"github.com/cilium/cilium/cilium-cli/k8s"
	"github.com/cilium/cilium/pkg/versioncheck"
)

func (k *K8sInstaller) getHelmValues() (map[string]any, error) {
	helmMapOpts := map[string]string{}

	switch {
	case versioncheck.MustCompile(">=1.14.0")(k.chartVersion):
		// TODO(aanm) to keep the previous behavior unchanged we will set the number
		// of the operator replicas to 1. Ideally this should be the default in the helm chart
		helmMapOpts["operator.replicas"] = "1"

		// Set nodeinit enabled option
		if needsNodeInit(k.flavor.Kind) {
			helmMapOpts["nodeinit.enabled"] = "true"
		}

		// Set Helm options specific to the detected Kubernetes cluster type
		switch k.flavor.Kind {
		case k8s.KindKind:
			helmMapOpts["ipam.mode"] = ipamKubernetes
		}

		// Set Helm options specific to the detected / selected datapath mode
		switch k.params.DatapathMode {
		case DatapathTunnel:
			helmMapOpts["routingMode"] = routingModeTunnel
			helmMapOpts["tunnelProtocol"] = tunnelVxlan
		}

		if k.params.ClusterName != "" {
			helmMapOpts["cluster.name"] = k.params.ClusterName
		}

		// TODO: remove when removing "ipv4-native-routing-cidr" flag (marked as
		// deprecated), kept for backwards compatibility
		if k.params.IPv4NativeRoutingCIDR != "" {
			helmMapOpts["ipv4NativeRoutingCIDR"] = k.params.IPv4NativeRoutingCIDR
		}

	default:
		return nil, fmt.Errorf("cilium version unsupported %s", k.chartVersion)
	}

	return helm.MergeVals(k.params.HelmOpts, helmMapOpts)
}
