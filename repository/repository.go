package repository

import (
	"com.banxiaoxiao.server/config"
	"com.banxiaoxiao.server/logger"
	"database/sql"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Repository interface {
	Model(value interface{}) *gorm.DB
	Select(query interface{}, args ...interface{}) *gorm.DB
	Find(out interface{}, where ...interface{}) *gorm.DB
	Exec(sql string, values ...interface{}) *gorm.DB
	First(out interface{}, where ...interface{}) *gorm.DB
	Raw(sql string, values ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Update(value interface{}) *gorm.DB
	Delete(value interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	Limit(limit interface{}) *gorm.DB
	Preload(column string, conditions ...interface{}) *gorm.DB
	Scopes(funcs ...func(*gorm.DB) *gorm.DB) *gorm.DB
	ScanRows(rows *sql.Rows, result interface{}) error
	Transaction(fc func(tx Repository) error) (err error)
	Close() error
	DropTableIfExists(value interface{}) *gorm.DB
	AutoMigrate(value interface{}) *gorm.DB
}

type repository struct {
	db *gorm.DB
}

type accountRepository struct {
	*repository
}

func NewProxyRepository(conf *config.Config, l *logger.Logger) Repository {
	logger.Log().Infof("Try database connection")
	connection := getConnection(conf)
	logger.Log().Infof("%s", connection)
	db, err := gorm.Open(conf.Database.Dialect, connection)
	if err != nil {
		logger.Log().Errorf("Failure database connection: %s", err.Error())
		os.Exit(2)
	}
	logger.Log().Infof("Success database connection, %s", conf.Database.Host)
	db.LogMode(true)
	db.SetLogger(l)
	return &accountRepository{&repository{db: db}}
}

const (
	// SQLITE represents SQLite3
	SQLITE = "sqlite3"
	// POSTGRES represents PostgreSQL
	POSTGRES = "postgres"
	// MYSQL represents MySQL
	MYSQL = "mysql"
)

func getConnection(config *config.Config) string {
	if config.Database.Dialect == POSTGRES {
		return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", config.Database.Host, config.Database.Port, config.Database.Username, config.Database.Dbname, config.Database.Password)
	} else if config.Database.Dialect == MYSQL {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&character_set_server=utf8mb4",
			config.Database.Username, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Dbname)
	}
	return config.Database.Host
}

func (rep *repository) Model(value interface{}) *gorm.DB {
	return rep.db.Model(value)
}

func (rep *repository) Select(query interface{}, args ...interface{}) *gorm.DB {
	return rep.db.Select(query, args...)
}

func (rep *repository) Find(out interface{}, where ...interface{}) *gorm.DB {
	return rep.db.Find(out, where...)
}

func (rep *repository) Exec(sql string, values ...interface{}) *gorm.DB {
	return rep.db.Exec(sql, values...)
}

func (rep *repository) First(out interface{}, where ...interface{}) *gorm.DB {
	return rep.db.First(out, where...)
}

func (rep *repository) Raw(sql string, values ...interface{}) *gorm.DB {
	return rep.db.Raw(sql, values...)
}

func (rep *repository) Create(value interface{}) *gorm.DB {
	return rep.db.Create(value)
}

func (rep *repository) Save(value interface{}) *gorm.DB {
	return rep.db.Save(value)
}

func (rep *repository) Update(value interface{}) *gorm.DB {
	return rep.db.Update(value)
}

func (rep *repository) Delete(value interface{}) *gorm.DB {
	return rep.db.Delete(value)
}

func (rep *repository) Where(query interface{}, args ...interface{}) *gorm.DB {
	return rep.db.Where(query, args...)
}

func (rep *repository) Limit(limit interface{}) *gorm.DB {
	return rep.db.Limit(limit)
}

func (rep *repository) Preload(column string, conditions ...interface{}) *gorm.DB {
	return rep.db.Preload(column, conditions...)
}

func (rep *repository) Scopes(funcs ...func(*gorm.DB) *gorm.DB) *gorm.DB {
	return rep.db.Scopes(funcs...)
}

func (rep *repository) ScanRows(rows *sql.Rows, result interface{}) error {
	return rep.db.ScanRows(rows, result)
}

func (rep *repository) Close() error {
	return rep.db.Close()
}

func (rep *repository) DropTableIfExists(value interface{}) *gorm.DB {
	return rep.db.DropTableIfExists(value)
}

func (rep *repository) AutoMigrate(value interface{}) *gorm.DB {
	return rep.db.AutoMigrate(value)
}

func (rep *repository) Transaction(fc func(tx Repository) error) (err error) {
	panicked := true
	tx := rep.db.Begin()
	defer func() {
		if panicked || err != nil {
			tx.Rollback()
		}
	}()

	txRep := &repository{}
	txRep.db = tx
	err = fc(txRep)

	if err == nil {
		err = tx.Commit().Error
	}

	panicked = false
	return
}
