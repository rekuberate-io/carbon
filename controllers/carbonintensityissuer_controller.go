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
	"fmt"
	"github.com/go-logr/logr"
	"github.com/rekuberate-io/carbon/pkg/providers"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"os"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	carbonv1alpha1 "github.com/rekuberate-io/carbon/api/v1alpha1"
	carboninfluxdb "github.com/rekuberate-io/carbon/pkg/influxdb"
)

const (
	labelProviderInstance = "core.rekuberate.io/carbon-issuer-instance"
	labelProviderType     = "core.rekuberate.io/carbon-issuer-type"
	labelProviderZone     = "core.rekuberate.io/carbon-issuer-zone"
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

// CarbonIntensityIssuerReconciler reconciles a CarbonIntensityIssuer object
type CarbonIntensityIssuerReconciler struct {
	client.Client
	Scheme          *runtime.Scheme
	Recorder        record.EventRecorder
	InfluxDb2Client influxdb2.Client
}

//+kubebuilder:rbac:groups=core.rekuberate.io,resources=carbonintensityissuers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.rekuberate.io,resources=carbonintensityissuers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core.rekuberate.io,resources=carbonintensityissuers/finalizers,verbs=update
//+kubebuilder:rbac:groups=core.rekuberate.io,resources=electricitymaps,verbs=get;list;watch
//+kubebuilder:rbac:groups=core.rekuberate.io,resources=watttimes,verbs=get;list;watch
//+kubebuilder:rbac:groups=core.rekuberate.io,resources=simulators,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *CarbonIntensityIssuerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger = log.FromContext(ctx).WithName("carbon-controller")

	// get carbon intensity provider resource
	before := &carbonv1alpha1.CarbonIntensityIssuer{}
	if err := r.Get(ctx, req.NamespacedName, before); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		logger.V(dbglvl).Error(err, "unable to fetch carbon intensity provider")
		return ctrl.Result{}, err
	}

	after := before.DeepCopy()

	// initialize status conditions
	if before.Status.Conditions == nil {
		conditions := carbonv1alpha1.GetConditions()
		for _, condition := range conditions {
			meta.SetStatusCondition(&after.Status.Conditions, condition)
		}

		res, err := r.UpdateStatus(ctx, before, after)
		if err != nil {
			return res, err
		}
	}

	// get a concrete provider
	providerRef := before.Spec.ProviderRef
	if providerRef.Namespace == "" {
		providerRef.Namespace = req.Namespace
	}

	provider, err := providers.GetProvider(ctx, req, r.Client, providerRef)
	if err != nil {
		condition := carbonv1alpha1.ConditionHealthy.DeepCopy()
		condition.Status = metav1.ConditionFalse
		condition.Reason = carbonv1alpha1.ProviderInitFailed
		condition.Message = err.Error()
		meta.SetStatusCondition(&after.Status.Conditions, *condition)

		logger.Error(err, "unable to get provider", "providerKind", providerRef.Kind)
		r.UpdateStatus(ctx, before, after)

		return ctrl.Result{}, err
	}

	// set the condition healthy as true because we initialized the concrete provider
	condition := carbonv1alpha1.ConditionHealthy.DeepCopy()
	condition.Status = metav1.ConditionTrue
	condition.Reason = carbonv1alpha1.ProviderInitFinished
	condition.Message = fmt.Sprintf("Initialized Provider '%s', (%s)", providerRef.Name, providerRef.Kind)
	meta.SetStatusCondition(&after.Status.Conditions, *condition)

	// get current carbon intensity
	carbonIntensity, err := provider.GetCurrent(ctx, before.Spec.Zone)
	if err != nil {
		logger.Error(err, "unable to get carbon intensity", "providerKind", providerRef.Kind, "provider", providerRef.Name)
		return ctrl.Result{}, err
	}

	// get carbon intensity forecast
	// TODO: change to time.Hours
	forecast := make(map[time.Time]float64)
	if before.Status.LastForecast == nil ||
		before.Status.LastForecast.Add(time.Duration(before.Spec.ForecastRefreshIntervalInHours)*time.Minute).Before(time.Now()) {
		forecast, err = provider.GetForecast(ctx, before.Spec.Zone)
		if err != nil {
			logger.Error(err, "unable to get carbon intensity forecast", "providerKind", providerRef.Kind, "provider", providerRef.Name)
			//return ctrl.Result{}, err
		}

		after.Status.LastForecast = &metav1.Time{Time: time.Now()}
	}

	// update status of custom resource, push carbon intensity measurements to influxdb
	if carbonIntensity > 0 {
		carbonIntensityAsString := fmt.Sprintf("%.2f", carbonIntensity)
		after.Status.CarbonIntensity = &carbonIntensityAsString
	} else {
		notAvailable := "-"
		after.Status.CarbonIntensity = &notAvailable
	}

	// TODO: change to time.Hours
	requeueAfter := time.Minute * time.Duration(before.Spec.LiveRefreshIntervalInHours)
	now := time.Now()

	after.Status.NextUpdate = &metav1.Time{Time: now.Add(requeueAfter)}
	after.Status.LastUpdate = &metav1.Time{Time: now}

	result, err := r.UpdateStatus(ctx, before, after)
	if err != nil {
		return result, err
	}

	// TODO: set N/A value as well in metric
	//if carbonIntensity > 0 {
	//	metrics.CipLiveCarbonIntensityMetric.WithLabelValues(
	//		providerRef.Kind,
	//		req.String(),
	//		before.Spec.Zone,
	//	).Set(carbonIntensity)
	//}

	orgName := os.Getenv("INFLUXDB2_ORG")
	bucketName := os.Getenv("INFLUXDB2_BUCKET")
	err = carboninfluxdb.InitDb(ctx, &r.InfluxDb2Client, orgName, bucketName)
	if err != nil {
		return ctrl.Result{}, err
	}

	tags := map[string]string{
		"providerKind": providerRef.Kind,
		"provider":     providerRef.Name,
		"zone":         before.Spec.Zone,
	}

	err = carboninfluxdb.PushMeasurements(
		ctx,
		&r.InfluxDb2Client,
		orgName,
		bucketName,
		"carbonIntensity",
		req.String(),
		tags,
		map[time.Time]float64{time.Now(): carbonIntensity},
	)
	if err != nil {
		logger.Error(err, "unable to push ci measurement to influxdb", "providerKind", providerRef.Kind, "provider", providerRef.Name, "bucket", bucketName)
	}
	logger.Info("pushed ci measurement to influxdb", "providerKind", providerRef.Kind, "provider", providerRef.Name, "bucket", os.Getenv("INFLUXDB2_BUCKET"))

	if len(forecast) > 0 {
		err = carboninfluxdb.PushMeasurements(
			ctx,
			&r.InfluxDb2Client,
			orgName,
			bucketName,
			"carbonIntensity",
			fmt.Sprintf("%s_%s", req.String(), "forecast"),
			tags,
			forecast,
		)
		if err != nil {
			logger.Error(err, "unable to push forecasts to influxdb", "providerKind", providerRef.Kind, "provider", providerRef.Name, "bucket", os.Getenv("INFLUXDB2_BUCKET"))
		}
		logger.Info("pushed ci forecast to influxdb", "providerKind", providerRef.Kind, "provider", providerRef.Name, "bucket", os.Getenv("INFLUXDB2_BUCKET"))
	}

	result.RequeueAfter = requeueAfter
	return result, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CarbonIntensityIssuerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&carbonv1alpha1.CarbonIntensityIssuer{}, eventFilters).
		Complete(r)
}

func (r *CarbonIntensityIssuerReconciler) UpdateStatus(
	ctx context.Context,
	current *carbonv1alpha1.CarbonIntensityIssuer,
	desired *carbonv1alpha1.CarbonIntensityIssuer,
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
