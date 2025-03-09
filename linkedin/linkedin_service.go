package linkedin

import (
	"Linkedin-Video-Posting/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// *********************** STEP: 1 - INITIALIZING THE VIDEO *****************************//
func InitializeLinkedInVideoUpload(accessToken, owner string, fileSize int) (*utils.UploadResponse, error) {
	url := "https://api.linkedin.com/rest/videos?action=initializeUpload"

	// Prepare request body
	reqBody := map[string]utils.InitializeUploadRequest{
		"initializeUploadRequest": {
			Owner:          owner,
			FileSizeBytes:  fileSize,
			UploadCaptions: true,
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")
	req.Header.Set("LinkedIn-Version", "202502")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Print raw response for debugging
	// fmt.Println("Raw API Response:", string(body))

	// Parse JSON response
	var uploadResp utils.UploadResponse
	err = json.Unmarshal(body, &uploadResp)
	if err != nil {
		return nil, err
	}

	// Check if upload instructions exist
	if len(uploadResp.Value.UploadInstructions) == 0 {
		return nil, fmt.Errorf("no upload instructions received. Response: %s", string(body))
	}

	return &uploadResp, nil
}

// *********************** STEP: 2 - UPLOADING THE VIDEO IN CHUNKS ***********************//
func UploadVideoFile(filePath string, uploadInstructions []utils.UploadInstruction, accessToken string) ([]string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	client := &http.Client{}
	var uploadPartIds []string

	// Loop through each upload instruction
	for _, instruction := range uploadInstructions {
		// Read chunk data based on first and last byte
		chunkSize := instruction.LastByte - instruction.FirstByte
		chunk := make([]byte, chunkSize)

		// Move to the first byte of the chunk to read
		_, err := file.Seek(int64(instruction.FirstByte), 0)
		if err != nil {
			return nil, fmt.Errorf("failed to seek to byte %d: %v", instruction.FirstByte, err)
		}

		// Read the chunk
		_, err = file.Read(chunk)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("failed to read chunk: %v", err)
		}

		// Create PUT request
		req, err := http.NewRequest("PUT", instruction.UploadURL, bytes.NewReader(chunk))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}

		// Set headers
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("LinkedIn-Version", "202502")

		// Execute request
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("upload failed: %v", err)
		}
		defer resp.Body.Close()

		// Read ETag from response
		etag := resp.Header.Get("ETag")
		if etag != "" {
			uploadPartIds = append(uploadPartIds, etag)
		}
	}

	return uploadPartIds, nil
}

// *********************** STEP: 3 - FINALIZING THE VIDEO ******************************//
func FinalizeVideoUpload(accessToken, urn string, parts []string) error {
	// Prepare the request body
	payload := utils.FinalizeUploadPayload{
		FinalizeUploadRequest: utils.FinalizeUploadRequest{
			Video:           urn,
			UploadToken:     "",
			UploadedPartIds: parts,
		},
	}

	// Convert the payload to JSON
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Create the POST request
	req, err := http.NewRequest("POST", "https://api.linkedin.com/rest/videos?action=finalizeUpload", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set the required headers
	req.Header.Set("LinkedIn-Version", "202502") // Replace with the version number in the format YYYYMM
	req.Header.Set("X-RestLi-Protocol-Version", "2.0.0")
	req.Header.Set("Authorization", "Bearer "+accessToken) // Replace {INSERT_TOKEN} with your actual Bearer token
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response: %v", resp.Status)
	}

	log.Println("Request successful, video finalized.")
	return nil
}

// *********************** STEP: 4 - CREATING THE POST *************************//
func CreatePost(accessToken, videoUrn, authorUrn, title, commentary string) error {
	// Prepare the post data in the new format
	postData := map[string]interface{}{
		"author":     authorUrn, // e.g., "urn:li:organization:5515715"
		"commentary": commentary,
		"visibility": "PUBLIC",
		"distribution": map[string]interface{}{
			"feedDistribution":               "MAIN_FEED",
			"targetEntities":                 []interface{}{},
			"thirdPartyDistributionChannels": []interface{}{},
		},
		"content": map[string]interface{}{
			"media": map[string]interface{}{
				"title": title,
				"id":    videoUrn, // e.g., "urn:li:video:C5F10AQGKQg_6y2a4sQ"
			},
		},
		"lifecycleState":            "PUBLISHED",
		"isReshareDisabledByAuthor": false,
	}

	// Convert post data to JSON
	postDataJSON, err := json.Marshal(postData)
	if err != nil {
		return fmt.Errorf("error marshalling post data: %v", err)
	}

	// Create the HTTP request
	url := "https://api.linkedin.com/rest/posts"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postDataJSON))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")
	req.Header.Set("LinkedIn-Version", "202502") // Use the appropriate API version
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making the request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	// Log the response for debugging
	log.Printf("Response Status: %d\n", resp.StatusCode)
	log.Printf("Response Body: %s\n", string(body))

	// Check if the status is not OK
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error creating post: %d %s", resp.StatusCode, string(body))
	}

	// Get the post ID from the response header
	postID := resp.Header.Get("x-restli-id")
	log.Printf("Post created successfully with ID: %s", postID)

	return nil
}
