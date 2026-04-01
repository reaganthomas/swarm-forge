#!/bin/bash
# Usage: ./notify-agent.sh <target-pane-index> "message"
# Panes: 0=Architect, 1=E2E Interpreter, 2=Coder, 3=Metrics
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')
echo "[$TIMESTAMP] [pane $1] $2" >> logs/agent_messages.log
tmux send-keys -t swarmforge:swarm.$1 "$2" Enter
