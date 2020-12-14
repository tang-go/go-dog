package model

import "time"

type Watermark struct {
	Timestamp time.Time `json:"timestamp"`
	Appid     string    `json:"appid"`
}

type UserWX struct {
	OpenId    string      `json:"openId"`
	UnionId   string      `json:"unionId"`
	NickName  string      `json:"nickName"`
	Gender    int         `json:"gender"`
	Language  string      `json:"language"`
	City      string      `json:"city"`
	Province  string      `json:"province"`
	Country   string      `json:"country"`
	AvatarUrl string      `json:"avatarUrl"`
	Watermark []Watermark `json:"watermark"`
}

type UserPhone struct {
	PhoneNumber     string    `json:"phoneNumber"`
	PurePhoneNumber string    `json:"purePhoneNumber"`
	CountryCode     int       `json:"countryCode"`
	Watermark       Watermark `json:"watermark"`
}
