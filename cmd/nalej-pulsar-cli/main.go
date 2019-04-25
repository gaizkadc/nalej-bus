/*
 *  Copyright (C) 2019 Nalej Group - All Rights Reserved
 */

package main

import (
    "github.com/nalej/nalej-bus/cmd/nalej-pulsar-cli/cmd"
    "github.com/nalej/nalej-bus/version"
)

var MainVersion string

var MainCommit string

func main() {
    version.AppVersion = MainVersion
    version.Commit = MainCommit
    cmd.Execute()
}