package controllers

import (
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ciscov1 "operator-cisco/api/v1"
)

func (r *CiscoCRDReconciler) createIngress(obj *ciscov1.CiscoCRD, labels map[string]string) *networkingv1.Ingress {
	ingressclassname := "nginx"
	prefix := networkingv1.PathType("Prefix")
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      labels["name"] + "-ingress",
			Namespace: labels["namespace"],
			Labels:    labels,
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion: obj.APIVersion,
				Kind:       obj.Kind,
				Name:       obj.Name,
				UID:        obj.UID,
			}},
		},
		//You can customize by adding extra fields from this package : https://pkg.go.dev/k8s.io/api/networking/v1
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: obj.Spec.Host,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &prefix,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: labels["name"],
											Port: networkingv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			IngressClassName: &ingressclassname,
		},
	}
	return ingress

}
