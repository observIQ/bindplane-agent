package logcountprocessor

import (
	"errors"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
)

// Expression is an Expression used to evaluate values.
type Expression struct {
	*vm.Program
}

// Match checks if an expression matches the supplied environment.
func (e *Expression) Match(env map[string]interface{}) (bool, error) {
	matches, err := e.Evaluate(env)
	if err != nil {
		return false, err
	}

	matchesBool, ok := matches.(bool)
	if !ok {
		return false, errors.New("expression did not return a boolean")
	}

	return matchesBool, nil
}

// Evaluate evaluates an expression against the supplied environment.
func (e *Expression) Evaluate(env map[string]interface{}) (interface{}, error) {
	return vm.Run(e.Program, env)
}

// NewExpression creates an expression from a string.
func NewExpression(str string, opts ...expr.Option) (*Expression, error) {
	program, err := expr.Compile(str, opts...)
	if err != nil {
		return nil, err
	}

	return &Expression{program}, nil
}
