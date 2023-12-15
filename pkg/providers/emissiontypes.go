package providers

type EmissionsType string

const (
	Average  EmissionsType = "average"
	Marginal EmissionsType = "marginal"
)

func GetSupportedEmissionsTypes() []EmissionsType {
	return supportedEmissionTypes
}
