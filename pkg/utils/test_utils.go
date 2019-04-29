/*
 *  Copyright (C) 2019 Nalej Group - All Rights Reserved
 */

package utils


import "os"

// RunIntegrationTests checks whether integration tests should be executed.
func RunIntegrationTests() bool {
    var runIntegration = os.Getenv("RUN_INTEGRATION_TEST")
    return runIntegration == "true"
}
