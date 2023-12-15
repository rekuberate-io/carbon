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
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/rekuberate-io/carbon/pkg/providers"
	"github.com/rekuberate-io/carbon/pkg/providers/electricitymaps"
	"github.com/rekuberate-io/carbon/pkg/providers/simulator"
	"github.com/rekuberate-io/carbon/pkg/providers/watttime"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	carbonv1alpha1 "github.com/rekuberate-io/carbon/api/v1alpha1"
)

const (
	labelProviderInstance = "core.rekuberate.io/carbon-provider-instance"
	labelProviderType     = "core.rekuberate.io/carbon-provider-type"
	labelProviderZone     = "core.rekuberate.io/carbon-provider-zone"
)

var (
	eventFilters = builder.WithPredicates(predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			// We only need to check generation changes here, because it is only
			// updated on spec changes. On the other hand RevisionVersion
			// changes also on status changes. We want to omit reconciliation
			// for status updates.
			return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			// DeleteStateUnknown evaluates to false only if the object
			// has been confirmed as deleted by the api server.
			return !e.DeleteStateUnknown
		},
	})
	logger logr.Logger
	dbglvl int = 5
)

// CarbonIntensityProviderReconciler reconciles a CarbonIntensityProvider object
type CarbonIntensityProviderReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=core.rekuberate.io,resources=carbonintensityproviders,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.rekuberate.io,resources=carbonintensityproviders/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core.rekuberate.io,resources=carbonintensityproviders/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *CarbonIntensityProviderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger = log.FromContext(ctx).WithName("carbon-controller")

	current := &carbonv1alpha1.CarbonIntensityProvider{}
	if err := r.Get(ctx, req.NamespacedName, current); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		logger.V(dbglvl).Error(err, "unable to fetch carbon intensity provider")
		return ctrl.Result{}, err
	}

	desired := current.DeepCopy()

	if current.Status.Conditions == nil {
		conditions := carbonv1alpha1.GetConditions()
		for _, condition := range conditions {
			meta.SetStatusCondition(&desired.Status.Conditions, condition)
		}

		res, err := r.updateStatus(ctx, current, desired)
		if err != nil {
			return res, err
		}
	}

	providerRef := current.Spec.ProviderRef
	providerRefKind := strings.ToLower(providerRef.Kind)

	if providerRefKind == "" {
		err := fmt.Errorf("carbon intensity provider is missing")
		logger.Error(err, "terminate reconciliation")
		return ctrl.Result{}, nil
	}

	if !providers.IsSupported(providerRefKind) {
		err := fmt.Errorf("not supported carbon intensity provider")
		logger.Error(err, "terminate reconciliation", "providerKind", providerRef.Kind)
		return ctrl.Result{}, nil
	}

	providerRefNamespace := req.Namespace
	if providerRef.Namespace != "" {
		providerRefNamespace = providerRef.Namespace
	}

	objectKey := client.ObjectKey{Name: providerRef.Name, Namespace: providerRefNamespace}
	var providerType client.Object
	var provider providers.Provider

	switch providerRefKind {
	case string(providers.Simulator):
		providerType = &carbonv1alpha1.Simulator{}
		provider = &simulator.Simulator{}
	case string(providers.WattTime):
		providerType = &carbonv1alpha1.WattTime{}
		provider = &watttime.WattTimeProvider{}
	case string(providers.ElectricityMaps):
		providerType = &carbonv1alpha1.ElectricityMaps{}
		provider = &electricitymaps.ElectricityMapsProvider{}
	}

	if err := r.Get(ctx, objectKey, providerType); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		logger.Error(err, "unable to fetch provider", "provider", objectKey)
		return ctrl.Result{}, err
	}

	return r.updateStatus(ctx, current, desired)
}

// SetupWithManager sets up the controller with the Manager.
func (r *CarbonIntensityProviderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&carbonv1alpha1.CarbonIntensityProvider{}, eventFilters).
		Complete(r)
}

func (r *CarbonIntensityProviderReconciler) updateStatus(
	ctx context.Context,
	current *carbonv1alpha1.CarbonIntensityProvider,
	desired *carbonv1alpha1.CarbonIntensityProvider,
) (ctrl.Result, error) {
	if !reflect.DeepEqual(current, desired) {
		err := r.Status().Update(ctx, desired)
		if err != nil {
			logger.Error(err, "unable to update carbon intensity provider status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *CarbonIntensityProviderReconciler) prepareConfigMap(
	req ctrl.Request,
	forecast []providers.Forecast,
	zone string,
	pointTime time.Time,
	providerType providers.ProviderType,
	immutable bool,
) (*corev1.ConfigMap, error) {

	jsonData, err := json.Marshal(forecast)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(jsonData); err != nil {
		return nil, err
	}

	data := map[string]string{
		"provider":  string(providerType),
		"zone":      zone,
		"pointTime": pointTime.String(),
	}

	configMapName := fmt.Sprintf("%s-forecast", req.Name)

	labels := map[string]string{
		"app.kubernetes.io/name":       "carbonintensityprovider",
		"app.kubernetes.io/instance":   configMapName,
		"app.kubernetes.io/component":  "forecast",
		"app.kubernetes.io/part-of":    "carbon",
		"app.kubernetes.io/managed-by": "controller",
		"app.kubernetes.io/created-by": "carbon",
		labelProviderInstance:          req.Name,
		labelProviderType:              string(providerType),
		labelProviderZone:              zone,
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: req.Namespace,
			Labels:    labels,
		},
		Immutable: &immutable,
		Data:      data,
		BinaryData: map[string][]byte{
			"BinaryData": buffer.Bytes(),
		},
	}

	return configMap, nil
}
