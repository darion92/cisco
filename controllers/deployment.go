package controllers

import (
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	ciscov1 "operator-cisco/api/v1"
)

//You can customize by adding extra fields from this package : https://pkg.go.dev/
func (r *CiscoCRDReconciler) createDeployment(obj *ciscov1.CiscoCRD, labels map[string]string) *v1.Deployment {
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      labels["name"],
			Namespace: labels["namespace"],
			Labels:    labels,
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion: obj.APIVersion,
				Kind:       obj.Kind,
				Name:       obj.Name,
				UID:        obj.UID,
			}},
		},
		Spec: v1.DeploymentSpec{
			Replicas: obj.Spec.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: v12.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: v12.PodSpec{
					Containers: []v12.Container{{
						Name:  labels["name"],
						Image: obj.Spec.ContainerImage,
					}},
				},
			},
		},
	}
	ctrl.SetControllerReference(obj, deployment, r.Scheme)
	return deployment

}
