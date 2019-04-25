/*
 *  Copyright (C) 2019 Nalej Group - All Rights Reserved
 */

package cmd

import (
    "github.com/nalej/nalej-bus/version"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "github.com/spf13/cobra"
    "os"
)

// Pulsar host address E.G. localhost:6550
var pulsarAddress string
// Default log flag
var debugLevel bool
// Set console logging format
var consoleLogging bool


var RootCmd = &cobra.Command{
    Use: "nalej-pulsar-cli",
    Short: "Pulsar client",
    Long: "Example of a Nalej client for a Pulsar message queue",
    Version: "unknown-version",
}


func Execute(){
    SetupLogging()
    RootCmd.SetVersionTemplate(version.GetVersionInfo())
    if err := RootCmd.Execute(); err != nil {
        log.Error().Msg(err.Error())
    }
}


func init() {

    RootCmd.PersistentFlags().StringVar(&pulsarAddress, "pulsarAddress", "localhost:6650", "Address for the pulsar server")
    RootCmd.PersistentFlags().BoolVar(&debugLevel, "debug", false, "enable debugLevel mode")
    RootCmd.PersistentFlags().BoolVar(&consoleLogging, "consoleLogging", false, "Pretty print logging")
}

// SetupLogging sets the debugLevel level and console logging if required.
func SetupLogging() {
    zerolog.TimeFieldFormat = ""
    zerolog.SetGlobalLevel(zerolog.InfoLevel)
    if debugLevel {
        zerolog.SetGlobalLevel(zerolog.DebugLevel)
    }

    if consoleLogging {
        log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
    }
}