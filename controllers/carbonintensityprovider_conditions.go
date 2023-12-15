package controllers

const (
	OperatorSucceededCondition               = "Ready"
	OperatorResourceNotAvailableReason       = "ResourceNotAvailable"
	OperatorResourceStatusUpdateFailedReason = "ResourceStatusUpdateFailed"
	OperatorInitializeProviderFailedReason   = "InitializeCarbonIntensityProviderFailed"
	OperatorConfigMapDeploymentFailedReason  = "ConfigMapDeploymentFailed"
	OperatorReconcileSucceededReason         = "ReconcileSucceeded"
	OperatorReconcileFailedReason            = "ReconcileFailed"

	ConditionHealthy string = "Healthy"
	UnknownReason    string = "Unknown"
)
