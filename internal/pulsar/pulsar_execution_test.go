/*
 *  Copyright (C) 2019 Nalej Group - All Rights Reserved
 */

package pulsar

import (
    "github.com/apache/pulsar/pulsar-client-go/pulsar"
    "github.com/nalej/derrors"
    "github.com/nalej/nalej-bus/internal/bus"
    "github.com/nalej/nalej-bus/internal/utils"
    "github.com/onsi/ginkgo"
    "github.com/onsi/gomega"
    "github.com/rs/zerolog/log"
    "os"
)

var _ = ginkgo.Describe("Test execution of Pulsar wrappers in Nalej", func() {

    var isReady bool
    var PulsarAddress string
    var client pulsar.Client
    var err derrors.Error

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


        ginkgo.BeforeEach(func(){
            //create client
            client, err = NewClient(PulsarAddress,5,1)
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
            log.Debug().Msg("Before producer")
            prod,err := NewPulsarProducer(client, "prod1", "public/default/topic")
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
            gomega.Expect(prod).ShouldNot(gomega.BeNil())

            // create consumer
            cons, err := NewPulsarConsumer(client, "cons1", "public/default/topic", PulsarExclusiveConsumer)
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
            gomega.Expect(prod).ShouldNot(gomega.BeNil())

            msg := "test message"

            received_ch := make(chan []byte,1)

            go receive(cons, received_ch)

            // produce something
            err = prod.Send([]byte(msg))
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred(), "error sending message")

            // this channel should return something interesting
            received := <-received_ch
            gomega.Expect(string(received)).Should(gomega.Equal(msg))

            err = prod.Close()
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred(),"error closing producer")
            err = cons.Close()
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred(),"error closing consumer")

        })
    })
})

// Helping function using a channel to return results
func receive(cons bus.NalejConsumer, ch chan <- []byte)  {
    received, err := cons.Receive()
    gomega.Expect(err).ShouldNot(gomega.HaveOccurred(), "error receiving message")

    ch <- received
}