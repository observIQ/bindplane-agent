package azure

import (
	"context"
	"fmt"
	"sync"

	azhub "github.com/Azure/azure-event-hubs-go/v3"
	"go.uber.org/zap"
)

// Eventhub provides methods for reading events from Azure Event Hub.
type EventHub struct {
	AzureConfig
	Persist *Persister
	WG      sync.WaitGroup
	Handler func(context.Context, *azhub.Event) error

	hub *azhub.Hub
}

// StartConsumers starts an Azure Event Hub handler for each partition_id.
func (e *EventHub) StartConsumers(ctx context.Context) error {
	if err := e.Connect(ctx); err != nil {
		return err
	}

	runtimeInfo, err := e.hub.GetRuntimeInformation(ctx)
	if err != nil {
		return err
	}

	for _, partitionID := range runtimeInfo.PartitionIDs {
		if err := e.startConsumer(ctx, partitionID, e.hub); err != nil {
			return err
		}
		e.Debugw(fmt.Sprintf("Successfully connected to Azure Event Hub '%s' partition_id '%s'", e.Name, partitionID))
	}
	return nil
}

// StopConsumers closes connections to Azure Event Hub.
func (e *EventHub) StopConsumers() error {
	if e.hub == nil {
		return nil
	}
	e.WG.Wait()
	if err := e.hub.Close(context.Background()); err != nil {
		return err
	}
	e.Debugw(fmt.Sprintf("Closed all connections to Azure Event Hub '%s'", e.Name))
	return nil
}

// startConsumer starts polling an Azure Event Hub partition id for new events
func (e *EventHub) startConsumer(ctx context.Context, partitionID string, hub *azhub.Hub) error {
	if e.startAtBeginning {
		_, err := hub.Receive(
			ctx, partitionID, e.Handler, azhub.ReceiveWithStartingOffset(""),
			azhub.ReceiveWithPrefetchCount(e.PrefetchCount))
		return err
	}

	offset, err := e.Persist.Read(e.Namespace, e.Name, e.Group, partitionID)
	if err != nil {
		x := fmt.Sprintf("Error while reading offset for partition_id %s", partitionID)
		e.Debugw(x, zap.Error(err))
	}

	// start at end and no offset was found
	if offset.Offset == "" {
		e.Debugw("No offset found, starting from end")
		_, err := hub.Receive(
			ctx, partitionID, e.Handler, azhub.ReceiveWithLatestOffset(),
			azhub.ReceiveWithPrefetchCount(e.PrefetchCount))
		return err
	}

	// start at end and offset exists
	e.Debugw(fmt.Sprintf("Starting with offset '%s'", offset.Offset))
	_, err = hub.Receive(
		ctx, partitionID, e.Handler, azhub.ReceiveWithStartingOffset(offset.Offset),
		azhub.ReceiveWithPrefetchCount(e.PrefetchCount))
	return err
}

// Connect initializes the connection to Azure Event Hub ensures the input parameters are valid
func (e *EventHub) Connect(ctx context.Context) error {
	var err error
	e.hub, err = azhub.NewHubFromConnectionString(e.ConnectionString, azhub.HubWithOffsetPersistence(e.Persist))
	return err
}
