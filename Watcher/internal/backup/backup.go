package backup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/WeatherGod3218/weather-reels-watcher/internal/logging"
	"github.com/WeatherGod3218/weather-reels-watcher/internal/models"
)

func AbortUpload(fileName string, uploadId string) error {

	abortJson := models.VideoAbortBackupRequest{
		Filename:        fileName,
		VideoS3UploadID: uploadId,
	}

	abortBytes, err := json.Marshal(abortJson)
	if err != nil {
		return err
	}

	abortReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/videos/abort", os.Getenv("SERVER_URL")), bytes.NewBuffer(abortBytes))
	if err != nil {
		return err
	}
	abortReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(abortReq)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func BackupFile(fileName string) error {
	fileLength, err := getMP4Length(fileName)
	if err != nil {
		return err
	}

	if fileLength < MINIMUM_VIDEO_LENGTH {
		logging.Logger.Info(fmt.Sprintf("Recieved video that was less than %d seconds, not backing up!", MINIMUM_VIDEO_LENGTH))
		return nil
	}

	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	jsonData := models.VideoStartBackupRequest{
		FileName:   fileName,
		FileLength: fileLength,
		FileSize:   fileSize,
	}

	jsonBytes, err := json.Marshal(&jsonData)
	if err != nil {
		return err
	}

	createReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/videos/backup", os.Getenv("SERVER_URL")), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	createReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(createReq)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to start video request!")
	}

	defer res.Body.Close()

	var resp models.VideoStartBackupResponse

	bodyBytes, _ := io.ReadAll(res.Body)
	if err = json.Unmarshal(bodyBytes, &resp); err != nil {
		return err
	}

	uploadCompleted := false
	defer func() {
		if !uploadCompleted {
			AbortUpload(fileName, resp.VideoS3UploadID)
		}
	}()

	if resp.VideoS3UploadID == "" {
		return fmt.Errorf("server returned empty upload ID, raw body: %s", bodyBytes)
	}

	completeParts, err := UploadVideo(resp, fileName)
	if err != nil {
		return err
	}

	completeBackup := models.VideoCompleteBackupRequest{
		RowID:           resp.RowID,
		VideoS3UploadID: resp.VideoS3UploadID,
		CompletedParts:  completeParts,
		Filename:        fileName,
	}
	completeBackupBytes, err := json.Marshal(&completeBackup)
	if err != nil {
		return err
	}

	completeReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/videos/complete", os.Getenv("SERVER_URL")), bytes.NewBuffer(completeBackupBytes))
	if err != nil {
		return err
	}

	completeReq.Header.Set("Content-Type", "application/json")

	completeRes, err := client.Do(completeReq)
	if err != nil {
		return err
	}
	defer completeRes.Body.Close()

	if completeRes.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to start video request!")
	}

	logging.Logger.Info(fmt.Sprintf("Successfully backed up Video: %s", fileName))
	uploadCompleted = true
	return nil
}
