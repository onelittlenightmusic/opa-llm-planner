package types

// Action represents a single action in a plan.
type Action struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Status      string                 `json:"status"`
}

// Plan represents a generated execution plan.
type Plan struct {
	PlanID  string   `json:"plan_id"`
	GoalID  string   `json:"goal_id,omitempty"`
	Actions []Action `json:"actions"`
}
