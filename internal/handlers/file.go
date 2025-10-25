package handlers

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/okoye-dev/oss-archive/internal/storage"
	"github.com/okoye-dev/oss-archive/internal/transport/rest"
)

type FilesResponse struct {
	Files []FileResponse `json:"files"`
}

type FileResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	StorageKey string `json:"storage_key"`
	Size       int64  `json:"size"`
}

type FileDownloadResponse struct {
	URL        string `json:"url"`
	ExpiresIn  int    `json:"expires_in"`
	Download   bool   `json:"download"`
}

type FileHandler struct {
	storage storage.StorageInterface
}

func NewFileHandler(storage storage.StorageInterface) *FileHandler {
	return &FileHandler{
		storage: storage,
	}
}

func (h *FileHandler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		rest.BadRequest(c, "No file provided")
		return
	}
	defer file.Close()

	// Generate unique ID and storage key
	fileID := uuid.New().String()
	storageKey := fmt.Sprintf("%s_%s", fileID, header.Filename)
	
	// Get content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload to storage using storage key
	err = h.storage.UploadFile(storageKey, file, header.Size, contentType)
	if err != nil {
		rest.InternalError(c, err)
		return
	}

	// Return file info
	fileData := map[string]interface{}{
		"id":          fileID,
		"name":        header.Filename,
		"storage_key": storageKey,
		"file_size":   header.Size,
		"file_type":   contentType,
		"created_at":  time.Now().Format(time.RFC3339),
		"updated_at":  time.Now().Format(time.RFC3339),
	}

	rest.Success(c, fileData)
}


func (h *FileHandler) GetFiles(c *gin.Context) {
	files, err := h.storage.ListFiles()
	if err != nil {
		rest.InternalError(c, err)
		return
	}

	var fileList []FileResponse
	for _, storageKey := range files {
		// Extract ID and name from storage key (format: "timestamp_originalname.ext")
		parts := strings.SplitN(storageKey, "_", 2)
		var fileID, fileName string
		
		if len(parts) == 2 {
			fileID = parts[0]
			fileName = parts[1]
		} else {
			// Fallback for files without timestamp prefix
			fileID = storageKey
			fileName = filepath.Base(storageKey)
		}
		
		// Get file size
		fileSize, err := h.storage.GetFileSize(storageKey)
		if err != nil {
			log.Printf("Failed to get file size for %s: %v", storageKey, err)
			fileSize = 0
		}
		
		fileList = append(fileList, FileResponse{
			ID:         fileID,
			Name:       fileName,
			StorageKey: storageKey,
			Size:       fileSize,
		})
	}

	rest.Success(c, FilesResponse{
		Files: fileList,
	})
}

func (h *FileHandler) GetFile(c *gin.Context) {
	filename := c.Param("id")
	if filename == "" {
		rest.BadRequest(c, "File ID required")
		return
	}

	forceDownload := c.Query("download") == "true"

	// Generate presigned URL
	presignedURL, err := h.storage.GetPresignedURL(filename, forceDownload)
	if err != nil {
		rest.NotFound(c, "File not found")
		return
	}

	rest.Success(c, FileDownloadResponse{
		URL:        presignedURL,
		Download:   forceDownload,
		ExpiresIn:  3600,
	})
}

func (h *FileHandler) DeleteFile(c *gin.Context) {
	filename := c.Param("id")
	if filename == "" {
		rest.BadRequest(c, "File ID required")
		return
	}

	err := h.storage.DeleteFile(filename)
	if err != nil {
		rest.InternalError(c, err)
		return
	}

	rest.Success(c, FileResponse{
		ID:         filename,
		Name:       filename,
		StorageKey: filename,
		Size:       0,
	})
}