package core_test

import (
	"github.com/ghtak/golang.grpc.base/internal/core"
	"github.com/magiconair/properties/assert"
	"testing"
)

type Item struct {
	Value int
}

func TestRoundRobin(t *testing.T) {
	rr := core.NewRoundRobin(
		&Item{Value: 0},
		&Item{Value: 1},
		&Item{Value: 2},
		&Item{Value: 3},
		&Item{Value: 4},
	)
	assert.Equal(t, 0, rr.Next().Value, "")
	assert.Equal(t, 1, rr.Next().Value, "")
	assert.Equal(t, 2, rr.Next().Value, "")
	assert.Equal(t, 3, rr.Next().Value, "")
	assert.Equal(t, 4, rr.Next().Value, "")
	assert.Equal(t, 0, rr.Next().Value, "")
	assert.Equal(t, 1, rr.Next().Value, "")
	assert.Equal(t, 2, rr.Next().Value, "")
	assert.Equal(t, 3, rr.Next().Value, "")
	assert.Equal(t, 4, rr.Next().Value, "")
}
