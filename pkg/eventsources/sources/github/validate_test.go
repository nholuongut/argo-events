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

package github

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/nholuongut/argo-events/pkg/apis/events/v1alpha1"
	"github.com/nholuongut/argo-events/pkg/eventsources/sources"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
)

func TestValidateEventSource(t *testing.T) {
	listener := &EventListener{}
	err := listener.ValidateEventSource(context.Background())
	assert.Error(t, err)
	assert.Equal(t, "either repositories or organizations is required", err.Error())

	content, err := os.ReadFile(fmt.Sprintf("%s/%s", sources.EventSourceDir, "github.yaml"))
	assert.Nil(t, err)

	var eventSource *v1alpha1.EventSource
	err = yaml.Unmarshal(content, &eventSource)
	assert.Nil(t, err)
	assert.NotNil(t, eventSource.Spec.Github)

	for name, value := range eventSource.Spec.Github {
		fmt.Println(name)
		l := &EventListener{
			GithubEventSource: value,
		}
		err := l.ValidateEventSource(context.Background())
		assert.NoError(t, err)
	}
}
