package model

import "time"

type IPFSResult struct {
	Hash   string `json:"Hash,omitempty"`
	Name   string `json:"Name,omitempty"`
	Size   string `json:"Size,omitempty"`
	AwsURL string `json:"awsUrl,omitempty"`
}

type Backups struct {
	ID           int64     `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
	Filehash     string    `json:"filehash,omitempty" validate:"required"`
	UserId       int64     `json:"user_id,omitempty" validate:"required"`
	BackupNature string    `json:"backup_nature,omitempty" validate:"required"`
	FileName     string    `json:"file_name,omitempty" validate:"required"`
	FileSize     string    `json:"file_size,omitempty" validate:"required"`
}
