package upload

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/WeatherGod3218/weather-reels-watcher/internal/models"
	"github.com/alfg/mp4"
	"github.com/bdragon300/tusgo"
)

const MINIMUM_VIDEO_LENGTH = 5

var client *tusgo.Client

func getMP4Length(fileName string) (float64, error) {
	file, err := mp4.Open(fileName)
	if err != nil {
		return 0, err
	}

	if file.Moov == nil || file.Moov.Mvhd == nil {
		return 0, fmt.Errorf("file does not contain metadata")
	}

	durationSeconds := (float64(file.Moov.Mvhd.Duration) / float64(file.Moov.Mvhd.Timescale))
	return durationSeconds, nil
}

func createUploadFromFile(config models.Config, file *os.File, fileName string, baseDir string) (*tusgo.Upload, error) {
	finfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	length, err := getMP4Length(fileName)
	if err != nil {
		return nil, err
	}

	if !config.IncludeDirectoryPath {
		fileName = strings.Trim(fileName, baseDir)
	}

	metaData := map[string]string{
		"filename":      fileName,
		"file_mod_date": strconv.FormatInt(finfo.ModTime().UnixMicro(), 10),
		"video_length":  strconv.FormatFloat(length, 'f', 2, 64),
	}

	upload := &tusgo.Upload{}
	if _, err = client.CreateUpload(upload, finfo.Size(), false, metaData); err != nil {
		return nil, err
	}

	return upload, nil
}

func uploadWithRetry(destination *tusgo.UploadStream, file *os.File) error {

	if _, err := destination.Sync(); err != nil {
		return err
	}
	if _, err := file.Seek(destination.Tell(), io.SeekStart); err != nil {
		return err
	}

	_, err := io.Copy(destination, file)
	attempts := 10
	for err != nil && attempts > 0 {
		if _, ok := err.(net.Error); !ok && !errors.Is(err, tusgo.ErrChecksumMismatch) {
			return err // Permanent error, no luck
		}
		time.Sleep(5 * time.Second)
		attempts--
		_, err = io.Copy(destination, file) // Try to resume the transfer again
	}
	if attempts == 0 {
		return errors.New("too many attempts to upload the data")
	}
	return nil
}

func InitTusio(creds models.Credentials) {
	baseUrl, _ := url.Parse(fmt.Sprintf("%s/api/v2/files/", creds.ServerURL))

	client = tusgo.NewClient(http.DefaultClient, baseUrl)
	client.GetRequest = func(method, url string, body io.Reader, tusClient *tusgo.Client, httpClient *http.Client) (*http.Request, error) {
		req, err := http.NewRequest(method, url, body)
		if err != nil {
			return nil, err
		}

		req.Header.Set("X-User-Id", creds.UserId) //TODO: Add users? would be nifty
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", creds.ServerAccessCode))

		return req, nil
	}
}

func UploadVideo(config models.Config, fileName string, baseDir string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	upload, err := createUploadFromFile(config, file, fileName, baseDir)
	if err != nil {
		return err
	}

	stream := tusgo.NewUploadStream(client, upload)
	if err := uploadWithRetry(stream, file); err != nil {
		return err
	}

	return nil
}

// func BackupFile(fileName string) error {
// 	fileLength, err := getMP4Length(fileName)
// 	if err != nil {
// 		return err
// 	}

// 	if fileLength < MINIMUM_VIDEO_LENGTH {
// 		logging.Logger.Info(fmt.Sprintf("Recieved video that was less than %d seconds, not backing up!", MINIMUM_VIDEO_LENGTH))
// 		return nil
// 	}

// 	fileInfo, err := os.Stat(fileName)
// 	if err != nil {
// 		return err
// 	}
// 	fileSize := fileInfo.Size()

// 	jsonData := models.VideoStartBackupRequest{
// 		FileName:   fileName,
// 		FileLength: fileLength,
// 		FileSize:   fileSize,
// 	}

// 	jsonBytes, err := json.Marshal(&jsonData)
// 	if err != nil {
// 		return err
// 	}

// 	createReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/videos/backup", os.Getenv("SERVER_URL")), bytes.NewBuffer(jsonBytes))
// 	if err != nil {
// 		return err
// 	}
// 	createReq.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	res, err := client.Do(createReq)
// 	if err != nil {
// 		return err
// 	}
// 	if res.StatusCode != http.StatusOK {
// 		return fmt.Errorf("Failed to start video request!")
// 	}

// 	defer res.Body.Close()

// 	var resp models.VideoStartBackupResponse

// 	bodyBytes, _ := io.ReadAll(res.Body)
// 	if err = json.Unmarshal(bodyBytes, &resp); err != nil {
// 		return err
// 	}

// 	uploadCompleted := false
// 	defer func() {
// 		if !uploadCompleted {
// 			AbortUpload(fileName, resp.VideoS3UploadID)
// 		}
// 	}()

// 	if resp.VideoS3UploadID == "" {
// 		return fmt.Errorf("server returned empty upload ID, raw body: %s", bodyBytes)
// 	}

// 	completeParts, err := UploadVideo(resp, fileName)
// 	if err != nil {
// 		return err
// 	}

// 	completeBackup := models.VideoCompleteBackupRequest{
// 		RowID:           resp.RowID,
// 		VideoS3UploadID: resp.VideoS3UploadID,
// 		CompletedParts:  completeParts,
// 		Filename:        fileName,
// 	}
// 	completeBackupBytes, err := json.Marshal(&completeBackup)
// 	if err != nil {
// 		return err
// 	}

// 	completeReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/videos/complete", os.Getenv("SERVER_URL")), bytes.NewBuffer(completeBackupBytes))
// 	if err != nil {
// 		return err
// 	}

// 	completeReq.Header.Set("Content-Type", "application/json")

// 	completeRes, err := client.Do(completeReq)
// 	if err != nil {
// 		return err
// 	}
// 	defer completeRes.Body.Close()

// 	if completeRes.StatusCode != http.StatusOK {
// 		return fmt.Errorf("Failed to start video request!")
// 	}

// 	logging.Logger.Info(fmt.Sprintf("Successfully backed up Video: %s", fileName))
// 	uploadCompleted = true
// 	return nil
// }
