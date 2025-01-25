package cccorrelation

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_GetCorrelationID_CorrelationKeyFound_ShouldReturnCorrelationID(t *testing.T) {
	// Arrange
	correlationID := uuid.New().String()
	expected := correlationID

	ctx := context.WithValue(context.Background(), Key, correlationID)

	// Act
	got := GetCorrelationID(ctx)
	assert.Equal(t, expected, got)

}

func Test_GetCorrelationID_CorrelationKeyNotFound_ShouldReturnEmptyString(t *testing.T) {
	expected := ""

	// Arrange
	ctx := context.Background()

	// Act
	got := GetCorrelationID(ctx)
	assert.Equal(t, expected, got)
}
