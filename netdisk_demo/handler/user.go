package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net_disk_demo/db"
	"net_disk_demo/utils"
	"time"
)

const (
	pwdSalt   = "joey"
	tokenSalt = "joey"
)

// SignUpHandler:注册用户
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("failed to read data from file"))
			return
		}
		w.Write(data)
		return
	}

	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	if len(username) < 3 || len(passwd) < 5 {
		w.Write([]byte("invalid parameter"))
		return
	}

	encPasswd := utils.Sha1([]byte(passwd + pwdSalt))
	suc := db.UserSignUp(username, encPasswd)
	if suc {
		w.Write([]byte("sign Up successfully"))
	} else {
		w.Write([]byte("failed to sign up"))
	}
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	encPasswd := utils.Sha1([]byte(password + pwdSalt))
	// 1. 校验用户名及密码
	ok := db.UserSignIn(username, encPasswd)
	if !ok {
		w.Write([]byte("failed to login"))
		return
	}

	// 2. 生成访问凭证(token)
	token := GenToken(username)
	if ok := db.UpdateToken(username, token); !ok {
		w.Write([]byte("failed to login"))
		return
	}
	// 3. 登录成功后重定向到首页
	w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
}

// GenToken:生成token
func GenToken(username string) string {
	// md5(username + timestamp + token_salt) + timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	return utils.MD5([]byte(username+ts+tokenSalt)) + ts[:8]
}
