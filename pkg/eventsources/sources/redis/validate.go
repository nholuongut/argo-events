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
package redis

import (
	"context"
	"fmt"

	"github.com/nholuongut/argo-events/pkg/apis/events/v1alpha1"
)

// ValidateEventSource validates nats event source
func (listener *EventListener) ValidateEventSource(ctx context.Context) error {
	return validate(&listener.RedisEventSource)
}

func validate(eventSource *v1alpha1.RedisEventSource) error {
	if eventSource == nil {
		return v1alpha1.ErrNilEventSource
	}
	if eventSource.HostAddress == "" {
		return fmt.Errorf("host address must be specified")
	}
	if eventSource.Channels == nil {
		return fmt.Errorf("channel/s must be specified")
	}
	if eventSource.Password != nil && eventSource.Namespace == "" {
		return fmt.Errorf("namespace must be defined in order to retrieve the password from the secret")
	}
	if eventSource.TLS != nil {
		return v1alpha1.ValidateTLSConfig(eventSource.TLS)
	}
	return nil
}
