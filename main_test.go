package main

import (
	"strings"
	"testing"
	"time"
)

func TestValidateFlagsRejectsInvalidValues(t *testing.T) {
	originalCPUPercent := *FlagCPUPercent
	originalCPU := *FlagCPU
	originalMemory := *FlagMemory
	originalNetwork := *FlagNetwork
	originalNetworkConnectionCount := *FlagNetworkConnectionCount

	t.Cleanup(func() {
		*FlagCPUPercent = originalCPUPercent
		*FlagCPU = originalCPU
		*FlagMemory = originalMemory
		*FlagNetwork = originalNetwork
		*FlagNetworkConnectionCount = originalNetworkConnectionCount
	})

	tests := []struct {
		name      string
		setup     func()
		wantError string
	}{
		{
			name: "negative memory",
			setup: func() {
				*FlagMemory = -1
			},
			wantError: "-m must be 0 or greater",
		},
		{
			name: "negative cpu interval",
			setup: func() {
				*FlagCPU = -time.Second
			},
			wantError: "-c must be 0 or greater",
		},
		{
			name: "negative network interval",
			setup: func() {
				*FlagNetwork = -time.Second
			},
			wantError: "-n must be 0 or greater",
		},
		{
			name: "invalid network connection count",
			setup: func() {
				*FlagNetworkConnectionCount = 0
			},
			wantError: "-t must be greater than 0",
		},
		{
			name: "invalid cpu ratio",
			setup: func() {
				*FlagCPUPercent = 1.5
			},
			wantError: "-cp must be a ratio between 0 and 1",
		},
		{
			name: "cpu modes are mutually exclusive",
			setup: func() {
				*FlagCPU = time.Second
				*FlagCPUPercent = 0.2
			},
			wantError: "-c and -cp cannot be used together",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetValidationFlags()
			tt.setup()

			err := validateFlags()
			if err == nil {
				t.Fatalf("validateFlags() error = nil, want %q", tt.wantError)
			}
			if !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("validateFlags() error = %q, want substring %q", err.Error(), tt.wantError)
			}
		})
	}
}

func TestValidateFlagsAcceptsValidValues(t *testing.T) {
	originalCPUPercent := *FlagCPUPercent
	originalCPU := *FlagCPU
	originalMemory := *FlagMemory
	originalNetwork := *FlagNetwork
	originalNetworkConnectionCount := *FlagNetworkConnectionCount

	t.Cleanup(func() {
		*FlagCPUPercent = originalCPUPercent
		*FlagCPU = originalCPU
		*FlagMemory = originalMemory
		*FlagNetwork = originalNetwork
		*FlagNetworkConnectionCount = originalNetworkConnectionCount
	})

	resetValidationFlags()
	*FlagCPU = time.Second
	*FlagNetwork = time.Second
	*FlagNetworkConnectionCount = 4

	if err := validateFlags(); err != nil {
		t.Fatalf("validateFlags() error = %v, want nil", err)
	}
}

func resetValidationFlags() {
	*FlagCPUPercent = 0
	*FlagCPU = 0
	*FlagMemory = 0
	*FlagNetwork = 0
	*FlagNetworkConnectionCount = 10
}
