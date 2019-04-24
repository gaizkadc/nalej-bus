/*
 *  Copyright (C) 2019 Nalej Group - All Rights Reserved
 */

package bus

import (
    "github.com/nalej/derrors"
)


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

    // Close the producer. This operation must close any connection with brokers with an estabalished connection.
    // return:
    //  error if any
    Close() derrors.Error
}