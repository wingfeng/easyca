package model

import "time"

type ServerCa struct {
	Id        string        `xorm:"varchar(36) pk notnull 'id'"`
	UserId    string        `xorm:"varchar(36)    notnull 'user_id'"`
	Dn        string        `xorm:"varchar(36)    notnull 'dn'"`
	Created   time.Time     `xorm:"datetime   created   'created'"`
	Updated   time.Time     `xorm:"datetime   updated   'updated'"`
	Deleted   time.Time     `xorm:"datetime   deleted   'deleted'"`
}