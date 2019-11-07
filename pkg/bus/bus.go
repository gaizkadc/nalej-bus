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

package bus

import (
	"context"
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
	// params:
	//  ctx context
	// return:
	//  message payload
	//  error if any
	Receive(ctx context.Context) ([]byte, derrors.Error)

	// Close the consumer
	// params:
	//  ctx context
	// return:
	//  error if any
	Close(ctx context.Context) derrors.Error
}

type NalejProducer interface {

	// Send a new message to the topic of this producer.
	// params:
	//  msg message to be sent
	//  ctx context
	// return:
	//  error if any
	Send(ctx context.Context, msg []byte) derrors.Error

	// Close the producer. This operation must close any connection with brokers with an established connection.
	// params:
	//  ctx context
	// return:
	//  error if any
	Close(ctx context.Context) derrors.Error
}
