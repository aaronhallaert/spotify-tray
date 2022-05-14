#!/usr/bin/env bash

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source "$CURRENT_DIR/helpers.sh"

progress() {
read -r -d '' SCRIPT <<END
set theApp to "$APP"

if application theApp is running then
  tell application "$APP"
    return player position as string
  end tell
end if
END

osascript -e "${SCRIPT}"
}

print_progress() {
  echo "$(progress)"
}

main() {
  print_progress
}

main
