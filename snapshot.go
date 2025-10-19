package owl

import (
	"errors"
	"os"
	"time"
)

// SnapShot is a snap shot of the file meta data
type SnapShot struct {
	Path    string
	Exists  bool
	Size    int64
	ModTime time.Time
	INO     uint64
	mapid   string
	DEV     uint64
}

// ReadSnap is a snap shot of the contents of the file
type ReadSnap struct {
	Path    string
	ModTime time.Time
	Content []byte
}

// takesnap takes a snap shot of the file state
func takesnap(path string) (SnapShot, error) {
	snap := SnapShot{Path: path}

	info, e := os.Stat(path)
	if e != nil {
		if errors.Is(e, os.ErrNotExist) {
			return snap, nil
		}
		return snap, e
	}
	snap.Exists = true
	snap.Size = info.Size()
	snap.ModTime = info.ModTime()

	fillSysInfo(&snap, info.Sys())

	return snap, nil
}
