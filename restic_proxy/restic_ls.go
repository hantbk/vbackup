package resticProxy

import (
	"os"
	"time"

	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/restic"
)

type lsSnapshot struct {
	*restic.Snapshot
	ID         *restic.ID `json:"id"`
	ShortID    string     `json:"short_id"`
	StructType string     `json:"struct_type"` // "snapshot"
}

type lsNode struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Path        string      `json:"path"`
	UID         uint32      `json:"uid"`
	GID         uint32      `json:"gid"`
	Size        string      `json:"size,omitempty"`
	Mode        os.FileMode `json:"mode,omitempty"`
	Permissions string      `json:"permissions,omitempty"`
	ModTime     time.Time   `json:"mtime,omitempty"`
	AccessTime  time.Time   `json:"atime,omitempty"`
	ChangeTime  time.Time   `json:"ctime,omitempty"`
	StructType  string      `json:"struct_type"` // "node"

	size uint64 // Target for Size pointer.
}
