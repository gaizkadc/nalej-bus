/*
 *  Copyright (C) 2019 Nalej Group - All Rights Reserved
 */

package pulsar

import (
    "github.com/apache/pulsar/pulsar-client-go/pulsar"
    "github.com/nalej/nalej-bus/internal/utils"
    "github.com/onsi/ginkgo"
    "github.com/onsi/gomega"
    "os"
)

var _ = ginkgo.Describe("Test execution of Pulsar wrappers in Nalej", func() {

    var isReady bool
    var PulsarAddress string

    ginkgo.BeforeSuite(func() {
        isReady = false
        if utils.RunIntegrationTests() {
            PulsarAddress = os.Getenv(utils.IT_PULSAR_ADDRESS)
            if PulsarAddress != "" {
                isReady = true
            }
        }

        if !isReady {
            return
        }
    })

    ginkgo.Context("Evaluate consumers", func(){
        var client pulsar.Client

        ginkgo.BeforeEach(func(){
            //create client
            client, err := NewClient(PulsarAddress,5,1)
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
            gomega.Expect(client).ShouldNot(gomega.BeNil())
        })

        ginkgo.AfterEach(func(){
            // close client
            err := client.Close()
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
        })

        ginkgo.It("produce and consume with the same client", func(){
            // create producer
            prod,err := NewPulsarProducer(client, "prod1", "test/topic")
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
            gomega.Expect(prod).ShouldNot(gomega.BeNil())

            // create consumer
            cons, err := NewPulsarConsumer(client, "cons1", "test/topic", PulsarExclusiveConsumer)
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
            gomega.Expect(prod).ShouldNot(gomega.BeNil())

            // produce something
            msg := "test message"
            err = prod.Send([]byte(msg))
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred(), "error sending message")

            //consume something
            received, err := cons.Receive()
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred(), "error receiving message")
            gomega.Expect(received).Should(gomega.Equal(msg))

            err = prod.Close()
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred(),"error closing producer")
            err = cons.Close()
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred(),"error closing consumer")

        })
    })
})
