// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	gops "github.com/google/gops/agent"

	"github.com/cilium/cilium/cilium-cli/cli"
	_ "github.com/cilium/cilium/cilium-cli/logging" // necessary to disable unwanted cfssl log messages
)

func init() {
	os.Setenv("KUBECONFIG", filepath.Join(os.Getenv("HOME"), ".kube", "cilium"))
	os.Args = append(os.Args, []string{

		"install",
		"--version=v1.19.0-pre.0",
		"--namespace=cilium-system",
		"--set", "routing-mode=native",
		"--set", "ipv4NativeRoutingCIDR=10.0.0.0/8",
		"--set", "kubeProxyReplacement=true",
		"--set", "bpf.masquerade=true",
		"--set", "nodePort.enabled=true",
		"--set", "bandwidthManager.enabled=true",
		"--set", "bandwidthManager.bbr=true",
		"--set", "bandwidthManager.bbrHostNamespaceOnly=true",
		"--set", "encryption.enabled=true",
		"--set", "encryption.type=wireguard",
		"--set", "encryption.nodeEncryption=true",
		"--set", "ingressController.enabled=true",
		"--set", "ingressController.loadbalancerMode=dedicated",
		"--set", "localRedirectPolicies.enabled=true",
		"--set", "image.pullPolicy=IfNotPresent --set monitorAggregation=none",
		"--set", "debug.enabled=true",
		"--set", "debug.verbose=flow agent envoy daemon monitor kvstore ipam config datapath",
		"--set", "hubble.metrics.enabled={dns,drop,tcp,flow,port-distribution,icmp,http}",
		"--set", "monitor.enabled=true",
		"--set", "hubble.enabled=true",
		"--set", "hubble.relay.enabled=true",
		"--set", "hubble.relay.image.repository=registry.cn-hangzhou.aliyuncs.com/acejilam/hubble-relay",
		"--set", "hubble.ui.enabled=true",
		"--set", "hubble.ui.frontend.image.repository=registry.cn-hangzhou.aliyuncs.com/acejilam/hubble-ui",
		"--set", "hubble.ui.backend.image.repository=registry.cn-hangzhou.aliyuncs.com/acejilam/hubble-ui-backend",
		"--set", "image.repository=registry.cn-hangzhou.aliyuncs.com/acejilam/cilium",
		"--set", "envoy.image.repository=registry.cn-hangzhou.aliyuncs.com/acejilam/cilium-envoy",
		"--set", "preflight.image.repository=registry.cn-hangzhou.aliyuncs.com/acejilam/cilium-ci",
		"--set", "preflight.envoy.image.repository=registry.cn-hangzhou.aliyuncs.com/acejilam/cilium-envoy",
		"--set", "operator.image.repository=registry.cn-hangzhou.aliyuncs.com/acejilam/operator",
	}...)
}

func main() {
	if err := gops.Listen(gops.Options{}); err != nil {
		log.Printf("Unable to start gops: %s", err)
	}

	if err := cli.NewDefaultCiliumCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
