package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var db *sql.DB

func init() {
	dsn := "root:joeytest123...@tcp(127.0.0.1:3306)/netdisk_demo?charset=utf8"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("failed to open mysql, err:%s\n", err.Error())
		os.Exit(1)
		return
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("failed to connect to mysql, err:%s\n", err.Error())
		os.Exit(1)
		return
	}
	db.SetMaxOpenConns(1000)
	fmt.Println("start mysql service successfully!!!")
}

// DBConn : 返回数据库连接
func DBConn() *sql.DB {
	return db
}

func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		// 将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Fatal(err)
			panic(err)
		}

		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records
}
