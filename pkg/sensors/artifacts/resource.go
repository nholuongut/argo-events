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

package artifacts

import (
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/nholuongut/argo-events/pkg/apis/events/v1alpha1"
	"github.com/nholuongut/argo-events/pkg/shared/logging"
)

// ResourceReader implements the ArtifactReader interface for resource artifacts
type ResourceReader struct {
	resourceArtifact *unstructured.Unstructured
}

// NewResourceReader creates a new ArtifactReader for resource
func NewResourceReader(resourceArtifact *v1alpha1.K8SResource) (ArtifactReader, error) {
	if resourceArtifact == nil {
		return nil, fmt.Errorf("ResourceArtifact does not exist")
	}
	data, err := json.Marshal(resourceArtifact)
	if err != nil {
		return nil, err
	}
	object := make(map[string]interface{})
	err = json.Unmarshal(data, &object)
	if err != nil {
		return nil, err
	}
	un := &unstructured.Unstructured{Object: object}
	return &ResourceReader{un}, nil
}

func (reader *ResourceReader) Read() ([]byte, error) {
	log := logging.NewArgoEventsLogger()
	log.Debugw("reading artifact from resource template", "resource", reader.resourceArtifact.Object)
	return yaml.Marshal(reader.resourceArtifact.Object)
}
