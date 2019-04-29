/*
 *  Copyright (C) 2019 Nalej Group - All Rights Reserved
 */

package bus

import (
    "github.com/nalej/derrors"
)

// Interface for nalej clients.
type NalejClient interface {

    // Build a producer using this client
    // params:
    //  name of the producer
    //  topic this producer will use
    // return:
    //  resulting producer
    //  error if any
    BuildProducer(name string, topic string) (NalejProducer, derrors.Error)

    // Build a consumer using this client
    // params:
    //  name of the producer
    //  topic this producer will use
    //  exclusive sets this consumer as the only one consuming data from the topic
    // return:
    //  resulting consumer
    //  error if any
    BuildConsumer(name string, topic string, exclusive bool) (NalejConsumer, derrors.Error)
}

// Main interfaces to be implemented by any producer or consumer in the Nalej platform
type NalejConsumer interface {
    // Receive a message from a subscribed entry
    // return:
    //  message payload
    //  error if any
    Receive() ([]byte, derrors.Error)

    // Close the consumer
    // return:
    //  error if any
    Close() derrors.Error
}

type NalejProducer interface {

    // Send a new message to the topic of this producer.
    // params:
    //  msg message to be sent
    // return:
    //  error if any
    Send(msg []byte) derrors.Error

    // Close the producer. This operation must close any connection with brokers with an established connection.
    // return:
    //  error if any
    Close() derrors.Error
}