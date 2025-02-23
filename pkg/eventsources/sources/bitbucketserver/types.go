/*

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

package bitbucketserver

import (
	"github.com/nholuongut/argo-events/pkg/apis/events/v1alpha1"
	"github.com/nholuongut/argo-events/pkg/eventsources/common/webhook"
	metrics "github.com/nholuongut/argo-events/pkg/metrics"
	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
)

// EventListener implements ConfigExecutor
type EventListener struct {
	EventSourceName            string
	EventName                  string
	BitbucketServerEventSource v1alpha1.BitbucketServerEventSource
	Metrics                    *metrics.Metrics
}

// GetEventSourceName returns name of event source
func (el *EventListener) GetEventSourceName() string {
	return el.EventSourceName
}

// GetEventName returns name of event
func (el *EventListener) GetEventName() string {
	return el.EventName
}

// GetEventSourceType return type of event server
func (el *EventListener) GetEventSourceType() v1alpha1.EventSourceType {
	return v1alpha1.BitbucketServerEvent
}

// Router contains the configuration information for a route
type Router struct {
	// route contains information about a API endpoint
	route *webhook.Route
	// client is the bitbucket server client
	client *bitbucketv1.APIClient
	// customClient is a custom bitbucket server client which implements a method the gfleury/go-bitbucket-v1 client is missing
	customClient *customBitbucketServerClient
	// deleteClient is used to delete webhooks. This client does not contain the cancelable context of the default client
	deleteClient *bitbucketv1.APIClient
	// hookIDs is a map of webhook IDs
	// (projectKey + "," + repoSlug) -> hook ID
	// Bitbucket Server API docs:
	// https://developer.atlassian.com/server/bitbucket/reference/rest-api/
	hookIDs map[string]int
	// hookSecret is a Bitbucket Server webhook secret
	hookSecret string
	// bitbucketServerEventSource is the event source that contains configuration necessary to consume events from Bitbucket Server
	bitbucketServerEventSource *v1alpha1.BitbucketServerEventSource
}
