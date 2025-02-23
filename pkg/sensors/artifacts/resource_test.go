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
	"testing"

	"github.com/nholuongut/argo-events/pkg/apis/events/v1alpha1"
	"github.com/smartystreets/goconvey/convey"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestNewResourceReader(t *testing.T) {
	convey.Convey("Given a resource, get new reader", t, func() {
		un := unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Secret",
				"metadata": map[string]string{
					"name": "mysecret",
				},
				"type": "Opaque",
				"data": map[string]string{
					"access": "c2VjcmV0",
					"secret": "c2VjcmV0",
				},
			},
		}
		artifact := v1alpha1.NewK8SResource(un)
		reader, err := NewResourceReader(&artifact)
		convey.So(err, convey.ShouldBeNil)
		convey.So(reader, convey.ShouldNotBeNil)

		data, err := reader.Read()
		convey.So(err, convey.ShouldBeNil)
		convey.So(data, convey.ShouldNotBeNil)
	})
}
