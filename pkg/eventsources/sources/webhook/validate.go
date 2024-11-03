/*
Copyright 2018 The Nho Luong

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

package webhook

import (
	"context"

	"github.com/nholuongut/argo-events/pkg/apis/events/v1alpha1"
	"github.com/nholuongut/argo-events/pkg/eventsources/common/webhook"
)

// ValidateEventSource validates webhook event source
func (listener *EventListener) ValidateEventSource(ctx context.Context) error {
	return validate(&listener.Webhook)
}

func validate(webhookEventSource *v1alpha1.WebhookEventSource) error {
	if webhookEventSource == nil {
		return v1alpha1.ErrNilEventSource
	}
	return webhook.ValidateWebhookContext(&webhookEventSource.WebhookContext)
}
