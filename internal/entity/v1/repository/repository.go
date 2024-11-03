package repository

import "github.com/hantbk/vbackup/internal/entity/v1/common"

// Repository type
const (
	S3    = 1
	Sftp  = 3
	Local = 4
	Rest  = 5
)

const (
	StatusNone = 1
	StatusRun  = 2
	StatusErr  = 3
)

// Repository struct for backup storage, inspired by Restic
type Repository struct {
	common.BaseModel `storm:"inline"`
	Name             string `json:"name"`               // Repository name
	Type             int    `json:"type"`               // Repository type identifier
	Endpoint         string `json:"endPoint"`           // Endpoint URL for the repository

	// AWS-specific fields
	Region string `json:"region"`   // AWS default region
	Bucket string `json:"bucket"`   // S3 bucket name
	KeyId  string `json:"keyId"`    // AWS access key ID
	Secret string `json:"secret"`   // AWS secret access key

	// Google Cloud-specific fields
	ProjectID string `json:"projectId"` // Google Cloud project ID

	// Azure-specific fields
	AccountName string `json:"accountName"` // Azure account name
	AccountKey  string `json:"accountKey"`  // Azure account key, B2 account key

	// Backblaze B2-specific fields
	AccountID string `json:"accountId"` // B2 account ID

	Password          string `json:"password"`          // Repository password
	Status            int    `json:"status"`            // Status of the repository connection
	Errmsg            string `json:"errmsg"`            // Error message for failed operations
	RepositoryVersion string `json:"repositoryVersion"` // Version of the repository format
	Compression       int    `json:"compression"`       // Compression mode: auto(0), off(1), max(2)
	PackSize          int    `json:"packSize"`          // Suggested pack size for repository data
}
