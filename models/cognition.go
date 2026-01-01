package models

type CognitionResponse struct {
	Input struct {
		RawText   string `json:"raw_text"`
		Language  string `json:"language"`
		Timestamp string `json:"timestamp"`
	} `json:"input"`

	Analysis struct {
		Ideas          []map[string]any `json:"ideas"`
		Assumptions    []map[string]any `json:"assumptions"`
		Emotions       []map[string]any `json:"emotions"`
		Conflicts      []map[string]any `json:"conflicts"`
		Contradictions []map[string]any `json:"contradictions"`
	} `json:"analysis"`

	Graph struct {
		Nodes []map[string]any `json:"nodes"`
		Edges []map[string]any `json:"edges"`
	} `json:"graph"`

	Reflection struct {
		Summary             string   `json:"summary"`
		CoreTension         string   `json:"core_tension"`
		UnresolvedQuestions []string `json:"unresolved_questions"`
	} `json:"reflection"`

	Confidence struct {
		AnalysisConfidence float64 `json:"analysis_confidence"`
	} `json:"confidence"`
}
