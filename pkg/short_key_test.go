package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeShortKey(t *testing.T) {
	tests := []struct {
		name    string
		randKey string
		id      uint64
		want    string
	}{
		{name: "ID 123", randKey: "ab", id: 123, want: "aB9b"},
		{name: "ID 0", randKey: "xy", id: 0, want: "xAy"},
		{name: "ID 61 (Z)", randKey: "qr", id: 61, want: "q9r"},
		{name: "ID 62 (10)", randKey: "ef", id: 62, want: "eBAf"},
		{name: "Large ID 3843 (ZZ)", randKey: "gh", id: 3843, want: "g99h"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeShortKey(tt.randKey, tt.id)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSplitShortKey(t *testing.T) {
	tests := []struct {
		name        string
		shortKey    string
		wantID      uint64
		wantRandKey string
	}{
		{name: "Valid key: a1Zb (ID 123)", shortKey: "aB9b", wantID: 123, wantRandKey: "ab"},
		{name: "Valid key: x0y (ID 0)", shortKey: "xAy", wantID: 0, wantRandKey: "xy"},
		{name: "Valid key: qZr (ID 61)", shortKey: "q9r", wantID: 61, wantRandKey: "qr"},
		{name: "Valid key: e10f (ID 62)", shortKey: "eBAf", wantID: 62, wantRandKey: "ef"},
		{name: "Valid key: gZZh (ID 3843)", shortKey: "g99h", wantID: 3843, wantRandKey: "gh"},
		{name: "Shortest valid key: a0b (ID 0)", shortKey: "aAb", wantID: 0, wantRandKey: "ab"},
		{
			name:        "Key with empty middle part (e.g., direct merge of randKey if base62 of ID is empty - though 0 is '0')",
			shortKey:    "ab",
			wantID:      0,
			wantRandKey: "ab",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotRandKey := SplitShortKey(tt.shortKey)
			assert.Equal(t, tt.wantID, gotID)
			assert.Equal(t, tt.wantRandKey, gotRandKey)
		})
	}
}

func TestRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		id      uint64
		randKey string
	}{
		{name: "Small ID", id: 42, randKey: "ab"},
		{name: "Large ID", id: 9999999, randKey: "xy"},
		{name: "ID 0", id: 0, randKey: "cd"},
		{name: "Special boundary value ID 61", id: 61, randKey: "ef"},
		{name: "Special boundary value ID 62", id: 62, randKey: "gh"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortKey := MergeShortKey(tt.randKey, tt.id)
			gotID, gotRandKey := SplitShortKey(shortKey)
			assert.Equal(t, tt.id, gotID, "ID should match")
			assert.Equal(t, tt.randKey, gotRandKey, "Random key should match")
		})
	}
}
