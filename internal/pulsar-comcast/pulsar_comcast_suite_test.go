/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package pulsar_comcast

import (
    "testing"
    "github.com/onsi/ginkgo"
    "github.com/onsi/gomega"
)

func TestPulsar(t *testing.T) {
    gomega.RegisterFailHandler(ginkgo.Fail)
    ginkgo.RunSpecs(t, "Pulsar comcast executor Suite")
}