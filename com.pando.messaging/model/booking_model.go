package model

import "time"

type BookingDetails struct {
	ID          int64     `json:"id,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	UserId      int64     `json:"user_id,omitempty" validate:"required,numeric"`
	TicketCode  string    `json:"ticket_code,omitempty"`
	ConcertID   int64     `json:"concert_id,omitempty" validate:"required"`
	TicketPrice int64     `json:"ticket_price,omitempty" validate:"required"`
	Status      string    `json:"status,omitempty"`
	ViewerId    int64     `json:"viewer_id,omitempty"`
	ReceiverId  int64     `json:"receiver_id,omitempty"`
}

type BookingDetailDto struct {
	UserId         int64  `json:"user_id,omitempty"`
	TotalTickets   int64  `json:"total_tickets,omitempty" validate:"required"`
	ConcertID      int64  `json:"concert_id,omitempty" validate:"required"`
	WalletPassword string `json:"wallet_password,omitempty"  validate:"required"`
}
type TicketDetail struct {
	TicketCode       string `json:"ticketCode,omitempty"`
	ReceiverUsername string `json:"receiverUsername,omitempty"`
	ReceiverUserId   int64  `json:"receiverId,omitempty"`
	SenderUsername   string `json:"senderUsername,omitempty"`
	SenderUserId     int64  `json:"senderUserId,omitempty"`
	Status           string `json:"status,omitempty"`
	TimeZone         string `json:"timezone,omitempty"`
}
type Refund struct {
	WalletId string `json:"wallet_id,omitempty"`
	Amount   string `json:"amount,omitempty"`
	UserId   int64  `json:"user_id,omitempty"`
}
type Concerts struct {
	Id                  int64     `json:"id,omitempty"`
	CreatedAt           time.Time `json:"created_at,omitempty"`
	UpdatedAt           time.Time `json:"updated_at,omitempty"`
	ArtistName          string    `json:"artist_name,omitempty"`
	ConcertDate         time.Time `json:"concert_date,omitempty"`
	Description         string    `json:"description,omitempty"`
	Thumbnail_file_hash string    `json:"thumbnail_file_hash,omitempty"`
	TicketPrice         float64   `json:"ticket_price,omitempty"`
	Title               string    `json:"title,omitempty"`
	UserId              int64     `json:"user_id,omitempty"`
	Viewing_experience  string    `json:"viewing_experience,omitempty"`
	Awss3url            string    `json:"awss3url,omitempty"`
	ShowType            string    `json:"show_type,omitempty"`
	Language            string    `json:"language,omitempty"`
	Status              string    `json:"status,omitempty"`
	PaymentStatus       string    `json:"payment_status,omitempty"`
}
type ConcertDto struct {
	Id                  int64          `json:"id,omitempty"`
	CreatedAt           time.Time      `json:"createdAt,omitempty"`
	UpdatedAt           time.Time      `json:"updatedAt,omitempty"`
	ArtistName          string         `json:"artistName,omitempty"`
	ConcertDate         time.Time      `json:"concertDate,omitempty"`
	Description         string         `json:"description,omitempty"`
	Thumbnail_file_hash string         `json:"thumbnailFileHash,omitempty"`
	TicketPrice         float64        `json:"ticketPrice,omitempty"`
	Title               string         `json:"title,omitempty"`
	UserId              int64          `json:"userId,omitempty"`
	Viewing_experience  string         `json:"viewingExperience,omitempty"`
	Awss3url            string         `json:"thumbnailAwsS3Url,omitempty"`
	ShowType            string         `json:"showType,omitempty"`
	Language            string         `json:"language,omitempty"`
	Status              string         `json:"status,omitempty"`
	PaymentStatus       string         `json:"paymentStatus,omitempty"`
	TicketDetail        []TicketDetail `json:"ticketDetails,omitempty"`
}

type VerifyTicketCodeDto struct {
	ArtistName         string    `json:"artistName,omitempty"`
	ConcertDate        time.Time `json:"concertDate,omitempty"`
	StreamId           string    `json:"streamId,omitempty"`
	LiveStreamingUrl   string    `json:"liveStreamingUrl,omitempty"`
	Viewing_experience string    `json:"viewingExperience,omitempty"`
	Liked              bool      `json:"liked,omitempty"`
}

type ConcertDetail struct {
	Title       string    `json:"title,omitempty"`
	ConcertDate time.Time `json:"concert_date,omitempty"`
	TicketPrice float64   `json:"ticket_price,omitempty"`
}
