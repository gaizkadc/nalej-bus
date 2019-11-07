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
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/nalej-bus/pkg/bus"
	"github.com/nalej/nalej-bus/pkg/queue"
	"github.com/rs/zerolog/log"
)

const (
	InventoryOpsTopic = "nalej/inventory/ops"
)

// Producer
type InventoryOpsProducer struct {
	producer bus.NalejProducer
}

// Create a new producer for the application operations topic
// params:
//  client to be used
//  name of the producer
// return:
//  built producer
func NewInventoryOpsProducer(client bus.NalejClient, name string) (*InventoryOpsProducer, derrors.Error) {
	prod, err := client.BuildProducer(name, InventoryOpsTopic)
	if err != nil {
		return nil, err
	}
	return &InventoryOpsProducer{producer: prod}, nil
}

// Generic function that decides what to send depending on the object type.
func (m InventoryOpsProducer) Send(ctx context.Context, msg proto.Message) derrors.Error {

	var wrapper grpc_bus_go.InventoryOps

	switch x := msg.(type) {
	case *grpc_inventory_manager_go.AgentOpResponse:
		wrapper = grpc_bus_go.InventoryOps{Operation: &grpc_bus_go.InventoryOps_AgentOpResponse{x}}
	case *grpc_inventory_manager_go.EdgeControllerOpResponse:
		wrapper = grpc_bus_go.InventoryOps{Operation: &grpc_bus_go.InventoryOps_EdgeControllerOpResponse{x}}
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

// Consumer

type InventoryOpsConsumer struct {
	Consumer bus.NalejConsumer
	Config   ConfigInventoryOpsConsumer
}

func NewInventoryOpsConsumer(client bus.NalejClient, name string, exclusive bool, config ConfigInventoryOpsConsumer) (*InventoryOpsConsumer, derrors.Error) {
	consumer, err := client.BuildConsumer(name, InventoryOpsTopic, exclusive)
	if err != nil {
		return nil, err
	}

	return &InventoryOpsConsumer{Consumer: consumer, Config: config}, nil
}

// Data struct indicating what data structures available in this topic will be accepted.
type ConsumableStructsInventoryOpsConsumer struct {
	// AgentOpResponse to indicate that the consumer wants to receive AgentOpResponse messages.
	AgentOpResponse bool
	// EdgeControllerOpResponse to indicate that the consumer wants to receive EdgeControllerOpResponse messages.
	EdgeControllerOpResponse bool
}

// Struct designed to config a consumer defining what actions to perform depending on the incoming object.
type ConfigInventoryOpsConsumer struct {
	// ChAgentOpResponse channel to receive operation Responses
	ChAgentOpResponse chan *grpc_inventory_manager_go.AgentOpResponse
	// ChEdgeControllerOpResponse channel to receive operation Responses
	ChEdgeControllerOpResponse chan *grpc_inventory_manager_go.EdgeControllerOpResponse
	// object types to be considered for consumption
	ToConsume ConsumableStructsInventoryOpsConsumer
}

// Create a new configuration structure for a given channel size
// params:
//  size of channels
// return:
//  instance of a configuration object
func NewConfigInventoryOpsConsumer(size int, toConsume ConsumableStructsInventoryOpsConsumer) ConfigInventoryOpsConsumer {
	chOpResponse := make(chan *grpc_inventory_manager_go.AgentOpResponse, size)
	chECOpResponse := make(chan *grpc_inventory_manager_go.EdgeControllerOpResponse, size)

	return ConfigInventoryOpsConsumer{
		ChAgentOpResponse:          chOpResponse,
		ChEdgeControllerOpResponse: chECOpResponse,
		ToConsume:                  toConsume,
	}
}

// Consume any of the possible objects that can be sent to this queue and send it to the corresponding channel.
func (c InventoryOpsConsumer) Consume(ctx context.Context) derrors.Error {
	msg, err := c.Consumer.Receive(ctx)
	if err != nil {
		return err
	}

	target := &grpc_bus_go.InventoryOps{}

	derr := queue.UnmarshallPbMsg(msg, target)
	if derr != nil {
		return derr
	}

	switch x := target.Operation.(type) {
	case *grpc_bus_go.InventoryOps_AgentOpResponse:
		if c.Config.ToConsume.AgentOpResponse {
			c.Config.ChAgentOpResponse <- x.AgentOpResponse
		}
	case *grpc_bus_go.InventoryOps_EdgeControllerOpResponse:
		if c.Config.ToConsume.EdgeControllerOpResponse {
			c.Config.ChEdgeControllerOpResponse <- x.EdgeControllerOpResponse
		}
	case nil:
		errMsg := "received nil entry"
		log.Error().Msg(errMsg)
		return derrors.NewInvalidArgumentError(errMsg)
	default:
		errMsg := fmt.Sprintf("unknown object type in %s", InventoryOpsTopic)
		log.Error().Interface("type", x).Msg(errMsg)
		return derrors.NewInvalidArgumentError(errMsg)
	}
	return nil
}
