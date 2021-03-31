package snowflake

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnowFlakeGenerate(t *testing.T) {
	s := NewSnowFlake(999)
	assert.NotNil(t, s)
	var lastID int64
	for i := 0; i < 10000; i++ {
		id, err := s.GenID()
		assert.Nil(t, err)
		assert.Greater(t, id, lastID)
		lastID = id
	}
}

func BenchmarkSnowFlakeGenerate(b *testing.B) {
	s := NewSnowFlake(999)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.GenID()
	}
}
