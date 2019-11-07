/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
    Use: "nalej-pulsar-cli-legacy",
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