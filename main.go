package main

import (
	"net/http"
	"strings"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/api/admin/search"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
)

func uploadFile(c *gin.Context) {

	// Create our instance
	cld, _ := cloudinary.NewFromURL("cloudinary://775819266531562:0vvHyxdwaqAHBjT44neW6h572HY@djmi67ke1")

	// Get the preferred name of the file if its not supplied
	fileName := c.PostForm("name")

	// Add tags
	fileTags := c.PostForm("tags")
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "Failed to upload",
		})
	}

	result, err := cld.Upload.Upload(c, file, uploader.UploadParams{
		PublicID: fileName,
		// Split the tags by comma's
		Tags: strings.Split(",", fileTags),
	})

	if err != nil {
		c.String(http.StatusConflict, "Upload to cloudinary failed")
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Successfully uploaded the file",
		"secureURL": result.SecureURL,
		"publicURL": result.URL,
	})
}

func getAllUploadedAssets(c *gin.Context) {
	cld, _ := cloudinary.NewFromURL("cloudinary://775819266531562:0vvHyxdwaqAHBjT44neW6h572HY@djmi67ke1")
	var urls []string

	searchQ := search.Query{
		Expression: "resource_type:image AND uploaded_at>1d AND bytes<1m",
		SortBy:     []search.SortByField{{"created_at": search.Descending}},
		MaxResults: 10,
	}

	results, err := cld.Admin.Search(c, searchQ)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   err,
			"message": "Failed to query your files",
		})
	}

	for _, asset := range results.Assets {
		urls = append(urls, asset.SecureURL)
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    urls,
	})
}

func getUploadedFile(c *gin.Context) {
	// Create our instance
	cld, _ := cloudinary.NewFromURL("cloudinary://775819266531562:0vvHyxdwaqAHBjT44neW6h572HY@djmi67ke1")
	fileName := c.Param("name")

	// Access the filename using a desired filename
	result, err := cld.Admin.Asset(c, admin.AssetParams{PublicID: fileName})
	if err != nil {
		c.String(http.StatusNotFound, "We were unable to find the file requested")
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":    "Successfully found the image",
		"secureURL":  result.SecureURL,
		"publicURL":  result.URL,
		"created_at": result.CreatedAt.String(),
	})
}

func updateFile(c *gin.Context) {
	cld, _ := cloudinary.NewFromURL("cloudinary://775819266531562:0vvHyxdwaqAHBjT44neW6h572HY@djmi67ke1")
	fileId := c.Param("publicId")
	newFileName := c.PostForm("fileName")

	// Access the filename using a desired filename
	result, err := cld.Upload.Rename(c, uploader.RenameParams{
		FromPublicID: fileId,
		ToPublicID:   newFileName,
	})
	if err != nil {
		c.String(http.StatusNotFound, "We were unable to find the file requested")
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":    "Successfully found the image",
		"secureURL":  result.SecureURL,
		"publicURL":  result.URL,
		"created_at": result.CreatedAt.String(),
	})
}

func deleteAsset(c *gin.Context) {
	cld, _ := cloudinary.NewFromURL("cloudinary://775819266531562:0vvHyxdwaqAHBjT44neW6h572HY@djmi67ke1")
	fileId := c.Param("assetId")
	result, err := cld.Upload.Destroy(c, uploader.DestroyParams{PublicID: fileId})

	if err != nil {
		c.String(http.StatusBadRequest, "File could not be deleted")
	}

	c.JSON(http.StatusContinue, result.Result)
}

func home(c *gin.Context) {
	c.String(200, "We are blessed")
}

func routes() {
	router := gin.Default()
	router.GET("/", home)
	router.POST("/upload", uploadFile)
	router.POST("/get-files", getAllUploadedAssets)
	router.GET("/get-upload/:assetId", getUploadedFile)
	router.PUT("/update-file/:name", updateFile)
	router.DELETE("/delete-file/:assetId", deleteAsset)

	router.Run(":1243")
}

func main() {
	routes()
}
