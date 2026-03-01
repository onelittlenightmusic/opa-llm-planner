# opa-llm-planner

A Go CLI tool that combines [OPA (Open Policy Agent)](https://www.openpolicyagent.org/) Rego policies with LLM to automatically generate execution action plans from goal and current state JSON files.

[日本語](#日本語)

---

## Overview

`opa-llm-planner` evaluates what actions are **missing** between your goal state and current state using OPA Rego rules, then optionally uses an LLM (Anthropic Claude or OpenAI GPT-4o) to enrich those actions with descriptions and parameters.

```
goal.json + current.json + policies/*.rego  →  plan.json
```

## Commands

| Command | Description |
|---------|-------------|
| `plan` | Generate an action plan from goal/current state using OPA policies |
| `consider` | Generate new Rego rules for missing actions using LLM |

## Installation

```bash
git clone https://github.com/onelittlenightmusic/opa-llm-planner.git
cd opa-llm-planner
go build -o opa-llm-planner .
```

## Usage

### `plan` — Generate an action plan

```bash
# OPA only (no LLM)
./opa-llm-planner plan \
  --goal examples/goal.json \
  --current examples/current.json \
  --policy ./policies \
  --out plan.json

# With LLM enrichment (Anthropic)
ANTHROPIC_API_KEY=xxx ./opa-llm-planner plan \
  --goal examples/goal.json \
  --current examples/current.json \
  --policy ./policies \
  --out plan.json \
  --llm --llm-provider anthropic

# With LLM enrichment (OpenAI)
OPENAI_API_KEY=xxx ./opa-llm-planner plan \
  --goal examples/goal.json \
  --current examples/current.json \
  --policy ./policies \
  --out plan.json \
  --llm --llm-provider openai
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--goal` | `examples/goal.json` | Path to goal JSON file |
| `--current` | `examples/current.json` | Path to current state JSON file |
| `--policy` | `policies` | Directory containing Rego policy files |
| `--out` | stdout | Output file for plan JSON |
| `--llm` | false | Enrich actions using LLM |
| `--llm-provider` | `$LLM_PROVIDER` | `anthropic` or `openai` |

**Example output (`plan.json`):**

```json
{
  "plan_id": "ff734a1f-5865-4271-a8c7-c3c9338d58c1",
  "goal_id": "goal-001",
  "actions": [
    {
      "type": "reserve_hotel",
      "description": "Reserve a hotel room in Tokyo for the trip",
      "parameters": { "destination": "Tokyo" },
      "status": "pending"
    },
    {
      "type": "reserve_dinner",
      "description": "Make a dinner reservation in Tokyo",
      "parameters": { "destination": "Tokyo" },
      "status": "pending"
    }
  ]
}
```

### `consider` — Generate new Rego rules

When actions are missing and no policy rule exists for them, use `consider` to generate new Rego rules via LLM.

```bash
# Dry-run: print generated rules without writing
ANTHROPIC_API_KEY=xxx ./opa-llm-planner consider \
  --goal examples/goal.json \
  --current examples/current.json \
  --policy ./policies \
  --dry-run \
  --llm-provider anthropic

# Write to a new file
ANTHROPIC_API_KEY=xxx ./opa-llm-planner consider \
  --goal examples/goal.json \
  --current examples/current.json \
  --policy ./policies \
  --out new_rules.rego \
  --llm-provider anthropic

# Append to an existing policy file
ANTHROPIC_API_KEY=xxx ./opa-llm-planner consider \
  --goal examples/goal.json \
  --current examples/current.json \
  --policy ./policies \
  --out policies/planner.rego \
  --append \
  --llm-provider anthropic
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--goal` | `examples/goal.json` | Path to goal JSON file |
| `--current` | `examples/current.json` | Path to current state JSON file |
| `--policy` | `policies` | Directory containing Rego policy files |
| `--out` | stdout | Output file for generated Rego rules |
| `--append` | false | Append to output file instead of overwriting |
| `--dry-run` | false | Print rules without writing to file |
| `--llm-provider` | `$LLM_PROVIDER` | `anthropic` or `openai` |

## Writing Rego Policies

Policies live in the `--policy` directory. Each `.rego` file is loaded automatically.

Rules must be in `package planner` and contribute to the `missing` set:

```rego
package planner

missing[action] {
  input.goal.trip.require_hotel
  not input.current.hotel_reserved
  action := "reserve_hotel"
}

missing[action] {
  input.goal.trip.require_dinner
  not input.current.dinner_reserved
  action := "reserve_dinner"
}
```

The input available in Rego is:
```json
{
  "goal":    { ...contents of goal.json... },
  "current": { ...contents of current.json... }
}
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `ANTHROPIC_API_KEY` | Anthropic API key |
| `OPENAI_API_KEY` | OpenAI API key |
| `LLM_PROVIDER` | Default LLM provider (`anthropic` or `openai`) |

## Examples

**`examples/goal.json`:**
```json
{
  "id": "goal-001",
  "trip": {
    "destination": "Tokyo",
    "require_hotel": true,
    "require_dinner": true
  }
}
```

**`examples/current.json`:**
```json
{
  "hotel_reserved": false,
  "dinner_reserved": false
}
```

## Architecture

```
opa-llm-planner/
├── main.go
├── cmd/
│   ├── root.go       # Cobra root command
│   ├── plan.go       # plan command
│   └── consider.go   # consider command
├── internal/
│   ├── types/        # Action, Plan types
│   ├── opa/          # OPA SDK wrapper
│   ├── llm/          # LLMClient interface + Anthropic/OpenAI implementations
│   └── planner/      # plan and consider business logic
├── policies/         # Rego policy files
└── examples/         # Example input files
```

---

## 日本語

`opa-llm-planner` は OPA (Open Policy Agent) の Rego ルールと LLM を組み合わせ、goal/current JSON から実行アクション列を生成する Go CLI ツールです。

### 概要

goal（目標状態）と current（現在状態）の差分を Rego ポリシーで評価し、不足しているアクションを特定します。LLM（Anthropic Claude または OpenAI GPT-4o）を使って各アクションの説明やパラメータを自動補完することもできます。

### コマンド

| コマンド | 説明 |
|---------|------|
| `plan` | goal/current と Rego ポリシーからアクションプランを生成 |
| `consider` | 不足アクションに対応する Rego ルールを LLM で生成 |

### インストール

```bash
git clone https://github.com/onelittlenightmusic/opa-llm-planner.git
cd opa-llm-planner
go build -o opa-llm-planner .
```

### 基本的な使い方

```bash
# plan コマンド（OPAのみ）
./opa-llm-planner plan \
  --goal examples/goal.json \
  --current examples/current.json \
  --policy ./policies

# plan コマンド（LLM でアクションを補完）
ANTHROPIC_API_KEY=xxx ./opa-llm-planner plan \
  --goal examples/goal.json \
  --current examples/current.json \
  --policy ./policies \
  --llm --llm-provider anthropic

# consider コマンド（不足ルールを生成・ドライラン）
ANTHROPIC_API_KEY=xxx ./opa-llm-planner consider \
  --goal examples/goal.json \
  --current examples/current.json \
  --policy ./policies \
  --dry-run \
  --llm-provider anthropic
```

### Rego ポリシーの書き方

`package planner` の `missing` セットにアクション名を追加するルールを記述します。

```rego
package planner

missing[action] {
  input.goal.trip.require_hotel
  not input.current.hotel_reserved
  action := "reserve_hotel"
}
```

Rego 内で参照できる `input` は以下の構造です：

```json
{
  "goal":    { /* goal.json の内容 */ },
  "current": { /* current.json の内容 */ }
}
```

### 環境変数

| 変数 | 説明 |
|------|------|
| `ANTHROPIC_API_KEY` | Anthropic API キー |
| `OPENAI_API_KEY` | OpenAI API キー |
| `LLM_PROVIDER` | デフォルトの LLM プロバイダ（`anthropic` または `openai`） |
