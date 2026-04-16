package logs

import (
	"io"
	"os"
	"strings"
)

const readBuf = 32 * 1024 // 32 KB look-back window

// TailN returns the last n lines of the file at path.
// Returns an empty slice (no error) when the file does not exist yet.
func TailN(path string, n int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	defer f.Close()

	size, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}

	bufSize := int64(readBuf)
	if size < bufSize {
		bufSize = size
	}

	if _, err := f.Seek(-bufSize, io.SeekEnd); err != nil {
		return nil, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	raw := strings.Split(string(data), "\n")
	// Drop a leading partial line when we seeked into the middle of the file
	if bufSize < size && len(raw) > 1 {
		raw = raw[1:]
	}
	// Remove trailing empty entry from a file ending in \n
	if len(raw) > 0 && raw[len(raw)-1] == "" {
		raw = raw[:len(raw)-1]
	}

	if len(raw) <= n {
		return raw, nil
	}
	return raw[len(raw)-n:], nil
}
