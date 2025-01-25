package cccopy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Test struct {
	Field1 int
	Field2 bool
	Field3 *int64
}

func TestDeepCopy_GivenSrcTestAndDestTestType_WhenSameStructFields_ThenAllFieldDeepCopy(t *testing.T) {

	field3 := int64(19)

	src := Test{
		Field1: 1,
		Field2: true,
		Field3: &field3,
	}

	var dest Test
	err := DeepCopy(src, &dest)
	assert.NoError(t, err)
	assert.Equal(t, src, dest)
	assert.Equal(t, *src.Field3, int64(19))
}
