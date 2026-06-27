package handler

import "testing"

func TestGetHello(t *testing.T) {

}

func TestValidateEstateInput(t *testing.T) {
	tests := []struct {
		name   string
		length int
		width  int
		wantErr bool
	}{
		{name: "valid", length: 10, width: 5, wantErr: false},
		{name: "zero length", length: 0, width: 5, wantErr: true},
		{name: "negative width", length: 10, width: -1, wantErr: true},
		{name: "too large", length: 50001, width: 1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEstateInput(tt.length, tt.width)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error for length=%d width=%d", tt.length, tt.width)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
