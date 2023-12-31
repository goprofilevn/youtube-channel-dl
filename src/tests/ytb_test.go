package tests

import (
	"fmt"
	"io"
	"os"
	"strings"

	"testing"

	"github.com/kkdai/youtube/v2"
)

// ExampleDownload : Example code for how to use this package for download video.
func TestClient(t *testing.T) {
	videoID := "P5esQ3PreGs"
	t.Log("Download video ID: ", videoID)
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		t.Error(err)
	}
	formats := video.Formats.WithAudioChannels()
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		t.Error(err)
	}
	defer stream.Close()

	file, err := os.Create("video.mp4")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		t.Error(err)
	}
}

// Example usage for playlists: downloading and checking information.
func TestPlaylist(t *testing.T) {
	playlistID := "UULFhmJqLUWbSHvZRFZ_DJk4Kw"
	client := youtube.Client{}

	playlist, err := client.GetPlaylist(playlistID)
	if err != nil {
		t.Error(err)
	}

	/* ----- Enumerating playlist videos ----- */
	header := fmt.Sprintf("Playlist %s by %s", playlist.Title, playlist.Author)
	t.Log(header)
	t.Log(strings.Repeat("=", len(header)) + "\n")

	for k, v := range playlist.Videos {
		fmt.Printf("(%d) %s - '%s' - %s\n", k+1, v.Author, v.Title, v.ID)
	}

	/* ----- Downloading the 1st video ----- */
	entry := playlist.Videos[0]
	video, err := client.VideoFromPlaylistEntry(entry)
	if err != nil {
		t.Error(err)
	}
	// Now it's fully loaded.

	fmt.Printf("Downloading %s by '%s'!\n", video.Title, video.Author)

	stream, _, err := client.GetStream(video, &video.Formats[0])
	if err != nil {
		t.Error(err)
	}

	file, err := os.Create("video.mp4")

	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	_, err = io.Copy(file, stream)

	if err != nil {
		t.Error(err)
	}

	t.Log("Downloaded /video.mp4")
}
