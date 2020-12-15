package gorm

import (
	"context"
	"reflect"
	"strings"

	"github.com/go-kratos/kratos/pkg/database/sql"
	"gorm.io/gorm"
)

// DB DB
type DB struct {
	*sql.DB
	Builder *gorm.DB
}

// NewMySQL NewMySQL
func NewMySQL(c *sql.Config) *DB {
	db := sql.NewMySQL(c)
	return &DB{DB: db, Builder: NewBuilder()}
}

// Build 返回用于生成sql的 DB
func (d *DB) Build() *gorm.DB {
	return d.Builder.Session(&gorm.Session{DryRun: true})
}

// Execute Execute
func (d *DB) Execute(ctx context.Context, db *gorm.DB) error {
	if db.Error == nil {
		sql := db.Statement.SQL.String()

		switch {
		case strings.HasPrefix(sql, "SELECT"):
			d.doQuery(ctx, db)
		case strings.HasPrefix(sql, "INSERT"):
			d.doCreate(ctx, db)
		default:
			d.doExec(ctx, db)
		}
	}

	return db.Error
}

// doQuery doQuery
func (d *DB) doQuery(ctx context.Context, db *gorm.DB) {
	if db.Error == nil {
		rows, err := d.DB.Query(ctx, db.Statement.SQL.String(), db.Statement.Vars...)
		if err != nil {
			db.AddError(err)
			return
		}

		ScanClose(rows, db, false)
	}
	db.Statement.SQL.Reset()
	db.Statement.Vars = nil

}

// doExec doExec/Update
func (d *DB) doExec(ctx context.Context, db *gorm.DB) {
	if db.Error == nil {
		result, err := d.DB.Exec(ctx, db.Statement.SQL.String(), db.Statement.Vars...)
		if err != nil {
			db.AddError(err)
		} else {
			db.RowsAffected, _ = result.RowsAffected()
		}
	}
	db.Statement.SQL.Reset()
	db.Statement.Vars = nil
}

// doCreate doCreate
func (d *DB) doCreate(ctx context.Context, db *gorm.DB) {
	if db.Error == nil {
		result, err := d.DB.Exec(ctx, db.Statement.SQL.String(), db.Statement.Vars...)
		if err == nil {
			db.RowsAffected, _ = result.RowsAffected()

			if db.RowsAffected > 0 {
				if db.Statement.Schema != nil && db.Statement.Schema.PrioritizedPrimaryField != nil && db.Statement.Schema.PrioritizedPrimaryField.HasDefaultValue {
					if insertID, err := result.LastInsertId(); err == nil && insertID > 0 {
						switch db.Statement.ReflectValue.Kind() {
						case reflect.Slice, reflect.Array:
							// if config.LastInsertIDReversed {
							if false {
								for i := db.Statement.ReflectValue.Len() - 1; i >= 0; i-- {
									rv := db.Statement.ReflectValue.Index(i)
									if reflect.Indirect(rv).Kind() != reflect.Struct {
										break
									}

									_, isZero := db.Statement.Schema.PrioritizedPrimaryField.ValueOf(rv)
									if isZero {
										db.Statement.Schema.PrioritizedPrimaryField.Set(rv, insertID)
										insertID--
									}
								}
							} else {
								for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
									rv := db.Statement.ReflectValue.Index(i)
									if reflect.Indirect(rv).Kind() != reflect.Struct {
										break
									}

									if _, isZero := db.Statement.Schema.PrioritizedPrimaryField.ValueOf(rv); isZero {
										db.Statement.Schema.PrioritizedPrimaryField.Set(rv, insertID)
										insertID++
									}
								}
							}
						case reflect.Struct:
							if _, isZero := db.Statement.Schema.PrioritizedPrimaryField.ValueOf(db.Statement.ReflectValue); isZero {
								db.Statement.Schema.PrioritizedPrimaryField.Set(db.Statement.ReflectValue, insertID)
							}
						}
					} else {
						db.AddError(err)
					}
				}
			}
		} else {
			db.AddError(err)
		}
	}

	db.Statement.SQL.Reset()
	db.Statement.Vars = nil
}
