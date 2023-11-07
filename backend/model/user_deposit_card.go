package model

import "database/sql"

type UserDepositCard struct {
	Id         uint          `default:"0" json:"id"`
	CardKey    string        `default:"" json:"card_key"`
	CardAmount float64       `default:"" json:"card_amount"`
	UsedUserId sql.NullInt32 `default:"0" json:"-"`
	UsedTime   sql.NullTime  `default:"0000-00-00 00:00:00" json:"-"`
}
