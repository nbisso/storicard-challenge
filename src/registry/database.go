package registry

import (
	"sync"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"github.com/nbisso/storicard-challenge/infrastracture/conf"
)

var dbinstance *sqlx.DB
var dbonce sync.Once

func (r *register) NewDatabase() *sqlx.DB {

	dbonce.Do(func() {
		dbinstance = sqlx.MustConnect(conf.Instance.Database.Driver, conf.Instance.Database.DSN)
	})

	return dbinstance
}
