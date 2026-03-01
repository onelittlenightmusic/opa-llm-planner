package opa

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/open-policy-agent/opa/rego"
)

// Evaluator wraps the OPA SDK for Rego policy evaluation.
type Evaluator struct {
	policyDir string
}

// NewEvaluator creates an Evaluator that loads policies from the given directory.
func NewEvaluator(policyDir string) *Evaluator {
	return &Evaluator{policyDir: policyDir}
}

// EvaluateMissing evaluates data.planner.missing with {goal, current} as input
// and returns the list of missing action names.
func (e *Evaluator) EvaluateMissing(ctx context.Context, goal, current map[string]interface{}) ([]string, error) {
	modules, err := e.loadModules()
	if err != nil {
		return nil, fmt.Errorf("loading policies: %w", err)
	}

	options := []func(*rego.Rego){
		rego.Query("data.planner.missing"),
		rego.Input(map[string]interface{}{
			"goal":    goal,
			"current": current,
		}),
	}
	for name, src := range modules {
		options = append(options, rego.Module(name, src))
	}

	r := rego.New(options...)
	rs, err := r.Eval(ctx)
	if err != nil {
		return nil, fmt.Errorf("evaluating policy: %w", err)
	}

	if len(rs) == 0 || len(rs[0].Expressions) == 0 {
		return nil, nil
	}

	raw := rs[0].Expressions[0].Value
	set, ok := raw.([]interface{})
	if !ok {
		// OPA returns sets as map[string]interface{} when using Eval
		setMap, ok2 := raw.(map[string]interface{})
		if !ok2 {
			return nil, fmt.Errorf("unexpected result type: %T", raw)
		}
		var actions []string
		for k := range setMap {
			actions = append(actions, k)
		}
		return actions, nil
	}

	var actions []string
	for _, v := range set {
		s, ok := v.(string)
		if !ok {
			continue
		}
		actions = append(actions, s)
	}
	return actions, nil
}

// ValidateRego checks that a Rego source string is syntactically valid.
func ValidateRego(ctx context.Context, name, src string) error {
	r := rego.New(
		rego.Query("data"),
		rego.Module(name, src),
	)
	_, err := r.PrepareForEval(ctx)
	return err
}

func (e *Evaluator) loadModules() (map[string]string, error) {
	modules := make(map[string]string)
	entries, err := os.ReadDir(e.policyDir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".rego" {
			continue
		}
		path := filepath.Join(e.policyDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		modules[entry.Name()] = string(data)
	}
	return modules, nil
}
