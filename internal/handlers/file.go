package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/okoye-dev/oss-archive/internal/storage"
	"github.com/okoye-dev/oss-archive/internal/transport/rest"
)

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

	// Generate unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
	
	// Get content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload to storage
	err = h.storage.UploadFile(filename, file, header.Size, contentType)
	if err != nil {
		rest.InternalError(c, err)
		return
	}

	// Return file info
	fileData := map[string]interface{}{
		"id":         filename,
		"file_name":  header.Filename,
		"file_path":  filename,
		"file_size":  header.Size,
		"file_type":  contentType,
		"created_at": time.Now().Format(time.RFC3339),
		"updated_at": time.Now().Format(time.RFC3339),
	}

	rest.Success(c, fileData)
}

func (h *FileHandler) GetFiles(c *gin.Context) {
	files, err := h.storage.ListFiles()
	if err != nil {
		rest.InternalError(c, err)
		return
	}

	var fileList []map[string]interface{}
	for _, filename := range files {
		fileList = append(fileList, map[string]interface{}{
			"id":        filename,
			"file_name": filepath.Base(filename),
			"file_path": filename,
		})
	}

	rest.Success(c, fileList)
}

func (h *FileHandler) GetFile(c *gin.Context) {
	filename := c.Param("id")
	if filename == "" {
		rest.BadRequest(c, "File ID required")
		return
	}

	reader, err := h.storage.GetFile(filename)
	if err != nil {
		rest.NotFound(c, "File not found")
		return
	}
	defer reader.Close()

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(filename)))
	c.DataFromReader(http.StatusOK, -1, "application/octet-stream", reader, nil)
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

	rest.Success(c, map[string]string{"message": "File deleted successfully"})
}