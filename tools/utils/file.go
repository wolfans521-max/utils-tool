package utils

import (
	"os"
	"syscall"
)

// GetInode 获取文件的 inode
func GetInode(f *os.File) uint64 {
	if f == nil {
		return 0
	}
	stat, err := f.Stat()
	if err != nil {
		return 0
	}
	if sysStat, ok := stat.Sys().(*syscall.Stat_t); ok {
		return sysStat.Ino
	}
	return 0
}

// GetPathInode 获取路径的 inode
func GetPathInode(path string) uint64 {
	stat, err := os.Stat(path)
	if err != nil {
		return 0
	}
	if sysStat, ok := stat.Sys().(*syscall.Stat_t); ok {
		return sysStat.Ino
	}
	return 0
}
