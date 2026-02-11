package attribution

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculate(t *testing.T) {
	attr := Calculate(146, 54)
	assert.Equal(t, 73.0, attr.AgentPercent)
	assert.Equal(t, 146, attr.AgentLines)
	assert.Equal(t, 200, attr.TotalLines)
}

func TestCalculateAllAgent(t *testing.T) {
	attr := Calculate(100, 0)
	assert.Equal(t, 100.0, attr.AgentPercent)
	assert.Equal(t, 100, attr.TotalLines)
}

func TestCalculateAllHuman(t *testing.T) {
	attr := Calculate(0, 100)
	assert.Equal(t, 0.0, attr.AgentPercent)
	assert.Equal(t, 100, attr.TotalLines)
}

func TestCalculateZero(t *testing.T) {
	attr := Calculate(0, 0)
	assert.Equal(t, 0.0, attr.AgentPercent)
	assert.Equal(t, 0, attr.TotalLines)
}
