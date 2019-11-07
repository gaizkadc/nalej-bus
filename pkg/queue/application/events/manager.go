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

package events

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-bus-go"
	"github.com/nalej/grpc-conductor-go"
	"github.com/nalej/nalej-bus/pkg/bus"
	"github.com/nalej/nalej-bus/pkg/queue"
	"github.com/rs/zerolog/log"
)

const (
	ApplicationEventsTopic="nalej/application/events"
)

type ApplicationEventsProducer struct {
	producer bus.NalejProducer
}

// Create a new producer for the infrastructure events topic
// params:
//  client to be used
//  name of the producer
// return:
//  built producer
func NewApplicationEventsProducer (client bus.NalejClient, name string) (*ApplicationEventsProducer, derrors.Error) {
	prod, err := client.BuildProducer(name, ApplicationEventsTopic)
	if err != nil {
		return nil, err
	}
	return &ApplicationEventsProducer{producer: prod}, nil
}

// Generic function that decides what to send depending on the object type.
func (m ApplicationEventsProducer) Send(ctx context.Context, msg proto.Message) derrors.Error {

	var  wrapper grpc_bus_go.ApplicationEvents

	switch x := msg.(type) {
	case *grpc_conductor_go.DeploymentServiceUpdateRequest:
		wrapper = grpc_bus_go.ApplicationEvents{Event: &grpc_bus_go.ApplicationEvents_DeploymentServiceStatusUpdateRequest{x}}
	default:
		return derrors.NewInvalidArgumentError("invalid proto message type")
	}

	payload, err := queue.MarshallPbMsg(&wrapper)
	if err != nil {
		return err
	}

	err = m.producer.Send(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

// Data struct indicating what data structures available in this topic will be accepted.
type ConsumableStructsApplicationEventsConsumer struct {
	// Consume DeploymentServiceUpdateRequest
	DeploymentServiceUpdateRequest bool
}

// Struct designed to config a consumer defining what actions to perform depending on the incoming object.
type ConfigApplicationEventsConsumer struct {
	// channel to receive DeploymentServiceStatusUpdateRequest
	ChDeploymentServiceStatusUpdateRequest chan *grpc_conductor_go.DeploymentServiceUpdateRequest
	// object types to be considered for consumption
	ToConsume ConsumableStructsApplicationEventsConsumer
}

// Create a new configuration structure for a given channel size
// params:
//  size of channels
// return:
//  instance of a configuration object
func NewConfigApplicationEventsConsumer(size int, toConsume ConsumableStructsApplicationEventsConsumer) ConfigApplicationEventsConsumer {
	chDeploymentServiceStatusUpdateRequest := make(chan *grpc_conductor_go.DeploymentServiceUpdateRequest, size)

	return ConfigApplicationEventsConsumer{
		ChDeploymentServiceStatusUpdateRequest: chDeploymentServiceStatusUpdateRequest,
		ToConsume: toConsume,
	}
}

// Application consumer
type ApplicationEventsConsumer struct {
	Consumer bus.NalejConsumer
	Config   ConfigApplicationEventsConsumer
}

func NewApplicationEventsConsumer (client bus.NalejClient, name string, exclusive bool, config ConfigApplicationEventsConsumer) (*ApplicationEventsConsumer, derrors.Error) {
	consumer, err := client.BuildConsumer(name, ApplicationEventsTopic, exclusive)
	if err != nil {
		return nil, err
	}

	return &ApplicationEventsConsumer{Consumer: consumer,Config: config}, nil
}

// Consume any of the possible objects that can be sent to this queue and send it to the corresponding channel.
func (c ApplicationEventsConsumer) Consume(ctx context.Context) derrors.Error{
	msg, err := c.Consumer.Receive(ctx)
	if err != nil {
		return err
	}

	target := &grpc_bus_go.ApplicationEvents{}

	derr := queue.UnmarshallPbMsg(msg, target)
	if derr != nil {
		return derr
	}

	switch x := target.Event.(type) {
	case *grpc_bus_go.ApplicationEvents_DeploymentServiceStatusUpdateRequest:
		if c.Config.ToConsume.DeploymentServiceUpdateRequest {
			c.Config.ChDeploymentServiceStatusUpdateRequest <- x.DeploymentServiceStatusUpdateRequest
		}
	case nil:
		errMsg := "received nil entry"
		log.Error().Msg(errMsg)
		return derrors.NewInvalidArgumentError(errMsg)
	default:
		errMsg := fmt.Sprintf("unknown object type in %s",ApplicationEventsTopic)
		log.Error().Interface("type",x).Msg(errMsg)
		return derrors.NewInvalidArgumentError(errMsg)
	}
	return nil
}
