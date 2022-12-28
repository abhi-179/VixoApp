package model

import (
	"time"
)

type UserInfo struct {
	Username        string `json:"username,omitempty"`
	Email           string `json:"email,omitempty"`
	Profile_pic_url string `json:"profile_pic_url,omitempty"`
	WalletId        string `json:"wallet_id,omitempty"`
}
type TicketInfo struct {
	WalletId    string `json:"wallet_id,omitempty"`
	TicketPrice string `json:"ticket_price,omitempty"`
	Username    string `json:"username,omitempty"`
	UserId      int64  `json:"user_id,omitempty"`
}
type WalletDetails struct {
	ID        int64     `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	UserId    int64     `json:"user_id,omitempty" validate:"required"`
	WalletId  string    `json:"wallet_id,omitempty" validate:"required"`
	Username  string    `json:"username,omitempty"`
	Balance   string    `json:"balance,omitempty"`
}
type WalletReq struct {
	UserId   int64  `json:"user_id,omitempty" validate:"required"`
	Password string `json:"password,omitempty"`
}
type Response struct {
	Status                bool                    `json:"isSuccess"`
	ResponseCode          int64                   `json:"statusCode,omitempty"`
	Msg                   string                  `json:"message,omitempty"`
	ValidationErrors      []string                `json:",omitempty"`
	WalletAddress         string                  `json:"walletId,omitempty"`
	UserDetail            *UserDto                `json:"userDetail,omitempty"`
	Balance               string                  `json:"balance,omitempty"`
	TransactionHash       string                  `json:"transactionHash,omitempty"`
	Transactions          *[]TransactionsDto      `json:"transactions,omitempty"`
	Data                  *[]SpendAnalytics       `json:"data,omitempty"`
	Users                 *[]UserDto              `json:"usersInfo,omitempty"`
	RequestToken          *[]RequestTokenDto      `json:"requests,omitempty"`
	ConcertDetail         *[]Concerts             `json:"concertDetails,omitempty"`
	ConcertDetails        *ConcertDto             `json:"concertDetail,omitempty"`
	LiveStreamInfo        *VerifyTicketCodeDto    `json:"liveStreamInfo,omitempty"`
	Reviews               *[]ReviewDto            `json:"reviews,omitempty"`
	Rating                float64                 `json:"rating,omitempty"`
	Groups                *Group                  `json:"groups,omitempty"`
	GroupDetail           *[]Group                `json:"groupsDetail,omitempty"`
	CallDetails           *[]CallDetail           `json:"callDetails,omitempty"`
	Details               *CallDetail             `json:"details,omitempty"`
	UserDetails           *map[int64]interface{}  `json:"userDetails,omitempty"`
	GroupDetails          *map[int64]interface{}  `json:"groupDetails,omitempty"`
	BlockedContactDetails *map[string]interface{} `json:"blockedContactDetails,omitempty"`
	User_details          *[]User                 `json:"userInfo,omitempty"`
	Statues               *[]Statuses             `json:"yourStatuses,omitempty"`
	Statues1              *map[int64][]Statuses   `json:"othersStatuses,omitempty"`
	Status2               *map[string][]Statuses  `json:"statuses,omitempty"`
	URL                   string                  `json:"url,omitempty"`
	Data1                 *Data                   `json:"wallpapers,omitempty"`
	GroupChatSetting      *string                 `json:"groupChatType,omitempty"`
	BlockedContacts       *[]BlockedContact       `json:"blockedContacts,omitempty"`
	Backups               *Backups                `json:"backupInfo,omitempty"`
	NewAdminId            int64                   `json:"newAdminId,omitempty"`
}
type Transactions struct {
	ID                  int64     `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt           time.Time `json:"created_at,omitempty"`
	SenderWalletId      string    `json:"sender_wallet_id,omitempty"`
	SenderUsername      string    `json:"sender_username,omitempty"`
	ReceiverWalletId    string    `json:"receiver_wallet_id,omitempty"`
	ReceiverUsername    string    `json:"receiver_username,omitempty"`
	Amount              string    `json:"amount,omitempty"`
	Message             string    `json:"message,omitempty"`
	Status              string    `json:"status,omitempty"`
	TransactionHash     string    `json:"transaction_hash,omitempty"`
	Time                string    `json:"time,omitempty"`
	TransactionCategory string    `json:"transaction_category,omitempty"`
}
type User struct {
	ID              int64     `json:"id,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	Username        string    `json:"username,omitempty" validate:"required"`
	Profile_Pic_Url string    `json:"profile_pic_url,omitempty"`
	Email           string    `json:"email,omitempty"`
}
type RequestTokens struct {
	ID                    int64     `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt             time.Time `json:"created_at,omitempty"`
	RequestStatus         string    `json:"request_status,omitempty"`
	RequestedFromUserId   int64     `json:"requested_from_user_id,omitempty" validate:"required,numeric"`
	RequestedByUserId     int64     `json:"requested_by_user_id,omitempty" validate:"required,numeric"`
	RequestedFromWalletId string    `json:"requested_from_wallet_id,omitempty" validate:"required"`
	RequestedByWalletId   string    `json:"requested_by_wallet_id,omitempty" validate:"required"`
	Amount                string    `json:"amount,omitempty" validate:"required"`
	Message               string    `json:"message,omitempty"`
}
type RequestTokenDto struct {
	ID                    int64     `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt             time.Time `json:"created_at,omitempty"`
	RequestStatus         string    `json:"request_status,omitempty"`
	RequestedFromUserId   int64     `json:"requested_from_user_id,omitempty" validate:"required,numeric"`
	RequestedByUserId     int64     `json:"requested_by_user_id,omitempty" validate:"required,numeric"`
	RequestedFromWalletId string    `json:"requested_from_wallet_id,omitempty" validate:"required"`
	RequestedByWalletId   string    `json:"requested_by_wallet_id,omitempty" validate:"required"`
	Amount                string    `json:"amount,omitempty" validate:"required"`
	Message               string    `json:"message,omitempty"`
	Username              string    `json:"requested_by_username,omitempty"`
	Profile_Pic_Url       string    `json:"profile_pic_url"`
}
type Notifications struct {
	SenderUserId      int64  `json:"senderUserId,omitempty"`
	ReceiverUserId    int64  `json:"receiverUserId,omitempty"`
	Message           string `json:"message,omitempty"`
	SenderUsername    string `json:"senderUsername,omitempty"`
	ReceiverUsername  string `json:"receiverUsername,omitempty"`
	NotificationTitle string `json:"title,omitempty"`
	NotificationType  string `json:"type,omitempty"`
	Profile_Pic_Url   string `json:"senderProfilePicUrl,omitempty"`
}
type TransactionsDto struct {
	Date            string
	Username        string
	TransactionHash string
	WalletId        string
	Amount          string
	Type            string
}
type SpendAnalyticsReq struct {
	Title  string  `json:"title,omitempty"`
	Amount float64 `json:"amount,omitempty"`
	Month  string  `json:"month,omitempty"`
	Total  float64 `json:"total,omitempty"`
}

type SpendAnalytics struct {
	Month     string      `json:"month,omitempty"`
	Total     float64     `json:"total,omitempty"`
	SpendList []SpendList `json:"spendList,omitempty"`
}
type SpendList struct {
	Title  string  `json:"title,omitempty"`
	Amount float64 `json:"amount,omitempty"`
}

type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}
type UserDto struct {
	ID              int64  `json:"id,omitempty"`
	Username        string `json:"username,omitempty" validate:"required"`
	Profile_Pic_Url string `json:"profilePicUrl,omitempty"`
	WalletId        string `json:"wallet_id,omitempty"`
}

type WalletDto struct {
	UserId          int64  `json:"id,omitempty" validate:"required"`
	Username        string `json:"username,omitempty"`
	Profile_Pic_Url string `json:"profile_pic_url,omitempty"`
	WalletId        string `json:"wallet_id,omitempty" validate:"required"`
}
type UserInnfo struct {
	Username string `json:"username,omitempty"`
	WalletId string `json:"wallet_id,omitempty"`
	Email    string `json:"email,omitempty"`
}

type UserWalletInfo struct {
	Username string `json:"username,omitempty"`
	WalletId string `json:"wallet_id,omitempty"`
}
type ConcertInfo struct {
	Title       string    `json:"title,omitempty"`
	ConcertDate time.Time `json:"concert_date,omitempty"`
	TicketPrice float64   `json:"ticket_price,omitempty"`
}
type ConcertID struct {
	Id       int64  `json:"id,omitempty"`
	ShowType string `json:"show_type,omitempty"`
}
type WalletStatement struct {
	WalletId    string `json:"wallet_id,omitempty" validate:"required"`
	StartDate   string `json:"start_date,omitempty"`
	EndDate     string `json:"end_date,omitempty"`
	TotalMonths string `json:"total_months,omitempty"`
	QueryType   string `json:"query_type,omitempty" validate:"required"`
}

type BookingInfo struct {
	WalletId    string `json:"wallet_id,omitempty"`
	TicketPrice string `json:"ticket_price,omitempty"`
	Username    string `json:"username,omitempty"`
	UserId      int64  `json:"user_id,omitempty"`
	Title       string `json:"title,omitempty"`
	TicketCode  string `json:"ticket_code,omitempty"`
}
