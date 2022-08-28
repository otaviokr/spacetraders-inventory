package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/otaviokr/spacetraders-inventory/component"
	"github.com/otaviokr/spacetraders-inventory/web"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"

	traceSdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	TracerName          = "spacetrader-inventory"
	collectionFrequency = 300000 // 5 min = 300000 ms
)

// main is just the starting point, but we keep just the bare minimum here.
//
// https://pace.dev/blog/2020/02/12/why-you-shouldnt-use-func-main-in-golang-by-mat-ryer.html
func main() {
	token := os.Getenv("USER_TOKEN")
	jaegerUrl := os.Getenv("JAEGER_URL")

	metricsPort := os.Getenv("METRICS_PORT")
	if len(metricsPort) < 1 {
		metricsPort = "9091"
	}

	// This is function to expose the metrics to Prometheus.
	go exposeMetrics(metricsPort)

	// The main loop is actually inside the run function.
	if err := run(token, jaegerUrl); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// run contains the main loop of the program. It will collect data from the Space Traders game and
// expose them to Prometheus.
func run(token, jaegerUrl string) error {
	bgCtx := context.Background()
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerUrl)))
	if err != nil {
		log.Fatal(err)
	}

	tp := traceSdk.NewTracerProvider(
		traceSdk.WithBatcher(exp),
		traceSdk.WithResource(newResource()))
	defer func() {
		if err := tp.Shutdown(bgCtx); err != nil {
			log.Fatal(err)
		}
	}()
	otel.SetTracerProvider(tp)

	tracer := otel.Tracer(TracerName)

	user, err := component.NewUser(bgCtx, tracer, token)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("User logged in: %s\n", user.Details.Username)

	for {
		startTime := time.Now()

		// Check game status.
		status, err := user.GetStatus(bgCtx)
		if err != nil {
			log.Fatal(err)
		}
		web.GameStatus.WithLabelValues(user.Details.Username).Set(float64(status))
		log.Printf("%+v\n", status)

		// Get user details (credits, ship count, structure count, joined date).
		err = user.GetDetails(bgCtx)
		if err != nil {
			log.Fatal(err)
		}
		web.Credits.WithLabelValues(user.Details.Username).Set(float64(user.Details.Credits))
		web.ShipCount.WithLabelValues(user.Details.Username).Set(float64(user.Details.ShipCount))
		web.StructureCount.WithLabelValues(user.Details.Username).Set(float64(user.Details.StructureCount))

		log.Printf("%+v\n", user.Details)

		// TODO // Get current loans.
		// loans, err := user.GetLoans(bgCtx)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// Get leaderboard.
		board, err := user.GetLeaderboard(bgCtx)
		if err != nil {
			log.Fatal(err)
		}

		web.UserRank.WithLabelValues(user.Details.Username).Set(float64(board.UserNetWorth.Rank))
		log.Printf("%+v\n", board.UserNetWorth)

		// Get details of all ships.
		ships, err := user.GetShips(bgCtx)
		if err != nil {
			log.Fatal(err)
		}

		for _, ship := range ships.Ships {
			web.ShipLoad.WithLabelValues(
				user.Details.Username,
				ship.Id,
				ship.Class,
				ship.Manufacturer,
				ship.Type,
				fmt.Sprintf("%d", ship.MaxCargo),
				fmt.Sprintf("%d", ship.Plating),
				fmt.Sprintf("%d", ship.Speed),
				fmt.Sprintf("%d", ship.Weapons)).Set(float64(ship.SpaceAvailable))
			log.Printf("%+v\n", ship)
		}

		// TODO // Get details of all structures.

		// Wait before collecting new data.
		wait := time.Since(startTime).Milliseconds()
		log.Printf("Duration: %d / Wait: %d\n", wait, collectionFrequency-wait)
		time.Sleep(time.Duration(collectionFrequency-wait) * time.Millisecond)
	}
}

// exposeMetrics is a very simple web server that Prometheus can access to collect the metrics.
//
// port is the port where the web server is listening.
func exposeMetrics(port string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// newResource returns a resource describing this application to be used in the traces for Jaeger.
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("spacetraders-inventory"),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}
