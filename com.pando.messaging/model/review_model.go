package model

import "time"

type Reviews struct {
	Id        int64     `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UserId    int64     `json:"user_id,omitempty"`
	ConcertId int64     `json:"concert_id,omitempty"`
	Rating    int64     `json:"rating,omitempty"`
	Feedback  string    `json:"feedback,omitempty"`
}
type ReviewDto struct {
	Id              int64     `json:"id,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	UserId          int64     `json:"user_id,omitempty"`
	ConcertId       int64     `json:"concert_id,omitempty"`
	Rating          int64     `json:"rating,omitempty"`
	Feedback        string    `json:"feedback,omitempty"`
	Username        string    `json:"username,omitempty"`
	Profile_Pic_Url string    `json:"profile_pic_url"`
}
