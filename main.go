package main

import (
	"Linkedin-Video-Posting/linkedin"
	"Linkedin-Video-Posting/utils"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	accessToken := os.Getenv("AccessToken")                //add your Access Token in .env
	owner := "urn:li:organization:<your_organization_urn>" // replace with your oraganization URN

	filePath := "media/temp.mp4"

	fileSize, err := utils.GetFileSize(filePath) //Getting File Size
	if err != nil {
		fmt.Println("Failed to fetch file size: ", err)
		return
	}
	title := "Test post title"
	caption := "test Caption"

	//*********** STEP-1 Initializing Video **************//
	fmt.Println("::::::::::::::::STEP-1 Initializing the Video:::::::::::::::::")
	fmt.Println("")

	response, err := linkedin.InitializeLinkedInVideoUpload(accessToken, owner, int(fileSize))
	if err != nil {
		fmt.Println("Error In Initialize:", err)
		return
	}
	fmt.Printf("VIDEO URN: %s\n", response.Value.Video)
	VideoURN := response.Value.Video
	fmt.Println("")
	fmt.Printf("Upload URL: %s\n", response.Value.UploadInstructions[0].UploadURL)
	fmt.Println("")

	//************ STEP-2 Uploading the video *************//
	fmt.Println("::::::::::::::STEP-2 Uploading the video:::::::::::::::::::::")
	fmt.Println("")

	uploadInstructions := response.Value.UploadInstructions
	etags, err := linkedin.UploadVideoFile(filePath, uploadInstructions, accessToken)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Video Parts: %s\n", etags)
	fmt.Println("")

	//************ STEP-3 Finalizing the video *************//

	fmt.Println(":::::::::::::::STEP-3 Finalizing the video:::::::::::::::::::")
	fmt.Println("")
	err = linkedin.FinalizeVideoUpload(accessToken, VideoURN, etags)
	if err != nil {
		fmt.Println("Error in Finalize:", err)
		return
	}
	fmt.Println("")

	//*********** STEP-4 Publishing the video ************//
	fmt.Println(":::::::::::::: STEP-4 Publishing the video:::::::::::::::::::")
	fmt.Println("")
	err = linkedin.CreatePost(accessToken, VideoURN, owner, title, caption)
	if err != nil {
		fmt.Println("Error in Creating the post: ", err)
		return
	}
	fmt.Println("")
	fmt.Println("Video SuccessFully Posted !")

}
