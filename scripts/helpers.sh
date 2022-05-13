APP=${MUSIC_APP:-"Spotify"}
get_tmux_option() {
  local option=$1
  local default_value=$2
  local option_value=$(tmux show-option -gqv "$option")
  if [ -z "$option_value" ]; then
    echo "$default_value"
  else
    echo "$option_value"
  fi
}

set_tmux_option() {
  local option="$1"
  local value="$2"
  tmux set-option -gq "$option" "$value"
}

current_track_property() {
  local prop="${1}"
read -r -d '' SCRIPT <<END
set theApp to "$APP"

if application theApp is running then
  tell application "$APP"
    return %s of current track as string
  end tell
end if
END

osascript -e "$(printf "${SCRIPT}" "$prop")"
}
