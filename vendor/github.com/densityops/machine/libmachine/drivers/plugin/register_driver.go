package plugin

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	"github.com/densityops/machine/libmachine/drivers"
	"github.com/densityops/machine/libmachine/drivers/plugin/localbinary"
	rpcdriver "github.com/densityops/machine/libmachine/drivers/rpc"
	"github.com/densityops/machine/libmachine/version"
	log "github.com/sirupsen/logrus"
)

var (
	heartbeatTimeout = 10 * time.Second
)

func RegisterDriver(d drivers.Driver) {
	if os.Getenv(localbinary.PluginEnvKey) != localbinary.PluginEnvVal {
		fmt.Fprintf(os.Stderr, `This is a hypervisor plugin binary.
(Driver version: %s, API version: %d)
`, d.DriverVersion(),
			version.APIVersion)
		os.Exit(1)
	}

	log.SetLevel(log.DebugLevel)
	os.Setenv("MACHINE_DEBUG", "1")

	rpcd := rpcdriver.NewRPCServerDriver(d)
	if err := rpc.RegisterName(rpcdriver.RPCServiceNameV0, rpcd); err != nil {
		log.Error(err)
	}
	if err := rpc.RegisterName(rpcdriver.RPCServiceNameV1, rpcd); err != nil {
		log.Error(err)
	}
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading RPC server: %s\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println(listener.Addr())

	go func() {
		_ = http.Serve(listener, nil)
	}()

	for {
		select {
		case <-rpcd.CloseCh:
			log.Debug("Closing plugin on server side")
			os.Exit(0) //nolint
		case <-rpcd.HeartbeatCh:
			continue
		case <-time.After(heartbeatTimeout):
			// TODO: Add heartbeat retry logic
			os.Exit(1)
		}
	}
}
