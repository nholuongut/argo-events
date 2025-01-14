/*
Copyright 2020 The Nho Luong

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package azureservicebus

import (
	"context"
	"encoding/json"
	"fmt"

	servicebus "github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"go.uber.org/zap"

	"github.com/nholuongut/argo-events/pkg/apis/events/v1alpha1"
	"github.com/nholuongut/argo-events/pkg/sensors/triggers"
	"github.com/nholuongut/argo-events/pkg/shared/logging"
	sharedutil "github.com/nholuongut/argo-events/pkg/shared/util"
)

// AzureServiceBusTrigger describes the trigger to send messages to a Service Bus
type AzureServiceBusTrigger struct {
	// Sensor object
	Sensor *v1alpha1.Sensor
	// Trigger reference
	Trigger *v1alpha1.Trigger
	// Sender refers to Azure Service Bus Sender struct
	Sender *servicebus.Sender
	// Logger to log stuff
	Logger *zap.SugaredLogger
}

func NewAzureServiceBusTrigger(sensor *v1alpha1.Sensor, trigger *v1alpha1.Trigger, azureServiceBusClients sharedutil.StringKeyedMap[*servicebus.Sender], logger *zap.SugaredLogger) (*AzureServiceBusTrigger, error) {
	triggerLogger := logger.With(logging.LabelTriggerType, v1alpha1.TriggerTypeAzureServiceBus)
	azureServiceBusTrigger := trigger.Template.AzureServiceBus

	sender, ok := azureServiceBusClients.Load(trigger.Template.Name)

	if !ok {
		connStr, err := sharedutil.GetSecretFromVolume(azureServiceBusTrigger.ConnectionString)
		if err != nil {
			triggerLogger.With("connection-string", azureServiceBusTrigger.ConnectionString.Name).Errorw("failed to retrieve connection string from secret", zap.Error(err))
			return nil, err
		}

		triggerLogger.Info("connecting to the service bus...")
		clientOptions := servicebus.ClientOptions{}
		if azureServiceBusTrigger.TLS != nil {
			tlsConfig, err := sharedutil.GetTLSConfig(azureServiceBusTrigger.TLS)
			if err != nil {
				triggerLogger.Errorw("failed to get the tls configuration", zap.Error(err))
				return nil, err
			}
			clientOptions.TLSConfig = tlsConfig
		}

		client, err := servicebus.NewClientFromConnectionString(connStr, &clientOptions)
		if err != nil {
			triggerLogger.Errorw("failed to create a service bus client", zap.Error(err))
			return nil, err
		}

		// Set queueOrTopicName to be azureServiceBusTrigger.QueueName or azureServiceBusTrigger.TopicName
		var queueOrTopicName string
		switch {
		case azureServiceBusTrigger.QueueName != "":
			queueOrTopicName = azureServiceBusTrigger.QueueName
		case azureServiceBusTrigger.TopicName != "":
			queueOrTopicName = azureServiceBusTrigger.TopicName
		default:
			return nil, fmt.Errorf("neither queue name nor topic name is specified")
		}

		logger.With("queueOrTopicName", queueOrTopicName).Info("creating a new sender...")

		sender, err = client.NewSender(queueOrTopicName, &servicebus.NewSenderOptions{})
		if err != nil {
			triggerLogger.Errorw("failed to create a service bus sender", zap.Error(err))
			return nil, err
		}

		azureServiceBusClients.Store(trigger.Template.Name, sender)
	}

	return &AzureServiceBusTrigger{
		Sensor:  sensor,
		Trigger: trigger,
		Sender:  sender,
		Logger:  triggerLogger,
	}, nil
}

// GetTriggerType returns the type of the trigger
func (t *AzureServiceBusTrigger) GetTriggerType() v1alpha1.TriggerType {
	return v1alpha1.TriggerTypeAzureServiceBus
}

// FetchResource fetches the trigger resource
func (t *AzureServiceBusTrigger) FetchResource(ctx context.Context) (interface{}, error) {
	return t.Trigger.Template.AzureServiceBus, nil
}

// ApplyResourceParameters applies parameters to the trigger resource
func (t *AzureServiceBusTrigger) ApplyResourceParameters(events map[string]*v1alpha1.Event, resource interface{}) (interface{}, error) {
	fetchedResource, ok := resource.(*v1alpha1.AzureServiceBusTrigger)
	if !ok {
		return nil, fmt.Errorf("failed to interpret the fetched trigger resource")
	}

	resourceBytes, err := json.Marshal(fetchedResource)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal the azure service bus trigger resource, %w", err)
	}

	parameters := fetchedResource.Parameters
	if parameters != nil {
		updatedResourceBytes, err := triggers.ApplyParams(resourceBytes, parameters, events)
		if err != nil {
			return nil, err
		}
		var sbTrigger *v1alpha1.AzureServiceBusTrigger
		if err := json.Unmarshal(updatedResourceBytes, &sbTrigger); err != nil {
			return nil, fmt.Errorf("failed to unmarshal the updated azure service bus trigger resource after applying resource parameters, %w", err)
		}
		return sbTrigger, nil
	}

	return resource, nil
}

// Execute executes the trigger
func (t *AzureServiceBusTrigger) Execute(ctx context.Context, events map[string]*v1alpha1.Event, resource interface{}) (interface{}, error) {
	trigger, ok := resource.(*v1alpha1.AzureServiceBusTrigger)
	if !ok {
		return nil, fmt.Errorf("failed to interpret the trigger resource")
	}

	if trigger.Payload == nil {
		return nil, fmt.Errorf("payload parameters are not specified")
	}

	payload, err := triggers.ConstructPayload(events, trigger.Payload)
	if err != nil {
		return nil, err
	}

	message := &servicebus.Message{
		Body: payload,
	}

	err = t.Sender.SendMessage(ctx, message, &servicebus.SendMessageOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to send a message to the service bus, %w", err)
	}
	t.Logger.Info("successfully sent message to the service bus")

	return nil, nil
}

// ApplyPolicy applies the trigger policy
func (t *AzureServiceBusTrigger) ApplyPolicy(ctx context.Context, resource interface{}) error {
	return nil
}
