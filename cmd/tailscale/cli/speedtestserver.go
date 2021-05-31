package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"tailscale.com/client/tailscale"
	"tailscale.com/speedtest"

	"github.com/peterbourgon/ff/v2/ffcli"
)

var speedtestCmd = &ffcli.Command{
	Name:       "speedtest",
	ShortUsage: "speedtest <server|client> ...",
	ShortHelp:  "Run a speed test",
	Subcommands: []*ffcli.Command{
		speedtestServerCmd,
		speedtestClientCmd,
	},
	Exec: func(context.Context, []string) error {
		return errors.New("subcommand required; run 'tailscale speedtest -h' for details")
	},
}

var speedtestServerCmd = &ffcli.Command{
	Name:       "server",
	ShortUsage: "speedtest server <host> <port>",
	ShortHelp:  "Start a speed test server",
	Exec:       runServer,
	FlagSet: (func() *flag.FlagSet {
		fs := flag.NewFlagSet("server", flag.ExitOnError)
		fs.IntVar(&serverArgs.port, "port", 0, "port to listen on")
		return fs
	})(),
}

var speedtestClientCmd = &ffcli.Command{
	Name:       "client",
	ShortUsage: "speedtest client HOST:PORT <download|upload>",
	ShortHelp:  "Start a speed test client and connect to a speed test server",
	Exec:       runClient,
}

var serverArgs struct {
	port int
}

func runServer(ctx context.Context, args []string) error {
	if serverArgs.port == 0 {
		return errors.New("port needs to be provided")
	}
	portString := fmt.Sprint(serverArgs.port)
	st, err := tailscale.Status(ctx)
	if err != nil {
		return err
	}
	ips := st.TailscaleIPs
	if len(ips) == 0 {
		return errors.New("no tailscale ips found")
	}
	var host string
	for _, ip := range ips {
		if ip.Is4() {
			host = ip.String()
		}
	}
	err = speedtest.StartServer(host, portString)
	return err
}

func runClient(ctx context.Context, args []string) error {
	return errors.New("not implemented")
}
