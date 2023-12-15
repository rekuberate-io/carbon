package providers

import (
	"context"
	"fmt"
	carbonv1alpha1 "github.com/rekuberate-io/carbon/api/v1alpha1"
	"github.com/rekuberate-io/carbon/pkg/providers/electricitymaps"
	"github.com/rekuberate-io/carbon/pkg/providers/simulator"
	"github.com/rekuberate-io/carbon/pkg/providers/watttime"
	v1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"slices"
	"strings"
	"time"
)

const (
	WattTime        string = "watttime"
	ElectricityMaps string = "electricitymaps"
	Simulator       string = "simulator"
)

var (
	supportedProviders     = []string{WattTime, ElectricityMaps, Simulator}
	supportedEmissionTypes = []EmissionsType{Average, Marginal}
)

type Provider interface {
	GetCurrent(ctx context.Context) (float64, error)
	GetForecast(ctx context.Context) (map[time.Time]float64, error)
	Region() string
}

//
//type Forecast struct {
//	PointTime       time.Time `json:"pointTime"`
//	CarbonIntensity float64   `json:"carbonIntensity"`
//}

func GetProvider(
	ctx context.Context,
	req ctrl.Request,
	kClient client.Client,
	providerRef *v1.ObjectReference,
) (Provider, error) {
	providerRefKind := strings.ToLower(providerRef.Kind)
	if providerRefKind == "" {
		err := fmt.Errorf("carbon intensity provider is missing")
		return nil, err
	}

	if !IsSupported(providerRefKind) {
		err := fmt.Errorf("not supported carbon intensity provider")
		return nil, err
	}

	providerRefNamespace := req.Namespace
	if providerRef.Namespace != "" {
		providerRefNamespace = providerRef.Namespace
	}

	objectKey := client.ObjectKey{Name: providerRef.Name, Namespace: providerRefNamespace}

	switch providerRefKind {
	case Simulator:
		po := &carbonv1alpha1.Simulator{}
		if err := kClient.Get(ctx, objectKey, po); err != nil {
			return nil, err
		}

		p, err := simulator.NewProvider(*po)
		if err != nil {
			return nil, err
		}

		return Provider(p), nil
	case WattTime:
		po := &carbonv1alpha1.WattTime{}
		if err := kClient.Get(ctx, objectKey, po); err != nil {
			return nil, err
		}

		p, err := watttime.NewProvider(ctx, kClient, *po)
		if err != nil {
			return nil, err
		}

		return Provider(p), nil
	case ElectricityMaps:
		po := &carbonv1alpha1.ElectricityMaps{}
		if err := kClient.Get(ctx, objectKey, po); err != nil {
			return nil, err
		}

		p, err := electricitymaps.NewProvider(ctx, kClient, *po)
		if err != nil {
			return nil, err
		}

		return Provider(p), nil
	}

	return nil, nil
}

func GetSupportedProviders() []string {
	return supportedProviders
}

func IsSupported(providerType string) bool {
	if !slices.Contains(supportedProviders, strings.ToLower(providerType)) {
		return false
	}

	return true
}
