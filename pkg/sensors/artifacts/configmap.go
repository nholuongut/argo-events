package artifacts

import (
	sharedutil "github.com/nholuongut/argo-events/pkg/shared/util"
	corev1 "k8s.io/api/core/v1"
)

// ConfigMapReader implements the ArtifactReader interface for K8s configmap
type ConfigMapReader struct {
	configmapArtifact *corev1.ConfigMapKeySelector
}

// NewConfigMapReader returns a new configmap reader
func NewConfigMapReader(configmapArtifact *corev1.ConfigMapKeySelector) (*ConfigMapReader, error) {
	return &ConfigMapReader{
		configmapArtifact: configmapArtifact,
	}, nil
}

func (c *ConfigMapReader) Read() (body []byte, err error) {
	cm, err := sharedutil.GetConfigMapFromVolume(c.configmapArtifact)
	if err != nil {
		return nil, err
	}
	return []byte(cm), nil
}
