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
	"errors"
	"fmt"
	"github.com/rekuberate-io/carbon/controllers/metrics"
	"github.com/rekuberate-io/carbon/providers"
	"github.com/rekuberate-io/carbon/providers/simulator"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

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
	logger := log.FromContext(ctx).WithName("controller")
	var err error

	var cip carbonv1alpha1.CarbonIntensityProvider
	if err = r.Get(ctx, req.NamespacedName, &cip); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		metrics.CipReconciliationLoopErrorsTotal.WithLabelValues(req.NamespacedName.String()).Inc()

		msg := "unable to fetch carbon intensity provider"
		logger.V(5).Error(err, "unable to fetch carbon intensity provider")
		r.Recorder.Event(&cip, "Warning", OperatorResourceNotAvailableReason, msg)

		return r.SetCondition(
			ctx,
			&cip,
			OperatorSucceededCondition,
			OperatorResourceNotAvailableReason,
			metav1.ConditionFalse,
			client.IgnoreNotFound(err),
			msg,
			nil,
		)
	}

	var provider providers.Provider
	var zone string
	providerType := providers.ProviderType(cip.Spec.Provider)
	patch := client.MergeFrom(cip.DeepCopy())

	switch providerType {
	case providers.WattTime:
		if cip.Spec.WattTimeConfiguration == nil {
			err = errors.New("missing configuration in yaml")
			break
		}

		zone = cip.Spec.WattTimeConfiguration.Region

		passwordRef := cip.Spec.WattTimeConfiguration.Password
		objectKey := client.ObjectKey{
			Namespace: req.Namespace,
			Name:      passwordRef.Name,
		}

		var secret corev1.Secret
		if err := r.Get(ctx, objectKey, &secret); err != nil {
			msg := "finding secret failed"
			logger.Error(err, msg)
			r.Recorder.Event(&cip, "Warning", OperatorResourceNotAvailableReason, msg)

			metrics.CipReconciliationLoopErrorsTotal.WithLabelValues(req.NamespacedName.String()).Inc()

			return r.SetCondition(
				ctx,
				&cip,
				OperatorSucceededCondition,
				OperatorResourceNotAvailableReason,
				metav1.ConditionFalse,
				client.IgnoreNotFound(err),
				msg,
				nil,
			)
		}

		password := string(secret.Data["password"])
		provider, err = providers.NewWattTimeProvider(ctx, cip.Spec.WattTimeConfiguration.Username, password)
	case providers.ElectricityMaps:
		if cip.Spec.ElectricityMapsConfiguration == nil {
			err = errors.New("missing configuration in yaml")
			break
		}

		zone = cip.Spec.ElectricityMapsConfiguration.Zone

		apiKeyRef := cip.Spec.ElectricityMapsConfiguration.ApiKey
		objectKey := client.ObjectKey{
			Namespace: req.Namespace,
			Name:      apiKeyRef.Name,
		}

		var secret corev1.Secret
		if err := r.Get(ctx, objectKey, &secret); err != nil {
			msg := "finding secret failed"
			logger.Error(err, msg)
			r.Recorder.Event(&cip, "Warning", OperatorResourceNotAvailableReason, msg)

			metrics.CipReconciliationLoopErrorsTotal.WithLabelValues(req.NamespacedName.String()).Inc()

			return r.SetCondition(
				ctx,
				&cip,
				OperatorSucceededCondition,
				OperatorResourceNotAvailableReason,
				metav1.ConditionFalse,
				client.IgnoreNotFound(err),
				msg,
				nil,
			)
		}

		apiKey := string(secret.Data["apiKey"])
		subscriptionType := providers.SubscriptionType(cip.Spec.ElectricityMapsConfiguration.Subscription)

		switch subscriptionType {
		case providers.Commercial:
			provider, err = providers.NewElectricityMapsProvider(apiKey)
		case providers.CommercialTrial:
			provider, err = providers.NewElectricityMapsCommercialTrialProvider(
				apiKey,
				cip.Spec.ElectricityMapsConfiguration.CommercialTrialEndpoint,
			)
		case providers.FreeTier:
			provider, err = providers.NewElectricityMapsFreeTierProvider(apiKey)
		}
	case providers.Simulator:
		if cip.Spec.SimulatorConfiguration == nil {
			err = errors.New("missing configuration in yaml")
			break
		}

		zone = cip.Spec.SimulatorConfiguration.Zone
		provider, err = simulator.NewCarbonIntensityProviderSimulator(
			zone,
			*cip.Spec.SimulatorConfiguration.Randomize,
		)
	}

	if err != nil {
		msg := "unable to initialize provider"
		logger.Error(err, msg, msg, "providerType", providerType)
		r.Recorder.Event(&cip, "Warning", OperatorInitializeProviderFailedReason, msg)

		metrics.CipReconciliationLoopErrorsTotal.WithLabelValues(req.NamespacedName.String()).Inc()

		return r.SetCondition(
			ctx,
			&cip,
			OperatorSucceededCondition,
			OperatorInitializeProviderFailedReason,
			metav1.ConditionFalse,
			client.IgnoreNotFound(err),
			msg,
			nil,
		)

	}

	carbonIntensity, err := provider.GetCurrent(ctx, zone)
	if err != nil {
		msg := "request to provider failed"
		logger.Error(err, msg, "providerType", providerType)
		r.Recorder.Event(&cip, "Warning", OperatorReconcileFailedReason, msg)

		carbonIntensity = -1
	}

	var carbonIntensityAsString string
	if carbonIntensity < 0 {
		carbonIntensity = 0
		carbonIntensityAsString = "n/a"
	} else {
		carbonIntensityAsString = fmt.Sprintf("%.2f", carbonIntensity)
	}

	objectKey := client.ObjectKey{
		Namespace: req.Namespace,
		Name:      fmt.Sprintf("%s-forecast", req.Name),
	}

	var createConfigMap bool = false //cip.Status.Provider == nil || cip.Status.Zone == nil || cip.Spec.Provider != *cip.Status.Provider || zone != cip.Status.Zone
	var deleteConfigMap bool = true
	var updateForecast bool = false

	if cip.Status.Provider == nil || cip.Status.Zone == nil {
		createConfigMap = true
	} else {
		if cip.Spec.Provider != *cip.Status.Provider || zone != *cip.Status.Zone {
			createConfigMap = true
		}
	}

	configMap := &corev1.ConfigMap{}
	err = r.Get(ctx, objectKey, configMap)
	if err != nil {
		if apierrors.IsNotFound(err) {
			deleteConfigMap = false
			createConfigMap = true
		}
	}

	timestamp := time.Now()

	updateForecast = createConfigMap || (cip.Status.LastForecast == nil || cip.Status.LastForecast.Add(time.Duration(cip.Spec.ForecastRefreshIntervalInHours)*time.Minute).Before(time.Now()))
	if updateForecast {
		lastForecast := cip.Status.LastForecast
		cip.Status.LastForecast = &metav1.Time{Time: timestamp}

		forecast, err := provider.GetForecast(ctx, zone)
		if err != nil {
			msg := "request to provider failed"
			logger.Error(err, msg, "providerType", providerType)
			r.Recorder.Event(&cip, "Warning", OperatorReconcileFailedReason, msg)

			cip.Status.LastForecast = lastForecast
		}

		if deleteConfigMap {
			err := r.Delete(ctx, configMap)
			if err != nil {
				msg := "deleting configmap failed"
				logger.Error(err, msg, "objectKey", objectKey)
				r.Recorder.Event(&cip, "Warning", OperatorConfigMapDeploymentFailedReason, msg)

				metrics.CipReconciliationLoopErrorsTotal.WithLabelValues(req.NamespacedName.String()).Inc()
				return ctrl.Result{}, err
			}
		}

		if forecast != nil {
			configMap, err = r.PrepareConfigMap(req, forecast, zone, cip.Status.LastForecast.Time, providerType, true)
			if err != nil {
				msg := "preparing configmap failed"
				logger.Error(err, msg, "objectKey", objectKey)
				r.Recorder.Event(&cip, "Warning", OperatorConfigMapDeploymentFailedReason, msg)

				metrics.CipReconciliationLoopErrorsTotal.WithLabelValues(req.NamespacedName.String()).Inc()
				return ctrl.Result{}, err
			}

			err = r.Create(ctx, configMap)
			if err != nil {
				msg := "creating configmap failed"
				logger.Error(err, msg, "objectKey", objectKey)
				r.Recorder.Event(&cip, "Warning", OperatorConfigMapDeploymentFailedReason, msg)

				metrics.CipReconciliationLoopErrorsTotal.WithLabelValues(req.NamespacedName.String()).Inc()

				return r.SetCondition(
					ctx,
					&cip,
					OperatorSucceededCondition,
					OperatorConfigMapDeploymentFailedReason,
					metav1.ConditionFalse,
					client.IgnoreNotFound(err),
					msg,
					nil,
				)
			}

			err = controllerutil.SetOwnerReference(&cip, configMap, r.Scheme)
			if err != nil {
				msg := "setting owner reference"
				logger.Error(err, msg, "configmap", configMap.Name)
				r.Recorder.Event(&cip, "Warning", OperatorConfigMapDeploymentFailedReason, msg)

				metrics.CipReconciliationLoopErrorsTotal.WithLabelValues(req.NamespacedName.String()).Inc()

				return r.SetCondition(
					ctx,
					&cip,
					OperatorSucceededCondition,
					OperatorConfigMapDeploymentFailedReason,
					metav1.ConditionFalse,
					client.IgnoreNotFound(err),
					msg,
					nil,
				)
			}
		}
	}

	requeueAfter := time.Minute * time.Duration(cip.Spec.LiveRefreshIntervalInHours)

	cip.Status.Zone = &zone
	cip.Status.Provider = &cip.Spec.Provider
	cip.Status.LastUpdate = &metav1.Time{Time: timestamp}
	cip.Status.NextUpdate = &metav1.Time{Time: timestamp.Add(requeueAfter)}
	cip.Status.CarbonIntensity = &carbonIntensityAsString
	err = r.Status().Patch(ctx, &cip, patch)
	if err != nil {
		namespacedName := fmt.Sprintf("%v/%v", cip.Namespace, cip.Name)
		msg := "failed to patch provider's status"
		logger.Error(err, msg, "provider", namespacedName)
		r.Recorder.Event(&cip, "Warning", OperatorResourceStatusUpdateFailedReason, msg)

		metrics.CipReconciliationLoopErrorsTotal.WithLabelValues(req.NamespacedName.String()).Inc()

		return r.SetCondition(
			ctx,
			&cip,
			OperatorSucceededCondition,
			OperatorResourceStatusUpdateFailedReason,
			metav1.ConditionFalse,
			client.IgnoreNotFound(err),
			msg,
			nil,
		)
	}

	metrics.CipLiveCarbonIntensityMetric.WithLabelValues(string(providerType), zone).Set(carbonIntensity)
	metrics.CipReconciliationLoopsTotal.WithLabelValues(req.NamespacedName.String()).Inc()

	msg := "operator successfully reconciled"
	r.Recorder.Event(&cip, "Normal", OperatorReconcileSucceededReason, msg)

	return r.SetCondition(
		ctx,
		&cip,
		OperatorSucceededCondition,
		OperatorReconcileSucceededReason,
		metav1.ConditionTrue,
		nil,
		msg,
		&requeueAfter,
	)
}

// SetupWithManager sets up the controller with the Manager.
func (r *CarbonIntensityProviderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&carbonv1alpha1.CarbonIntensityProvider{}).
		WithEventFilter(ignorePredicates()).
		Complete(r)
}

func ignorePredicates() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			// Ignore updates to CR status in which case metadata.Generation does not change
			return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			// Evaluates to false if the object has been confirmed deleted.
			return !e.DeleteStateUnknown
		},
	}
}

func (r *CarbonIntensityProviderReconciler) PrepareConfigMap(
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

func (r *CarbonIntensityProviderReconciler) SetCondition(
	ctx context.Context,
	cip *carbonv1alpha1.CarbonIntensityProvider,
	condition string,
	reason string,
	status metav1.ConditionStatus,
	err error,
	msg string,
	requeueAfter *time.Duration,
) (ctrl.Result, error) {
	meta.SetStatusCondition(&cip.Status.Conditions,
		metav1.Condition{
			Type:               condition,
			Status:             status,
			LastTransitionTime: metav1.Time{},
			Reason:             reason,
			Message:            msg,
		})

	var result reconcile.Result

	if requeueAfter == nil {
		result = ctrl.Result{}
	} else {
		result = ctrl.Result{RequeueAfter: *requeueAfter}
	}

	return result, utilerrors.NewAggregate([]error{err, r.Status().Update(ctx, cip)})
}
