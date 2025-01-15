package tokens

import (
	"time"
)

type tokens struct {
	Id         int       `json:"id"`
	Token      string    `json:"token"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}
