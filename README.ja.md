# opa-llm-planner

OPA (Open Policy Agent) の Rego ルールと LLM を組み合わせ、goal/current JSON から実行アクション列を生成する Go CLI ツールです。

[English](README.md)

---

## 概要

goal（目標状態）と current（現在状態）の差分を Rego ポリシーで評価し、不足しているアクションを特定します。LLM（Anthropic Claude または OpenAI GPT-4o）を使って各アクションの説明やパラメータを自動補完することもできます。

```
goal.json + current.json + policies/*.rego  →  plan.json
```

## コマンド

| コマンド | 説明 |
|---------|------|
| `plan` | goal/current と Rego ポリシーからアクションプランを生成 |
| `consider` | 不足アクションに対応する Rego ルールを LLM で生成 |
| `explain` | OPA トレースでなぜアクションが missing なのかを表示 |

## インストール

```bash
git clone https://github.com/onelittlenightmusic/opa-llm-planner.git
cd opa-llm-planner
go build -o opa-llm-planner .
```

## 使い方

### `plan` — アクションプランの生成

```bash
# OPA のみ（LLM なし）
./opa-llm-planner plan \
  --goal examples/goal.json \
  --current examples/current.json \
  --policy ./policies \
  --out plan.json

# LLM でアクションを補完（Anthropic）
ANTHROPIC_API_KEY=xxx ./opa-llm-planner plan \
  --goal examples/goal.json \
  --current examples/current.json \
  --policy ./policies \
  --out plan.json \
  --llm --llm-provider anthropic

# LLM でアクションを補完（OpenAI）
OPENAI_API_KEY=xxx ./opa-llm-planner plan \
  --goal examples/goal.json \
  --current examples/current.json \
  --policy ./policies \
  --out plan.json \
  --llm --llm-provider openai
```

**フラグ一覧：**

| フラグ | デフォルト | 説明 |
|-------|-----------|------|
| `--goal` | `examples/goal.json` | goal JSON ファイルのパス |
| `--current` | `examples/current.json` | current JSON ファイルのパス |
| `--policy` | `policies` | Rego ポリシーファイルのディレクトリ |
| `--out` | 標準出力 | プラン JSON の出力先ファイル |
| `--llm` | false | LLM でアクションを補完する |
| `--llm-provider` | `$LLM_PROVIDER` | `anthropic` または `openai` |

**出力例（`plan.json`）：**

```json
{
  "plan_id": "ff734a1f-5865-4271-a8c7-c3c9338d58c1",
  "goal_id": "goal-001",
  "actions": [
    {
      "type": "reserve_hotel",
      "description": "東京旅行のためにホテルを予約する",
      "parameters": { "destination": "Tokyo" },
      "status": "pending"
    },
    {
      "type": "reserve_dinner",
      "description": "東京でのディナーを予約する",
      "parameters": { "destination": "Tokyo" },
      "status": "pending"
    }
  ]
}
```

### `consider` — Rego ルールの生成

既存ポリシーでカバーされていないアクションがある場合、`consider` コマンドで LLM に新しい Rego ルールを生成させます。

```bash
# ドライラン（ファイルに書かず画面に表示）
ANTHROPIC_API_KEY=xxx ./opa-llm-planner consider \
  --goal examples/goal.json \
  --current examples/current.json \
  --policy ./policies \
  --dry-run \
  --llm-provider anthropic

# 新しいファイルに書き出す
ANTHROPIC_API_KEY=xxx ./opa-llm-planner consider \
  --goal examples/goal.json \
  --current examples/current.json \
  --policy ./policies \
  --out new_rules.rego \
  --llm-provider anthropic

# 既存ポリシーファイルに追記する
ANTHROPIC_API_KEY=xxx ./opa-llm-planner consider \
  --goal examples/goal.json \
  --current examples/current.json \
  --policy ./policies \
  --out policies/planner.rego \
  --append \
  --llm-provider anthropic
```

**フラグ一覧：**

| フラグ | デフォルト | 説明 |
|-------|-----------|------|
| `--goal` | `examples/goal.json` | goal JSON ファイルのパス |
| `--current` | `examples/current.json` | current JSON ファイルのパス |
| `--policy` | `policies` | Rego ポリシーファイルのディレクトリ |
| `--out` | 標準出力 | 生成 Rego ルールの出力先ファイル |
| `--append` | false | 既存ファイルに追記する |
| `--dry-run` | false | ファイルに書かず画面に表示のみ |
| `--llm-provider` | `$LLM_PROVIDER` | `anthropic` または `openai` |

### `explain` — なぜ missing なのかを OPA トレースで表示

OPA のトレース機能を使って、どのルールが評価されたか・なぜアクションが missing になったかを詳細に表示します。

```bash
./opa-llm-planner explain \
  --goal examples/goal.json \
  --current examples/current.json \
  --policy ./policies

# ソースファイルと行番号も表示する
./opa-llm-planner explain \
  --goal examples/goal.json \
  --current examples/current.json \
  --policy ./policies \
  --location
```

**フラグ一覧：**

| フラグ | デフォルト | 説明 |
|-------|-----------|------|
| `--goal` | `examples/goal.json` | goal JSON ファイルのパス |
| `--current` | `examples/current.json` | current JSON ファイルのパス |
| `--policy` | `policies` | Rego ポリシーファイルのディレクトリ |
| `--location` | false | ソースファイル・行番号をトレースに含める |

**出力例：**

```
=== OPA Trace: why actions are missing ===

Enter data.planner.missing = _
| Eval data.planner.missing = _
| Index data.planner.missing (matched 2 rules)
| Enter data.planner.missing
| | Eval input.goal.trip.require_hotel
| | Eval not input.current.hotel_reserved
| | | Fail input.current.hotel_reserved    ← hotel_reserved が false → missing!
| | Eval action = "reserve_hotel"
| | Exit data.planner.missing
...

=== Result ===
Missing actions: [reserve_hotel, reserve_dinner]
```

## Rego ポリシーの書き方

`--policy` ディレクトリ内の `.rego` ファイルはすべて自動的に読み込まれます。

`package planner` の `missing` セットにアクション名を追加するルールを記述します：

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

Rego 内で参照できる `input` の構造：

```json
{
  "goal":    { /* goal.json の内容 */ },
  "current": { /* current.json の内容 */ }
}
```

## 環境変数

| 変数 | 説明 |
|------|------|
| `ANTHROPIC_API_KEY` | Anthropic API キー |
| `OPENAI_API_KEY` | OpenAI API キー |
| `LLM_PROVIDER` | デフォルトの LLM プロバイダ（`anthropic` または `openai`） |

## サンプルファイル

**`examples/goal.json`：**
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

**`examples/current.json`：**
```json
{
  "hotel_reserved": false,
  "dinner_reserved": false
}
```

## アーキテクチャ

```
opa-llm-planner/
├── main.go
├── cmd/
│   ├── root.go       # Cobra ルートコマンド
│   ├── plan.go       # plan コマンド
│   └── consider.go   # consider コマンド
├── internal/
│   ├── types/        # Action, Plan 型定義
│   ├── opa/          # OPA SDK ラッパー
│   ├── llm/          # LLMClient インターフェース + Anthropic/OpenAI 実装
│   └── planner/      # plan・consider のビジネスロジック
├── policies/         # Rego ポリシーファイル
└── examples/         # サンプル入力ファイル
```
