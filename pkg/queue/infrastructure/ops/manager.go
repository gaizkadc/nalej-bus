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

package ops

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
	InfrastructureOpsTopic = "nalej/infrastructure/ops"
)

type InfrastructureOpsProducer struct {
	producer bus.NalejProducer
}

// Create a new producer for the application operations topic
// params:
//  client to be used
//  name of the producer
// return:
//  built producer
func NewInfrastructureOpsProducer(client bus.NalejClient, name string) (*InfrastructureOpsProducer, derrors.Error) {
	prod, err := client.BuildProducer(name, InfrastructureOpsTopic)
	if err != nil {
		return nil, err
	}
	return &InfrastructureOpsProducer{producer: prod}, nil
}

// Generic function that decides what to send depending on the object type.
func (m InfrastructureOpsProducer) Send(ctx context.Context, msg proto.Message) derrors.Error {

	var wrapper grpc_bus_go.InfrastructureOps

	switch x := msg.(type) {
	case *grpc_conductor_go.DrainClusterRequest:
		wrapper = grpc_bus_go.InfrastructureOps{Operation: &grpc_bus_go.InfrastructureOps_DrainRequest{x}}
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

// Application consumer
type InfrastructureOpsConsumer struct {
	Consumer bus.NalejConsumer
	Config   ConfigInfrastructureOpsConsumer
}

// Struct designed to config a consumer defining what actions to perform depending on the incoming object.
type ConfigInfrastructureOpsConsumer struct {
	// channel to receive drain requests
	ChDrainRequest chan *grpc_conductor_go.DrainClusterRequest
	// object types to be considered for consumption
	ToConsume ConsumableStructsInfrastructureOpsConsumer
}

// Create a new configuration structure for a given channel size
// params:
//  size of channels
// return:
//  instance of a configuration object
func NewConfigInfrastructureOpsConsumer(size int, toConsume ConsumableStructsInfrastructureOpsConsumer) ConfigInfrastructureOpsConsumer {
	chDrainRequest := make(chan *grpc_conductor_go.DrainClusterRequest, size)

	return ConfigInfrastructureOpsConsumer{
		ChDrainRequest: chDrainRequest,
		ToConsume:      toConsume,
	}
}

// Data struct indicating what data structures available in this topic will be accepted.
type ConsumableStructsInfrastructureOpsConsumer struct {
	// Consume drain requests
	DrainRequest bool
}

func NewInfrastructureOpsConsumer(client bus.NalejClient, name string, exclusive bool, config ConfigInfrastructureOpsConsumer) (*InfrastructureOpsConsumer, derrors.Error) {
	consumer, err := client.BuildConsumer(name, InfrastructureOpsTopic, exclusive)
	if err != nil {
		return nil, err
	}

	return &InfrastructureOpsConsumer{Consumer: consumer, Config: config}, nil
}

// Consume any of the possible objects that can be sent to this queue and send it to the corresponding channel.
func (c InfrastructureOpsConsumer) Consume(ctx context.Context) derrors.Error {
	msg, err := c.Consumer.Receive(ctx)
	if err != nil {
		return err
	}

	target := &grpc_bus_go.InfrastructureOps{}

	derr := queue.UnmarshallPbMsg(msg, target)
	if derr != nil {
		return derr
	}

	switch x := target.Operation.(type) {
	case *grpc_bus_go.InfrastructureOps_DrainRequest:
		if c.Config.ToConsume.DrainRequest {
			c.Config.ChDrainRequest <- x.DrainRequest
		}
	case nil:
		errMsg := "received nil entry"
		log.Error().Msg(errMsg)
		return derrors.NewInvalidArgumentError(errMsg)
	default:
		errMsg := fmt.Sprintf("unknown object type in %s", InfrastructureOpsTopic)
		log.Error().Interface("type", x).Msg(errMsg)
		return derrors.NewInvalidArgumentError(errMsg)
	}
	return nil
}
