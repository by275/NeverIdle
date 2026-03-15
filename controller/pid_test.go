package controller

import "testing"

func TestNormalizeReferenceSignal(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  float64
	}{
		{
			name:  "zero ratio",
			input: 0,
			want:  0,
		},
		{
			name:  "valid ratio",
			input: 0.2,
			want:  20,
		},
		{
			name:  "full ratio",
			input: 1,
			want:  100,
		},
		{
			name:  "negative ratio falls back",
			input: -0.1,
			want:  15,
		},
		{
			name:  "oversized ratio falls back",
			input: 1.5,
			want:  15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeReferenceSignal(tt.input)
			if got != tt.want {
				t.Fatalf("normalizeReferenceSignal(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
