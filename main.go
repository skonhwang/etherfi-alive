package main

import (
	"bytes"
	"flag"
	"fmt"
	"os/exec"
	"strings"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/sdk"

	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	// Default values.
	defaultAddr     = "127.0.0.1"
	defaultPort     = 9002
	dirSyncClientV2 = "/root/bin/etherfi/sync-client-v2/"
	nameSyncClient  = "etherfi-sync-clientv2"

	pluginName = "vatz-plugin-etherfi-alive"
)

var (
	addr     string
	port     int
	dirSync  string
	nameSync string
)

func init() {
	flag.StringVar(&addr, "addr", defaultAddr, "IP Address(e.g. 0.0.0.0, 127.0.0.1)")
	flag.IntVar(&port, "port", defaultPort, "Port number, default 9091")
	flag.StringVar(&dirSync, "dirsync", dirSyncClientV2, "Location of etherfi-sync-clientv2")
	flag.StringVar(&nameSync, "namesync", nameSyncClient, "Process name of etherfi-sync-client")

	flag.Parse()
}

func main() {

	p := sdk.NewPlugin(pluginName)
	p.Register(pluginFeature)

	ctx := context.Background()
	if err := p.Start(ctx, addr, port); err != nil {
		fmt.Println("exit")
	}
}

func checkEtherfiSyncClientv2() bool {
	retBool := false
	//cmd := exec.Command("pgrep", "-x", "etherfi_sc_v2")
	fmt.Println("nameSync: ", nameSync)
	cmd := exec.Command("pgrep", "-x", nameSync)
	out, err := cmd.Output()
	if err != nil && err.Error() != "exit status 1" {
		fmt.Println("Error: ", err)
		return retBool
	}

	if len(out) == 0 {
		fmt.Println("etherfi_sc_v2 is not running")
		process := nameSync
		fmt.Println("Process name = ", nameSync)
		//cmd := exec.Command("./etherfi_sc_v2", "listen")
		cmd := exec.Command(process, "listen")
		//cmd.Dir = "/Users/hwangseungkon/dsrv/2022/etherfi/sync-client-v2/"
		cmd.Dir = dirSync
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error starting etherfi_sc_v2: ", err)
			fmt.Println("Stderr:", stderr.String())
			return retBool
		}
		return retBool
	} else {
		fmt.Println("etherfi_sc_v2 is running PID = ", strings.Split(string(out), "\n"))
		retBool = true
		return retBool
	}

}

func pluginFeature(info, option map[string]*structpb.Value) (sdk.CallResponse, error) {
	severity := pluginpb.SEVERITY_INFO
	state := pluginpb.STATE_NONE
	// TODO: Fill here.
	fmt.Println("Check process!!")
	chkProcess := checkEtherfiSyncClientv2()
	var str string
	if chkProcess == false {
		str = "etherfi sync client v2 is not running, so it will restart!"
		severity = pluginpb.SEVERITY_CRITICAL
		state = pluginpb.STATE_SUCCESS
	} else {
		str = "Sync client is okay"
		severity = pluginpb.SEVERITY_INFO
		state = pluginpb.STATE_SUCCESS
	}
	ret := sdk.CallResponse{
		FuncName:   "etherfi_alive_func",
		Message:    str,
		Severity:   severity,
		State:      state,
		AlertTypes: []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
	}

	return ret, nil
}
