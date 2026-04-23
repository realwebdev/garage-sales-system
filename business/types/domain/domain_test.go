package domain_test

import (
	"testing"

	"github.com/realwebdev/garage-sales-system/business/types/domain"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid User", "USER", false},
		{"Valid Product", "PRODUCT", false},
		{"Invalid String", "GHOST", true},
		{"EMPTY String", "", true},
		{"Too Long", "VERY_LONG_INVALID_DOMAIN_NAME", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
