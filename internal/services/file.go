package services

import "github.com/okoye-dev/oss-archive/internal/models"

func GetFiles() ([]models.File, error) {
	return []models.File{
		{
			ID: "1",
			FileName: "File 1",
			FilePath: "path/to/file1.txt",
			FileSize: 100,
			FileType: "text/plain",
		},
	}, nil
}