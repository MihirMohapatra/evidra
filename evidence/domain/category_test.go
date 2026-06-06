package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategory_Valid(t *testing.T) {
	tests := []struct {
		category Category
		valid    bool
	}{
		{CategoryPolicy, true},
		{CategoryAnswer, true},
		{CategoryClaim, true},
		{CategoryCertification, true},
		{CategoryArchitecture, true},
		{"unknown", false},
		{"", false},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.valid, tt.category.Valid(), "category=%q", tt.category)
	}
}

func TestValidCategories(t *testing.T) {
	cats := ValidCategories()
	assert.ElementsMatch(t, []Category{CategoryPolicy, CategoryAnswer, CategoryClaim, CategoryCertification, CategoryArchitecture}, cats)
}

func TestParseCategory(t *testing.T) {
	cat, ok := ParseCategory("policy")
	assert.True(t, ok)
	assert.Equal(t, CategoryPolicy, cat)

	_, ok = ParseCategory("unknown")
	assert.False(t, ok)

	_, ok = ParseCategory("")
	assert.False(t, ok)
}
