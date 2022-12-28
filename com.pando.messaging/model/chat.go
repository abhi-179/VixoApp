package model

import (
	"time"
)

type Status struct {
	ID         int64      `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt  time.Time  `json:"created_at,omitempty"`
	User_Id    int64      `json:"user_id,omitempty" validate:"required"`
	StatusType string     `json:"status_type,omitempty" validate:"required"`
	Message    string     `json:"message,omitempty"`
	File_hash  []FileHash `json:"file_hash,omitempty"`
	AwsUrl     []AwsUrl   `json:"aws_url,omitempty"`
}
type FileHash struct {
	File_hash string `json:"file_hash,omitempty"`
}
type AwsUrl struct {
	AwsUrl string `json:"aws_url,omitempty"`
}
type Statuses struct {
	ID         int64     `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	User_Id    int64     `json:"user_id,omitempty" validate:"required"`
	Username   string    `json:"username,omitempty"`
	StatusType string    `json:"status_type,omitempty" validate:"required"`
	Message    string    `json:"message,omitempty"`
	File_hash  string    `json:"file_hash,omitempty"`
	AwsUrl     string    `json:"aws_url,omitempty"`
}
type Reports struct {
	ID          int64     `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	Reporter_Id int64     `json:"reporter_id,omitempty"`
	Reportee_Id int64     `json:"reportee_id,omitempty"`
	GroupId     int64     `json:"group_id,omitempty"`
}

type Wallpapers struct {
	ID            int64     `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	WallpaperURL  string    `json:"wallpaper_url,omitempty"`
	WallpaperType string    `json:"wallpaper_type,omitempty"`
}
type Wallpaper struct {
	ID            int64   `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt     float64 `json:"created_at,omitempty"`
	WallpaperURL  string  `json:"wallpaper_url,omitempty"`
	WallpaperType string  `json:"wallpaper_type,omitempty"`
}
type BlockedContacts struct {
	ID        int64     `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	BlockeeId int64     `json:"blockee_id,omitempty"`
	BlockerId int64     `json:"blocker_id,omitempty"`
}
type BlockedContact struct {
	ID        int64   `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt float64 `json:"created_at,omitempty"`
	BlockeeId int64   `json:"blockee_id,omitempty"`
	BlockerId int64   `json:"blocker_id,omitempty"`
}

type ChatSettings struct {
	ID              int64     `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	User_Id         int64     `json:"user_id,omitempty"`
	Group_Chat_Type string    `json:"group_chat_type,omitempty"`
}
