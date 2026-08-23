//go:build windows

package installer

import (
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

// fileRenameInformation mirrors the NT FILE_RENAME_INFORMATION layout consumed by
// NtSetInformationFile. Its leading field carries either the BOOLEAN replace flag or, on Windows 10
// and later, the FILE_RENAME_* flag set. Go aligns RootDirectory (a pointer-sized handle) to its
// natural boundary, so the field offsets match the C ABI on both 32- and 64-bit Windows.
type fileRenameInformation struct {
	ReplaceIfExists uint32
	RootDirectory   windows.Handle
	FileNameLength  uint32
	FileName        [1]uint16
}

// rootedRename atomically renames oldName to newName within the directory anchored by directoryRoot.
//
// Windows has no directory-handle-relative rename in the Win32 surface, so the rename uses NT
// native APIs. The staged file is reopened relative to the anchored directory handle with DELETE
// access (NtCreateFile), then renamed in place (NtSetInformationFile with FileRenameInformation),
// naming the destination relative to the same handle. Every name resolves against the pinned
// directory descriptor, so a parent swap performed after anchoring cannot redirect the rename
// outside the repository. Go 1.24 os.Root lacks Rename, so there is no portable os.Root operation.
func rootedRename(
	directoryRoot *os.Root,
	_ string,
	oldName string,
	newName string,
) error {
	dir, err := directoryRoot.Open(".")
	if err != nil {
		return err
	}
	defer func() {
		_ = dir.Close() // closing a read-only directory cannot change the completed rename
	}()
	dirHandle := windows.Handle(dir.Fd())

	objectName, err := windows.NewNTUnicodeString(oldName)
	if err != nil {
		return err
	}
	attributes := &windows.OBJECT_ATTRIBUTES{
		RootDirectory: dirHandle,
		ObjectName:    objectName,
	}
	attributes.Length = uint32(unsafe.Sizeof(*attributes))

	var ioStatus windows.IO_STATUS_BLOCK
	var staged windows.Handle
	if err = windows.NtCreateFile(
		&staged,
		windows.DELETE,
		attributes,
		&ioStatus,
		nil,
		0,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		windows.FILE_OPEN,
		0,
		0,
		0,
	); err != nil {
		return err
	}
	defer windows.CloseHandle(staged)

	newNameUTF16, err := windows.UTF16FromString(newName)
	if err != nil {
		return err
	}
	nameBytes := (len(newNameUTF16) - 1) * 2

	var info fileRenameInformation
	infoSize := int(unsafe.Offsetof(info.FileName)) + nameBytes
	buffer := make([]byte, infoSize)
	header := (*fileRenameInformation)(unsafe.Pointer(&buffer[0]))
	header.ReplaceIfExists = windows.FILE_RENAME_REPLACE_IF_EXISTS
	header.RootDirectory = dirHandle
	header.FileNameLength = uint32(nameBytes)
	copy(
		(*[windows.MAX_LONG_PATH]uint16)(unsafe.Pointer(&header.FileName[0]))[:nameBytes/2:nameBytes/2],
		newNameUTF16,
	)

	return windows.NtSetInformationFile(
		staged,
		&ioStatus,
		&buffer[0],
		uint32(infoSize),
		windows.FileRenameInformation,
	)
}
