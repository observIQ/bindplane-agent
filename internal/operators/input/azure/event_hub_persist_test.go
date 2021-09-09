package azure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPersistenceKey(t *testing.T) {
	type TestKey struct {
		namespace     string
		name          string
		consumerGroup string
		partitionID   string
	}

	cases := []struct {
		name     string
		input    TestKey
		expected string
	}{
		{
			"basic",
			TestKey{
				namespace:     "stanza",
				name:          "devel",
				consumerGroup: "$Default",
				partitionID:   "0",
			},
			"stanza-devel-$Default-0",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := Persister{}
			out := p.persistenceKey(tc.input.namespace, tc.input.name, tc.input.consumerGroup, tc.input.partitionID)
			require.Equal(t, tc.expected, out)
		})
	}
}
