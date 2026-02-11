package attribution

import "github.com/yibudak/open-entire/pkg/types"

// Calculate computes the attribution for a set of changed files.
func Calculate(agentLines, humanLines int) types.Attribution {
	total := agentLines + humanLines
	if total == 0 {
		return types.Attribution{}
	}

	return types.Attribution{
		AgentPercent: float64(agentLines) / float64(total) * 100,
		AgentLines:   agentLines,
		TotalLines:   total,
	}
}
