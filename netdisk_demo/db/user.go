package db

import (
	"fmt"
	myDB "net_disk_demo/db/mysql"
)

// UserSignUp: 用户注册
func UserSignUp(username string, passwd string) bool {
	sql := "insert ignore into tbl_user (`user_name`,`user_pwd`) values(?,?)"
	stmt, err := myDB.DBConn().Prepare(sql)
	if err != nil {
		fmt.Printf("failed to prepare the statement, err:%s\n", err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Printf("failed to exec the statment, err:%s\n", err.Error())
		return false
	}
	affected, err := ret.RowsAffected()
	if err != nil {
		return false
	}
	if affected <= 0 {
		fmt.Println("user has been already existed")
		return false
	}
	return true
}

func UserSignIn(username string, encpwd string) bool {
	sql := "select * from tbl_user where user_name = ? limit 1"
	stmt, err := myDB.DBConn().Prepare(sql)
	if err != nil {
		fmt.Printf("failed to prepare the statement, err:%s\n", err)
		return false
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("cannot find the username:" + username)
		return false
	}

	pRows := myDB.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false
}

// UpdateToken:刷新用户登录的token
func UpdateToken(username string, token string) bool {
	sql := "replace into tbl_user_token (`user_name`,`user_token`) values (?,?)"
	stmt, err := myDB.DBConn().Prepare(sql)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

type UserToken struct {
	UserName string
	Token    string
}

// GetUserToken:获取用户token
func GetUserToken(username string) string {
	sql := "select user_name,user_token from tbl_user_token where user_name = ? limit 1"
	stmt, err := myDB.DBConn().Prepare(sql)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer stmt.Close()

	userToken := &UserToken{}
	err = stmt.QueryRow(username).Scan(&userToken.UserName, &userToken.Token)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return userToken.Token
}

type UserInfo struct {
	UserName     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

// GetUserInfo:获取用户信息
func GetUserInfo(username string) (UserInfo, error) {
	sql := "select user_name, email, phone,signup_at, last_active, status from tbl_user where user_name = ? limit 1"
	stmt, err := myDB.DBConn().Prepare(sql)
	if err != nil {
		fmt.Println(err)
		return UserInfo{}, err
	}
	defer stmt.Close()

	userInfo := UserInfo{}
	err = stmt.QueryRow(username).Scan(&userInfo.UserName, &userInfo.Email, &userInfo.Phone, &userInfo.SignupAt,
		&userInfo.LastActiveAt, &userInfo.Status)
	if err != nil {
		fmt.Println(err)
		return UserInfo{}, err
	}
	return userInfo, nil
}
