package model

import (
	"time"

	"github.com/lib/pq"
)

type CallDetails struct {
	ID                  int64         `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt           time.Time     `json:"created_at,omitempty"`
	CallId              int64         `json:"call_id,omitempty"`
	CallerId            int64         `json:"caller_id,omitempty"  validate:"required,numeric"`
	CallDuration        time.Time     `json:"call_duration" validate:"numeric"`
	Filehash            string        `json:"filehash,omitempty" validate:"required"`
	AwsUrl              string        `json:"aws_url,omitempty"`
	StartTime           time.Time     `json:"start_time" validate:"numeric"`
	EndTime             time.Time     `json:"end_time" validate:"numeric"`
	User_ids            pq.Int64Array `gorm:"type:integer[]" json:"user_ids,omitempty"`
	Deleted_by_user_ids pq.Int64Array `gorm:"type:integer[]" json:"deleted_by_user_ids,omitempty"`
	IsAudioCall         bool          `json:"is_audio_call"`
	IsMissedCall        bool          `json:"is_missed_call"`
	IsGroupCall         bool          `json:"is_group_call"`
}
type CallDetail struct {
	ID                  int64         `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt           float64       `json:"created_at,omitempty"`
	CallId              int64         `json:"call_id,omitempty"`
	CallerId            int64         `json:"caller_id,omitempty" validate:"required,numeric"`
	CallDuration        float64       `json:"call_duration,omitempty" validate:"numeric"`
	Filehash            string        `json:"filehash,omitempty" validate:"required"`
	AwsUrl              string        `json:"aws_url,omitempty"`
	StartTime           float64       `json:"start_time,omitempty" validate:"numeric"`
	EndTime             float64       `json:"end_time,omitempty" validate:"numeric"`
	User_ids            pq.Int64Array `gorm:"type:integer[]" json:"user_ids,omitempty"`
	Deleted_by_user_ids pq.Int64Array `gorm:"type:integer[]" json:"deleted_by_user_ids,omitempty"`
	IsAudioCall         bool          `json:"is_audio_call,omitempty"`
	IsMissedCall        bool          `json:"is_missed_call,omitempty"`
	IsGroupCall         bool          `json:"is_group_call,omitempty"`
}
type GroupId struct {
	GroupId int `json:"group_id,omitempty"`
}
