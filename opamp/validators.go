package opamp

import "gopkg.in/yaml.v3"

// ValidatorFunc function that takes in a config contents and validates
// Returns true if valid
type ValidatorFunc func([]byte) bool

// NewYamlValidator creates a new Validator that checks does a yaml unmarshal against the target interface{}
func NewYamlValidator(target interface{}) ValidatorFunc {
	return func(b []byte) bool {
		if err := yaml.Unmarshal(b, target); err != nil {
			return false
		}

		return true
	}
}
