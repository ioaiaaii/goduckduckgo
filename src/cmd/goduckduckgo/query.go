package main

import (
	v1 "goduckduckgo/pkg/api/query"
	"goduckduckgo/pkg/config"
	"goduckduckgo/pkg/store"

	"goduckduckgo/pkg/server/http"

	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const sortHelpQuery = "query, spawns a new server exposing the query API. Call it --help flag to check the defaults and the available parameters"
const exampleQuery = "goduckduckgo query --help"

// queryCmd represents the check command
var queryCmd = &cobra.Command{
	Use:     "query [flags]",
	Short:   sortHelpQuery,
	Example: exampleQuery,
	RunE:    runQuery,
}

// init, is a standard method for cobra framework to init and register commands flags and params.
// It provides as flags options all the server configuration. Also, it overwrites values from the given ENV Vars.
// e.g for server-port, the equivalent env var is SERVER_PORT.
func init() {
	var err error
	rootCmd.AddCommand(queryCmd)

	queryCmd.Flags().String("server-port", config.DefaultQueryHTTPPort, "The port for communication")
	err = viper.BindPFlag("server-port", queryCmd.Flags().Lookup("server-port"))
	if err != nil {
		log.Warn().Err(err).Msgf("Flagging server-port failed with error: %v", err.Error())
	}

	queryCmd.Flags().Duration("server-timeout", config.DefaultQueryHTTPServerTimeout, "Timeout for HTTP Server's Requests and Responses")
	err = viper.BindPFlag("server-timeout", queryCmd.Flags().Lookup("server-timeout"))
	if err != nil {
		log.Warn().Err(err).Msgf("Flagging server-timeout failed with error: %v", err.Error())
	}

	queryCmd.Flags().Duration("shutdown-timeout", config.DefaultQueryHTTPServerShutdownTimeout, "Timeout for HTTP Server graceful timeout")
	err = viper.BindPFlag("shutdown-timeout", queryCmd.Flags().Lookup("shutdown-timeout"))
	if err != nil {
		log.Warn().Err(err).Msgf("Flagging shutdown-timeout failed with error: %v", err.Error())
	}

	queryCmd.Flags().String("store-endpoint", config.DefaultQueryEndpoint, "gRPC Endpoint for Query")
	err = viper.BindPFlag("store-endpoint", queryCmd.Flags().Lookup("store-endpoint"))
	if err != nil {
		log.Warn().Err(err).Msgf("Flagging store-endpoint failed with error: %v", err.Error())
	}
}

// runQuery, contains the logic to run the query API.
// It creates a new router, and it registers to the new server.
// It creates a cancel channel, runs the server to run.Group and manages its lifecycle.
func runQuery(cmd *cobra.Command, args []string) error {

	var (
		err error
		g   run.Group
		cfg config.Config
	)

	//Init config
	cfg.QueryConfig.HTTPPort = viper.GetString("server-port")
	cfg.QueryConfig.HTTPServerShutdownTimeout = viper.GetDuration("shutdown-timeout")
	cfg.QueryConfig.HTTPServerTimeout = viper.GetDuration("server-timeout")
	cfg.QueryConfig.StoreEndpoint = viper.GetString("store-endpoint")

	// gRPC Client
	sc := store.RegisterStoreClient(cfg.QueryConfig.StoreEndpoint)
	client, err := sc.NewStoreClient()
	if err != nil {
		log.Fatal().Err(err)
	}

	//HTTP Server
	router := http.NewRouter()
	apiV1 := v1.NewQueryAPI(client)
	apiV1.Register(router.WithPrefix("/api/v1"))

	//Create new server and register handler with new router
	srv := http.NewServer(http.HTTPPort(viper.GetString("server-port")), http.HTTPServerTimeout(viper.GetDuration("server-timeout")), http.HTTPServerShutdownTimeout(viper.GetDuration("shutdown-timeout")))
	srv.HandleAPIPath("/", router)
	srv.HandleAPIPath("/metrics", promhttp.Handler())

	//add function to run group
	//start serving and gracefully shutdown the server
	{
		g.Add(func() error {
			return srv.ListenAndServe()
		}, func(err error) {
			srv.Shutdown(err)
		})

	}

	//create and monitor the cancel channel.
	//if an interruption signal found, call interrupt to stop serving
	{
		cancel := make(chan struct{})
		g.Add(func() error {
			return interrupt(cancel)
		}, func(error) {
			close(cancel)
		})
	}

	//Run all actors (functions) concurrently.
	//When the first actor returns, all others are interrupted.
	if err := g.Run(); err != nil {
		log.Warn().Err(err).Msgf("Running run Groups failed with error: %v", err.Error())
		return err
	}

	return nil
}
