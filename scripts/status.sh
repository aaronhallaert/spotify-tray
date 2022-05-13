#!/usr/bin/env bash

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source "$CURRENT_DIR/helpers.sh"

playing_icon="▶︎"
paused_icon="❚❚"
stopped_icon="■"

music_status() {
read -r -d '' SCRIPT <<END
set theApp to "$APP"

if application theApp is running then
  tell application "$APP"
    return player state as string
  end tell
end if
END

osascript -e "${SCRIPT}"
}

print_music_status() {
  local status=$(music_status)

  if [[ "$status" == "playing" ]]; then
    echo "${playing_icon}"
  elif [[ "$status" == "paused" ]]; then
    echo "${paused_icon}"
  else
    echo "${stopped_icon}"
  fi
}

main() {
  print_music_status
}

main


