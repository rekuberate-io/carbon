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
	"github.com/rekuberate-io/carbon/providers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"net"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1alpha1 "github.com/rekuberate-io/carbon/api/v1alpha1"
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

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *CarbonIntensityProviderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithName("controller")

	var cip corev1alpha1.CarbonIntensityProvider
	if err := r.Get(ctx, req.NamespacedName, &cip); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		logger.V(5).Error(err, "unable to fetch carbon intensity provider")
		return ctrl.Result{}, err
	}

	providerType := providers.ProviderType(cip.Spec.Provider)
	emissionsSignal := providers.EmissionsSignal(cip.Spec.Signal)
	provider, err := providers.NewProvider(providerType, emissionsSignal)
	if err != nil {
		logger.Error(err, "unable to initialize provider", "providerType", providerType)
		return ctrl.Result{}, nil
	}

	provider.GetCurrent()

	requeueAfter := time.Hour * time.Duration(*cip.Spec.RefreshIntervalInHours)

	patch := client.MergeFrom(cip.DeepCopy())
	cip.Status.LastUpdate = &metav1.Time{Time: time.Now()}
	cip.Status.NextUpdate = &metav1.Time{Time: time.Now().Add(requeueAfter)}
	err = r.Status().Patch(ctx, &cip, patch)
	if err != nil {
		namespacedName := fmt.Sprintf("%v/%v", cip.Namespace, cip.Name)
		logger.Error(err, "failed to patch provider's status", "provider", namespacedName)
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: requeueAfter}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CarbonIntensityProviderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.CarbonIntensityProvider{}).
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

func (r *CarbonIntensityProviderReconciler) getIpAddress() (net.IP, error) {
	localAddress := "127.0.0.1"
	ipAddress := net.ParseIP(localAddress)

	return ipAddress, nil
}
