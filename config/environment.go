package config

type Environment string

func (e Environment) String() string {
	return string(e)
}

const (
	Staging    Environment = "staging"
	Production Environment = "production"
)
