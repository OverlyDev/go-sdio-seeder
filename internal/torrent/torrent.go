package torrent

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/OverlyDev/go-sdio-seeder/internal/logger"
	"github.com/cenkalti/rain/torrent"
)

const downloadDir = "downloads"

var config torrent.Config
var session *torrent.Session

// Creates download dir, config, and session
func Setup() {
	// Create download dir
	if err := os.MkdirAll(downloadDir, 0750); err != nil {
		logger.ErrorLogger.Println(err)
		os.Exit(1)
	}

	torrent.DisableLogging()

	config = torrent.DefaultConfig
	config.DataDir = downloadDir
	config.Database = filepath.Join(downloadDir, "session.db")
	config.RPCEnabled = false

	session, _ = torrent.NewSession(config)
}

// Returns torrent ID
func Download(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}

	t, _ := session.AddTorrent(f, nil)

	return t.ID(), nil
}

// Returns a string of the current status of a torrent by ID
func CheckStatus(torrentId string) (string, error) {
	t := session.GetTorrent(torrentId)
	if t == nil {
		return "", errors.New(fmt.Sprintf("Failed to get torrent: %s", torrentId))
	}

	stats := t.Stats()

	var msg string
	switch stats.Status.String() {
	case "Downloading":
		msg = fmt.Sprintf("Status: %s\tPeers: (In %d | Out %d | Total %d)\tETA: %s", stats.Status, stats.Peers.Incoming, stats.Peers.Outgoing, stats.Peers.Total, stats.ETA)
	case "Seeding":
		msg = fmt.Sprintf("Status: %s (%d)\tPeers: (In %d | Out %d | Total %d)", stats.Status, stats.SeededFor, stats.Peers.Incoming, stats.Peers.Outgoing, stats.Peers.Total)
	case "Stopping":
		msg = "Stopping (new torrent will start soon)"
	case "Stopped":
		msg = "Stopped (new torrent will start soon)"
	default:
		msg = "Some other status (check debug logs)"
		logger.DebugLogger.Printf("%+v\n", stats)
	}

	return msg, nil
}

// Stops a torrent by ID and removes its files
func StopTorrent(torrentId string) error {
	// Get torrent
	t := session.GetTorrent(torrentId)
	if t == nil {
		return errors.New(fmt.Sprintf("Failed to get torrent: %s", torrentId))
	}

	// Stop torrent
	logger.DebugLogger.Println("Stopping torrent:", torrentId)
	t.Stop()
	logger.DebugLogger.Println("Waiting 5s for torrent to stop")
	time.Sleep(5 * time.Second)

	// Delete torrent
	if err := session.RemoveTorrent(torrentId); err != nil {
		return err
	}
	logger.DebugLogger.Println("Removed torrent:", torrentId)

	return nil
}
