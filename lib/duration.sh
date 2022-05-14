#!/usr/bin/env bash

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source "$CURRENT_DIR/helpers.sh"

print_duration() {
  current_track_property "duration" $APP
}

main() {
  print_duration $APP
}

main


