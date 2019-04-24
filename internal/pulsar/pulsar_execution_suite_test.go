/*
 *  Copyright (C) 2019 Nalej Group - All Rights Reserved
 */

package pulsar


import (
    "testing"
    "github.com/onsi/ginkgo"
    "github.com/onsi/gomega"
)

func TestPulsar(t *testing.T) {
    gomega.RegisterFailHandler(ginkgo.Fail)
    ginkgo.RunSpecs(t, "Pulsar executor Suite")
}


