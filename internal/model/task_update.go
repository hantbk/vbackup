package model

// StatusUpdate Progress
type StatusUpdate struct {
	MessageType      string   `json:"messageType"`        // "status"
	SecondsElapsed   string   `json:"secondsElapsed"`     // Elapsed time
	SecondsRemaining string   `json:"secondsRemaining"`   // Remaining time
	PercentDone      float64  `json:"percentDone"`        // Progress
	TotalFiles       uint64   `json:"totalFiles"`         // Total number of files
	FilesDone        uint64   `json:"filesDone"`          // Number of completed files
	TotalBytes       string   `json:"totalBytes"`         // Total size of files
	BytesDone        string   `json:"bytesDone"`          // Size of completed files
	ErrorCount       uint     `json:"errorCount"`         // Number of errors
	CurrentFiles     []string `json:"currentFiles"`       // List of current files
	AvgSpeed         string   `json:"avgSpeed,omitempty"` // Average speed
}

// ErrorUpdate Error
type ErrorUpdate struct {
	MessageType string `json:"messageType"` // "error"
	Error       string `json:"error"`       // Error message
	During      string `json:"during"`      // Type of error
	Item        string `json:"item"`        // Target of the error
}

// VerboseUpdate Completion Item
type VerboseUpdate struct {
	MessageType  string `json:"messageType"` // "verbose_status"
	Action       string `json:"action"`
	Item         string `json:"item"`
	Duration     string `json:"duration"`     // Duration in seconds
	DataSize     string `json:"dataSize"`     // Amount of data
	MetadataSize string `json:"metadataSize"` // Amount of metadata
	TotalFiles   uint   `json:"totalFiles"`   // Total number of files
}

// SummaryOutput Backup Completion Summary
type SummaryOutput struct {
	MessageType         string `json:"messageType"` // "summary"
	FilesNew            uint   `json:"filesNew"`
	FilesChanged        uint   `json:"filesChanged"`
	FilesUnmodified     uint   `json:"filesUnmodified"`
	DirsNew             uint   `json:"dirsNew"`
	DirsChanged         uint   `json:"dirsChanged"`
	DirsUnmodified      uint   `json:"dirsUnmodified"`
	DataBlobs           int    `json:"dataBlobs"`
	TreeBlobs           int    `json:"treeBlobs"`
	DataAdded           string `json:"dataAdded"`            // Size of newly added files
	TotalFilesProcessed uint   `json:"totalFiles_processed"` // Total number of files
	TotalBytesProcessed string `json:"totalBytes_processed"` // Total size of files
	TotalDuration       string `json:"totalDuration"`        // Total duration in seconds
	SnapshotID          string `json:"snapshotId"`
	DryRun              bool   `json:"dryRun,omitempty"` // Dry run flag
}
