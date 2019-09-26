package util

import (
	"errors"
	"github.com/subosito/gotenv"
	"net"
	"os"
)

const (
	metadataFile = "/run/metadata/coreos"
)

// From: https://github.com/coreos/container-linux-config-transpiler/blob/master/config/templating/templating.go
var (
	envNamesV4Private []string = []string{
		"COREOS_AZURE_IPV4_DYNAMIC",
		"COREOS_DIGITALOCEAN_IPV4_PRIVATE_0",
		"COREOS_EC2_IPV4_LOCAL",
		"COREOS_GCE_IP_LOCAL_0",
		"COREOS_PACKET_IPV4_PRIVATE_0",
		"COREOS_OPENSTACK_IPV4_LOCAL",
		"COREOS_VAGRANT_VIRTUALBOX_PRIVATE_IPV4",
		"COREOS_CUSTOM_PRIVATE_IPV4",
	}

	envNamesV4Public []string = []string{
		"COREOS_AZURE_IPV4_VIRTUAL",
		"COREOS_DIGITALOCEAN_IPV4_PUBLIC_0",
		"COREOS_EC2_IPV4_PUBLIC",
		"COREOS_GCE_IP_EXTERNAL_0",
		"COREOS_PACKET_IPV6_PUBLIC_0",
		"COREOS_OPENSTACK_IPV4_PUBLIC",
		"COREOS_VAGRANT_VIRTUALBOX_PRIVATE_IPV4",
		"COREOS_CUSTOM_PUBLIC_IPV6",
	}
)

func LookupIpV4Address(public bool) (net.IP, error) {
	gotenv.Load(metadataFile)
	var envNames []string = envNamesV4Private
	if public {
		envNames = envNamesV4Public
	}

	for _, envName := range envNames {
		v := os.Getenv(envName)
		if v != "" {
			return net.ParseIP(v), nil
		}
	}
	return net.IP{}, errors.New("Ip address lookup failed")
}
