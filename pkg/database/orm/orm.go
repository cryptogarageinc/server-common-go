package orm

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm/schema"
	"time"

	"reflect"
	"unicode"

	"github.com/cryptogarageinc/server-common-go/pkg/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NoLimit is used when searching without an upper limit on the number of
// returned records.
const NoLimit = -1

// ORM represent an Object Relational Mapper instance.
type ORM struct {
	config        *Config
	log           *log.Log
	connectionStr string
	enableLog     bool
	logger        *logrus.Logger
	initialized   bool
	db            *gorm.DB
	sqldb *sql.DB
}

// NewORM creates a new ORM structure with the given parameters.
func NewORM(config *Config, l *log.Log) *ORM {
	return &ORM{
		config:      config,
		log:         l,
		initialized: false,
	}
}

// Initialize initializes the ORM structure.
func (o *ORM) Initialize() error {

	if o.initialized {
		return nil
	}

	o.log.Logger.Info("ORM initialization starts")
	defer o.log.Logger.Info("ORM initialization end")

	enableLog := o.config.EnableLogging

	o.enableLog = enableLog
	o.logger = o.log.Logger

	var dbDialector gorm.Dialector

	if o.config.InMemory {
		o.log.Logger.Info("InMemory flag detected : Using Sqlite Inmemory DB")
		o.connectionStr = ":memory:"
		dbDialector = sqlite.Open(o.connectionStr)
	} else {
		// postgres db
		o.connectionStr = fmt.Sprintf(
			"host=%s port=%s dbname=%s user=%s password=%s %s",
			o.config.Host,
			o.config.Port,
			o.config.DbName,
			o.config.DbUser,
			o.config.DbPassword,
			o.config.ConnectionParams)
		dbDialector = postgres.Open(o.connectionStr)
	}

	opened, err := gorm.Open(dbDialector, &gorm.Config{
		SkipDefaultTransaction:                   false,
		NamingStrategy:                           nil,
		FullSaveAssociations:                     false,
		Logger:                                   nil,
		NowFunc:                                  nil,
		DryRun:                                   false,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		AllowGlobalUpdate:                        false,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		Dialector:                                nil,
		Plugins:                                  nil,
	})
	if err != nil {
		o.log.Logger.Error(err, "Could not open database.")
		return errors.Wrap(err, "failed to open database")
	}

	o.db = opened
	sqldb, err := o.db.DB()
	if err != nil {
		err = errors.WithMessage(err, "Could not access sub sql database.")
		o.log.Logger.Error(err)
		return err
	}
	o.sqldb = sqldb
	o.sqldb.SetConnMaxLifetime(time.Hour)
	o.initialized = true

	return nil
}

// IsInitialized returns whether the orm is initialized.
func (o *ORM) IsInitialized() bool {
	return o.initialized
}

// Finalize releases the resources held by the orm.
func (o *ORM) Finalize() error {
	err := o.sqldb.Close()
	if err != nil {
		return errors.Errorf("failed to close database connection")
	}
	return nil
}

// GetDB returns the DB instance associated with the orm object. Panics if the
// object is not initialized.
func (o *ORM) GetDB() *gorm.DB {
	if !o.IsInitialized() {
		panic("Trying to access uninitialized ORM object.")
	}

	return o.db
}

// GetColumnNames returns the name of the columns for the given model.
func (o *ORM) GetColumnNames(model interface{}) []string {
	result := make([]string, 0)
	t := reflect.TypeOf(model)
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		first := rune(name[0])
		if unicode.IsUpper(first) {
			result = append(
				result,
				o.db.NamingStrategy.ColumnName("", name))
		}
	}
	return result
}

// GetTableName returns the name of the table for the given model.
// Assumes that the globalDB is initialized, returns empty string if not
func (o *ORM) GetTableName(model interface{}) string {
	if o.initialized {
		stmt := gorm.Statement{DB: o.db}
		stmt.Parse(model)
		return stmt.Schema.Table
	}
	structName := reflect.TypeOf(model).Elem().Name()
	return schema.NamingStrategy{}.TableName(structName)
}

// IsRecordNotFoundError returns whether the given error is due to a requested
// record not present in the DB.
func IsRecordNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// NewRecordNotFoundError returns a ErrRecordNotFoundError.
func NewRecordNotFoundError() error {
	return gorm.ErrRecordNotFound
}
