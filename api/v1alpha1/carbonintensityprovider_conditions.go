package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ConditionHealthy = metav1.Condition{
		Type:   "Available",
		Status: metav1.ConditionUnknown,
		Reason: "InitCarbonIntensityProvider",
	}
)

func GetConditions() []metav1.Condition {
	conditions := []metav1.Condition{
		ConditionHealthy,
	}

	return conditions
}
