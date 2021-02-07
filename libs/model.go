package libs

import "time"

type ClipRecord struct {
	User      UserCode  `form:"user" json:"user"`
	Payload   string    `form:"payload" json:"payload"`
	UpdatedAt time.Time `form:"updated_at" json:"updated_at"`
}
