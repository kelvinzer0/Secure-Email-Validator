package checker

import (
	"testing"

	"github.com/kelvinzer0/secure-email-validator/internal/config"
)

func TestEmailChecker_isValidEmailFormat(t *testing.T) {
	ec := NewEmailChecker(config.DefaultConfig())

	tests := []struct {
		email string
		want  bool
	}{
		{"test@example.com", true},
		{"user.name@domain.co.uk", true},
		{"user+tag@example.org", true},
		{"invalid-email", false},
		{"@example.com", false},
		{"test@", false},
		{"", false},
		{"test@.com", false},
		{"test@com.", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			if got := ec.isValidEmailFormat(tt.email); got != tt.want {
				t.Errorf("isValidEmailFormat(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestEmailChecker_normalizeEmail(t *testing.T) {
	ec := NewEmailChecker(config.DefaultConfig())

	tests := []struct {
		input string
		want  string
	}{
		{"test@gmail.com", "test@gmail.com"},
		{"te.st@gmail.com", "test@gmail.com"},
		{"test+tag@gmail.com", "test@gmail.com"},
		{"te.st+tag@gmail.com", "test@gmail.com"},
		{"test@googlemail.com", "test@gmail.com"},
		{"test@example.com", "test@example.com"},
		{"Test@Example.Com", "test@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := ec.normalizeEmail(tt.input); got != tt.want {
				t.Errorf("normalizeEmail(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestEmailChecker_extractDomain(t *testing.T) {
	ec := NewEmailChecker(config.DefaultConfig())

	tests := []struct {
		email string
		want  string
	}{
		{"test@example.com", "example.com"},
		{"user@DOMAIN.COM", "domain.com"},
		{"invalid-email", ""},
		{"@example.com", "example.com"},
		{"test@", ""},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			if got := ec.extractDomain(tt.email); got != tt.want {
				t.Errorf("extractDomain(%q) = %q, want %q", tt.email, got, tt.want)
			}
		})
	}
}
