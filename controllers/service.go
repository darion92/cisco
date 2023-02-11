package controllers

import (
	ciscov1 "operator-cisco/api/v1"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// You can customize by adding extra fields from this package : https://pkg.go.dev/
func (r *CiscoCRDReconciler) createService(obj *ciscov1.CiscoCRD, labels map[string]string) *v1.Service {
	service := &v1.Service{
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
		Spec: v1.ServiceSpec{
			Selector: labels,
			Ports: []v1.ServicePort{
				{
					Name:       "httpport",
					Port:       80,
					TargetPort: intstr.FromInt(80),
				},
			},
			//if you're in a Prod cluster set it to LoadBalancer
			Type: v1.ServiceTypeClusterIP,
		},
	}
	return service
}
