/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package pulsar_comcast

import (
    "context"
    "github.com/nalej/derrors"
    "github.com/nalej/nalej-bus/pkg/bus"
    "github.com/nalej/nalej-bus/pkg/utils"
    "github.com/onsi/ginkgo"
    "github.com/onsi/gomega"
    "os"
    "time"
)

const(
    SendTimeout = 2
)

var _ = ginkgo.Describe("Test execution of Pulsar wrappers in Nalej", func() {

    var isReady bool
    var PulsarAddress string
    var client PulsarClient
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
            if !isReady {
                ginkgo.Skip("no integration test set")
                return
            }

            //create client
            client = NewClient(PulsarAddress, nil).(PulsarClient)
            gomega.Expect(client).ShouldNot(gomega.BeNil())
        })


        ginkgo.It("produce and consume with the same client", func(){

            if !isReady {
                ginkgo.Skip("no integration test set")
                return
            }

            // create producer
            prod := NewPulsarProducer(client, "prod1", "public/default/topic")
            gomega.Expect(prod).ShouldNot(gomega.BeNil())

            // create consumer
            cons := NewPulsarConsumer(client, "cons1", "public/default/topic", true)
            gomega.Expect(prod).ShouldNot(gomega.BeNil())

            msg := "test message"

            received_ch := make(chan []byte,1)

            go receive(cons, received_ch)

            // produce something
            ctx,cancel := context.WithTimeout(context.Background(), time.Second * SendTimeout)
            err = prod.Send(ctx, []byte(msg))
            cancel()
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred(), "error sending message")

            // In one second this should arrive
            time.Sleep(time.Second)
            // this channel should return something interesting
            received := <-received_ch
            gomega.Expect(string(received)).Should(gomega.Equal(msg))

            ctx,cancel = context.WithTimeout(context.Background(), time.Second * SendTimeout)
            err = prod.Close(ctx)
            cancel()
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred(),"error closing producer")

            ctx,cancel = context.WithTimeout(context.Background(), time.Second * SendTimeout)
            err = cons.Close(ctx)
            cancel()
            gomega.Expect(err).ShouldNot(gomega.HaveOccurred(),"error closing consumer")

        })
    })
})

// Helping function using a channel to return results
func receive(cons bus.NalejConsumer, ch chan <- []byte)  {
    ctx,cancel := context.WithTimeout(context.Background(), time.Second * SendTimeout)
    received, err := cons.Receive(ctx)
    cancel()
    gomega.Expect(err).ShouldNot(gomega.HaveOccurred(), "error receiving message")

    ch <- received
}
