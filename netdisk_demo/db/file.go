package db

import (
	"database/sql"
	"fmt"
	myDB "net_disk_demo/db/mysql"
)

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

func GetFileMeta(fileHash string) (*TableFile, error) {
	sql := "select `file_sha1`,`file_name`,`file_size`,`file_addr` " +
		"from tbl_file where `file_sha1` = ? and status =1 limit 1"
	stmt, err := myDB.DBConn().Prepare(sql)
	if err != nil {
		fmt.Printf("failed to prepare statment, err:%s\n", err.Error())
		return nil, err
	}
	fileData := &TableFile{}
	err = stmt.QueryRow(fileHash).Scan(&fileData.FileHash, &fileData.FileName, &fileData.FileSize, &fileData.FileAddr)
	if err != nil {
		return nil, err
	}
	return fileData, nil
}

// OnFileUploadFinished : 文件上传完成，保存meta
func OnFileUploadFinished(fileHash, fileName, fileAddr string, fileSize int64) bool {
	sql := "insert ignore into tbl_file " +
		"(`file_sha1`, `file_name`, `file_size`, `file_addr`,`status`) " +
		"values (?,?,?,?,1)"
	stmt, err := myDB.DBConn().Prepare(sql)
	if err != nil {
		fmt.Printf("failed to prepare statment, err:%s\n", err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(fileHash, fileName, fileSize, fileAddr)
	if err != nil {
		fmt.Printf("failed to exec statment, err:%s\n", err.Error())
		return false
	}
	rf, err := ret.RowsAffected()
	if err != nil {
		return false
	}
	if rf <= 0 {
		fmt.Printf("file with hash: %s has been uploaded before", fileHash)
	}
	return true
}
