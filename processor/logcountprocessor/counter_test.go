package logcountprocessor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogCounter(t *testing.T) {
	counter := NewLogCounter()
	resourceMap1 := map[string]interface{}{"resource1": "value1"}
	resourceMap2 := map[string]interface{}{"resource2": "value2"}
	attrMap1 := map[string]interface{}{"attr1": "value1"}
	attrMap2 := map[string]interface{}{"attr2": "value2"}

	for i := 0; i < 10; i++ {
		counter.Add(resourceMap1, attrMap1)
	}

	for i := 0; i < 5; i++ {
		counter.Add(resourceMap1, attrMap2)
	}

	counter.Add(resourceMap2, attrMap1)

	resourceKey1 := getDimensionKey(resourceMap1)
	resourceKey2 := getDimensionKey(resourceMap2)
	attrKey1 := getDimensionKey(attrMap1)
	attrKey2 := getDimensionKey(attrMap2)

	require.Equal(t, 10, counter.resources[resourceKey1].attributes[attrKey1].count)
	require.Equal(t, 5, counter.resources[resourceKey1].attributes[attrKey2].count)
	require.Equal(t, 1, counter.resources[resourceKey2].attributes[attrKey1].count)

	counter.Reset()
	require.Len(t, counter.resources, 0)
}
