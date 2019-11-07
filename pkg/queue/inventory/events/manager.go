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
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-bus-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/nalej-bus/pkg/bus"
	"github.com/nalej/nalej-bus/pkg/queue"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

const (
	InventoryEventsTopic = "nalej/inventory/events"
)

type InventoryEventsProducer struct {
	producer bus.NalejProducer
}

// Create a new producer for the application operations topic
// params:
//  client to be used
//  name of the producer
// return:
//  built producer
func NewInventoryEventsProducer(client bus.NalejClient, name string) (*InventoryEventsProducer, derrors.Error) {
	prod, err := client.BuildProducer(name, InventoryEventsTopic)
	if err != nil {
		return nil, err
	}
	return &InventoryEventsProducer{producer: prod}, nil
}

// Generic function that decides what to send depending on the object type.
func (m InventoryEventsProducer) Send(ctx context.Context, msg proto.Message) derrors.Error {

	var wrapper grpc_bus_go.InventoryEvents

	switch x := msg.(type) {
	case *grpc_inventory_manager_go.AgentsAlive:
		wrapper = grpc_bus_go.InventoryEvents{
			Event: &grpc_bus_go.InventoryEvents_AgentsAlive{x}}
	case *grpc_inventory_go.EdgeControllerId:
		wrapper = grpc_bus_go.InventoryEvents{
			Event: &grpc_bus_go.InventoryEvents_EdgeControllerId{x}}
	case *grpc_inventory_manager_go.EICStartInfo:
		wrapper = grpc_bus_go.InventoryEvents{
			Event: &grpc_bus_go.InventoryEvents_EicStart{x}}
	case *grpc_inventory_go.AssetUninstalledId:
		wrapper = grpc_bus_go.InventoryEvents{
			Event: &grpc_bus_go.InventoryEvents_AssetUninstalledId{x}}
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

// InventoryEventsConsumer to receive message from the InventoryEvents topic
type InventoryEventsConsumer struct {
	Consumer bus.NalejConsumer
	Config   ConfigInventoryEventsConsumer
}

// Struct designed to config a consumer defining what actions to perform depending on the incoming object.
type ConfigInventoryEventsConsumer struct {
	// ChAgentsIds with the channel to receive agent ids
	ChAgentsAlive chan *grpc_inventory_manager_go.AgentsAlive
	// ChEdgeControllerId with the channel to receive edge inventory ids
	ChEdgeControllerId chan *grpc_inventory_go.EdgeControllerId
	// ChEICStart with the channel to receive EIC start notifications.
	ChEICStart chan *grpc_inventory_manager_go.EICStartInfo
	// ChUninstalledAssetId with the channel to receive asset ids
	ChUninstalledAssetId chan *grpc_inventory_go.AssetUninstalledId
	// object types to be considered for consumption
	ToConsume ConsumableStructsInventoryEventsConsumer
}

// Data struct indicating what data structures available in this topic will be accepted.
type ConsumableStructsInventoryEventsConsumer struct {
	// AgentsAlive consumption
	AgentsAclive bool
	// EdgeControllerId consumption
	EdgeControllerId bool
	// EICStartInfo consumption
	EICStartInfo bool
	// UninstalledAssetId consumption
	UninstalledAssetId bool
}

// NewConfigInventoryEventsConsumer creates the underlying channels with a given configuration.
func NewConfigInventoryEventsConsumer(size int, toConsume ConsumableStructsInventoryEventsConsumer) ConfigInventoryEventsConsumer {

	chAgentsAlive := make(chan *grpc_inventory_manager_go.AgentsAlive, size)
	chEdgeControllerId := make(chan *grpc_inventory_go.EdgeControllerId, size)
	chEICStart := make(chan *grpc_inventory_manager_go.EICStartInfo, size)
	ChUninstalledAssetId := make(chan *grpc_inventory_go.AssetUninstalledId, size)

	return ConfigInventoryEventsConsumer{
		ChAgentsAlive:        chAgentsAlive,
		ChEdgeControllerId:   chEdgeControllerId,
		ChEICStart:           chEICStart,
		ChUninstalledAssetId: ChUninstalledAssetId,
		ToConsume:            toConsume,
	}
}

// NewInventoryEventsConsumer creates a consumer with a given string and config. Exclusivity determines how many consumer of a given
// name will receive the message.
func NewInventoryEventsConsumer(client bus.NalejClient, name string, exclusive bool, config ConfigInventoryEventsConsumer) (*InventoryEventsConsumer, derrors.Error) {
	consumer, err := client.BuildConsumer(name, InventoryEventsTopic, exclusive)
	if err != nil {
		return nil, err
	}

	return &InventoryEventsConsumer{Consumer: consumer, Config: config}, nil
}

// Consume any of the possible objects that can be sent to this queue and send it to the corresponding channel.
func (c InventoryEventsConsumer) Consume(ctx context.Context) derrors.Error {
	msg, err := c.Consumer.Receive(ctx)
	if err != nil {
		return err
	}

	target := &grpc_bus_go.InventoryEvents{}

	derr := queue.UnmarshallPbMsg(msg, target)
	if derr != nil {
		return derr
	}

	switch x := target.Event.(type) {
	case *grpc_bus_go.InventoryEvents_AgentsAlive:
		if c.Config.ToConsume.AgentsAclive {
			c.Config.ChAgentsAlive <- x.AgentsAlive
		}
	case *grpc_bus_go.InventoryEvents_EdgeControllerId:
		if c.Config.ToConsume.EdgeControllerId {
			c.Config.ChEdgeControllerId <- x.EdgeControllerId
		}
	case *grpc_bus_go.InventoryEvents_EicStart:
		if c.Config.ToConsume.EICStartInfo {
			c.Config.ChEICStart <- x.EicStart
		}
	case *grpc_bus_go.InventoryEvents_AssetUninstalledId:
		if c.Config.ToConsume.UninstalledAssetId {
			c.Config.ChUninstalledAssetId <- x.AssetUninstalledId
		}

	case nil:
		errMsg := "received nil entry"
		log.Error().Msg(errMsg)
		return derrors.NewInvalidArgumentError(errMsg)
	default:
		errMsg := fmt.Sprintf("unknown object type in %s", InventoryEventsTopic)
		log.Error().Interface("type", x).Msg(errMsg)
		return derrors.NewInvalidArgumentError(errMsg)
	}
	return nil
}
