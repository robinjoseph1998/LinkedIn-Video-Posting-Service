package utils

import "os"

type InitializeUploadRequest struct {
	Owner          string `json:"owner"`
	FileSizeBytes  int    `json:"fileSizeBytes"`
	UploadCaptions bool   `json:"uploadCaptions"`
}

type UploadInstruction struct {
	UploadURL string `json:"uploadUrl"`
	FirstByte int    `json:"firstByte"`
	LastByte  int    `json:"lastByte"`
}

type UploadResponse struct {
	Value struct {
		UploadUrlsExpireAt int64               `json:"uploadUrlsExpireAt"`
		Video              string              `json:"video"`
		UploadInstructions []UploadInstruction `json:"uploadInstructions"`
		UploadToken        string              `json:"uploadToken"`
	} `json:"value"`
}

// Struct for the finalize video upload parameters
type FinalizeUploadRequest struct {
	Video           string   `json:"video"`
	UploadToken     string   `json:"uploadToken"`
	UploadedPartIds []string `json:"uploadedPartIds"`
}

type FinalizeUploadPayload struct {
	FinalizeUploadRequest FinalizeUploadRequest `json:"finalizeUploadRequest"`
}

func GetFileSize(filePath string) (int64, error) {
	// Open the file
	file, err := os.Stat(filePath)
	if err != nil {
		return 0, err // return 0 and error if the file doesn't exist or cannot be accessed
	}
	return file.Size(), nil // return the file size
}
