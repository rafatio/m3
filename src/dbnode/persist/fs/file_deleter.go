// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package fs

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	xunsafe "github.com/m3db/m3/src/x/unsafe"
)

// FileDeleter deletes files.
type FileDeleter interface {
	Delete(pattern string) error
}

// EfficientFileDeleter holds reusable memory for filenames to be deleted. This allows
// batch file deletions to be performedf with minimal memory usage.
type EfficientFileDeleter struct {
	fileNames []byte
}

// SimpleFileDeleter deletes files using native file deletion code.
type SimpleFileDeleter struct {
}

// NewEfficientFileDeleter returns a new deleter that will avoid allocating unnecessary memory. Specifically
// it batch processes the matching of the filenames to be deleted in a fixed buffer with the provided filenameBatchByteSize.
// It then allocates memory only for the matched string filenames to be deleted, rather than all files in a given directory.
func NewEfficientFileDeleter(filenameBatchByteSize int) FileDeleter {
	return &EfficientFileDeleter{
		fileNames: make([]byte, maxMemoryInBytes),
	}
}

// NewSimpleFileDeleter returns a new deleter that performs standard native GO file deletion.
func NewSimpleFileDeleter() FileDeleter {
	return &SimpleFileDeleter{}
}

// Delete deletes all of the files that match the pattern.
func (d *SimpleFileDeleter) Delete(pattern string) error {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	DeleteFiles(files)
	return nil
}

// Delete deletes all of the files that match the pattern.
func (d *EfficientFileDeleter) Delete(pattern string) error {
	dir, file := filepath.Split(pattern)
	openDir, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer openDir.Close()

	var namesToDelete []string
	for i := 0; ; i++ {
		n, err := syscall.ReadDirent(int(openDir.Fd()), d.fileNames)
		if err != nil {
			return err
		}
		namesToDelete = append(namesToDelete, matchFilesToDelete(dir, d.fileNames[:n], file)...)
		if n <= 0 {
			break
		}
	}

	for _, nameToDelete := range namesToDelete {
		// Even if we hit an err, we continue deletion attempts in case it is file specific.
		err = os.Remove(nameToDelete)
	}
	return err
}

func matchFilesToDelete(dir string, fileNames []byte, filePattern string) []string {
	var namesToDelete []string
	for len(fileNames) > 0 {
		reclen, ok := direntReclen(fileNames)
		if !ok || reclen > uint64(len(fileNames)) {
			return namesToDelete
		}

		rec := fileNames[:reclen]
		fileNames = fileNames[reclen:]

		ino, ok := direntIno(rec)
		if !ok {
			break
		}
		if ino == 0 { // File absent in directory.
			continue
		}
		const namoff = uint64(unsafe.Offsetof(syscall.Dirent{}.Name))
		namlen, ok := direntNamlen(rec)
		if !ok || namoff+namlen > uint64(len(rec)) {
			break
		}
		name := rec[namoff : namoff+namlen]
		for i, c := range name {
			if c == 0 {
				name = name[:i]
				break
			}
		}

		if bytes.Equal(name, dotBytes) || bytes.Equal(name, doubleDotBytes) {
			// Check for useless names before allocating a string.
			continue
		}

		// TODO: consider implementing the wildcard pattern matching on the
		// raw bytes to avoid this yolo cast to call the Match func.
		matched := false
		xunsafe.WithString(name, func(s string) {
			matched, _ = filepath.Match(filePattern, s)
		})
		if !matched {
			continue
		}

		// Use only []byte instead of string until we know for sure we will delete the filename.
		// That way we avoid unnecessary string allocation for filenames we are not going to delete.
		var pathToDelete strings.Builder
		pathToDelete.WriteString(dir)
		pathToDelete.Write(name)
		namesToDelete = append(namesToDelete, pathToDelete.String())
	}
	return namesToDelete
}

func direntReclen(buf []byte) (uint64, bool) {
	return readInt(buf, unsafe.Offsetof(syscall.Dirent{}.Reclen), unsafe.Sizeof(syscall.Dirent{}.Reclen))
}

func direntIno(buf []byte) (uint64, bool) {
	return readInt(buf, unsafe.Offsetof(syscall.Dirent{}.Ino), unsafe.Sizeof(syscall.Dirent{}.Ino))
}

func direntNamlen(buf []byte) (uint64, bool) {
	reclen, ok := direntReclen(buf)
	if !ok {
		return 0, false
	}
	return reclen - uint64(unsafe.Offsetof(syscall.Dirent{}.Name)), true
}

// readInt returns the size-bytes unsigned integer in native byte order at offset off.
func readInt(b []byte, off, size uintptr) (u uint64, ok bool) {
	if len(b) < int(off+size) {
		return 0, false
	}
	return readIntLE(b[off:], size), true
}

func readIntLE(b []byte, size uintptr) uint64 {
	switch size {
	case 1:
		return uint64(b[0])
	case 2:
		_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
		return uint64(b[0]) | uint64(b[1])<<8
	case 4:
		_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
		return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24
	case 8:
		_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
		return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
			uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
	default:
		panic("syscall: readInt with unsupported size")
	}
}