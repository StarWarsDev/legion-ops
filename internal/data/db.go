package data

import (
	"github.com/StarWarsDev/legion-ops/internal/orm"
	"github.com/StarWarsDev/legion-ops/internal/util"
	"github.com/jinzhu/gorm"
)

func NewDB(orm *orm.ORM) *gorm.DB {
	db := orm.DB.New()
	if util.DebugEnabled() {
		db = db.Debug()
	}
	return db
}
