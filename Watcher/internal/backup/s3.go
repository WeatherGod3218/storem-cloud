package backup

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"sync"

	"github.com/WeatherGod3218/weather-reels-watcher/internal/logging"
	"github.com/WeatherGod3218/weather-reels-watcher/internal/models"
)

const GLOBAL_UPLOAD_CAP = 32
const PER_FILE_UPLOAD_CAP = 8

var uploadSemaphore = make(chan struct{}, GLOBAL_UPLOAD_CAP)

func uploadPart(ctx context.Context, url string, data *io.SectionReader, length int64) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, data)
	if err != nil {
		return "", err
	}
	req.ContentLength = length

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed: %s: %s", resp.Status, body)
	}

	etag := resp.Header.Get("ETag")
	return etag, nil
}

func UploadVideo(resp models.VideoStartBackupResponse, fileName string) ([]models.VideoCompletedPart, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logging.Logger.Info(fmt.Sprintf("Starting upload on file %s", fileName))
	var (
		wg        sync.WaitGroup
		lock      sync.Mutex
		uploadErr error

		completed        map[int]string = map[int]string{}
		perFileSemaphore                = make(chan struct{}, PER_FILE_UPLOAD_CAP)
	)

	modify := func(partNum int, etag string) {
		lock.Lock()
		defer lock.Unlock()
		completed[partNum] = etag
	}

	setErr := func(err error) {
		lock.Lock()
		defer lock.Unlock()
		if uploadErr == nil {
			uploadErr = err
			cancel()
		}
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	for _, part := range resp.VideoS3URLs {
		wg.Add(1)
		go func(part models.VideoURLPart) {
			defer wg.Done()

			if ctx.Err() != nil {
				return
			}

			uploadSemaphore <- struct{}{}
			defer func() { <-uploadSemaphore }()

			perFileSemaphore <- struct{}{}
			defer func() { <-perFileSemaphore }()

			reader := io.NewSectionReader(
				file,
				part.Offset,
				part.Size,
			)

			etag, err := uploadPart(ctx, part.RequestURL, reader, part.Size)
			if err != nil {
				setErr(err)
				return
			}

			modify(int(part.PartNumber), etag)
		}(part)
	}
	wg.Wait()

	if uploadErr != nil {
		return nil, uploadErr
	}

	partList := make([]models.VideoCompletedPart, len(completed))

	i := 0
	for partNum, etag := range completed {
		partList[i] = models.VideoCompletedPart{
			PartNumber: partNum,
			ETag:       etag,
		}
		i++
	}

	sort.Slice(partList, func(i, j int) bool {
		return partList[i].PartNumber < partList[j].PartNumber
	})

	logging.Logger.Info(fmt.Sprintf("Completed upload on %s", fileName))

	return partList, nil
}
