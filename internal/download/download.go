package download

import (
	"io"
	"os"
	"path/filepath"

	"github.com/OverlyDev/go-sdio-seeder/internal/logger"
)

// Returns downloaded file name on success, otherwise empty string and error
func DownloadFile(url string, destDir string) (string, error) {
	// Get filename from url
	filename := filepath.Base(url)

	// Make request
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create file handle
	out, err := os.Create(filepath.Join(destDir, filename))
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Write to file
	size, err := io.Copy(out, resp.Body)
	if err != nil {
		logger.DebugLogger.Printf("Incomplete: %s Size: %d\n", filename, size)
		return "", err
	}

	return filename, nil

}
