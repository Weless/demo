package db

import (
	"fmt"
	myDB "net_disk_demo/db/mysql"
	"time"
)

// UserFile: 用户文件表结构体
type UserFile struct {
	UserName    string
	FileHash    string
	FileName    string
	FileSize    int64
	UploadAt    string
	LastUpdated string
}

func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) bool {
	sql := "insert ignore into tbl_user_file (user_name, file_sha1, file_name, file_size, upload_at)" +
		" values(?,?,?,?,?)"
	stmt, err := myDB.DBConn().Prepare(sql)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash, filename, filesize, time.Now())
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	sql := "select file_sha1,file_name,file_size,upload_at,last_update from tbl_user_file " +
		"where user_name = ? limit ?"
	stmt, err := myDB.DBConn().Prepare(sql)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, limit)
	if err != nil {
		return nil, err
	}

	var userFiles []UserFile
	for rows.Next() {
		ufile := UserFile{}
		err = rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdated)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		ufile.UserName = username
		userFiles = append(userFiles, ufile)
	}
	return userFiles, nil
}
