// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package varnishreceiver // import "github.com/observiq/observiq-otel-collector/receiver/varnishreceiver"

import (
	"go.opentelemetry.io/collector/model/pdata"

	"github.com/observiq/observiq-otel-collector/receiver/varnishreceiver/internal/metadata"
)

// FullStats holds stats from a 6.5+ response.
type FullStats struct {
	Version   int    `json:"version"`
	Timestamp string `json:"timestamp"`
	Stats     Stats
}

// Stats holds the metric stats.
type Stats struct {
	Stat struct {
		MAINBackendConn struct {
			Value int64 `json:"value"`
		} `json:"MAIN.backend_conn"`
		MAINBackendUnhealthy struct {
			Value int64 `json:"value"`
		} `json:"MAIN.backend_unhealthy"`
		MAINBackendBusy struct {
			Value int64 `json:"value"`
		} `json:"MAIN.backend_busy"`
		MAINBackendFail struct {
			Value int64 `json:"value"`
		} `json:"MAIN.backend_fail"`
		MAINBackendReuse struct {
			Value int64 `json:"value"`
		} `json:"MAIN.backend_reuse"`
		MAINBackendRecycle struct {
			Value int64 `json:"value"`
		} `json:"MAIN.backend_recycle"`
		MAINBackendRetry struct {
			Value int64 `json:"value"`
		} `json:"MAIN.backend_retry"`
		MAINCacheHit struct {
			Value int64 `json:"value"`
		} `json:"MAIN.cache_hit"`
		MAINCacheHitpass struct {
			Value int64 `json:"value"`
		} `json:"MAIN.cache_hitpass"`
		MAINCacheMiss struct {
			Value int64 `json:"value"`
		} `json:"MAIN.cache_miss"`
		MAINThreadsCreated struct {
			Value int64 `json:"value"`
		} `json:"MAIN.threads_created"`
		MAINThreadsDestroyed struct {
			Value int64 `json:"value"`
		} `json:"MAIN.threads_destroyed"`
		MAINThreadsFailed struct {
			Value int64 `json:"value"`
		} `json:"MAIN.threads_failed"`
		MAINSessConn struct {
			Value int64 `json:"value"`
		} `json:"MAIN.sess_conn"`
		MAINSessFail struct {
			Value int64 `json:"value"`
		} `json:"MAIN.sess_fail"`
		MAINSessDropped struct {
			Value int64 `json:"value"`
		} `json:"MAIN.sess_dropped"`
		MAINNObject struct {
			Value int64 `json:"value"`
		} `json:"MAIN.n_object"`
		MAINNExpired struct {
			Value int64 `json:"value"`
		} `json:"MAIN.n_expired"`
		MAINNLruNuked struct {
			Value int64 `json:"value"`
		} `json:"MAIN.n_lru_nuked"`
		MAINNLruMoved struct {
			Value int64 `json:"value"`
		} `json:"MAIN.n_lru_moved"`
		MAINClientReq struct {
			Value int64 `json:"value"`
		} `json:"MAIN.client_req"`
		MAINBackendReq struct {
			Value int64 `json:"value"`
		} `json:"MAIN.backend_req"`
	} `json:"counters"`
}

func (v *varnishScraper) recordVarnishBackendConnectionsCountDataPoint(now pdata.Timestamp, stats *Stats) {
	attributeMappings := map[string]int64{
		metadata.AttributeBackendConnectionType.Success:   stats.Stat.MAINBackendConn.Value,
		metadata.AttributeBackendConnectionType.Recycle:   stats.Stat.MAINBackendRecycle.Value,
		metadata.AttributeBackendConnectionType.Reuse:     stats.Stat.MAINBackendReuse.Value,
		metadata.AttributeBackendConnectionType.Fail:      stats.Stat.MAINBackendFail.Value,
		metadata.AttributeBackendConnectionType.Unhealthy: stats.Stat.MAINBackendUnhealthy.Value,
		metadata.AttributeBackendConnectionType.Busy:      stats.Stat.MAINBackendBusy.Value,
		metadata.AttributeBackendConnectionType.Retry:     stats.Stat.MAINBackendRetry.Value,
	}

	for attributeName, attributeValue := range attributeMappings {
		v.mb.RecordVarnishBackendConnectionsCountDataPoint(now, attributeValue, attributeName)
	}
}

func (v *varnishScraper) recordVarnishCacheOperationsCountDataPoint(now pdata.Timestamp, stats *Stats) {
	attributeMappings := map[string]int64{
		metadata.AttributeCacheOperations.Hit:     stats.Stat.MAINCacheHit.Value,
		metadata.AttributeCacheOperations.HitPass: stats.Stat.MAINCacheHitpass.Value,
		metadata.AttributeCacheOperations.Miss:    stats.Stat.MAINCacheMiss.Value,
	}

	for attributeName, attributeValue := range attributeMappings {
		v.mb.RecordVarnishCacheOperationsCountDataPoint(now, attributeValue, attributeName)
	}
}

func (v *varnishScraper) recordVarnishThreadOperationsCountDataPoint(now pdata.Timestamp, stats *Stats) {
	attributeMappings := map[string]int64{
		metadata.AttributeThreadOperations.Created:   stats.Stat.MAINThreadsCreated.Value,
		metadata.AttributeThreadOperations.Destroyed: stats.Stat.MAINThreadsDestroyed.Value,
		metadata.AttributeThreadOperations.Failed:    stats.Stat.MAINThreadsFailed.Value,
	}

	for attributeName, attributeValue := range attributeMappings {
		v.mb.RecordVarnishThreadOperationsCountDataPoint(now, attributeValue, attributeName)
	}
}

func (v *varnishScraper) recordVarnishSessionCountDataPoint(now pdata.Timestamp, stats *Stats) {
	attributeMappings := map[string]int64{
		metadata.AttributeSessionType.Accepted: stats.Stat.MAINSessConn.Value,
		metadata.AttributeSessionType.Dropped:  stats.Stat.MAINSessDropped.Value,
		metadata.AttributeSessionType.Failed:   stats.Stat.MAINSessFail.Value,
	}

	for attributeName, attributeValue := range attributeMappings {
		v.mb.RecordVarnishSessionCountDataPoint(now, attributeValue, attributeName)
	}
}
