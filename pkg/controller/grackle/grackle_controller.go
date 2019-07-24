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
	"context"
	"fmt"

	k8sv1alpha1 "github.com/jmckind/grackle-operator/pkg/apis/k8s/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8slabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_grackle")

// Add creates a new Grackle Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileGrackle{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("grackle-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Grackle
	err = c.Watch(&source.Kind{Type: &k8sv1alpha1.Grackle{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Grackle
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &k8sv1alpha1.Grackle{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileGrackle implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileGrackle{}

// ReconcileGrackle reconciles a Grackle object
type ReconcileGrackle struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Grackle object and makes changes based on the state read
// and what is in the Grackle.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileGrackle) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("grackle reconciliation started")

	// Fetch the Grackle instance
	grackle, err := r.fetchGrackleInstance(request)
	if err != nil {
		reqLogger.Error(err, "unable to fetch grackle instance")
		return reconcile.Result{}, err
	} else if grackle == nil {
		// Request object not found, could have been deleted after reconcile request.
		reqLogger.Info("grackle instance not found")
		return reconcile.Result{}, nil
	}

	// Set default values and requeue the request if needed.
	if r.setDefaults(grackle) {
		reqLogger.Info("grackle instance initialized")
		err = r.client.Update(context.TODO(), grackle)
		if err != nil {
			reqLogger.Error(err, "unable to update grackle instance")
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}

	// Reconcile Ingest process
	if err := r.reconcileIngest(grackle); err != nil {
		reqLogger.Error(err, "unable to reconcile ingest resources")
		return reconcile.Result{}, err
	}

	// Reconcile Web UI process
	if err := r.reconcileWeb(grackle); err != nil {
		reqLogger.Error(err, "unable to reconcile web ui resources")
		return reconcile.Result{}, err
	}

	reqLogger.Info("grackle reconciliation complete")
	return reconcile.Result{}, nil
}

// fetchGrackleInstance will retun the Grackle instance for the given request.
func (r *ReconcileGrackle) fetchGrackleInstance(request reconcile.Request) (*k8sv1alpha1.Grackle, error) {
	grackle := &k8sv1alpha1.Grackle{}
	err := r.client.Get(context.TODO(), request.NamespacedName, grackle)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Return nil with no error
			return nil, nil
		}
		// Error reading the object.
		return nil, err
	}
	return grackle, nil
}

// reconcileIngest will reconcile the Grackle ingest process.
func (r *ReconcileGrackle) reconcileIngest(cr *k8sv1alpha1.Grackle) error {
	ingest := cr.Spec.Ingest
	if ingest == nil || len(ingest.Track) <= 0 {
		return nil
	}

	labels := labelsForCluster(cr)
	labels["component"] = "ingest"
	labelSelector := k8slabels.SelectorFromSet(labels)
	listOpts := &client.ListOptions{Namespace: cr.Namespace, LabelSelector: labelSelector}

	podList := &corev1.PodList{}
	err := r.client.List(context.TODO(), listOpts, podList)
	if err != nil {
		// Unable to list pods for some reason...
		return err
	}

	// Generate hashes for each search track.
	tracks := make(map[string]*IngestTrack)
	for _, track := range ingest.Track {
		hash := hashValue(track)
		tracks[hash] = NewIngestTrack(track)
	}

	// Remove tracks that have already been assigned to pods.
	for _, pod := range podList.Items {
		delete(tracks, pod.ObjectMeta.Annotations[AnnotationIngestHash])
	}

	if len(tracks) <= 0 {
		// Nothing remains...
		return nil
	}

	// Create pod for next remaining hash.
	var track *IngestTrack
	for _, track = range tracks {
		break
	}

	pod := newGrackleIngestPod(cr, track)
	controllerutil.SetControllerReference(cr, pod, r.scheme)

	err = r.client.Create(context.TODO(), pod)
	if err != nil {
		// Unable to create pod for some reason...
		return err
	}

	return nil
}

// reconcileWeb will reconcile the Grackle web UI process.
func (r *ReconcileGrackle) reconcileWeb(cr *k8sv1alpha1.Grackle) error {
	depFound := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: fmt.Sprintf("%s-web", cr.Name), Namespace: cr.Namespace}, depFound)
	if err != nil && errors.IsNotFound(err) {
		// Not found, create a new deployment.
		deployment := newGrackleWebDeployment(cr)
		controllerutil.SetControllerReference(cr, deployment, r.scheme)

		err = r.client.Create(context.TODO(), deployment)
		if err != nil {
			// Unable to create deployment for some reason...
			return err
		}
		// Deployment created successfully...
		return nil
	} else if err != nil {
		// Unable to get deployment for some reason...
		return err
	}

	svcFound := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, svcFound)
	if err != nil && errors.IsNotFound(err) {
		// Not found, create a new service.
		service := newGrackleWebService(cr)
		controllerutil.SetControllerReference(cr, service, r.scheme)

		err = r.client.Create(context.TODO(), service)
		if err != nil {
			// Unable to create service for some reason...
			return err
		}
		// Service created successfully...
		return nil
	} else if err != nil {
		// Unable to get service for some reason...
		return err
	}

	// All resources exist...
	return nil
}

// setDefaults will set the default values for any "required" properties.
func (r *ReconcileGrackle) setDefaults(cr *k8sv1alpha1.Grackle) bool {
	changed := false

	// Defaults for Ingest
	if cr.Spec.Ingest == nil {
		cr.Spec.Ingest = &k8sv1alpha1.IngestSpec{
			Version: "latest",
		}
		changed = true
	} else {
		if len(cr.Spec.Ingest.Version) <= 0 {
			cr.Spec.Ingest.Version = "latest"
			changed = true
		}
	}

	// Defaults for Web UI
	if cr.Spec.Web == nil {
		cr.Spec.Web = &k8sv1alpha1.WebSpec{
			Replicas: &DefaultWebReplicas,
			Version:  "latest",
		}
		changed = true
	} else {
		if cr.Spec.Web.Replicas == nil {
			cr.Spec.Web.Replicas = &DefaultWebReplicas
			changed = true
		}
		if len(cr.Spec.Web.Version) <= 0 {
			cr.Spec.Ingest.Version = "latest"
			changed = true
		}
	}

	return changed
}
