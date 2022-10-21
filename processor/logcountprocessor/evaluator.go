package logcountprocessor

import (
	"github.com/antonmedv/expr/vm"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

// Evaluator is used to evaluate expressions.
type Evaluator struct {
	match      *vm.Program
	attributes map[string]*vm.Program
	logger     *zap.Logger
}

// NewEvaluator returns a new expression evaluator.
func NewEvaluator(match *vm.Program, attributes map[string]*vm.Program, logger *zap.Logger) *Evaluator {
	return &Evaluator{
		match,
		attributes,
		logger,
	}
}

// GetExpressionEnv returns an expression env for open telemetry logs.
func GetExpressionEnv(resource pcommon.Resource, log plog.LogRecord) map[string]interface{} {
	env := make(map[string]interface{})
	env["body"] = log.Body().AsRaw()
	env["resource"] = resource.Attributes().AsRaw()
	env["attributes"] = log.Attributes().AsRaw()
	env["severity"] = log.SeverityText()
	return env
}

// MatchesLog determines if a log matches the evaluator's filter.
func (e *Evaluator) MatchesLog(resource pcommon.Resource, log plog.LogRecord) bool {
	env := GetExpressionEnv(resource, log)
	matches, err := vm.Run(e.match, env)
	if err != nil {
		e.logger.Error("Failed to evaluate match expression", zap.Error(err))
		return false
	}

	matchesBool, ok := matches.(bool)
	if !ok {
		e.logger.Error("Match expression did not return a boolean", zap.Error(err))
		return false
	}

	return matchesBool
}

// GetAttributes returns attributes extracted from logs.
func (e *Evaluator) GetAttributes(resource pcommon.Resource, log plog.LogRecord) map[string]interface{} {
	attributes := make(map[string]interface{})
	env := GetExpressionEnv(resource, log)
	for key, expression := range e.attributes {
		value, err := vm.Run(expression, env)
		if err != nil {
			e.logger.Error("Failed to evaluate attributes expression", zap.Error(err))
			continue
		}

		attributes[key] = value
	}

	return attributes
}
