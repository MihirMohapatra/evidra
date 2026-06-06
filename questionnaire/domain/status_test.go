package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatus_Valid(t *testing.T) {
	tests := []struct {
		status Status
		valid  bool
	}{
		{StatusUploaded, true},
		{StatusQueued, true},
		{StatusParsing, true},
		{StatusParsed, true},
		{StatusFailed, true},
		{StatusAssigned, true},
		{"unknown", false},
		{"", false},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.valid, tt.status.Valid(), "status=%q", tt.status)
	}
}

func TestStatus_CanTransitionTo(t *testing.T) {
	tests := []struct {
		current Status
		next    Status
		allowed bool
	}{
		{StatusUploaded, StatusQueued, true},
		{StatusUploaded, StatusFailed, true},
		{StatusUploaded, StatusParsed, false},
		{StatusQueued, StatusParsing, true},
		{StatusQueued, StatusFailed, true},
		{StatusQueued, StatusParsed, false},
		{StatusParsing, StatusParsed, true},
		{StatusParsing, StatusFailed, true},
		{StatusParsing, StatusQueued, false},
		{StatusParsed, StatusAssigned, true},
		{StatusParsed, StatusFailed, false},
		{StatusFailed, StatusQueued, true},
		{StatusFailed, StatusParsed, false},
		{StatusAssigned, StatusParsed, true},
		{StatusAssigned, StatusFailed, false},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.allowed, tt.current.CanTransitionTo(tt.next), "%s -> %s", tt.current, tt.next)
	}
}
