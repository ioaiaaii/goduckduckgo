package main

import (
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd, represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goduckduckgo",
	Short: "goduckduckgo queries GDD API",
	Long:  "goduckduckgo is a cli tool that handles the service runtime.\n",
}

// execute, adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func execute() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// interrupt, handles the lifecycle of run groups, called by the sub-commands to spawn a http servers with concurrency.
// When a OS signal for interruption comes, it sends a cancel signal to the running groups.
func interrupt(cancel <-chan struct{}) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	select {
	case s := <-c:
		log.Info().Msgf("Caught signal %v", s)

		return nil
	case <-cancel:
		return errors.New("canceled")
	}
}
