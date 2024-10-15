package model

import (
	"com.banxiaoxiao.server/config"
	"com.banxiaoxiao.server/snow"
	"github.com/jinzhu/gorm"
)

const ZeroID = int64(0)

var node *snow.Node

func InitNode() {
	node, _ = snow.NewNode(config.Cfg.Etcd.NewId)
	if node == nil {
		panic("failed to init snowflake for distributed ID generation")
	}
}

func GenID() int64 {
	id := node.Generate().Int64()
	return id
}

type Model struct {
	ID int64 `gorm:"primaryKey" json:"id"`
}

func (m *Model) BeforeCreate(scope *gorm.Scope) (err error) {
	if m.ID <= ZeroID {
		id := GenID()
		return scope.SetColumn("ID", id)
	}
	return nil
}
