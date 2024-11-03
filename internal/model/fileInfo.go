package model

type FileInfo struct {
	IsDir      bool   `json:"isDir"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	Mode       string `json:"mode"`       // Permissions
	ModTime    string `json:"modTime"`    // Modification time
	CreateTime string `json:"createTime"` // Creation time
	Size       int64  `json:"size"`       // Size
	Gid        int    `json:"gid"`        // Group ID
	Uid        int    `json:"uid"`        // User ID
}
