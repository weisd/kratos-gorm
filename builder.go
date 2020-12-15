package gorm

import (
	"context"
	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _sqlbuilder *gorm.DB

func init() {
	_sqlbuilder = NewBuilder()
}

// Build Build
func Build() *gorm.DB {
	return _sqlbuilder.Session(&gorm.Session{DryRun: true})
}

// NewBuilder 用gorm来创建sql
func NewBuilder() *gorm.DB {
	db, _ := gorm.Open(mysql.New(mysql.Config{
		Conn:                      &TryConn{},
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	return db
}

// TryConn TryConn
type TryConn struct {
}

// PrepareContext PrepareContext
func (c *TryConn) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return nil, nil
}

// ExecContext ExecContext
func (c *TryConn) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

// QueryContext QueryContext
func (c *TryConn) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

// QueryRowContext QueryRowContext
func (c *TryConn) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return nil
}
