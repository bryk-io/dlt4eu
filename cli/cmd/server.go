package cmd

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"os"
	"syscall"

	"github.com/bryk-io/dlt4eu/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.bryk.io/x/cli"
	"go.bryk.io/x/net/graphql"
	"go.bryk.io/x/net/graphql/extension/tracing"
	"go.bryk.io/x/net/loader"
	"go.bryk.io/x/net/middleware"
	"go.bryk.io/x/observability"
)

var helper *loader.Helper

var serverCmd = &cobra.Command{
	Use:     "server",
	Short:   "Start a GraphQL server instance to handle incoming requests",
	Example: "dlt4eu server --config file.yaml",
	RunE:    runServer,
}

func init() {
	segments := []string{
		loader.SegmentHTTP,
		loader.SegmentGraphQL,
		loader.SegmentObservability,
		loader.SegmentMiddlewareCORS,
		loader.SegmentMiddlewareMetadata,
	}
	helper = loader.New()
	if err := cli.SetupCommandParams(serverCmd, helper.Params(segments...)); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(serverCmd)
}

func runServer(_ *cobra.Command, _ []string) error {
	// Load general configuration settings.
	if err := viper.Unmarshal(helper.Data); err != nil {
		return err
	}

	// Load service handler settings.
	conf := new(service.Config)
	if err := viper.UnmarshalKey("service", conf); err != nil {
		return err
	}

	// Get observability operator. Discard tracer output and keep
	// only logging and metrics.
	oopSettings := helper.Observability()
	oop, err := observability.NewOperator(oopSettings...)
	if err != nil {
		return err
	}

	// Get service handler instance.
	handler, err := service.New(conf)
	if err != nil {
		return err
	}

	// Get GraphQL server instance.
	srv, err := graphql.NewServer(handler.Schema(), helper.GraphQL()...)
	if err != nil {
		return err
	}

	// Apply HTTP middleware.
	srv.Use(middleware.GzipCompression(gzip.DefaultCompression))
	srv.Use(middleware.ProxyHeaders())
	srv.Use(middleware.CORS(helper.MiddlewareCORS()))
	srv.Use(middleware.ContextMetadata(helper.MiddlewareMetadata()))
	srv.Use(customHeaders())
	srv.Use(oop.HTTPServerMiddleware())

	// Apply GraphQL extensions.
	srv.UseHandlerExtension(tracing.FieldTracing(oop, false))

	// Start server
	router := http.NewServeMux()
	router.Handle("/graphql", srv.GetHandler())
	router.HandleFunc("/ping", pingHandler())
	oop.Infof("waiting for requests on: http://localhost:%d/graphql", helper.Data.HTTP.Port)
	go func() {
		_ = http.ListenAndServe(fmt.Sprintf(":%d", helper.Data.HTTP.Port), router)
	}()

	// wait for system signals
	<-cli.SignalsHandler([]os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	})

	// Close server
	oop.Info("closing server")
	handler.Shutdown()
	return nil
}

// Attach version information as HTTP headers on all responses
func customHeaders() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("x-dlt4eu-version", coreVersion)
			w.Header().Set("x-dlt4eu-build", buildCode)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// Manage a basic (not instrumented) ping request
func pingHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, _ *http.Request) {
		_, _ = res.Write([]byte("ok"))
	}
}
