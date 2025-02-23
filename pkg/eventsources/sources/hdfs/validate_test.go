package hdfs

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
	assert.Equal(t, "type is required", err.Error())

	content, err := os.ReadFile(fmt.Sprintf("%s/%s", sources.EventSourceDir, "hdfs.yaml"))
	assert.Nil(t, err)

	var eventSource *v1alpha1.EventSource
	err = yaml.Unmarshal(content, &eventSource)
	assert.Nil(t, err)
	assert.NotNil(t, eventSource.Spec.HDFS)

	for name, value := range eventSource.Spec.HDFS {
		fmt.Println(name)
		l := &EventListener{
			HDFSEventSource: value,
		}
		err := l.ValidateEventSource(context.Background())
		assert.NoError(t, err)
	}
}
