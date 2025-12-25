package ft

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDurationConversionPrecision(t *testing.T) {
	d := time.Second + 123*time.Nanosecond

	expectedMs := float64(d) / float64(time.Millisecond)
	expectedSeconds := float64(d) / float64(time.Second)

	assert.InDelta(t, expectedMs, durationToMillisecond(d), 1e-12)
	assert.InDelta(t, expectedSeconds, durationToSecond(d), 1e-12)
}
