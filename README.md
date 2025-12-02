# Spotifycli

[![CircleCI](https://circleci.com/gh/masroorhasan/spotifycli/tree/master.svg?style=svg)](https://circleci.com/gh/masroorhasan/spotifycli/tree/master)

A command line interface to manage Spotify playlists.

## Install

To use `spotifycli` you have to register the application on Spotify's developer platform:

1. Sign up or login at [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Create a new app
3. In your app settings, add the following **Redirect URI** (note: use `127.0.0.1` not `localhost`):
   ```
   http://127.0.0.1:8080/callback
   ```
4. Set the following environment variables with your app's Client ID and Client Secret:
   ```bash
   export SPOTIFY_ID=xxx
   export SPOTIFY_SECRET=xxx
   ```

## Usage

### Commands
List of available commands:
```
$ ./spotifycli --help
A command line interface to manage Spotify playlists.

Usage:
  spotifycli [command]

Available Commands:
  add         Add track by name to playlist
  aid         Add track by ID to playlist
  ato         Add currently playing track to playlist
  del         Delete a playlist
  help        Help about any command
  list        List tracks in playlist
  login       Login to authenticate Spotify account
  logout      Logout from Spotify account
  new         Create new playlist
  now         Displays the currently playing track
  playlists   Show all playlists
  rm          Remove track from playlist
  search      search tracks, albums, artists, playlists by name
  show        Display information about a track by ID

Flags:
  -h, --help   help for spotifycli

Use "spotifycli [command] --help" for more information about a command.
```

### Search
Search using query terms on top of tracks (`tr`), albums (`al`), artists (`ar`) or playlists (`pl`) by name.

```
./spotifycli search --help
search tracks, albums, artists, playlists by name

Usage:
  spotifycli search --t [SEARCH_TYPE] --q [SEARCH_QUERY] [flags]

Flags:
  -h, --help       help for search
      --q string   The search query term.
      --t string   The search type (tr, al, ar, pl).
```

Sample search for type `tr` (track).
```
./spotifycli search --t "tr" --q "one step closer - live"
```
