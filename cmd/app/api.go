package main

import (
	"context"
	"net/http"

	"github.com/cristalhq/acmd"
	"go.temporal.io/sdk/client"

	"github.com/xenking/temporal-apps/api"
)

// apiCmd represents the api command
var apiCmd = acmd.Command{
	Name:        "api",
	Description: "Run API Server",
	ExecFunc: func(ctx context.Context, args []string) error {
		c, err := client.DialContext(ctx, client.Options{
			//HostPort: "localhost:7233",
		})
		if err != nil {
			return err
		}
		defer c.Close()

		srv := &http.Server{
			Handler: api.Router(c),
			Addr:    "0.0.0.0:8084",
		}

		errCh := make(chan error, 1)
		go func() {
			errCh <- srv.ListenAndServe()
		}()

		select {
		case <-ctx.Done():
			return srv.Close()
		case err = <-errCh:
			return err
		}
	},
}

func init() {
	commands = append(commands, apiCmd)
}
