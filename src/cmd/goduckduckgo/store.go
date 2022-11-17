package main

import (
	"goduckduckgo/internal/db"
	"goduckduckgo/pkg/config"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const sortHelpStore = "WIP - store, stores data to DB"
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

}

// runstore, contains the logic to run the store API.
// It creates a new router, and it registers to the new server.
// It creates a cancel channel, runs the server to run.Group and manages its lifecycle.
func runStore(cmd *cobra.Command, args []string) error {

	var cfg config.Config

	cfg.StoreConfig.DBHost = viper.GetString("db-host")
	cfg.StoreConfig.DBPort = viper.GetString("db-port")
	cfg.StoreConfig.DBName = viper.GetString("db-name")
	cfg.StoreConfig.DBUser = viper.GetString("db-user")
	cfg.StoreConfig.DBPassword = viper.GetString("db-password")

	newDB := db.NewDB(cfg)

	err := newDB.AutoMigrateDB()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	return nil
}
