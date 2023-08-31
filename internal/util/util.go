package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
)

// Returns a byte slice of the contents of file at filepath
func ReadFile(filepath string) ([]byte, error) {
	// Open file
	f, fErr := os.Open(filepath)
	if fErr != nil {
		return nil, fErr
	}
	defer f.Close()

	// Get file size for buffer
	s, sErr := f.Stat()
	if sErr != nil {
		return nil, sErr
	}

	// Create buffer the size of file
	buf := make([]byte, s.Size())

	// Read file into buffer
	_, rErr := f.Read(buf)
	if rErr != nil {
		return nil, rErr
	}

	return buf, nil
}

// Calculates SHA256 hash from data.
// Returns hash string on success, otherwise empty string and error.
func ChecksumBytes(data []byte) (string, error) {
	r := bytes.NewReader(data)

	hasher := sha256.New()
	_, hErr := r.WriteTo(hasher)
	if hErr != nil {
		return "", hErr
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// Calculates SHA256 hash from reader.
// Rewinds reader upon finish.
// Returns hash string on success, otherwise empty string and error.
func ChecksumReader(reader *bytes.Reader) (string, error) {
	hasher := sha256.New()
	_, hErr := reader.WriteTo(hasher)
	if hErr != nil {
		return "", hErr
	}

	reader.Seek(0, 0)

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
