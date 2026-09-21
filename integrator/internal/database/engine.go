package database

import (
	"app-platform/internal/config"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// --------------------------------
// CARA PAKAI :
// --------------------------------
// SELECT
// row, err := db.Query(sql).Fetch(...?:vals) (select row 1)
// rows, err := db.Query(sql).FetchAll(..?:vals) (select banyak)
// rows, err := db.Query(sql).FetchGroup(...key1,key2) (key yg akan di group)

// INSERT
// err := db.Query(sql).Exec(...)
// id, err := db.Query(sql).LastInsertID(...)

// UPDATE / DELETE sederhana + jumlah row
// err := db.Query(sql).Exec(...)
// affected, err := db.Query(sql).RowsAffected(...)

type Connect struct {
	conn *sql.DB
}

func New() (*Connect, error) {
	dbparse := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		config.App.MySQLUser,
		config.App.MySQLPass,
		config.App.MySQLHost,
		config.App.MySQLPort,
		config.App.MySQLDB,
	)
	db, err := sql.Open("mysql", dbparse)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Connect{
		conn: db,
	}, nil

}

func (db *Connect) Close() {
	_ = db.conn.Close()
}

func (db *Connect) Query(sqlText string) *Query {
	return &Query{
		db:  db,
		sql: sqlText,
	}
}
