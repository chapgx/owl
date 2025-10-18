//go:build linux

package owl

import "syscall"

func fillSysInfo(s *SnapShot, info interface{}) {
	if stat, ok := info.(*syscall.Stat_t); ok {
		s.INO = stat.Ino
		s.DEV = stat.Dev
	}
}
