package model

import (
	"time"

	"github.com/lib/pq"
)

type Groups struct {
	ID                int64         `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt         time.Time     `json:"created_at,omitempty"`
	Group_name        string        `json:"group_name,omitempty"`
	Admin_ids         pq.Int64Array `gorm:"type:integer[]" json:"admin_ids,omitempty" validate:"required"`
	TotalUsers        int64         `json:"total_users,omitempty"`
	Status            string        `json:"status,omitempty"`
	BlockedTime       time.Time     `json:"blocked_time,omitempty"`
	Subject_Timestamp time.Time     `json:"subject_timestamp,omitempty"`
	Subject_Owner_Id  int64         `json:"subject_owner_id,omitempty"`
	Profile_Pic_Url   string        `json:"profile_pic_url,omitempty"`
	ChatId            int64         `json:"chat_id,omitempty"`
	User_ids          pq.Int64Array `gorm:"type:integer[]" json:"user_ids,omitempty" validate:"required"`
	Pending_User_ids  pq.Int64Array `gorm:"type:integer[]" json:"pending_user_ids,omitempty"`
}

type Group struct {
	ID                int64         `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt         float64       `json:"created_at,omitempty"`
	Group_name        string        `json:"group_name,omitempty"`
	Admin_ids         pq.Int64Array `gorm:"type:integer[]" json:"admin_ids,omitempty"`
	TotalUsers        int64         `json:"total_users,omitempty"`
	BlockedTime       float64       `json:"blocked_time,omitempty"`
	Status            string        `json:"status,omitempty"`
	Subject_Timestamp float64       `json:"subject_timestamp,omitempty"`
	Subject_Owner_Id  int64         `json:"subject_owner_id,omitempty"`
	Profile_Pic_Url   string        `json:"profile_pic_url,omitempty"`
	ChatId            int64         `json:"chat_id,omitempty"`
	User_ids          pq.Int64Array `gorm:"type:integer[]" json:"user_ids,omitempty" validate:"required"`
	Pending_User_ids  pq.Int64Array `gorm:"type:integer[]" json:"pending_user_ids,omitempty"`
}
type Data struct {
	Dark    []Wallpaper `json:"dark,omitempty"`
	Light   []Wallpaper `json:"light,omitempty"`
	Bright  []Wallpaper `json:"bright,omitempty"`
	Pattern []Wallpaper `json:"pattern,omitempty"`
}

type AcceptGroupInvitationDto struct {
	GroupId int64  `json:"group_id,omitempty"`
	UserId  int64  `json:"user_id,omitempty"`
	Type    string `json:"type,omitempty"`
}
type UserDtos struct {
	ID                int64   `json:"id,omitempty"`
	CreatedAt         float64 `json:"created_at,omitempty"`
	Username          string  `json:"username,omitempty" validate:"required"`
	Profile_Pic_Url   string  `json:"profile_pic_url,omitempty"`
	Email             string  `json:"email,omitempty"`
	Friendship_status string  `json:"friendship_status,omitempty"`
	Is_friend         bool    `json:"is_friend"`
}
