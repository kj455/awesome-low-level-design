#!/bin/bash
set -eo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
WORKSPACE_DIR="$ROOT_DIR/workspace/go"

# Problem definitions: "id|difficulty|name|diagram_name"
PROBLEMS=(
  "parking-lot|Easy|Design Parking Lot|parkinglot"
  "stack-overflow|Easy|Design Stack Overflow|stackoverflow"
  "vending-machine|Easy|Design Vending Machine|vendingmachine"
  "logging-framework|Easy|Design Logging Framework|loggingframework"
  "traffic-signal|Easy|Design Traffic Signal|trafficcontrolsystem"
  "coffee-vending-machine|Easy|Design Coffee Vending Machine|coffeevendingmachine"
  "task-management-system|Easy|Design Task Management System|taskmanagementsystem"
  "atm|Medium|Design ATM|atm"
  "linkedin|Medium|Design LinkedIn|linkedin"
  "lru-cache|Medium|Design LRU Cache|lrucache"
  "tic-tac-toe|Medium|Design Tic Tac Toe|tictactoe"
  "pub-sub-system|Medium|Design Pub Sub System|pubsubsystem"
  "elevator-system|Medium|Design Elevator System|elevatorsystem"
  "car-rental-system|Medium|Design Car Rental System|carrentalsystem"
  "online-auction-system|Medium|Design Online Auction System|onlineauctionsystem"
  "hotel-management-system|Medium|Design Hotel Management System|hotelmanagementsystem"
  "digital-wallet-service|Medium|Design Digital Wallet Service|digitalwalletservice"
  "airline-management-system|Medium|Design Airline Management System|airlinemanagementsystem"
  "library-management-system|Medium|Design Library Management System|librarymanagementsystem"
  "social-networking-service|Medium|Design Social Networking Service|socialnetworkingservice"
  "restaurant-management-system|Medium|Design Restaurant Management System|restaurantmanagementsystem"
  "concert-ticket-booking-system|Medium|Design Concert Ticket Booking System|concertticketbookingsystem"
  "cricinfo|Hard|Design Cricinfo|cricinfo"
  "splitwise|Hard|Design Splitwise|splitwise"
  "chess-game|Hard|Design Chess Game|chessgame"
  "snake-and-ladder|Hard|Design Snake and Ladder|snakeandladdergame"
  "ride-sharing-service|Hard|Design Ride Sharing Service|ridesharingservice"
  "course-registration-system|Hard|Design Course Registration System|courseregistrationsystem"
  "movie-ticket-booking-system|Hard|Design Movie Ticket Booking System|movieticketbookingsystem"
  "online-shopping-service|Hard|Design Online Shopping Service|onlineshoppingservice"
  "online-stock-brokerage-system|Hard|Design Online Stock Brokerage System|onlinestockbrokeragesystem"
  "music-streaming-service|Hard|Design Music Streaming Service|musicstreamingservice"
  "food-delivery-service|Hard|Design Food Delivery Service|fooddeliveryservice"
)

# Check if problem is solved (workspace directory exists)
is_solved() {
  [ -d "$WORKSPACE_DIR/$1" ]
}

# Show progress summary
show_progress() {
  local solved=0
  local total=${#PROBLEMS[@]}
  for entry in "${PROBLEMS[@]}"; do
    local id="${entry%%|*}"
    if is_solved "$id"; then
      solved=$((solved + 1))
    fi
  done
  echo "  $solved/$total completed"
}

# Generate list for fzf
generate_list() {
  for entry in "${PROBLEMS[@]}"; do
    local id="${entry%%|*}"
    local rest="${entry#*|}"
    local difficulty="${rest%%|*}"
    rest="${rest#*|}"
    local name="${rest%%|*}"
    local mark=""
    if is_solved "$id"; then
      mark=" ☑"
    fi
    printf "[%-6s] %s%s\n" "$difficulty" "$name" "$mark"
  done | sort
}

# Get problem entry by name
get_problem_entry() {
  local search_name="$1"
  for entry in "${PROBLEMS[@]}"; do
    local id="${entry%%|*}"
    local rest="${entry#*|}"
    local difficulty="${rest%%|*}"
    rest="${rest#*|}"
    local name="${rest%%|*}"
    if [ "$name" = "$search_name" ]; then
      echo "$entry"
      return 0
    fi
  done
  return 1
}

# Setup workspace for selected problem
setup_workspace() {
  local entry="$1"
  local id="${entry%%|*}"
  local rest="${entry#*|}"
  rest="${rest#*|}"
  rest="${rest#*|}"
  local diagram_name="${rest%%|*}"
  local dir="$WORKSPACE_DIR/$id"

  mkdir -p "$dir"

  # Copy problem.md (Requirements only, remove UML section and below)
  if [ -f "$ROOT_DIR/problems/$id.md" ]; then
    sed '/^## UML Class Diagram/,$d' "$ROOT_DIR/problems/$id.md" > "$dir/problem.md"
  fi

  # Create main.go template (only if doesn't exist)
  if [ ! -f "$dir/main.go" ]; then
    cat > "$dir/main.go" << 'EOF'
package main

import "fmt"

func main() {
	fmt.Println("TODO: Implement solution")
}
EOF
  fi

  # Initialize go.mod (if doesn't exist)
  if [ ! -f "$WORKSPACE_DIR/go.mod" ]; then
    (cd "$WORKSPACE_DIR" && go mod init workspace)
  fi

  echo ""
  echo "Created: $dir"
  echo ""
  echo "Next steps:"
  echo "  cd $dir"
  echo "  go run main.go"
}

# Main
main() {
  # Check fzf is installed
  if ! command -v fzf > /dev/null 2>&1; then
    echo "Error: fzf is not installed. Install it with: brew install fzf"
    exit 1
  fi

  show_progress
  echo ""

  # Run fzf and get selection
  selected=$(generate_list | fzf --height=40% --reverse) || true
  [ -z "$selected" ] && exit 0

  # Extract problem name from selection (remove difficulty badge and checkmark)
  name=$(echo "$selected" | sed 's/\[.*\] //' | sed 's/ ☑$//')

  # Get problem entry and setup workspace
  entry=$(get_problem_entry "$name")
  if [ -n "$entry" ]; then
    setup_workspace "$entry"
  else
    echo "Error: Problem not found"
    exit 1
  fi
}

main
