package measurements

// BindplaneAgentThroughputMeasurementsRegistry is the registry singleton used by bindplane agent to
// track throughput measurements
var BindplaneAgentThroughputMeasurementsRegistry = NewConcreteThroughputMeasurementsRegistry(false)
