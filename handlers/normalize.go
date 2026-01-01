package handlers

func normalizeCognition(raw map[string]any) map[string]any {
	// -----------------------------
	// Analysis
	// -----------------------------
	analysis := map[string]any{
		"ideas":          safeSlice(raw["ideas"]),
		"assumptions":    safeSlice(raw["assumptions"]),
		"contradictions": safeSlice(raw["contradictions"]),
	}

	// -----------------------------
	// Reflection
	// -----------------------------
	reflection := map[string]any{
		"summary":              safeString(raw["reflection"]),
		"core_tension":         safeString(raw["core_tension"]),
		"unresolved_questions": []any{},
	}

	// -----------------------------
	// Confidence (never null)
	// -----------------------------
	confidenceValue := 0.5 // honest uncertainty default
	if c, ok := raw["confidence"].(float64); ok {
		confidenceValue = c
	}

	confidence := map[string]any{
		"analysis_confidence": confidenceValue,
	}

	// -----------------------------
	// Graph (safe defaults)
	// -----------------------------
	graph := map[string]any{
		"nodes": []any{},
		"edges": []any{},
	}

	return map[string]any{
		"analysis":   analysis,
		"reflection": reflection,
		"confidence": confidence,
		"graph":      graph,
	}
}

// -----------------------------
// Helpers (defensive coding)
// -----------------------------
func safeSlice(v any) []any {
	if s, ok := v.([]any); ok {
		return s
	}
	return []any{}
}

func safeString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
