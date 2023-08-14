package model

import "time"

type ISP struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	//DeletedAt *time.Time `json:"-" gorm:"index"`
	ID int64 `json:"id"   gorm:"primary_key"`

	UserID       int          `json:"user_id"`
	Mobile       string       `json:"mobile"`
	ISPType      string       `json:"isp_type"      gorm:"conment:运营商类型，unicom,telecom"`
	Status       bool         `json:"status"        gorm:"default:false"`
	UnicomConfig UnicomConfig `json:"unicom_config" gorm:"embedded"`
}

type UnicomConfig struct {
	Version  string `json:"version" gorm:"default:iphone_c@10.5"`
	APPID    string `json:"app_id"`
	Cookie   string `json:"cookie"`
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}
