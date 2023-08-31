package sdio

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/OverlyDev/go-sdio-seeder/internal/download"
	"github.com/OverlyDev/go-sdio-seeder/internal/logger"
	"github.com/OverlyDev/go-sdio-seeder/internal/torrent"
	"github.com/OverlyDev/go-sdio-seeder/internal/util"
)

// Credit: https://www.technibble.com/forums/threads/snappy-driver-installer.60557/post-477761
const torrentUrl = "http://driveroff.net/SDI_Update.torrent"

// Creates new SdioHelper struct with default values
func NewHelper() SdioHelper {
	helper := SdioHelper{
		Url:     torrentUrl,
		DataDir: "data",
	}

	if err := helper.Setup(); err != nil {
		logger.ErrorLogger.Println(err)
		os.Exit(1)
	}
	return helper
}

// Contains variables and functions related to SDIO
type SdioHelper struct {
	Url           string // Url for driverpacks torrent file
	DataDir       string // Folder holding torrent files
	CurrentFile   string // Most recent torrent file (hashed)
	ActiveTorrent string // ID of most recent torrent
}

// Performs initial setup (just makes the data dir for now)
func (s *SdioHelper) Setup() error {
	// Create data dir
	err := os.MkdirAll(s.DataDir, 0750)
	if err != nil {
		return err
	}

	return nil
}

// Starts the show
func (s *SdioHelper) Start() error {
	// Download torrent file
	file, err := s.getTorrentFile()
	if err != nil {
		return err
	}
	logger.InfoLogger.Println("Downloaded:", file)

	// Do checksum of torrent file
	hash, err := s.hashTorrentFile(file)
	if err != nil {
		return err
	}
	logger.InfoLogger.Println("File:", file, "Hash:", hash)

	// Set current file
	if err := s.setCurrentFile(hash); err != nil {
		return err
	}
	logger.InfoLogger.Println("Current file:", s.CurrentFile)

	// Start downloading/seeding
	s.downloadTorrent()
	logger.InfoLogger.Println("Active torrent:", s.ActiveTorrent)

	return nil

}

func (s *SdioHelper) Update() {
	// Download new torrent file
	new, err := s.getTorrentFile()
	if err != nil {
		logger.ErrorLogger.Println(err)
		return
	}

	// Read in new torrent file
	data, err := util.ReadFile(filepath.Join(s.DataDir, new))
	if err != nil {
		logger.ErrorLogger.Println(err)
		return
	}

	// Checksum new torrent file
	sum, err := util.ChecksumBytes(data)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return
	}

	// Current file is up to date, no change
	if sum == s.CurrentFile {
		logger.DebugLogger.Println("Current file is up to date")
		os.Remove(filepath.Join(s.DataDir, new))
		return
	}

	// Got a new torrent file
	logger.DebugLogger.Println("Got new torrent file")

	// Stop active torrent
	if err := torrent.StopTorrent(s.ActiveTorrent); err != nil {
		logger.ErrorLogger.Println(err)
		return
	}

	// Delete old hashed torrent file
	os.Remove(filepath.Join(s.DataDir, s.CurrentFile))

	// Do checksum of torrent file
	hash, err := s.hashTorrentFile(new)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return
	}
	logger.DebugLogger.Println("File:", new, "Hash:", hash)

	// Set current file
	if err := s.setCurrentFile(hash); err != nil {
		logger.ErrorLogger.Println(err)
		return
	}
	logger.InfoLogger.Println("New current file:", s.CurrentFile)

	// Start downloading/seeding
	s.downloadTorrent()
	logger.InfoLogger.Println("New active torrent:", s.ActiveTorrent)
}

// Prints status of the active torrent
func (s *SdioHelper) Status() {
	status, err := torrent.CheckStatus(s.ActiveTorrent)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return
	}
	logger.InfoLogger.Println(status)
}

// Updates SdioHelper instance's CurrentFile field
func (s *SdioHelper) setCurrentFile(filename string) error {
	path := filepath.Join(s.DataDir, filename)
	if _, err := os.Stat(path); err != nil {
		return err
	}
	s.CurrentFile = filename

	return nil
}

// Returns downloaded file name on success, otherwise empty string and error
func (s *SdioHelper) getTorrentFile() (string, error) {
	file, err := download.DownloadFile(s.Url, s.DataDir)
	if err != nil {
		return "", nil
	}
	logger.DebugLogger.Printf("Got torrent file: %s Url: %s\n", file, s.Url)
	return file, nil
}

// Returns new filename on success, otherwise empty string and error
func (s *SdioHelper) hashTorrentFile(filename string) (string, error) {
	// Read file into buffer
	buf, fErr := util.ReadFile(filepath.Join(s.DataDir, filename))
	if fErr != nil {
		return "", fErr
	}

	// Create reader for buffer
	r := bytes.NewReader(buf)

	// Get checksum
	hash, rErr := util.ChecksumReader(r)
	if rErr != nil {
		return "", rErr
	}

	// Create new file with filename of hash string
	out, oErr := os.Create(filepath.Join(s.DataDir, hash))
	if oErr != nil {
		return "", oErr
	}
	defer out.Close()

	// Copy original bytes from buffer to new file
	_, cErr := io.Copy(out, r)
	if cErr != nil {
		return "", cErr
	}

	// Delete original file
	os.Remove(filepath.Join(s.DataDir, filename))

	return hash, nil
}

// Starts download of s.CurrentFile and saves the ID in s.ActiveTorrent
func (s *SdioHelper) downloadTorrent() error {
	torrentId, err := torrent.Download(filepath.Join(s.DataDir, s.CurrentFile))
	if err != nil {
		return err
	}
	s.ActiveTorrent = torrentId

	return nil
}
