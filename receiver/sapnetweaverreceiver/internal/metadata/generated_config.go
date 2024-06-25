// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/filter"
)

// MetricConfig provides common config for a particular metric.
type MetricConfig struct {
	Enabled bool `mapstructure:"enabled"`

	enabledSetByUser bool
}

func (ms *MetricConfig) Unmarshal(parser *confmap.Conf) error {
	if parser == nil {
		return nil
	}
	err := parser.Unmarshal(ms)
	if err != nil {
		return err
	}
	ms.enabledSetByUser = parser.IsSet("enabled")
	return nil
}

// MetricsConfig provides config for sapnetweaver metrics.
type MetricsConfig struct {
	SapnetweaverAbapRfcCount                MetricConfig `mapstructure:"sapnetweaver.abap.rfc.count"`
	SapnetweaverAbapSessionCount            MetricConfig `mapstructure:"sapnetweaver.abap.session.count"`
	SapnetweaverAbapUpdateStatus            MetricConfig `mapstructure:"sapnetweaver.abap.update.status"`
	SapnetweaverCacheEvictions              MetricConfig `mapstructure:"sapnetweaver.cache.evictions"`
	SapnetweaverCacheHits                   MetricConfig `mapstructure:"sapnetweaver.cache.hits"`
	SapnetweaverCertificateValidity         MetricConfig `mapstructure:"sapnetweaver.certificate.validity"`
	SapnetweaverConnectionErrorCount        MetricConfig `mapstructure:"sapnetweaver.connection.error.count"`
	SapnetweaverCPUSystemUtilization        MetricConfig `mapstructure:"sapnetweaver.cpu.system.utilization"`
	SapnetweaverCPUUtilization              MetricConfig `mapstructure:"sapnetweaver.cpu.utilization"`
	SapnetweaverDatabaseDialogRequestTime   MetricConfig `mapstructure:"sapnetweaver.database.dialog.request.time"`
	SapnetweaverHostMemoryVirtualOverhead   MetricConfig `mapstructure:"sapnetweaver.host.memory.virtual.overhead"`
	SapnetweaverHostMemoryVirtualSwap       MetricConfig `mapstructure:"sapnetweaver.host.memory.virtual.swap"`
	SapnetweaverHostSpoolListUtilization    MetricConfig `mapstructure:"sapnetweaver.host.spool_list.utilization"`
	SapnetweaverLocksDequeueErrorsCount     MetricConfig `mapstructure:"sapnetweaver.locks.dequeue.errors.count"`
	SapnetweaverLocksEnqueueCurrentCount    MetricConfig `mapstructure:"sapnetweaver.locks.enqueue.current.count"`
	SapnetweaverLocksEnqueueErrorsCount     MetricConfig `mapstructure:"sapnetweaver.locks.enqueue.errors.count"`
	SapnetweaverLocksEnqueueHighCount       MetricConfig `mapstructure:"sapnetweaver.locks.enqueue.high.count"`
	SapnetweaverLocksEnqueueLockTime        MetricConfig `mapstructure:"sapnetweaver.locks.enqueue.lock_time"`
	SapnetweaverLocksEnqueueLockWaitTime    MetricConfig `mapstructure:"sapnetweaver.locks.enqueue.lock_wait_time"`
	SapnetweaverLocksEnqueueMaxCount        MetricConfig `mapstructure:"sapnetweaver.locks.enqueue.max.count"`
	SapnetweaverMemoryConfigured            MetricConfig `mapstructure:"sapnetweaver.memory.configured"`
	SapnetweaverMemoryFree                  MetricConfig `mapstructure:"sapnetweaver.memory.free"`
	SapnetweaverMemorySwapSpaceUtilization  MetricConfig `mapstructure:"sapnetweaver.memory.swap_space.utilization"`
	SapnetweaverProcessAvailability         MetricConfig `mapstructure:"sapnetweaver.process_availability"`
	SapnetweaverQueueCount                  MetricConfig `mapstructure:"sapnetweaver.queue.count"`
	SapnetweaverQueueMaxCount               MetricConfig `mapstructure:"sapnetweaver.queue_max.count"`
	SapnetweaverQueuePeakCount              MetricConfig `mapstructure:"sapnetweaver.queue_peak.count"`
	SapnetweaverRequestCount                MetricConfig `mapstructure:"sapnetweaver.request.count"`
	SapnetweaverRequestTimeoutCount         MetricConfig `mapstructure:"sapnetweaver.request.timeout.count"`
	SapnetweaverResponseDuration            MetricConfig `mapstructure:"sapnetweaver.response.duration"`
	SapnetweaverSessionCount                MetricConfig `mapstructure:"sapnetweaver.session.count"`
	SapnetweaverSessionsBrowserCount        MetricConfig `mapstructure:"sapnetweaver.sessions.browser.count"`
	SapnetweaverSessionsEjbCount            MetricConfig `mapstructure:"sapnetweaver.sessions.ejb.count"`
	SapnetweaverSessionsHTTPCount           MetricConfig `mapstructure:"sapnetweaver.sessions.http.count"`
	SapnetweaverSessionsSecurityCount       MetricConfig `mapstructure:"sapnetweaver.sessions.security.count"`
	SapnetweaverSessionsWebCount            MetricConfig `mapstructure:"sapnetweaver.sessions.web.count"`
	SapnetweaverShortDumpsRate              MetricConfig `mapstructure:"sapnetweaver.short_dumps.rate"`
	SapnetweaverSpoolRequestErrorCount      MetricConfig `mapstructure:"sapnetweaver.spool.request.error.count"`
	SapnetweaverSystemInstanceAvailability  MetricConfig `mapstructure:"sapnetweaver.system.instance_availability"`
	SapnetweaverWorkProcessActiveCount      MetricConfig `mapstructure:"sapnetweaver.work_process.active.count"`
	SapnetweaverWorkProcessJobAbortedStatus MetricConfig `mapstructure:"sapnetweaver.work_process.job.aborted.status"`
}

func DefaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		SapnetweaverAbapRfcCount: MetricConfig{
			Enabled: true,
		},
		SapnetweaverAbapSessionCount: MetricConfig{
			Enabled: true,
		},
		SapnetweaverAbapUpdateStatus: MetricConfig{
			Enabled: true,
		},
		SapnetweaverCacheEvictions: MetricConfig{
			Enabled: true,
		},
		SapnetweaverCacheHits: MetricConfig{
			Enabled: true,
		},
		SapnetweaverCertificateValidity: MetricConfig{
			Enabled: true,
		},
		SapnetweaverConnectionErrorCount: MetricConfig{
			Enabled: true,
		},
		SapnetweaverCPUSystemUtilization: MetricConfig{
			Enabled: true,
		},
		SapnetweaverCPUUtilization: MetricConfig{
			Enabled: true,
		},
		SapnetweaverDatabaseDialogRequestTime: MetricConfig{
			Enabled: true,
		},
		SapnetweaverHostMemoryVirtualOverhead: MetricConfig{
			Enabled: true,
		},
		SapnetweaverHostMemoryVirtualSwap: MetricConfig{
			Enabled: true,
		},
		SapnetweaverHostSpoolListUtilization: MetricConfig{
			Enabled: true,
		},
		SapnetweaverLocksDequeueErrorsCount: MetricConfig{
			Enabled: true,
		},
		SapnetweaverLocksEnqueueCurrentCount: MetricConfig{
			Enabled: true,
		},
		SapnetweaverLocksEnqueueErrorsCount: MetricConfig{
			Enabled: true,
		},
		SapnetweaverLocksEnqueueHighCount: MetricConfig{
			Enabled: true,
		},
		SapnetweaverLocksEnqueueLockTime: MetricConfig{
			Enabled: true,
		},
		SapnetweaverLocksEnqueueLockWaitTime: MetricConfig{
			Enabled: true,
		},
		SapnetweaverLocksEnqueueMaxCount: MetricConfig{
			Enabled: true,
		},
		SapnetweaverMemoryConfigured: MetricConfig{
			Enabled: true,
		},
		SapnetweaverMemoryFree: MetricConfig{
			Enabled: true,
		},
		SapnetweaverMemorySwapSpaceUtilization: MetricConfig{
			Enabled: true,
		},
		SapnetweaverProcessAvailability: MetricConfig{
			Enabled: true,
		},
		SapnetweaverQueueCount: MetricConfig{
			Enabled: true,
		},
		SapnetweaverQueueMaxCount: MetricConfig{
			Enabled: true,
		},
		SapnetweaverQueuePeakCount: MetricConfig{
			Enabled: true,
		},
		SapnetweaverRequestCount: MetricConfig{
			Enabled: true,
		},
		SapnetweaverRequestTimeoutCount: MetricConfig{
			Enabled: true,
		},
		SapnetweaverResponseDuration: MetricConfig{
			Enabled: true,
		},
		SapnetweaverSessionCount: MetricConfig{
			Enabled: true,
		},
		SapnetweaverSessionsBrowserCount: MetricConfig{
			Enabled: true,
		},
		SapnetweaverSessionsEjbCount: MetricConfig{
			Enabled: true,
		},
		SapnetweaverSessionsHTTPCount: MetricConfig{
			Enabled: true,
		},
		SapnetweaverSessionsSecurityCount: MetricConfig{
			Enabled: true,
		},
		SapnetweaverSessionsWebCount: MetricConfig{
			Enabled: true,
		},
		SapnetweaverShortDumpsRate: MetricConfig{
			Enabled: true,
		},
		SapnetweaverSpoolRequestErrorCount: MetricConfig{
			Enabled: true,
		},
		SapnetweaverSystemInstanceAvailability: MetricConfig{
			Enabled: true,
		},
		SapnetweaverWorkProcessActiveCount: MetricConfig{
			Enabled: true,
		},
		SapnetweaverWorkProcessJobAbortedStatus: MetricConfig{
			Enabled: true,
		},
	}
}

// ResourceAttributeConfig provides common config for a particular resource attribute.
type ResourceAttributeConfig struct {
	Enabled bool `mapstructure:"enabled"`
	// Experimental: MetricsInclude defines a list of filters for attribute values.
	// If the list is not empty, only metrics with matching resource attribute values will be emitted.
	MetricsInclude []filter.Config `mapstructure:"metrics_include"`
	// Experimental: MetricsExclude defines a list of filters for attribute values.
	// If the list is not empty, metrics with matching resource attribute values will not be emitted.
	// MetricsInclude has higher priority than MetricsExclude.
	MetricsExclude []filter.Config `mapstructure:"metrics_exclude"`

	enabledSetByUser bool
}

func (rac *ResourceAttributeConfig) Unmarshal(parser *confmap.Conf) error {
	if parser == nil {
		return nil
	}
	err := parser.Unmarshal(rac)
	if err != nil {
		return err
	}
	rac.enabledSetByUser = parser.IsSet("enabled")
	return nil
}

// ResourceAttributesConfig provides config for sapnetweaver resource attributes.
type ResourceAttributesConfig struct {
	SapnetweaverSID      ResourceAttributeConfig `mapstructure:"sapnetweaver.SID"`
	SapnetweaverInstance ResourceAttributeConfig `mapstructure:"sapnetweaver.instance"`
	SapnetweaverNode     ResourceAttributeConfig `mapstructure:"sapnetweaver.node"`
}

func DefaultResourceAttributesConfig() ResourceAttributesConfig {
	return ResourceAttributesConfig{
		SapnetweaverSID: ResourceAttributeConfig{
			Enabled: true,
		},
		SapnetweaverInstance: ResourceAttributeConfig{
			Enabled: true,
		},
		SapnetweaverNode: ResourceAttributeConfig{
			Enabled: true,
		},
	}
}

// MetricsBuilderConfig is a configuration for sapnetweaver metrics builder.
type MetricsBuilderConfig struct {
	Metrics            MetricsConfig            `mapstructure:"metrics"`
	ResourceAttributes ResourceAttributesConfig `mapstructure:"resource_attributes"`
}

func DefaultMetricsBuilderConfig() MetricsBuilderConfig {
	return MetricsBuilderConfig{
		Metrics:            DefaultMetricsConfig(),
		ResourceAttributes: DefaultResourceAttributesConfig(),
	}
}
