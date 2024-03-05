// Copyright observIQ, Inc.
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

package telemetrygeneratorreceiver

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/golden"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest/pmetrictest"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest/ptracetest"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var expectedOTLPDir = filepath.Join("testdata", "expected_otlp")

func TestOTLPGenerator_Traces(t *testing.T) {
	test := []struct {
		name           string
		getCurrentTime func() time.Time
		cfg            GeneratorConfig
		expectedFile   string
	}{
		{
			name: "BPOP traces",
			getCurrentTime: func() time.Time {
				return time.Unix(0, 1706791445999459125)
			},
			cfg: GeneratorConfig{
				Type: generatorTypeOTLP,
				AdditionalConfig: map[string]any{
					"telemetry_type": "traces",
					"otlp_json": `{
						"resourceSpans": [
							{
								"resource": {
									"attributes": [
										{
											"key": "host.arch",
											"value": {
												"stringValue": "arm64"
											}
										},
										{
											"key": "host.name",
											"value": {
												"stringValue": "Sams-M1-Pro.local"
											}
										},
										{
											"key": "service.name",
											"value": {
												"stringValue": "bindplane"
											}
										},
										{
											"key": "service.version",
											"value": {
												"stringValue": "unknown"
											}
										}
									]
								},
								"scopeSpans": [
									{
										"scope": {},
										"spans": [
										
										
											{
												"endTimeUnixNano": "1706791445927702458",
												"kind": 1,
												"name": "pgstore/addTransitiveUpdates",
												"parentSpanId": "b6e1d82a58e8fd61",
												"spanId": "5a055056ad7713e5",
												"startTimeUnixNano": "1706791445927700000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445927704250",
												"kind": 1,
												"name": "pgstore/notify",
												"parentSpanId": "38ff7d679d77bdbd",
												"spanId": "b6e1d82a58e8fd61",
												"startTimeUnixNano": "1706791445927698000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
										
											{
												"endTimeUnixNano": "1706791445973101000",
												"kind": 1,
												"name": "pgstore/scanPostgresResource",
												"parentSpanId": "804ce3fb1b57be5d",
												"spanId": "3eadb90414f2cf22",
												"startTimeUnixNano": "1706791445971633000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445973109084",
												"kind": 1,
												"name": "pgstore/pgResourceInternal",
												"parentSpanId": "5236aa938eb9341c",
												"spanId": "804ce3fb1b57be5d",
												"startTimeUnixNano": "1706791445962593000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445973110334",
												"kind": 1,
												"name": "pgstore/pgResource",
												"parentSpanId": "ca3b8e53681a9e8d",
												"spanId": "5236aa938eb9341c",
												"startTimeUnixNano": "1706791445962592000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445973258875",
												"kind": 1,
												"name": "pgstore/pgEditConfiguration",
												"parentSpanId": "a894c9beda2d173a",
												"spanId": "ca3b8e53681a9e8d",
												"startTimeUnixNano": "1706791445962589000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445973576917",
												"kind": 1,
												"name": "pgstore/addTransitiveUpdates",
												"parentSpanId": "94557eef510d0814",
												"spanId": "c3ea8f993f9009ba",
												"startTimeUnixNano": "1706791445973575000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445973579333",
												"kind": 1,
												"name": "pgstore/notify",
												"parentSpanId": "a894c9beda2d173a",
												"spanId": "94557eef510d0814",
												"startTimeUnixNano": "1706791445973572000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445975115042",
												"kind": 1,
												"name": "pgstore/releaseAdvisoryLock",
												"parentSpanId": "a894c9beda2d173a",
												"spanId": "8906c43abedb6bd9",
												"startTimeUnixNano": "1706791445973581000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445975119416",
												"kind": 1,
												"name": "pgstore/UpdateRollout",
												"parentSpanId": "9d5a7e824fa7ba3b",
												"spanId": "a894c9beda2d173a",
												"startTimeUnixNano": "1706791445937536000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445975734417",
												"kind": 1,
												"name": "pgstore/acquireAdvisoryLock",
												"parentSpanId": "078ab0ab7eb707b2",
												"spanId": "f141ff98614cfc0c",
												"startTimeUnixNano": "1706791445975133000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445979767041",
												"kind": 1,
												"name": "pgstore/scanPostgresResource",
												"parentSpanId": "f06908731ed62890",
												"spanId": "e185a7e1c60473b8",
												"startTimeUnixNano": "1706791445976351000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445979799208",
												"kind": 1,
												"name": "pgstore/pgResourceInternal",
												"parentSpanId": "078ab0ab7eb707b2",
												"spanId": "f06908731ed62890",
												"startTimeUnixNano": "1706791445975818000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
										
											{
												"endTimeUnixNano": "1706791445999459125",
												"kind": 1,
												"name": "pgstore/UpdateAllRollouts",
												"parentSpanId": "",
												"spanId": "aeb2a416b8796cba",
												"startTimeUnixNano": "1706791445908223000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											}
										]
									},
									{
										"scope": {},
										"spans": [
											{
												"attributes": [
													{
														"key": "operation",
														"value": {
															"stringValue": "GetConfiguration"
														}
													}
												],
												"endTimeUnixNano": "1706791445376564375",
												"kind": 1,
												"name": "graphql/GetConfiguration/response",
												"parentSpanId": "723c3f6eb4457b5c",
												"spanId": "fd55f461239efdfc",
												"startTimeUnixNano": "1706791445359466000",
												"status": {},
												"traceId": "a3fbd5dc5db5e1734cb54419ca540b66"
											},
											{
												"attributes": [
													{
														"key": "operation",
														"value": {
															"stringValue": "GetConfiguration"
														}
													}
												],
												"endTimeUnixNano": "1706791445376589750",
												"kind": 1,
												"name": "graphql/GetConfiguration/response",
												"parentSpanId": "3e7909bbebcae0ba",
												"spanId": "01f07757e7bb6612",
												"startTimeUnixNano": "1706791445359560000",
												"status": {},
												"traceId": "d70c2b5eea8977bb8a0712f8c2a1fcb4"
											}
										]
									},
									{
										"scope": {},
										"spans": [
											{
												"attributes": [
													{
														"key": "http.method",
														"value": {
															"stringValue": "POST"
														}
													},
													{
														"key": "http.scheme",
														"value": {
															"stringValue": "http"
														}
													},
													{
														"key": "net.host.name",
														"value": {
															"stringValue": "bindplane"
														}
													},
													{
														"key": "net.host.port",
														"value": {
															"intValue": "3001"
														}
													},
													{
														"key": "net.sock.peer.addr",
														"value": {
															"stringValue": "127.0.0.1"
														}
													},
													{
														"key": "net.sock.peer.port",
														"value": {
															"intValue": "50141"
														}
													},
													{
														"key": "user_agent.original",
														"value": {
															"stringValue": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:123.0) Gecko/20100101 Firefox/123.0"
														}
													},
													{
														"key": "http.client_ip",
														"value": {
															"stringValue": "127.0.0.1"
														}
													},
													{
														"key": "http.target",
														"value": {
															"stringValue": "/v1/graphql"
														}
													},
													{
														"key": "net.protocol.version",
														"value": {
															"stringValue": "1.1"
														}
													},
													{
														"key": "http.route",
														"value": {
															"stringValue": "/v1/graphql"
														}
													},
													{
														"key": "http.status_code",
														"value": {
															"intValue": "200"
														}
													}
												],
												"endTimeUnixNano": "1706791445376694750",
												"kind": 2,
												"name": "/v1/graphql",
												"parentSpanId": "",
												"spanId": "723c3f6eb4457b5c",
												"startTimeUnixNano": "1706791445332980000",
												"status": {},
												"traceId": "a3fbd5dc5db5e1734cb54419ca540b66"
											},
											{
												"attributes": [
													{
														"key": "http.method",
														"value": {
															"stringValue": "POST"
														}
													},
													{
														"key": "http.scheme",
														"value": {
															"stringValue": "http"
														}
													},
													{
														"key": "net.host.name",
														"value": {
															"stringValue": "bindplane"
														}
													},
													{
														"key": "net.host.port",
														"value": {
															"intValue": "3001"
														}
													},
													{
														"key": "net.sock.peer.addr",
														"value": {
															"stringValue": "127.0.0.1"
														}
													},
													{
														"key": "net.sock.peer.port",
														"value": {
															"intValue": "50140"
														}
													},
													{
														"key": "user_agent.original",
														"value": {
															"stringValue": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:123.0) Gecko/20100101 Firefox/123.0"
														}
													},
													{
														"key": "http.client_ip",
														"value": {
															"stringValue": "127.0.0.1"
														}
													},
													{
														"key": "http.target",
														"value": {
															"stringValue": "/v1/graphql"
														}
													},
													{
														"key": "net.protocol.version",
														"value": {
															"stringValue": "1.1"
														}
													},
													{
														"key": "http.route",
														"value": {
															"stringValue": "/v1/graphql"
														}
													},
													{
														"key": "http.status_code",
														"value": {
															"intValue": "200"
														}
													}
												],
												"endTimeUnixNano": "1706791445376708291",
												"kind": 2,
												"name": "/v1/graphql",
												"parentSpanId": "",
												"spanId": "3e7909bbebcae0ba",
												"startTimeUnixNano": "1706791445332972000",
												"status": {},
												"traceId": "d70c2b5eea8977bb8a0712f8c2a1fcb4"
											}
										]
									},
									{
										"scope": {},
										"spans": [
											{
												"endTimeUnixNano": "1706791445913878000",
												"kind": 1,
												"name": "pgindex/Suggestions",
												"parentSpanId": "9d5a7e824fa7ba3b",
												"spanId": "4c2049c4cd14c987",
												"startTimeUnixNano": "1706791445912675000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445997017791",
												"kind": 1,
												"name": "pgindex/Suggestions",
												"parentSpanId": "96ae55c03e5146b3",
												"spanId": "aa69c45bc0970c2f",
												"startTimeUnixNano": "1706791445996229000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											}
										]
									}
								]
							}
						]
					}`,
				},
			},
			expectedFile: filepath.Join(expectedOTLPDir, "traces", "bpop_traces.yaml"),
		},
		{
			name: "BPOP traces 2",
			getCurrentTime: func() time.Time {
				return time.Unix(0, 1706791445999459839)
			},
			cfg: GeneratorConfig{
				Type: generatorTypeOTLP,
				AdditionalConfig: map[string]any{
					"telemetry_type": "traces",
					"otlp_json": `{
						"resourceSpans": [
							{
								"resource": {
									"attributes": [
										{
											"key": "host.arch",
											"value": {
												"stringValue": "arm64"
											}
										},
										{
											"key": "host.name",
											"value": {
												"stringValue": "Sams-M1-Pro.local"
											}
										},
										{
											"key": "service.name",
											"value": {
												"stringValue": "bindplane"
											}
										},
										{
											"key": "service.version",
											"value": {
												"stringValue": "unknown"
											}
										}
									]
								},
								"scopeSpans": [
									{
										"scope": {},
										"spans": [
										
										
											{
												"endTimeUnixNano": "1706791445927702458",
												"kind": 1,
												"name": "pgstore/addTransitiveUpdates",
												"parentSpanId": "b6e1d82a58e8fd61",
												"spanId": "5a055056ad7713e5",
												"startTimeUnixNano": "1706791445927700000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445927704250",
												"kind": 1,
												"name": "pgstore/notify",
												"parentSpanId": "38ff7d679d77bdbd",
												"spanId": "b6e1d82a58e8fd61",
												"startTimeUnixNano": "1706791445927698000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
										
											{
												"endTimeUnixNano": "1706791445973101000",
												"kind": 1,
												"name": "pgstore/scanPostgresResource",
												"parentSpanId": "804ce3fb1b57be5d",
												"spanId": "3eadb90414f2cf22",
												"startTimeUnixNano": "1706791445971633000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445973109084",
												"kind": 1,
												"name": "pgstore/pgResourceInternal",
												"parentSpanId": "5236aa938eb9341c",
												"spanId": "804ce3fb1b57be5d",
												"startTimeUnixNano": "1706791445962593000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445973110334",
												"kind": 1,
												"name": "pgstore/pgResource",
												"parentSpanId": "ca3b8e53681a9e8d",
												"spanId": "5236aa938eb9341c",
												"startTimeUnixNano": "1706791445962592000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445973258875",
												"kind": 1,
												"name": "pgstore/pgEditConfiguration",
												"parentSpanId": "a894c9beda2d173a",
												"spanId": "ca3b8e53681a9e8d",
												"startTimeUnixNano": "1706791445962589000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445973576917",
												"kind": 1,
												"name": "pgstore/addTransitiveUpdates",
												"parentSpanId": "94557eef510d0814",
												"spanId": "c3ea8f993f9009ba",
												"startTimeUnixNano": "1706791445973575000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445973579333",
												"kind": 1,
												"name": "pgstore/notify",
												"parentSpanId": "a894c9beda2d173a",
												"spanId": "94557eef510d0814",
												"startTimeUnixNano": "1706791445973572000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445975115042",
												"kind": 1,
												"name": "pgstore/releaseAdvisoryLock",
												"parentSpanId": "a894c9beda2d173a",
												"spanId": "8906c43abedb6bd9",
												"startTimeUnixNano": "1706791445973581000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445975119416",
												"kind": 1,
												"name": "pgstore/UpdateRollout",
												"parentSpanId": "9d5a7e824fa7ba3b",
												"spanId": "a894c9beda2d173a",
												"startTimeUnixNano": "1706791445937536000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445975734417",
												"kind": 1,
												"name": "pgstore/acquireAdvisoryLock",
												"parentSpanId": "078ab0ab7eb707b2",
												"spanId": "f141ff98614cfc0c",
												"startTimeUnixNano": "1706791445975133000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445979767041",
												"kind": 1,
												"name": "pgstore/scanPostgresResource",
												"parentSpanId": "f06908731ed62890",
												"spanId": "e185a7e1c60473b8",
												"startTimeUnixNano": "1706791445976351000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445979799208",
												"kind": 1,
												"name": "pgstore/pgResourceInternal",
												"parentSpanId": "078ab0ab7eb707b2",
												"spanId": "f06908731ed62890",
												"startTimeUnixNano": "1706791445975818000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
										
											{
												"endTimeUnixNano": "1706791445999459839",
												"kind": 1,
												"name": "pgstore/UpdateAllRollouts",
												"parentSpanId": "",
												"spanId": "aeb2a416b8796cba",
												"startTimeUnixNano": "1706791445908223000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											}
										]
									},									
									{
										"scope": {},
										"spans": [
											{
												"attributes": [
													{
														"key": "http.method",
														"value": {
															"stringValue": "POST"
														}
													},
													{
														"key": "http.scheme",
														"value": {
															"stringValue": "http"
														}
													},
													{
														"key": "net.host.name",
														"value": {
															"stringValue": "bindplane"
														}
													},
													{
														"key": "net.host.port",
														"value": {
															"intValue": "3001"
														}
													},
													{
														"key": "net.sock.peer.addr",
														"value": {
															"stringValue": "127.0.0.1"
														}
													},
													{
														"key": "net.sock.peer.port",
														"value": {
															"intValue": "50141"
														}
													},
													{
														"key": "user_agent.original",
														"value": {
															"stringValue": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:123.0) Gecko/20100101 Firefox/123.0"
														}
													},
													{
														"key": "http.client_ip",
														"value": {
															"stringValue": "127.0.0.1"
														}
													},
													{
														"key": "http.target",
														"value": {
															"stringValue": "/v1/graphql"
														}
													},
													{
														"key": "net.protocol.version",
														"value": {
															"stringValue": "1.1"
														}
													},
													{
														"key": "http.route",
														"value": {
															"stringValue": "/v1/graphql"
														}
													},
													{
														"key": "http.status_code",
														"value": {
															"intValue": "200"
														}
													}
												],
												"endTimeUnixNano": "1706791445376694750",
												"kind": 2,
												"name": "/v1/graphql",
												"parentSpanId": "",
												"spanId": "723c3f6eb4457b5c",
												"startTimeUnixNano": "1706791445332980000",
												"status": {},
												"traceId": "a3fbd5dc5db5e1734cb54419ca540b66"
											},
											{
												"attributes": [
													{
														"key": "http.method",
														"value": {
															"stringValue": "POST"
														}
													},
													{
														"key": "http.scheme",
														"value": {
															"stringValue": "http"
														}
													},
													{
														"key": "net.host.name",
														"value": {
															"stringValue": "bindplane"
														}
													},
													{
														"key": "net.host.port",
														"value": {
															"intValue": "3001"
														}
													},
													{
														"key": "net.sock.peer.addr",
														"value": {
															"stringValue": "127.0.0.1"
														}
													},
													{
														"key": "net.sock.peer.port",
														"value": {
															"intValue": "50140"
														}
													},
													{
														"key": "user_agent.original",
														"value": {
															"stringValue": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:123.0) Gecko/20100101 Firefox/123.0"
														}
													},
													{
														"key": "http.client_ip",
														"value": {
															"stringValue": "127.0.0.1"
														}
													},
													{
														"key": "http.target",
														"value": {
															"stringValue": "/v1/graphql"
														}
													},
													{
														"key": "net.protocol.version",
														"value": {
															"stringValue": "1.1"
														}
													},
													{
														"key": "http.route",
														"value": {
															"stringValue": "/v1/graphql"
														}
													},
													{
														"key": "http.status_code",
														"value": {
															"intValue": "200"
														}
													}
												],
												"endTimeUnixNano": "1706791445376708291",
												"kind": 2,
												"name": "/v1/graphql",
												"parentSpanId": "",
												"spanId": "3e7909bbebcae0ba",
												"startTimeUnixNano": "1706791445332972000",
												"status": {},
												"traceId": "d70c2b5eea8977bb8a0712f8c2a1fcb4"
											}
										]
									},
									{
										"scope": {},
										"spans": [
											{
												"endTimeUnixNano": "1706791445913878000",
												"kind": 1,
												"name": "pgindex/Suggestions",
												"parentSpanId": "9d5a7e824fa7ba3b",
												"spanId": "4c2049c4cd14c987",
												"startTimeUnixNano": "1706791445912675000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											},
											{
												"endTimeUnixNano": "1706791445997017791",
												"kind": 1,
												"name": "pgindex/Suggestions",
												"parentSpanId": "96ae55c03e5146b3",
												"spanId": "aa69c45bc0970c2f",
												"startTimeUnixNano": "1706791445996229000",
												"status": {},
												"traceId": "c7f3bb6aa9e7a7dce92d85d1566f2c31"
											}
										]
									}
								]
							}
						]
					}`,
				},
			},
			expectedFile: filepath.Join(expectedOTLPDir, "traces", "bpop_traces2.yaml"),
		},
	}
	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			getCurrentTime = tc.getCurrentTime
			err := tc.cfg.Validate()
			require.NoError(t, err)

			g := newOtlpGenerator(tc.cfg, zap.NewNop())
			traces := g.generateTraces()
			// clearTimeStamps(traces)

			// golden.WriteTraces(t, tc.expectedFile, traces)
			expectedTraces, err := golden.ReadTraces(tc.expectedFile)
			require.NoError(t, err)
			// unmarshaler := &ptrace.JSONMarshaler{}
			// fileBytes, _ := unmarshaler.MarshalTraces(traces)

			// os.WriteFile(filepath.Join(expectedOTLPDir, "bpop_traces.json"), fileBytes, 0600)
			// clearTimeStamps(expectedLogs)
			err = ptracetest.CompareTraces(expectedTraces, traces)
			require.NoError(t, err)

			// require.NoError(t, err)
			// require.NotNil(t, config)

		})
	}
}

func Test_findLastTraceEndTime(t *testing.T) {

	tests := []struct {
		name           string
		traceFile      string
		expectedTime   time.Time
		getCurrentTime func() time.Time
	}{

		{
			name:         "Traces 1",
			traceFile:    filepath.Join(expectedOTLPDir, "traces", "bpop_traces.yaml"),
			expectedTime: time.Date(2024, time.February, 1, 12, 44, 6, 622353875, time.UTC),
		},
		{
			name:         "Traces 2",
			traceFile:    filepath.Join(expectedOTLPDir, "traces", "bpop_traces2.yaml"),
			expectedTime: time.Date(2024, time.February, 1, 12, 44, 6, 622224928, time.UTC),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			getCurrentTime = tc.getCurrentTime
			traces, err := golden.ReadTraces(tc.traceFile)
			require.NoError(t, err)
			lastTime := findLastTraceEndTime(traces)
			require.Equal(t, tc.expectedTime, lastTime)
		})
	}
}

func TestOTLPGenerator_Metrics(t *testing.T) {

	tests := []struct {
		name           string
		getCurrentTime func() time.Time
		cfg            GeneratorConfig
		expectedFile   string
	}{
		{
			name: "exp histogram",
			getCurrentTime: func() time.Time {
				return time.Unix(0, 1706791445999459839)
			},
			cfg: GeneratorConfig{
				Type: generatorTypeOTLP,
				AdditionalConfig: map[string]any{
					"telemetry_type": "metrics",
					"otlp_json":      `{"resourceMetrics":[{"resource":{},"scopeMetrics":[{"scope":{},"metrics":[{"exponentialHistogram":{"dataPoints":[{"attributes":[{"key":"prod-machine","value":{"stringValue":"prod-1"}}],"count":"4","positive":{},"negative":{},"min":0,"max":100}]}}]}]}]}`,
				},
			},
			expectedFile: filepath.Join(expectedOTLPDir, "metrics", "exp_histogram.yaml"),
		},
		{
			name: "gauge",
			getCurrentTime: func() time.Time {
				return time.Unix(0, 1706791445999459839)
			},
			cfg: GeneratorConfig{
				Type: generatorTypeOTLP,
				AdditionalConfig: map[string]any{
					"telemetry_type": "metrics",
					"otlp_json":      `{"resourceMetrics":[{"resource":{"attributes":[{"key":"extra-resource-attr-key","value":{"stringValue":"extra-resource-attr-value"}},{"key":"host.name","value":{"stringValue":"Linux-Machine"}},{"key":"os.type","value":{"stringValue":"linux"}}]},"scopeMetrics":[{"scope":{},"metrics":[{"name":"system.cpu.load_average.1m","description":"Average CPU Load over 1 minute.","unit":"{thread}","gauge":{"dataPoints":[{"attributes":[{"key":"cool-attribute-key","value":{"stringValue":"cool-attribute-value"}}],"startTimeUnixNano":"1000000","timeUnixNano":"2000000","asDouble":3.71484375}]}}]}]}]}`,
				},
			},
			expectedFile: filepath.Join(expectedOTLPDir, "metrics", "gauge.yaml"),
		},
		{
			name: "histogram",
			getCurrentTime: func() time.Time {
				return time.Unix(0, 1706791445999459839)
			},
			cfg: GeneratorConfig{
				Type: generatorTypeOTLP,
				AdditionalConfig: map[string]any{
					"telemetry_type": "metrics",
					"otlp_json":      `{"resourceMetrics":[{"resource":{},"scopeMetrics":[{"scope":{},"metrics":[{"histogram":{"dataPoints":[{"attributes":[{"key":"prod-machine","value":{"stringValue":"prod-1"}}],"count":"4","bucketCounts":["0","2","2"],"explicitBounds":[0,50,100],"min":0,"max":100}]}}]}]}]}`,
				},
			},
			expectedFile: filepath.Join(expectedOTLPDir, "metrics", "histogram.yaml"),
		},
		{
			name: "sum",
			getCurrentTime: func() time.Time {
				return time.Unix(0, 1706791445999459839)
			},
			cfg: GeneratorConfig{
				Type: generatorTypeOTLP,
				AdditionalConfig: map[string]any{
					"telemetry_type": "metrics",
					"otlp_json":      `{"resourceMetrics":[{"resource":{"attributes":[{"key":"host.name","value":{"stringValue":"Sams-MBP"}},{"key":"os.type","value":{"stringValue":"darwin"}}]},"scopeMetrics":[{"scope":{},"metrics":[{"name":"system.filesystem.usage","description":"Filesystem bytes used.","unit":"By","sum":{"dataPoints":[{"attributes":[{"key":"device","value":{"stringValue":"/dev/disk4s1"}},{"key":"extra-sum-attr-key","value":{"stringValue":"extra-sum-attr-value"}},{"key":"mode","value":{"stringValue":"rw"}},{"key":"mountpoint","value":{"stringValue":"/Volumes/transfer"}},{"key":"state","value":{"stringValue":"free"}},{"key":"type","value":{"stringValue":"hfs"}}],"startTimeUnixNano":"1000000","timeUnixNano":"2000000","asInt":"8717185024"}]}}]}]}]}`,
				},
			},
			expectedFile: filepath.Join(expectedOTLPDir, "metrics", "sum.yaml"),
		},
		{
			name: "summary",
			getCurrentTime: func() time.Time {
				return time.Unix(0, 1706791445999459839)
			},
			cfg: GeneratorConfig{
				Type: generatorTypeOTLP,
				AdditionalConfig: map[string]any{
					"telemetry_type": "metrics",
					"otlp_json":      `{"resourceMetrics":[{"resource":{},"scopeMetrics":[{"scope":{},"metrics":[{"summary":{"dataPoints":[{"attributes":[{"key":"prod-machine","value":{"stringValue":"prod-1"}}],"count":"4","sum":111}]}}]}]}]}`,
				},
			},
			expectedFile: filepath.Join(expectedOTLPDir, "metrics", "summary.yaml"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			getCurrentTime = tc.getCurrentTime
			err := tc.cfg.Validate()
			require.NoError(t, err)

			g := newOtlpGenerator(tc.cfg, zap.NewNop())
			metrics := g.generateMetrics()

			// golden.WriteMetrics(t, tc.expectedFile, metrics)
			expectedMetrics, err := golden.ReadMetrics(tc.expectedFile)
			require.NoError(t, err)
			// unmarshaler := &pmetric.JSONMarshaler{}
			// fileBytes, _ := unmarshaler.MarshalMetrics(expectedMetrics)

			// os.WriteFile(tc.outFile, fileBytes, 0600)
			// clearTimeStamps(expectedLogs)
			err = pmetrictest.CompareMetrics(expectedMetrics, metrics)
			require.NoError(t, err)

			// require.NoError(t, err)
			// require.NotNil(t, config)

		})
	}
}
