package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOrganization(t *testing.T) {
	o := NewOrganization("Acme Corp", "acme-corp")
	require.NotNil(t, o)
	assert.Equal(t, "Acme Corp", o.Name)
	assert.Equal(t, "acme-corp", o.Slug)
	assert.NotZero(t, o.ID)
	assert.WithinDuration(t, time.Now(), o.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), o.UpdatedAt, time.Second)
}
