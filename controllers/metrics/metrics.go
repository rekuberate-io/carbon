package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	CipReconciliationLoopsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rekuberate_carbon_intensity_provider_reconciliations_total",
			Help: "Number of total reconciliation loops for carbon intensity provider controller",
		},
		[]string{"resource"},
	)

	CipReconciliationLoopErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rekuberate_carbon_intensity_provider_reconciliations_errors_total",
			Help: "Number of total failed reconciliation loops for carbon intensity provider controller",
		},
		[]string{"resource"},
	)

	CipLiveCarbonIntensityMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rekuberate_carbon_intensity_provider_live_gramsperkilowatthour",
			Help: "Carbon Intensity (grCO2eq/KWh)",
		},
		[]string{"provider", "zone"},
	)
)

func init() {
	metrics.Registry.MustRegister(CipReconciliationLoopsTotal)
	metrics.Registry.MustRegister(CipReconciliationLoopErrorsTotal)
	metrics.Registry.MustRegister(CipLiveCarbonIntensityMetric)
}
