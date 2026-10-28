#!/bin/bash
# Claude Code Status Line
# Displays: [Model (effort)] [Progress Bar] Context% | Git Branch | Lines Changed | Session Cost (if >$0)

# Read JSON input from Claude Code
input=$(cat)

# Extract values using jq
MODEL=$(echo "$input" | jq -r '.model.display_name // "Unknown"')
EFFORT=$(echo "$input" | jq -r '.effort.level // ""')
CONTEXT_SIZE=$(echo "$input" | jq -r '.context_window.context_window_size // 200000')
USAGE=$(echo "$input" | jq '.context_window.current_usage // null')
COST=$(echo "$input" | jq -r '.cost.total_cost_usd // 0')
LINES_ADDED=$(echo "$input" | jq -r '.cost.total_lines_added // 0')
LINES_REMOVED=$(echo "$input" | jq -r '.cost.total_lines_removed // 0')

# Calculate context percentage
if [ "$USAGE" != "null" ]; then
    CURRENT_TOKENS=$(echo "$USAGE" | jq '.input_tokens + .cache_creation_input_tokens + .cache_read_input_tokens')
    PERCENT_USED=$((CURRENT_TOKENS * 100 / CONTEXT_SIZE))
else
    PERCENT_USED=0
fi

# Build progress bar (10 chars)
FILLED=$((PERCENT_USED / 10))
EMPTY=$((10 - FILLED))
BAR=$(printf '█%.0s' $(seq 1 $FILLED 2>/dev/null) || true)$(printf '░%.0s' $(seq 1 $EMPTY 2>/dev/null) || true)

# Get git branch
GIT_BRANCH=""
if git rev-parse --git-dir > /dev/null 2>&1; then
    BRANCH=$(git branch --show-current 2>/dev/null)
    [ -n "$BRANCH" ] && GIT_BRANCH=" | 🌿 $BRANCH"
fi

# Build model label with optional effort level
if [ -n "$EFFORT" ]; then
    MODEL_LABEL="$MODEL ($EFFORT)"
else
    MODEL_LABEL="$MODEL"
fi

# Show cost when non-zero (API billing), otherwise show 5h rate limit (subscription)
USAGE_SECTION=""
if awk -v c="$COST" 'BEGIN { exit !(c + 0 > 0.001) }'; then
    COST_FMT=$(printf '%.2f' "$COST")
    USAGE_SECTION=" | \$$COST_FMT"
else
    RATE_5H=$(echo "$input" | jq -r '.rate_limits."5h".used_percentage // ""' 2>/dev/null)
    if [ -n "$RATE_5H" ] && [ "$RATE_5H" != "null" ]; then
        RATE_PCT=$(echo "$RATE_5H" | awk '{printf "%d", $1}')
        RATE_FILLED=$(( RATE_PCT / 10 ))
        RATE_EMPTY=$(( 10 - RATE_FILLED ))
        RATE_BAR=$(printf '█%.0s' $(seq 1 $RATE_FILLED 2>/dev/null) || true)$(printf '░%.0s' $(seq 1 $RATE_EMPTY 2>/dev/null) || true)
        USAGE_SECTION=" | 5h: [$RATE_BAR] ${RATE_PCT}%"
    fi
fi

# Output status line
echo "[$MODEL_LABEL] [$BAR] ${PERCENT_USED}%$GIT_BRANCH | +${LINES_ADDED}/-${LINES_REMOVED}$USAGE_SECTION"
