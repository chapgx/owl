//go:build darwin

package owl

import (
	"syscall"
)

func fillSysInfo(s *SnapShot, info interface{}) {
	if stat, ok := info.(*syscall.Stat_t); ok {
		s.INO = uint64(stat.Ino)
		s.DEV = uint64(stat.Dev)
	}
}
