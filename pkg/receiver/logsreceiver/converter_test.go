// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logsreceiver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/open-telemetry/opentelemetry-log-collection/entry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/model/pdata"
)

func BenchmarkConvertSimple(b *testing.B) {
	b.StopTimer()
	ent := entry.New()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		convert(ent, nil)
	}
}

func BenchmarkConvertComplex(b *testing.B) {
	b.StopTimer()
	ent := complexEntry()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		convert(ent, nil)
	}
}

func complexEntries(count int) []*entry.Entry {
	return complexEntriesForNDifferentHosts(count, 1)
}

func complexEntriesForNDifferentHosts(count int, n int) []*entry.Entry {
	ret := make([]*entry.Entry, count)
	for i := 0; i < count; i++ {
		e := entry.New()
		e.Severity = entry.Error
		e.AddResourceKey("type", "global")
		e.Resource = map[string]string{
			"host": fmt.Sprintf("host-%d", i%n),
		}
		e.Body = map[string]interface{}{
			"bool":   true,
			"int":    123,
			"double": 12.34,
			"string": "hello",
			"bytes":  []byte("asdf"),
			"object": map[string]interface{}{
				"bool":   true,
				"int":    123,
				"double": 12.34,
				"string": "hello",
				"bytes":  []byte("asdf"),
				"object": map[string]interface{}{
					"bool":   true,
					"int":    123,
					"double": 12.34,
					"string": "hello",
					"bytes":  []byte("asdf"),
				},
			},
		}
		ret[i] = e
	}
	return ret
}

func complexEntry() *entry.Entry {
	e := entry.New()
	e.Severity = entry.Error
	e.AddResourceKey("type", "global")
	e.AddAttribute("one", "two")
	e.AddAttribute("two", "three")
	e.Body = map[string]interface{}{
		"bool":   true,
		"int":    123,
		"double": 12.34,
		"string": "hello",
		// "bytes":  []byte("asdf"),
		"object": map[string]interface{}{
			"bool":   true,
			"int":    123,
			"double": 12.34,
			"string": "hello",
			// "bytes":  []byte("asdf"),
			"object": map[string]interface{}{
				"bool": true,
				"int":  123,
				// "double": 12.34,
				"string": "hello",
				// "bytes":  []byte("asdf"),
			},
		},
	}
	return e
}

/*func TestConvert(t *testing.T) {
	ent := func() *entry.Entry {
		e := entry.New()
		e.Severity = entry.Error
		e.AddResourceKey("type", "global")
		e.AddAttribute("one", "two")
		e.AddAttribute("two", "three")
		e.Body = map[string]interface{}{
			"bool":   true,
			"int":    123,
			"double": 12.34,
			"string": "hello",
			"bytes":  []byte("asdf"),
		}
		return e
	}()

	pLogs := Convert(ent)
	require.Equal(t, 1, pLogs.ResourceLogs().Len())
	rls := pLogs.ResourceLogs().At(0)
	require.Equal(t, 1, rls.Resource().Attributes().Len())
	{
		att, ok := rls.Resource().Attributes().Get("type")
		if assert.True(t, ok) {
			if assert.Equal(t, att.Type(), pdata.AttributeValueTypeString) {
				assert.Equal(t, att.StringVal(), "global")
			}
		}
	}

	ills := rls.InstrumentationLibraryLogs()
	require.Equal(t, 1, ills.Len())

	logs := ills.At(0).Logs()
	require.Equal(t, 1, logs.Len())

	lr := logs.At(0)

	assert.Equal(t, pdata.SeverityNumberERROR, lr.SeverityNumber())
	assert.Equal(t, "Error", lr.SeverityText())

	if atts := lr.Attributes(); assert.Equal(t, 2, atts.Len()) {
		m := pdata.NewAttributeMap()
		m.InitFromMap(map[string]pdata.AttributeValue{
			"one": pdata.NewAttributeValueString("two"),
			"two": pdata.NewAttributeValueString("three"),
		})
		assert.EqualValues(t, m.Sort(), atts.Sort())
	}

	if assert.Equal(t, pdata.AttributeValueTypeMap, lr.Body().Type()) {
		m := pdata.NewAttributeMap()
		m.InitFromMap(map[string]pdata.AttributeValue{
			"bool":   pdata.NewAttributeValueBool(true),
			"int":    pdata.NewAttributeValueInt(123),
			"double": pdata.NewAttributeValueDouble(12.34),
			"string": pdata.NewAttributeValueString("hello"),
			"bytes":  pdata.NewAttributeValueString("asdf"),
			// Don't include a nested object because AttributeValueMap sorting
			// doesn't sort recursively.
		})
		assert.EqualValues(t, m.Sort(), lr.Body().MapVal().Sort())
	}
}*/

func TestAllConvertedEntriesAreSentAndReceived(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		entries int
	}{
		{
			entries: 10,
		},
		{
			entries: 100,
		},
	}

	for i, tc := range testcases {
		tc := tc

		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			converter := NewConverter(
				WithWorkerCount(1),
			)
			converter.Start()
			defer converter.Stop()

			go func() {
				entries := complexEntries(tc.entries)
				assert.NoError(t, converter.Batch(entries))
			}()

			var (
				actualCount  int
				timeoutTimer = time.NewTimer(10 * time.Second)
				ch           = converter.OutChannel()
			)
			defer timeoutTimer.Stop()

		forLoop:
			for {
				if tc.entries == actualCount {
					break
				}

				select {
				case pLogs, ok := <-ch:
					if !ok {
						break forLoop
					}

					rLogs := pLogs.ResourceLogs()
					require.Equal(t, 1, rLogs.Len())

					rLog := rLogs.At(0)
					ills := rLog.InstrumentationLibraryLogs()
					require.Equal(t, 1, ills.Len())

					ill := ills.At(0)

					actualCount += ill.Logs().Len()
				case <-timeoutTimer.C:
					break forLoop
				}
			}

			assert.Equal(t, tc.entries, actualCount,
				"didn't receive expected number of entries after conversion",
			)
		})
	}
}

func TestAllConvertedEntriesAreSentAndReceivedWithinAnExpectedTimeDuration(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		entries    int
		hostsCount int
	}{
		{
			entries:    10,
			hostsCount: 1,
		},
		{
			entries:    50,
			hostsCount: 1,
		},
		{
			entries:    500,
			hostsCount: 1,
		},
		{
			entries:    500,
			hostsCount: 1,
		},
		{
			entries:    500,
			hostsCount: 4,
		},
	}

	for i, tc := range testcases {
		tc := tc

		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			converter := NewConverter(
				WithWorkerCount(1),
			)
			converter.Start()
			defer converter.Stop()

			go func() {
				entries := complexEntriesForNDifferentHosts(tc.entries, tc.hostsCount)
				assert.NoError(t, converter.Batch(entries))
			}()

			var (
				actualCount      int
				actualFlushCount int
				timeoutTimer     = time.NewTimer(10 * time.Second)
				ch               = converter.OutChannel()
			)
			defer timeoutTimer.Stop()

		forLoop:
			for {
				if tc.entries == actualCount {
					break
				}

				select {
				case pLogs, ok := <-ch:
					if !ok {
						break forLoop
					}

					actualFlushCount++

					rLogs := pLogs.ResourceLogs()
					require.Equal(t, 1, rLogs.Len())

					rLog := rLogs.At(0)
					ills := rLog.InstrumentationLibraryLogs()
					require.Equal(t, 1, ills.Len())

					ill := ills.At(0)

					actualCount += ill.Logs().Len()

				case <-timeoutTimer.C:
					break forLoop
				}
			}

			assert.Equal(t, tc.entries, actualCount,
				"didn't receive expected number of entries after conversion",
			)
		})
	}
}

func TestConverterCancelledContextCancellsTheFlush(t *testing.T) {
	converter := NewConverter()
	converter.Start()
	defer converter.Stop()
	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	go func() {
		defer wg.Done()
		pLogs := pdata.NewLogs()
		ills := pLogs.ResourceLogs().AppendEmpty().InstrumentationLibraryLogs().AppendEmpty()

		lr := convert(complexEntry(), nil)
		lr.CopyTo(ills.Logs().AppendEmpty())

		assert.Error(t, converter.flush(ctx, pLogs))
	}()
	wg.Wait()
}

func TestConvertMetadata(t *testing.T) {
	now := time.Now()

	e := entry.New()
	e.Timestamp = now
	e.Severity = entry.Error
	e.AddResourceKey("type", "global")
	e.AddAttribute("one", "two")
	e.Body = true

	result := convert(e, nil)

	atts := result.Attributes()
	require.Equal(t, 1, atts.Len(), "expected 1 attribute")
	attVal, ok := atts.Get("one")
	require.True(t, ok, "expected label with key 'one'")
	require.Equal(t, "two", attVal.StringVal(), "expected label to have value 'two'")

	bod := result.Body()
	require.Equal(t, pdata.AttributeValueTypeBool, bod.Type())
	require.True(t, bod.BoolVal())
}

func TestConvertSimpleBody(t *testing.T) {
	require.True(t, anyToBody(true).BoolVal())
	require.False(t, anyToBody(false).BoolVal())

	require.Equal(t, "string", anyToBody("string").StringVal())
	require.Equal(t, "bytes", anyToBody([]byte("bytes")).StringVal())

	require.Equal(t, int64(1), anyToBody(1).IntVal())
	require.Equal(t, int64(1), anyToBody(int8(1)).IntVal())
	require.Equal(t, int64(1), anyToBody(int16(1)).IntVal())
	require.Equal(t, int64(1), anyToBody(int32(1)).IntVal())
	require.Equal(t, int64(1), anyToBody(int64(1)).IntVal())

	require.Equal(t, int64(1), anyToBody(uint(1)).IntVal())
	require.Equal(t, int64(1), anyToBody(uint8(1)).IntVal())
	require.Equal(t, int64(1), anyToBody(uint16(1)).IntVal())
	require.Equal(t, int64(1), anyToBody(uint32(1)).IntVal())
	require.Equal(t, int64(1), anyToBody(uint64(1)).IntVal())

	require.Equal(t, float64(1), anyToBody(float32(1)).DoubleVal())
	require.Equal(t, float64(1), anyToBody(float64(1)).DoubleVal())
}

func TestConvertMapBody(t *testing.T) {
	structuredBody := map[string]interface{}{
		"true":    true,
		"false":   false,
		"string":  "string",
		"bytes":   []byte("bytes"),
		"int":     1,
		"int8":    int8(1),
		"int16":   int16(1),
		"int32":   int32(1),
		"int64":   int64(1),
		"uint":    uint(1),
		"uint8":   uint8(1),
		"uint16":  uint16(1),
		"uint32":  uint32(1),
		"uint64":  uint64(1),
		"float32": float32(1),
		"float64": float64(1),
	}

	result := anyToBody(structuredBody).MapVal()

	v, _ := result.Get("true")
	require.True(t, v.BoolVal())
	v, _ = result.Get("false")
	require.False(t, v.BoolVal())

	for _, k := range []string{"string", "bytes"} {
		v, _ = result.Get(k)
		require.Equal(t, k, v.StringVal())
	}
	for _, k := range []string{"int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64"} {
		v, _ = result.Get(k)
		require.Equal(t, int64(1), v.IntVal())
	}
	for _, k := range []string{"float32", "float64"} {
		v, _ = result.Get(k)
		require.Equal(t, float64(1), v.DoubleVal())
	}
}

func TestConvertArrayBody(t *testing.T) {
	structuredBody := []interface{}{
		true,
		false,
		"string",
		[]byte("bytes"),
		1,
		int8(1),
		int16(1),
		int32(1),
		int64(1),
		uint(1),
		uint8(1),
		uint16(1),
		uint32(1),
		uint64(1),
		float32(1),
		float64(1),
		[]interface{}{"string", 1},
		map[string]interface{}{"one": 1, "yes": true},
	}

	result := anyToBody(structuredBody).SliceVal()

	require.True(t, result.At(0).BoolVal())
	require.False(t, result.At(1).BoolVal())
	require.Equal(t, "string", result.At(2).StringVal())
	require.Equal(t, "bytes", result.At(3).StringVal())

	require.Equal(t, int64(1), result.At(4).IntVal())  // int
	require.Equal(t, int64(1), result.At(5).IntVal())  // int8
	require.Equal(t, int64(1), result.At(6).IntVal())  // int16
	require.Equal(t, int64(1), result.At(7).IntVal())  // int32
	require.Equal(t, int64(1), result.At(8).IntVal())  // int64
	require.Equal(t, int64(1), result.At(9).IntVal())  // uint
	require.Equal(t, int64(1), result.At(10).IntVal()) // uint8
	require.Equal(t, int64(1), result.At(11).IntVal()) // uint16
	require.Equal(t, int64(1), result.At(12).IntVal()) // uint32
	require.Equal(t, int64(1), result.At(13).IntVal()) // uint64

	require.Equal(t, float64(1), result.At(14).DoubleVal()) // float32
	require.Equal(t, float64(1), result.At(15).DoubleVal()) // float64

	nestedArr := result.At(16).SliceVal()
	require.Equal(t, "string", nestedArr.At(0).StringVal())
	require.Equal(t, int64(1), nestedArr.At(1).IntVal())

	nestedMap := result.At(17).MapVal()
	v, _ := nestedMap.Get("one")
	require.Equal(t, int64(1), v.IntVal())
	v, _ = nestedMap.Get("yes")
	require.True(t, v.BoolVal())
}

func TestConvertUnknownBody(t *testing.T) {
	unknownType := map[string]int{"0": 0, "1": 1}
	require.Equal(t, fmt.Sprintf("%v", unknownType), anyToBody(unknownType).StringVal())
}

func TestConvertNestedMapBody(t *testing.T) {
	unknownType := map[string]int{"0": 0, "1": 1}

	structuredBody := map[string]interface{}{
		"array":   []interface{}{0, 1},
		"map":     map[string]interface{}{"0": 0, "1": "one"},
		"unknown": unknownType,
	}

	result := anyToBody(structuredBody).MapVal()

	arrayAttVal, _ := result.Get("array")
	a := arrayAttVal.SliceVal()
	require.Equal(t, int64(0), a.At(0).IntVal())
	require.Equal(t, int64(1), a.At(1).IntVal())

	mapAttVal, _ := result.Get("map")
	m := mapAttVal.MapVal()
	v, _ := m.Get("0")
	require.Equal(t, int64(0), v.IntVal())
	v, _ = m.Get("1")
	require.Equal(t, "one", v.StringVal())

	unknownAttVal, _ := result.Get("unknown")
	require.Equal(t, fmt.Sprintf("%v", unknownType), unknownAttVal.StringVal())
}

func anyToBody(body interface{}) pdata.AttributeValue {
	entry := entry.New()
	entry.Body = body
	return convertAndDrill(entry).Body()
}

func convertAndDrill(entry *entry.Entry) pdata.LogRecord {
	return convert(entry, nil)
}

func TestConvertSeverity(t *testing.T) {
	cases := []struct {
		severity       entry.Severity
		expectedNumber pdata.SeverityNumber
		expectedText   string
	}{
		{entry.Default, pdata.SeverityNumberUNDEFINED, "Undefined"},
		{entry.Trace, pdata.SeverityNumberTRACE, "Trace"},
		{entry.Trace2, pdata.SeverityNumberTRACE2, "Trace2"},
		{entry.Trace3, pdata.SeverityNumberTRACE3, "Trace3"},
		{entry.Trace4, pdata.SeverityNumberTRACE4, "Trace4"},
		{entry.Debug, pdata.SeverityNumberDEBUG, "Debug"},
		{entry.Debug2, pdata.SeverityNumberDEBUG2, "Debug2"},
		{entry.Debug3, pdata.SeverityNumberDEBUG3, "Debug3"},
		{entry.Debug4, pdata.SeverityNumberDEBUG4, "Debug4"},
		{entry.Info, pdata.SeverityNumberINFO, "Info"},
		{entry.Info2, pdata.SeverityNumberINFO2, "Info2"},
		{entry.Info3, pdata.SeverityNumberINFO3, "Info3"},
		{entry.Info4, pdata.SeverityNumberINFO4, "Info4"},
		{entry.Warn, pdata.SeverityNumberWARN, "Warn"},
		{entry.Warn2, pdata.SeverityNumberWARN2, "Warn2"},
		{entry.Warn3, pdata.SeverityNumberWARN3, "Warn3"},
		{entry.Warn4, pdata.SeverityNumberWARN4, "Warn4"},
		{entry.Error, pdata.SeverityNumberERROR, "Error"},
		{entry.Error2, pdata.SeverityNumberERROR2, "Error2"},
		{entry.Error3, pdata.SeverityNumberERROR3, "Error3"},
		{entry.Error4, pdata.SeverityNumberERROR4, "Error4"},
		{entry.Fatal, pdata.SeverityNumberFATAL, "Fatal"},
		{entry.Fatal2, pdata.SeverityNumberFATAL2, "Fatal2"},
		{entry.Fatal3, pdata.SeverityNumberFATAL3, "Fatal3"},
		{entry.Fatal4, pdata.SeverityNumberFATAL4, "Fatal4"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.severity), func(t *testing.T) {
			entry := entry.New()
			entry.Severity = tc.severity
			log := convertAndDrill(entry)
			require.Equal(t, tc.expectedNumber, log.SeverityNumber())
			require.Equal(t, tc.expectedText, log.SeverityText())
		})
	}
}

func TestConvertTrace(t *testing.T) {
	record := convertAndDrill(&entry.Entry{
		TraceId: []byte{
			0x48, 0x01, 0x40, 0xf3, 0xd7, 0x70, 0xa5, 0xae, 0x32, 0xf0, 0xa2, 0x2b, 0x6a, 0x81, 0x2c, 0xff,
		},
		SpanId: []byte{
			0x32, 0xf0, 0xa2, 0x2b, 0x6a, 0x81, 0x2c, 0xff,
		},
		TraceFlags: []byte{
			0x01,
		}})

	require.Equal(t, pdata.NewTraceID(
		[16]byte{
			0x48, 0x01, 0x40, 0xf3, 0xd7, 0x70, 0xa5, 0xae, 0x32, 0xf0, 0xa2, 0x2b, 0x6a, 0x81, 0x2c, 0xff,
		}), record.TraceID())
	require.Equal(t, pdata.NewSpanID(
		[8]byte{
			0x32, 0xf0, 0xa2, 0x2b, 0x6a, 0x81, 0x2c, 0xff,
		}), record.SpanID())
	require.Equal(t, uint32(0x01), record.Flags())
}

func BenchmarkConverter(b *testing.B) {
	const (
		entryCount = 1_000_000
		hostsCount = 4
		batchCount = 200
	)

	var (
		workerCounts = []int{1, 2, 4, 6, 8}
		entries      = complexEntriesForNDifferentHosts(entryCount, hostsCount)
	)

	for _, wc := range workerCounts {
		b.Run(fmt.Sprintf("worker_count=%d", wc), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				converter := NewConverter(
					WithWorkerCount(wc),
				)
				converter.Start()
				defer converter.Stop()
				b.ResetTimer()

				go func() {
					for i := 0; i < entryCount; i += batchCount {
						if i+batchCount > entryCount {
							assert.NoError(b, converter.Batch(entries[i:entryCount]))
						} else {
							assert.NoError(b, converter.Batch(entries[i:i+batchCount]))
						}
					}
				}()

				var (
					timeoutTimer = time.NewTimer(10 * time.Second)
					ch           = converter.OutChannel()
				)
				defer timeoutTimer.Stop()

				var n int
			forLoop:
				for {
					if n == entryCount {
						break
					}

					select {
					case pLogs, ok := <-ch:
						if !ok {
							break forLoop
						}

						rLogs := pLogs.ResourceLogs()
						require.Equal(b, 1, rLogs.Len())

						rLog := rLogs.At(0)
						ills := rLog.InstrumentationLibraryLogs()
						require.Equal(b, 1, ills.Len())

						ill := ills.At(0)

						n += ill.Logs().Len()

					case <-timeoutTimer.C:
						break forLoop
					}
				}

				assert.Equal(b, entryCount, n,
					"didn't receive expected number of entries after conversion",
				)
			}
		})
	}
}

func BenchmarkGetResourceID(b *testing.B) {
	b.StopTimer()
	res := getResource()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		getResourceID(res)
	}
}

func BenchmarkGetResourceIDJSON(b *testing.B) {
	b.StopTimer()
	res := getResource()
	var underlyingBuffer [256]byte
	buf := bytes.NewBuffer(underlyingBuffer[:])
	enc := json.NewEncoder(buf)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err := enc.Encode(res)
		require.NoError(b, err)
	}
}

func getResource() map[string]string {
	return map[string]string{
		"file.name":        "filename.log",
		"file.directory":   "/some_directory",
		"host.name":        "localhost",
		"host.ip":          "192.168.1.12",
		"k8s.pod.name":     "test-pod-123zwe1",
		"k8s.node.name":    "aws-us-east-1.asfasf.aws.com",
		"k8s.container.id": "192end1yu823aocajsiocjnasd",
		"k8s.cluster.name": "my-cluster",
	}
}

type resourceIDOutput struct {
	name   string
	output uint64
}

type resourceIDOutputSlice []resourceIDOutput

func (o resourceIDOutputSlice) Len() int {
	return len(o)
}

func (x resourceIDOutputSlice) Less(i, j int) bool {
	return x[i].output < x[j].output
}

func (o resourceIDOutputSlice) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func TestGetResourceID(t *testing.T) {
	testCases := []struct {
		name  string
		input map[string]string
	}{
		{
			name:  "Typical Resource",
			input: getResource(),
		},
		{
			name: "Resource with non-utf bytes",
			input: map[string]string{
				"SomeKey":  "Value\xc0\xc1\xd4\xff\xfe",
				"\xff\xfe": "Ooops",
			},
		},
		{
			name: "Empty value/key",
			input: map[string]string{
				"SomeKey": "",
				"":        "Ooops",
			},
		},
		{
			name: "Empty value/key (reversed)",
			input: map[string]string{
				"":      "SomeKey",
				"Ooops": "",
			},
		},
		{
			name: "Ambiguous map 1",
			input: map[string]string{
				"AB": "CD",
				"EF": "G",
			},
		},
		{
			name: "Ambiguous map 2",
			input: map[string]string{
				"ABC": "DE",
				"F":   "G",
			},
		},
		{
			name: "Ambiguous map 3",
			input: map[string]string{
				"ABC": "DE\xfe",
				"F":   "G",
			},
		},
		{
			name: "Ambiguous map 4",
			input: map[string]string{
				"ABC":   "DE",
				"\xfeF": "G",
			},
		},
		{
			name:  "nil resource",
			input: nil,
		},
		{
			name: "Long resource value",
			input: map[string]string{
				"key": "This is a really long resource value; It's so long that the pre-allocated buffer size doesn't hold it.",
			},
		},
	}

	outputs := resourceIDOutputSlice{}
	for _, testCase := range testCases {
		outputs = append(outputs, resourceIDOutput{
			name:   testCase.name,
			output: getResourceID(testCase.input),
		})
	}

	// Ensure every output is unique
	sort.Sort(outputs)
	for i := 1; i < len(outputs); i++ {
		if outputs[i].output == outputs[i-1].output {
			t.Errorf("Test case %s and %s had the same output", outputs[i].name, outputs[i-1].name)
		}
	}
}

func TestGetResourceIDEmptyAndNilAreEqual(t *testing.T) {
	nilID := getResourceID(nil)
	emptyID := getResourceID(map[string]string{})
	require.Equal(t, nilID, emptyID)
}
