package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/containernetworking/cni/libcni"
)

const (
	EnvCNIPath        = "CNI_PATH"
	EnvNetDir         = "NETCONFPATH"
	EnvCapabilityArgs = "CAP_ARGS"

	DefaultNetDir = "/etc/cni/net.d"

	CmdAdd = "add"
	CmdDel = "del"
)

func main() {
	if len(os.Args) < 3 {
		usage()
		return
	}

	netdir := os.Getenv(EnvNetDir)
	if netdir == "" {
		netdir = DefaultNetDir
	}
	netconf, err := libcni.LoadConfList(netdir, os.Args[2])
	if err != nil {
		exit(err)
	}

	var capabilityArgs map[string]interface{}
	args := os.Getenv(EnvCapabilityArgs)
	if len(args) > 0 {
		if err = json.Unmarshal([]byte(args), &capabilityArgs); err != nil {
			exit(err)
		}
	}

	netns := os.Args[3]

	cninet := &libcni.CNIConfig{
		Path: filepath.SplitList(os.Getenv(EnvCNIPath)),
	}

	rt := &libcni.RuntimeConf{
		ContainerID:    "cni",
		NetNS:          netns,
		IfName:         "eth0",
		CapabilityArgs: capabilityArgs,
	}

	switch os.Args[1] {
	case CmdAdd:
		result, err := cninet.AddNetworkList(netconf, rt)
		if result != nil {
			_ = result.Print()
		}
		exit(err)
	case CmdDel:
		exit(cninet.DelNetworkList(netconf, rt))
	}
}

func usage() {
	exe := filepath.Base(os.Args[0])

	fmt.Fprintf(os.Stderr, "%s: Add or remove network interfaces from a network namespace\n", exe)
	fmt.Fprintf(os.Stderr, "  %s %s <net> <netns>\n", exe, CmdAdd)
	fmt.Fprintf(os.Stderr, "  %s %s <net> <netns>\n", exe, CmdDel)
	os.Exit(1)
}

func exit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
