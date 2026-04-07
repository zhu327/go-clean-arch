#!/bin/bash
# Session Start Hook - Load learned skills on new session
#
# This script runs when a new Cursor session starts.
# It reads learned skills and injects them as additional context.

set -e

# Read input JSON from stdin (required by Cursor hooks)
INPUT=$(cat)

# Get the workspace root from input JSON
WORKSPACE_ROOT=$(echo "$INPUT" | grep -o '"workspace_roots":\s*\[\s*"[^"]*"' | head -1 | sed 's/.*"\([^"]*\)"$/\1/')

# Fallback: use script directory to find project root
if [ -z "$WORKSPACE_ROOT" ]; then
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    WORKSPACE_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
fi

# Directory containing learned skills
SKILLS_DIR="$WORKSPACE_ROOT/.cursor/skills/learned"

# Build additional context
CONTEXT=""

# Check for learned skills
if [ -d "$SKILLS_DIR" ]; then
    # Find .md files (excluding .gitkeep)
    shopt -s nullglob
    SKILL_FILES=("$SKILLS_DIR"/*.md)
    SKILL_COUNT=${#SKILL_FILES[@]}
    
    if [ "$SKILL_COUNT" -gt 0 ]; then
        CONTEXT="## Learned Skills Available\n\n"
        CONTEXT+="以下是从历史会话中学习到的经验，遇到相关问题时请使用 Read 工具查看详情：\n\n"
        CONTEXT+="| Skill | Path |\n"
        CONTEXT+="|-------|------|\n"
        
        # List each skill file
        for skill in "${SKILL_FILES[@]}"; do
            if [ -f "$skill" ]; then
                # Get filename without extension
                name=$(basename "$skill" .md)
                # Use relative path for cleaner output
                rel_path=".cursor/skills/learned/$name.md"
                CONTEXT+="| $name | $rel_path |\n"
            fi
        done
        
        CONTEXT+="\n使用 /learn 命令可以从当前会话提取新的经验。\n"
    fi
fi

# Output JSON response (must output valid JSON to stdout)
if [ -n "$CONTEXT" ]; then
    # Use jq if available for proper JSON encoding, otherwise use simple approach
    if command -v jq &> /dev/null; then
        echo "{\"continue\": true, \"additional_context\": $(echo -e "$CONTEXT" | jq -Rs .)}"
    else
        # Escape special characters for JSON
        ESCAPED_CONTEXT=$(echo -e "$CONTEXT" | sed 's/\\/\\\\/g; s/"/\\"/g' | tr '\n' ' ' | sed 's/  */ /g')
        echo "{\"continue\": true, \"additional_context\": \"$ESCAPED_CONTEXT\"}"
    fi
else
    echo '{"continue": true}'
fi
