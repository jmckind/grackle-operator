// Copyright 2019 Grackle Operator authors

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grackle

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/jmckind/grackle-operator/pkg/apis/k8s/v1alpha1"
)

var (
	// AnnotationIngestHash is the annotation key for the ingest hash value
	AnnotationIngestHash = "k8s.mkz.io/grackle-ingest-hash"

	// DefaultWebReplicas is the number of Web UI pods to create by default.
	DefaultWebReplicas int32 = 2
)

// defaultLabels returns the default set of labels for the cluster.
func defaultLabels(cr *v1alpha1.Grackle) map[string]string {
	return map[string]string{
		"app":     "grackle",
		"cluster": cr.Name,
	}
}

// labelsForCluster returns the labels for all cluster resources.
func labelsForCluster(cr *v1alpha1.Grackle) map[string]string {
	labels := defaultLabels(cr)
	for key, val := range cr.ObjectMeta.Labels {
		labels[key] = val
	}
	return labels
}

// hashValue with return a URL encoded SHA256 hash of the given value.
func hashValue(value string) string {
	hasher := sha256.New()
	hasher.Write([]byte(value))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
