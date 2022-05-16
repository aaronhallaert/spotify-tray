APP=${MUSIC_APP:-"Spotify"}

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
