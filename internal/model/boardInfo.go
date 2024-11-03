package model

import "time"

type BoardInfo struct {
	PlanInfo       PlanInfo       `json:"planInfo"`
	RepositoryInfo RepositoryInfo `json:"repositoryInfo"`
	BackupInfo     BackupInfo     `json:"backupInfo"`  // Summary of backup data details.
	BackupInfos    []BackupInfo   `json:"backupInfos"` // Detailed backup data for each repository.
}

// PlanInfo contains schedule information.
type PlanInfo struct {
	Total        int `json:"total"`        // Total number of backup plans.
	RunningCount int `json:"runningCount"` // Number of backup plans currently running.
}

// RepositoryInfo contains repository information.
type RepositoryInfo struct {
	Total        int `json:"total"`        // Total number of repositories.
	RunningCount int `json:"runningCount"` // Number of repositories currently running.
}

// BackupInfo contains backup information.
type BackupInfo struct {
	RepositoryName           string    `json:"repositoryName"`           // Repository name
	SnapshotsNum             int       `json:"snapshotsNum"`             // Number of snapshots
	FileTotal                int       `json:"fileTotal"`                // Total number of files
	DataSize                 uint64    `json:"dataSize"`                 // Total data size
	DataSizeStr              string    `json:"dataSizeStr"`              // Total data size (string format)
	DataDay                  string    `json:"dataDay"`                  // Data retention days
	Time                     time.Time `json:"time"`                     // Statistics time
	Duration                 string    `json:"duration"`                 // Statistics duration
	CompressionSpaceSaving   string    `json:"compressionSpaceSaving"`   // Compression ratio
	TotalUncompressedSize    uint64    `json:"totalUncompressedSize"`    // Uncompressed data size
	TotalUncompressedSizeStr string    `json:"TotalUncompressedSizeStr"` // Uncompressed data size (string format)
}

