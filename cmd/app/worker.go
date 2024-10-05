package main

import (
	"context"
	"fmt"
	"github.com/xenking/temporal-apps/pkg/currency"
	"github.com/xenking/temporal-apps/pkg/currency/middleware"
	"log"
	"time"

	"github.com/cristalhq/acmd"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/uber-go/tally/v4"
	"github.com/uber-go/tally/v4/prometheus"
	"go.temporal.io/sdk/client"
	sdktally "go.temporal.io/sdk/contrib/tally"
	"go.temporal.io/sdk/worker"

	"github.com/xenking/temporal-apps/activities"
	"github.com/xenking/temporal-apps/workflows"
)

// workerCmd represents the worker command
var workerCmd = acmd.Command{
	Name:        "worker",
	Description: "Run worker",
	ExecFunc: func(ctx context.Context, args []string) error {
		cfg, err := parseConfig(args)
		if err != nil {
			return err
		}

		promScope, err := newPrometheusScope(prometheus.Configuration{
			ListenAddress: "0.0.0.0:9092",
			TimerType:     "histogram",
		})
		if err != nil {
			return err
		}

		c, err := client.DialContext(ctx, client.Options{
			HostPort:       "localhost:7233",
			MetricsHandler: sdktally.NewMetricsHandler(promScope),
		})
		if err != nil {
			return err
		}
		defer c.Close()

		cache.New()

		a := activities.Activities{
			Client: c,
			ConversionClient: currency.NewConversionClient(
				middleware.NewCaching(cacheProvider)(currency.NewClient(cfg.Currency.AppID)),
			),
		}

		w := worker.New(c, "currency-check", worker.Options{
			BackgroundActivityContext: ctx,
		})

		w.RegisterWorkflow(workflows.CurrencyVolatileCheckWorkflow)
		w.RegisterActivity(a.GetExchangeRate)

		err = w.Start()
		if err != nil {
			return err
		}

		<-ctx.Done()
		w.Stop()

		return nil
	},
}

func newPrometheusScope(c prometheus.Configuration) (tally.Scope, error) {
	reporter, err := c.NewReporter(
		prometheus.ConfigurationOptions{
			Registry: prom.NewRegistry(),
			OnError: func(err error) {
				log.Println("error in prometheus reporter", err)
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error creating prometheus reporter: %w", err)
	}

	scopeOpts := tally.ScopeOptions{
		CachedReporter:  reporter,
		Separator:       prometheus.DefaultSeparator,
		SanitizeOptions: &sdktally.PrometheusSanitizeOptions,
	}
	scope, _ := tally.NewRootScope(scopeOpts, time.Second)
	scope = sdktally.NewPrometheusNamingScope(scope)

	return scope, nil
}

func init() {
	commands = append(commands, workerCmd)
}
