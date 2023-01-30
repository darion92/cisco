/*
Copyright 2023.

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

package controllers

import (
	"context"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
	ciscov1 "operator-cisco/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	kindDeployment = "Deployment"
)

// CiscoCRDReconciler reconciles a CiscoCRD object
type CiscoCRDReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cisco.kind-kind,resources=ciscocrds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cisco.kind-kind,resources=ciscocrds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cisco.kind-kind,resources=ciscocrds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CiscoCRD object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *CiscoCRDReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Add a uuid for each reconciliation
	log := log.FromContext(ctx).WithValues("reconcileID", uuid.NewUUID())

	// Add the controller logger to the context
	ctx = ctrl.LoggerInto(ctx, log)

	obj := &ciscov1.CiscoCRD{}
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		log.Error(err, "unable to fetch AbstractWorkload")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	labels := map[string]string{
		"name":      req.NamespacedName.Name,
		"namespace": req.NamespacedName.Namespace,
	}

	// 1) Handle Deployment
	// Do we have a current lolcow instance deployed?
	existingD := &v1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: obj.Name, Namespace: obj.Namespace}, existingD)
	if err != nil {

		// Case 1: not found yet, check if deployment needs deletion
		if errors.IsNotFound(err) {
			dep := r.createDeployment(obj, labels)
			log.Info("✨ Creating a new Deployment ✨", "Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			err = r.Create(ctx, dep)
			if err != nil {
				log.Error(err, "❌ Failed to create new Deployment", "Namespace", dep.Namespace, "Name", dep.Name)
				return ctrl.Result{}, err
			}
			// Deployment created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Deployment")
			return ctrl.Result{}, err
		}
	}

	// If we get down here, we have a current instance deployed (not freshly created)
	replicas := obj.Spec.Replicas
	replicasSpec := existingD.Spec.Replicas
	if *replicas != *replicasSpec {
		log.Info("👋️ Update Deployment 👋️: ")

		existingD.Spec.Replicas = replicas
		err = r.Update(ctx, existingD)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Namespace", existingD.Namespace, "Name", existingD.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else {
		log.Info("👋️ No Change For deployment ! 👋️: ")
	}

	//2) Handle Service
	// Do we have a current service deployed?
	existingS := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: obj.Name, Namespace: obj.Namespace}, existingS)
	if err != nil {
		if errors.IsNotFound(err) {
			service := r.createService(obj, labels)
			log.Info("✨ Creating a new Service ✨", "Namespace", service.Namespace, "Name", service.Name)
			err = r.Create(ctx, service)
			if err != nil {
				log.Error(err, "❌ Failed to create new Service", "Namespace", service.Namespace, "Name", service.Name)
				return ctrl.Result{}, err
			}
			// Service created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Service")
			return ctrl.Result{}, err
		}

		// We found the service, but did the port change (and needs to be updated?)
	}

	//3) Handle Ingress
	// Do we have a current ingress deployed?
	ingress := &networkingv1.Ingress{}
	existingIngress := &networkingv1.Ingress{}
	err = r.Get(ctx, types.NamespacedName{Name: obj.Name + "-ingress", Namespace: obj.Namespace}, existingIngress)
	if err != nil {
		if errors.IsNotFound(err) {
			ingress := r.createIngress(obj, labels)
			log.Info("✨ Creating a new Ingress ✨", "Namespace", ingress.Namespace, "Name", ingress.Name)
			err = r.Create(ctx, ingress)
			if err != nil {
				log.Error(err, "❌ Failed to create new Ingress", "Namespace", ingress.Namespace, "Name", ingress.Name)
				return ctrl.Result{}, err
			}
			// Ingress created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Ingress")
			return ctrl.Result{}, err
		}
	}

	if obj.Spec.Host != existingIngress.Spec.Rules[0].Host {
		log.Info("👋️ Update Ingress 👋️: ")
		host := obj.Spec.Host
		existingIngress.Spec.Rules[0].Host = host
		existingIngress.ResourceVersion = ingress.ResourceVersion
		existingIngress.Finalizers = ingress.Finalizers
		r.Client.Update(ctx, existingIngress)
	} else {
		log.Info("👋️ No Change For Ingress ! 👋️: ")
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CiscoCRDReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ciscov1.CiscoCRD{}).
		Owns(&v1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&networkingv1.Ingress{}).
		Complete(r)
}
