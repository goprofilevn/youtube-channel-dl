package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"youtube-channel-dl/src/utils"

	"github.com/kkdai/youtube/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx            context.Context
	downloadCtx    context.Context
	downloadCancel context.CancelFunc
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type Folder struct {
	Folder string `json:"folder"`
	Err    error  `json:"err"`
}

func (a *App) SelectFolder() Folder {
	folder, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select a path save",
	})
	return Folder{Folder: folder, Err: err}
}

func (a *App) GetVideoFromPlaylist(playlistID string) ([]*youtube.PlaylistEntry, error) {
	client := youtube.Client{}

	playlist, err := client.GetPlaylist(playlistID)
	if err != nil {
		return nil, err
	}

	return playlist.Videos, nil
}

func (a *App) StartDownload(videoID []string, pathFolder string) error {
	ctx, cancel := context.WithCancel(context.Background())
	a.downloadCtx = ctx
	a.downloadCancel = cancel
	defer cancel()
	runtime.EventsEmit(a.ctx, "started")
	for i, id := range videoID {
		if ctx.Err() != nil {
			break
		}
		runtime.EventsEmit(a.ctx, "status", PrepareJSON(id, "downloading"))
		err := DownloadVideo(id, i+1, pathFolder)
		if err != nil {
			fmt.Println(err)
			runtime.EventsEmit(a.ctx, "status", PrepareJSON(id, "error"))
			runtime.EventsEmit(a.ctx, "error", err.Error())
		} else {
			runtime.EventsEmit(a.ctx, "status", PrepareJSON(id, "done"))
		}
	}
	runtime.EventsEmit(a.ctx, "stopped")
	return nil
}

func (a *App) StopDownload() {
	a.downloadCancel()
	runtime.EventsEmit(a.ctx, "stopped", 0)
}

type VideoJSON struct {
	ID     string `json:"ID"`
	Status string `json:"Status"`
}

func PrepareJSON(videoID string, status string) string {
	data := VideoJSON{
		ID:     videoID,
		Status: status,
	}
	jsonStr, _ := json.Marshal(data)
	return string(jsonStr)
}

func DownloadVideo(videoID string, stt int, pathFolder string) error {
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		return err
	}

	formats := video.Formats.WithAudioChannels()
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		return err
	}
	defer stream.Close()

	pathSave := fmt.Sprintf("%s/%d. %s.mp4", pathFolder, stt, utils.GetValidFileName(video.Title))

	file, err := os.Create(pathSave)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		return err
	}

	return nil
}
