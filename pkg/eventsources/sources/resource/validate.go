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

package resource

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/selection"

	"github.com/nholuongut/argo-events/pkg/apis/events/v1alpha1"
)

// ValidateEventSource validates a resource event source
func (listener *EventListener) ValidateEventSource(ctx context.Context) error {
	return validate(&listener.ResourceEventSource)
}

func validate(eventSource *v1alpha1.ResourceEventSource) error {
	if eventSource == nil {
		return v1alpha1.ErrNilEventSource
	}
	if eventSource.Version == "" {
		return fmt.Errorf("version must be specified")
	}
	if eventSource.Resource == "" {
		return fmt.Errorf("resource must be specified")
	}
	if eventSource.EventTypes == nil {
		return fmt.Errorf("event types must be specified")
	}
	if eventSource.Filter != nil {
		if eventSource.Filter.Labels != nil {
			if err := validateSelectors(eventSource.Filter.Labels); err != nil {
				return err
			}
		}
		if eventSource.Filter.Fields != nil {
			if err := validateSelectors(eventSource.Filter.Fields); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateSelectors(selectors []v1alpha1.Selector) error {
	for _, sel := range selectors {
		if sel.Key == "" {
			return fmt.Errorf("key can't be empty for selector")
		}
		if sel.Operation == "" {
			continue
		}
		if selection.Operator(sel.Operation) == "" {
			return fmt.Errorf("unknown selection operation %s", sel.Operation)
		}
	}
	return nil
}
