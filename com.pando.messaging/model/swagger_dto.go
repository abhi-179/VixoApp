package model

import (
	"time"
)

type VerifyTicketsDto struct {
	UserId     int64  `json:"user_id,omitempty"`
	TicketCode string `json:"ticket_code,omitempty"`
	ConcertID  int64  `json:"concert_id,omitempty"`
}

type SendTicketDto struct {
	SenderUserId   int64  `json:"sender_user_id,omitempty"`
	ReceiverUserId int64  `json:"receiver_user_id,omitempty"`
	TicketCode     string `json:"ticket_code,omitempty"`
}
type CancelTicketDto struct {
	UserId     int64  `json:"user_id,omitempty"`
	ConcertId  int64  `json:"concert_id,omitempty"`
	TicketCode string `json:"ticket_code,omitempty"`
}
type SaveReviewDto struct {
	UserId    int64  `json:"user_id,omitempty"`
	ConcertId int64  `json:"concert_id,omitempty"`
	Rating    int64  `json:"rating,omitempty"`
	Feedback  string `json:"feedback,omitempty"`
}

type RequestTokensDto struct {
	RequestStatus         string `json:"request_status,omitempty"`
	RequestedFromUserId   int64  `json:"requested_from_user_id,omitempty" validate:"required,numeric"`
	RequestedByUserId     int64  `json:"requested_by_user_id,omitempty" validate:"required,numeric"`
	RequestedFromWalletId string `json:"requested_from_wallet_id,omitempty" validate:"required"`
	RequestedByWalletId   string `json:"requested_by_wallet_id,omitempty" validate:"required"`
	Amount                string `json:"amount,omitempty" validate:"required"`
	Message               string `json:"message,omitempty"`
}

type RejectRequestDto struct {
	RequestId   int64  `json:"request_id,omitempty"`
	RequestType string `json:"request_type,omitempty"`
}
type SendTokenDto struct {
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	Amount   string `json:"amount,omitempty"`
	Password string `json:"password,omitempty"`
}
type WalletStatementDto struct {
	WalletId    string `json:"wallet_id,omitempty"`
	StartDate   string `json:"start_date,omitempty"`
	EndDate     string `json:"end_date,omitempty"`
	TotalMonths string `json:"total_months,omitempty"`
	QueryType   string `json:"query_type,omitempty"`
}

type ReportDto struct {
	Reporter_Id int64 `json:"reporter_id,omitempty"`
	Reportee_Id int64 `json:"reportee_id,omitempty"`
	GroupId     int64 `json:"group_id,omitempty"`
}

type BlockedContactDto struct {
	BlockeeId int64 `json:"blockee_id,omitempty"`
	BlockerId int64 `json:"blocker_id,omitempty"`
}

type ChatSettingDto struct {
	User_Id         int64  `json:"user_id,omitempty"`
	Group_Chat_Type string `json:"group_chat_type,omitempty"`
}

type CallDetailDto struct {
	CallerId     int64     `json:"caller_id,omitempty"  validate:"required,numeric"`
	CallDuration time.Time `json:"call_duration" validate:"numeric"`
	Filehash     string    `json:"filehash,omitempty" validate:"required"`
	AwsUrl       string    `json:"aws_url,omitempty"`
	StartTime    time.Time `json:"start_time" validate:"numeric"`
	EndTime      time.Time `json:"end_time" validate:"numeric"`
	User_ids     []int     `gorm:"type:integer[]" json:"user_ids,omitempty"`
	IsAudioCall  bool      `json:"is_audio_call"`
	IsMissedCall bool      `json:"is_missed_call"`
	IsGroupCall  bool      `json:"is_group_call"`
}
type GroupDto struct {
	Group_name        string  `json:"group_name,omitempty"`
	Admin_ids         []int64 `gorm:"type:integer[]" json:"admin_ids,omitempty"`
	Subject_Timestamp float64 `json:"subject_timestamp,omitempty"`
	Subject_Owner_Id  int64   `json:"subject_owner_id,omitempty"`
	Profile_Pic_Url   string  `json:"profile_pic_url,omitempty"`
	User_ids          []int64 `gorm:"type:integer[]" json:"user_ids,omitempty" validate:"required"`
}
type AddUserToGroupDto struct {
	Id      int64   `json:"id,omitempty"`
	UserId  []int64 `json:"user_ids,omitempty"`
	AdminId []int64 `json:"admin_ids,omitempty"`
}

type BackupDto struct {
	Filehash     string `json:"filehash,omitempty" validate:"required"`
	UserId       int64  `json:"user_id,omitempty" validate:"required"`
	BackupNature string `json:"backup_nature,omitempty" validate:"required"`
	FileName     string `json:"file_name,omitempty" validate:"required"`
	FileSize     string `json:"file_size,omitempty" validate:"required"`
}
