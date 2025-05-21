package factory

import (
	"gorm.io/gorm"

	"wakuwaku_nihongo/internals/pkg/database"

	"wakuwaku_nihongo/internals/pkg/redisutil"
)

type Factory struct {
	Db *gorm.DB

	Redis *redisutil.Redis
}

func NewFactory() *Factory {

	f := &Factory{}

	f.SetupDb()

	// f.SetupRedis()
	return f
}

func (f *Factory) SetupDb() {
	db := database.Connection()
	f.Db = db
}

func (f *Factory) SetupRedis() {
	f.Redis = redisutil.NewRedis()
}

func (f *Factory) SetupRepository() {
	if f.Db == nil {
		panic("Failed setup repository, db is undefined")
	}
}
