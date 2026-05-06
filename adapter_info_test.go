// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package gpucontext

import "testing"

func TestAdapterTypeString(t *testing.T) {
	tests := []struct {
		typ  AdapterType
		want string
	}{
		{AdapterTypeDiscrete, "Discrete"},
		{AdapterTypeIntegrated, "Integrated"},
		{AdapterTypeSoftware, "Software"},
		{AdapterTypeUnknown, "Unknown"},
		{AdapterType(99), "Unknown"},
	}
	for _, tt := range tests {
		if got := tt.typ.String(); got != tt.want {
			t.Errorf("AdapterType(%d).String() = %q, want %q", tt.typ, got, tt.want)
		}
	}
}

func TestAdapterInfoZeroValue(t *testing.T) {
	var info AdapterInfo
	if info.Name != "" {
		t.Errorf("zero AdapterInfo.Name = %q, want empty", info.Name)
	}
	if info.Type != AdapterTypeDiscrete {
		t.Errorf("zero AdapterInfo.Type = %d, want %d (Discrete)", info.Type, AdapterTypeDiscrete)
	}
}
