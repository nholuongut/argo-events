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
	"fmt"
	"os"

	"github.com/nholuongut/argo-events/pkg/apis/events/v1alpha1"
	"github.com/nholuongut/argo-events/pkg/shared/logging"
)

// FileReader implements the ArtifactReader interface for file artifacts
type FileReader struct {
	fileArtifact *v1alpha1.FileArtifact
}

// NewFileReader creates a new ArtifactReader for inline
func NewFileReader(fileArtifact *v1alpha1.FileArtifact) (ArtifactReader, error) {
	// This should never happen!
	if fileArtifact == nil {
		return nil, fmt.Errorf("FileArtifact cannot be empty")
	}
	return &FileReader{fileArtifact}, nil
}

func (reader *FileReader) Read() ([]byte, error) {
	content, err := os.ReadFile(reader.fileArtifact.Path)
	if err != nil {
		return nil, err
	}
	log := logging.NewArgoEventsLogger()
	log.Debugf("reading fileArtifact from %s", reader.fileArtifact.Path)
	return content, nil
}
