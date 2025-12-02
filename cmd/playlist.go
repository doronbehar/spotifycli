package cmd

import (
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zmb3/spotify"
)

var (
	addtoPlaylistName string
)

var (
	trackID string
)

var (
	addTrackID                 string
	addTrackByIDToPlaylistName string
)

var (
	addTrackName                 string
	addTrackByNameToPlaylistName string
)

var (
	rmTrackName             string
	rmTrackFromPlaylistName string
)

var (
	newPlaylistName string
)

var (
	delPlaylistName string
)

var (
	clearPlaylistName string
)

var (
	listPlaylistTracksName string
)

// isBase62 checks if a string contains only valid base62 characters (0-9, a-z, A-Z)
func isBase62(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
			return false
		}
	}
	return true
}

// parseTrackID extracts a Spotify track ID from either a URL or a plain ID string.
// If the input is a Spotify URL (e.g., https://open.spotify.com/track/2b7VhCSKWZAFDrDPKTJ1z2?si=xyz),
// it extracts the track ID from the path. Otherwise, it returns the input as-is.
func parseTrackID(input string) (string, error) {
	// Check if input looks like a URL
	if strings.Contains(input, "spotify.com") {
		parsedURL, err := url.Parse(input)
		if err != nil {
			return "", fmt.Errorf("invalid URL: %w", err)
		}

		// Extract track ID from path (e.g., /track/2b7VhCSKWZAFDrDPKTJ1z2)
		parts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
		if len(parts) >= 2 && parts[0] == "track" {
			return parts[1], nil
		}

		return "", fmt.Errorf("URL does not contain a valid track path")
	}

	// Return input as-is if it's not a URL
	return input, nil
}

func newCurrentTrackCmd() *cobra.Command {
	nowCmd := &cobra.Command{
		Use:   "now",
		Short: "Displays the currently playing track",
		RunE: func(cmd *cobra.Command, args []string) error {
			return displayCurrentTrack(cmd, args)
		},
	}
	return nowCmd
}

func newShowTrackCmd() *cobra.Command {
	addtoCmd := &cobra.Command{
		Use:   "show --tid [TRACK_ID]",
		Short: "Display information about a track by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			return displayTrackById(cmd, args)
		},
	}
	addtoCmd.Flags().StringVar(&trackID, "tid", "", "Id of track to display.")
	return addtoCmd
}

func newAddtoPlaylistCmd() *cobra.Command {
	addtoCmd := &cobra.Command{
		Use:   "ato --p [PLAYLIST_NAME]",
		Short: "Add currently playing track to playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return addto(cmd, args)
		},
	}
	addtoCmd.Flags().StringVar(&addtoPlaylistName, "p", "", "Add current track to specified playlist.")
	return addtoCmd
}

func newAddTrackByIDToPlaylistCmd() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "aid --tid [TRACK_ID] --p [PLAYLIST_NAME]",
		Short: "Add track by ID to playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return addTrackByIDToPlaylist(cmd, args)
		},
	}
	addCmd.Flags().StringVar(&addTrackID, "tid", "", "Id of track to add to playlist.")
	addCmd.Flags().StringVar(&addTrackByIDToPlaylistName, "p", "", "Name of playlist to add track to.")
	return addCmd
}

func newAddTrackByNameToPlaylistCmd() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add --t [TRACK_NAME] --p [PLAYLIST_NAME]",
		Short: "Add track by name to playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return addTrackByNameToPlaylist(cmd, args)
		},
	}
	addCmd.Flags().StringVar(&addTrackName, "t", "", "Name of track to add to playlist.")
	addCmd.Flags().StringVar(&addTrackByNameToPlaylistName, "p", "", "Name of playlist to add track to.")
	return addCmd
}

func newRemoveTrackFromPlaylistCmd() *cobra.Command {
	rmCmd := &cobra.Command{
		Use:   "rm --t [TRACK_NAME|TRACK_ID|TRACK_URL] --p [PLAYLIST_NAME]",
		Short: "Remove track from playlist by name, ID, or URL",
		RunE: func(cmd *cobra.Command, args []string) error {
			return rmTrackByNameFromPlaylist(cmd, args)
		},
	}
	rmCmd.Flags().StringVar(&rmTrackName, "t", "", "Name, ID, or URL of track to remove.")
	rmCmd.Flags().StringVar(&rmTrackFromPlaylistName, "p", "", "Name of playlist to remove track from.")
	return rmCmd
}

func newListPlaylistsCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "playlists",
		Short: "Show all playlists",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listPlaylists(cmd, args)
		},
	}
	return newCmd
}

func newCreatePlaylistCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "new --p [PLAYLIST_NAME]",
		Short: "Create new playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return newPlaylist(cmd, args)
		},
	}
	newCmd.Flags().StringVar(&newPlaylistName, "p", "", "Name of new playlist.")
	return newCmd
}

func newDeletePlaylistCmd() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "del --p [PLAYLIST_NAME]",
		Short: "Delete a playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return deletePlaylist(cmd, args)
		},
	}
	deleteCmd.Flags().StringVar(&delPlaylistName, "p", "", "Name of playlist to delete.")
	return deleteCmd
}

func newClearPlaylistCmd() *cobra.Command {
	clearCmd := &cobra.Command{
		Use:   "clear --p [PLAYLIST_NAME]",
		Short: "Remove all tracks from a playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return clearPlaylist(cmd, args)
		},
	}
	clearCmd.Flags().StringVar(&clearPlaylistName, "p", "", "Name of playlist to clear.")
	return clearCmd
}

func newListPlaylistTracksCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list --p [PLAYLIST_NAME]",
		Short: "List tracks in playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTracksFromPlaylist(cmd, args)
		},
	}
	listCmd.Flags().StringVar(&listPlaylistTracksName, "p", "", "Name of playlist to list tracks from.")
	return listCmd
}

func displayTrack(track *spotify.FullTrack) error {
	// format and display
	var data [][]interface{}
	item := []string{
		string(track.ID),
		track.Name,
		track.Album.Name,
		track.Artists[0].Name,
		(time.Duration(track.Duration) * time.Millisecond).Truncate(time.Second).String(),
		strconv.Itoa(track.Popularity),
		strconv.FormatBool(track.Explicit),
		track.PreviewURL,
	}
	row := make([]interface{}, len(item))
	for i, d := range item {
		row[i] = d
	}
	data = append(data, row)
	printSimple([]string{"ID", "Name", "Album", "Artist", "Duration", "Popularity", "Explicit", "Preview"}, data)
	return nil
}

func displayTrackById(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// parse track ID from URL or plain ID
	parsedTrackID, err := parseTrackID(trackID)
	if err != nil {
		return err
	}

	// get the track (check for existence)
	track, err := client.GetTrack(spotify.ID(parsedTrackID))
	if err != nil {
		return err
	}

  displayTrack(track)
  return nil
}

func displayCurrentTrack(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// get current playing song
	playing, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		return err
	}

  displayTrack(playing.Item)
	return nil
}

func addto(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// get my playlists
	pl, err := getPlaylistByName(addtoPlaylistName)
	if err != nil {
		return err
	}
	fmt.Println("Playlist: ", pl.Name)

	// get current playing song
	playing, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		return err
	}
	fmt.Println("Track: ", playing.Item.Name)

	// add track to playlist
	_, err = client.AddTracksToPlaylist(pl.ID, playing.Item.ID)
	if err != nil {
		return err
	}
	fmt.Printf("Added track \"%s\" to playlist \"%s\".\n", playing.Item.Name, pl.Name)
	return nil
}

func listPlaylists(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// get all playlists for the user
	playlists, err := getPlaylists()
	if err != nil {
		return err
	}

	// format resulting data
	var data [][]interface{}
	if playlists.Playlists != nil {
		for _, item := range playlists.Playlists {
			track := []string{
				string(item.ID),
				item.Name,
				item.Owner.DisplayName,
				strconv.FormatBool(item.IsPublic),
				strconv.FormatBool(item.Collaborative),
				strconv.FormatUint(uint64(item.Tracks.Total), 10)}
			row := make([]interface{}, len(track))
			for i, d := range track {
				row[i] = d
			}
			data = append(data, row)
		}
	}

	// pretty print track results
	printSimple([]string{"ID", "Name", "Owner", "Public", "Collaborative", "Tracks"}, data)
	return nil
}

func newPlaylist(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// create new playlist
	playlist, err := client.CreatePlaylistForUser(user.ID, newPlaylistName, "", true)
	if err != nil {
		return err
	}
	fmt.Println("Created public playlist: ", playlist.Name)
	return nil
}

func deletePlaylist(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// get the playlist
	pl, err := getPlaylistByName(delPlaylistName)
	if err != nil {
		return err
	}

	// unfollow and return
	// TODO: delete != unfollow?
	return client.UnfollowPlaylist(spotify.ID(user.ID), pl.ID)
}

func clearPlaylist(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// get the playlist
	pl, err := getPlaylistByName(clearPlaylistName)
	if err != nil {
		return err
	}
	fmt.Println("Playlist: ", pl.Name)

	// get all tracks from the playlist
	tracks, err := client.GetPlaylistTracks(pl.ID)
	if err != nil {
		return err
	}

	if len(tracks.Tracks) == 0 {
		fmt.Println("Playlist is already empty.")
		return nil
	}

	// collect all track IDs
	var trackIDs []spotify.ID
	for _, t := range tracks.Tracks {
		trackIDs = append(trackIDs, t.Track.ID)
	}

	// remove all tracks from the playlist
	_, err = client.RemoveTracksFromPlaylist(pl.ID, trackIDs...)
	if err != nil {
		return err
	}

	fmt.Printf("Removed %d track(s) from playlist \"%s\".\n", len(trackIDs), pl.Name)
	return nil
}

func addTrackByIDToPlaylist(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// get the playlist by name
	pl, err := getPlaylistByName(addTrackByIDToPlaylistName)
	if err != nil {
		return err
	}
	fmt.Println("Playlist: ", pl.Name)

	// parse track ID from URL or plain ID
	parsedTrackID, err := parseTrackID(addTrackID)
	if err != nil {
		return err
	}

	// get the track (check for existence)
	tr, err := client.GetTrack(spotify.ID(parsedTrackID))
	if err != nil {
		return err
	}
	fmt.Println("Track: ", tr.Name)

	// add track to playlist
	_, err = client.AddTracksToPlaylist(pl.ID, tr.ID)
	if err != nil {
		return err
	}
	fmt.Printf("Added track \"%s\" to playlist \"%s\".\n", tr.Name, pl.Name)
	return nil
}

func addTrackByNameToPlaylist(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// get the playlist by name
	pl, err := getPlaylistByName(addTrackByNameToPlaylistName)
	if err != nil {
		return err
	}
	fmt.Println("Playlist: ", pl.Name)

	// Search for the track
	results, err := client.Search(addTrackName, spotify.SearchTypeTrack)
	if err != nil {
		return err
	}

	// add most popular track to playlist from results
	if results.Tracks != nil {
		tracks := results.Tracks.Tracks[:]
		sort.Slice(tracks, func(i, j int) bool { return tracks[i].Popularity > tracks[j].Popularity })
		fmt.Println("Track: ", tracks[0].Name)

		// add track to playlist
		_, err = client.AddTracksToPlaylist(pl.ID, tracks[0].ID)
		if err != nil {
			return err
		}
		fmt.Printf("Added track \"%s\" to playlist \"%s\".\n", tracks[0].Name, pl.Name)
	} else {
		fmt.Printf("Track %s not found.\n", addTrackName)
	}
	return nil
}

func rmTrackByNameFromPlaylist(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// get the playlist by name
	pl, err := getPlaylistByName(rmTrackFromPlaylistName)
	if err != nil {
		return err
	}

	// get track in playlist and validate existence
	var matchedTrack spotify.SimpleTrack
	ptracks, err := client.GetPlaylistTracks(pl.ID)
	if err != nil {
		return err
	}

	// Check if input is a URL/ID or a track name
	parsedTrackID, parseErr := parseTrackID(rmTrackName)
	isIDOrURL := parseErr == nil && (strings.Contains(rmTrackName, "spotify.com") || isBase62(rmTrackName))

	// Search for track by ID or name
	for _, t := range ptracks.Tracks {
		if isIDOrURL {
			// Match by ID
			if string(t.Track.SimpleTrack.ID) == parsedTrackID {
				matchedTrack = t.Track.SimpleTrack
				break
			}
		} else {
			// Match by name
			if rmTrackName == t.Track.SimpleTrack.Name {
				matchedTrack = t.Track.SimpleTrack
				break
			}
		}
	}
	if reflect.DeepEqual(matchedTrack, spotify.SimpleTrack{}) {
		return fmt.Errorf("track %s not found in playlist %s", rmTrackName, rmTrackFromPlaylistName)
	}
	fmt.Println("Track: ", matchedTrack.Name)

	// remove track from playlist
	_, err = client.RemoveTracksFromPlaylist(pl.ID, matchedTrack.ID)
	if err != nil {
		return err
	}
	fmt.Printf("Removed track \"%s\" from playlist \"%s\".\n", matchedTrack.Name, rmTrackFromPlaylistName)
	return nil
}

func listTracksFromPlaylist(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	pl, err := getPlaylistByName(listPlaylistTracksName)
	if err != nil {
		return err
	}

	// get tracks from playlist
	tracks, err := client.GetPlaylistTracks(pl.ID)
	if err != nil {
		return err
	}

	// format resulting data
	var data [][]interface{}
	if tracks.Tracks != nil {
		for _, item := range tracks.Tracks {
			track := []string{
				string(item.Track.ID),
				item.Track.Name,
				item.Track.Album.Name,
				item.Track.Artists[0].Name,
				strconv.Itoa(item.Track.Popularity)}
			row := make([]interface{}, len(track))
			for i, d := range track {
				row[i] = d
			}
			data = append(data, row)
		}
	}

	// pretty print track results
	printSimple([]string{"ID", "Name", "Album", "Artist", "Popularity"}, data)
	return nil
}

func getPlaylists() (*spotify.SimplePlaylistPage, error) {
	playlists, err := client.CurrentUsersPlaylists()
	if err != nil {
		return &(spotify.SimplePlaylistPage{}), err
	}

  return playlists, nil
}

func getPlaylistByName(playlistName string) (spotify.SimplePlaylist, error) {
	// get current user's playlists
	playlists, err := client.CurrentUsersPlaylists()
	if err != nil {
		return spotify.SimplePlaylist{}, err
	}

	// match by name
	var matchPlaylist spotify.SimplePlaylist
	for _, p := range playlists.Playlists {
		if playlistName == p.Name {
			matchPlaylist = p
			break
		}
	}

	// check if found and return
	if reflect.DeepEqual(matchPlaylist, spotify.SimplePlaylist{}) {
		return spotify.SimplePlaylist{}, fmt.Errorf("playlist not found: %s", playlistName)
	}
	return matchPlaylist, nil
}
