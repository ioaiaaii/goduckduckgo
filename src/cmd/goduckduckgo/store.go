package main

import (
	"goduckduckgo/internal/storage/db"
	"goduckduckgo/pkg/config"
	"goduckduckgo/pkg/server/grpc"
	"goduckduckgo/pkg/store"

	"github.com/oklog/run"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const sortHelpStore = "store, stores data to DB"
const exampleStore = "goduckduckgo store --help"

// storeCmd represents the check command
var storeCmd = &cobra.Command{
	Use:     "store [flags]",
	Short:   sortHelpStore,
	Example: exampleStore,
	RunE:    runStore,
}

// init, is a standard method for cobra framework to init and register commands flags and params.
// It provides as flags options all the server configuration. Also, it overwrites values from the given ENV Vars.
// e.g for server-port, the equivalent env var is SERVER_PORT.
func init() {
	var err error
	rootCmd.AddCommand(storeCmd)

	storeCmd.Flags().String("db-host", config.DefaultStoreDBHost, "The DB Endpoint")
	err = viper.BindPFlag("db-host", storeCmd.Flags().Lookup("db-host"))
	if err != nil {
		log.Warn().Err(err).Msgf("Flagging db-host failed with error: %v", err.Error())
	}

	storeCmd.Flags().String("db-port", config.DefaultStoreDBPort, "The DB Port")
	err = viper.BindPFlag("db-port", storeCmd.Flags().Lookup("db-port"))
	if err != nil {
		log.Warn().Err(err).Msgf("Flagging db-port failed with error: %v", err.Error())
	}

	storeCmd.Flags().String("db-name", config.DefaultStoreDBName, "The DN Name")
	err = viper.BindPFlag("db-name", storeCmd.Flags().Lookup("db-name"))
	if err != nil {
		log.Warn().Err(err).Msgf("Flagging db-name failed with error: %v", err.Error())
	}

	storeCmd.Flags().String("db-user", config.DefaultStoreDBUser, "The DB Username")
	err = viper.BindPFlag("db-user", storeCmd.Flags().Lookup("db-user"))
	if err != nil {
		log.Warn().Err(err).Msgf("Flagging db-user failed with error: %v", err.Error())
	}

	storeCmd.Flags().String("db-password", config.DefaultStoreDBPassword, "The DB Password")
	err = viper.BindPFlag("db-password", storeCmd.Flags().Lookup("db-password"))
	if err != nil {
		log.Warn().Err(err).Msgf("Flagging db-password failed with error: %v", err.Error())
	}

	storeCmd.Flags().String("store-grpc-address", config.DefaultStoreGRPCAddress, "The gRPC Address for Store")
	err = viper.BindPFlag("store-grpc-address", storeCmd.Flags().Lookup("store-grpc-address"))
	if err != nil {
		log.Warn().Err(err).Msgf("Flagging store-grpc-address failed with error: %v", err.Error())
	}

	storeCmd.Flags().String("store-grpc-port", config.DefaultStoreGRPCPort, "The gRPC Port for Store")
	err = viper.BindPFlag("store-grpc-port", storeCmd.Flags().Lookup("store-grpc-port"))
	if err != nil {
		log.Warn().Err(err).Msgf("Flagging store-grpc-port failed with error: %v", err.Error())
	}

}

// runStore, contains the logic to run the store API.
// It creates a new router, and it registers to the new server.
// It creates a cancel channel, runs the server to run.Group and manages its lifecycle.
func runStore(cmd *cobra.Command, args []string) error {

	var cfg config.Config
	var g run.Group
	cfg.StoreConfig.DBHost = viper.GetString("db-host")
	cfg.StoreConfig.DBPort = viper.GetString("db-port")
	cfg.StoreConfig.DBName = viper.GetString("db-name")
	cfg.StoreConfig.DBUser = viper.GetString("db-user")
	cfg.StoreConfig.DBPassword = viper.GetString("db-password")
	cfg.StoreConfig.StoreGRPCAddress = viper.GetString("store-grpc-address")
	cfg.StoreConfig.StoreGRPCPort = viper.GetString("store-grpc-port")

	// DB Init
	storeDB := db.NewDB(cfg)

	err := storeDB.AutoMigrateDB()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	storeClient := store.NewStore(storeDB)

	//add function to run group
	//start serving and gracefully shutdown the server
	{

		s := grpc.NewServer(
			grpc.WithServer(store.RegisterStoreServer(storeClient)),
			grpc.WithListen(cfg.StoreConfig.StoreGRPCAddress+":"+cfg.StoreConfig.StoreGRPCPort),
		)

		g.Add(func() error {
			return s.ListenAndServe()
		}, func(err error) {
			s.Shutdown(err)
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
