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
	"fmt"

	"github.com/jmckind/grackle-operator/pkg/apis/k8s/v1alpha1"
	k8sv1alpha1 "github.com/jmckind/grackle-operator/pkg/apis/k8s/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IngestTrack defines the value to track for a Twitter stream,
// which is a string containing a comma-delimited list of search terms.
type IngestTrack struct {
	// Checksum holds the checksum for the value.
	Checksum string `json:"checksum,omitempty"`

	// Value holds the actual search track value.
	Value string `json:"value,omitempty"`
}

// NewIngestTrack will create a new IngestTrack from the given value.
func NewIngestTrack(track string) *IngestTrack {
	return &IngestTrack{
		Checksum: hashValue(track),
		Value:    track,
	}
}

// newGrackleIngestPod returns a grackle-ingest pod with the same name/namespace as the CR.
func newGrackleIngestPod(cr *k8sv1alpha1.Grackle, track *IngestTrack) *corev1.Pod {
	annotations := make(map[string]string)
	annotations[AnnotationIngestHash] = track.Checksum

	labels := labelsForCluster(cr)
	labels[LabelComponentKey] = LabelComponentIngest

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Annotations:  annotations,
			GenerateName: fmt.Sprintf("%s-%s-", cr.Name, LabelComponentIngest),
			Namespace:    cr.Namespace,
			Labels:       labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    LabelComponentIngest,
					Image:   fmt.Sprintf("%s:%s", DefaultGrackleImageName, cr.Spec.Ingest.Version),
					Command: []string{"/grackle-ingest"},
					Env: []corev1.EnvVar{{
						Name:  "GRK_RETHINKDB_HOST",
						Value: cr.Spec.Datastore.Host,
					}, {
						Name: "GRK_TWITTER_ACCESS_TOKEN",
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{Name: cr.Spec.Ingest.TwitterSecret},
								Key:                  "twitter-access-token",
							},
						},
					}, {
						Name: "GRK_TWITTER_ACCESS_SECRET",
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{Name: cr.Spec.Ingest.TwitterSecret},
								Key:                  "twitter-access-secret",
							},
						},
					}, {
						Name: "GRK_TWITTER_CONSUMER_KEY",
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{Name: cr.Spec.Ingest.TwitterSecret},
								Key:                  "twitter-consumer-key",
							},
						},
					}, {
						Name: "GRK_TWITTER_CONSUMER_SECRET",
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{Name: cr.Spec.Ingest.TwitterSecret},
								Key:                  "twitter-consumer-secret",
							},
						},
					}, {
						Name:  "GRK_TWITTER_TRACK",
						Value: track.Value,
					}},
				},
			},
		},
	}
}

// newIngestPodEvent returns an event for
func newIngestPodEvent(cr *v1alpha1.Grackle, track string) *corev1.Event {
	event := newEvent(cr)
	event.Type = corev1.EventTypeNormal
	event.Reason = "New Ingest Pod"
	event.Message = fmt.Sprintf("New ingest pod added to track: %s", track)
	return event
}
