package handler

import "net/http"

// HTTPInterceptor : http请求拦截器
func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		username := r.Form.Get("username")
		token := r.Form.Get("token")
		// todo: 判断token是否有效
		// 无效的话直接返回错误信息给客户端
		if len(username) < 3 || !IsTokenValid(username, token) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("failed to login"))
			return
		}
		// 有效的话则传递到具体的handler函数中进行逻辑处理 h(w,r)
		h(w, r)
	})
}
