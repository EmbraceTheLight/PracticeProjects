// 使用标准库的 HTTP 编程
package main

import (
	"fmt"
	"net/http"
)

// 中间件
func loggingMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("请求: %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		fmt.Println("请求结束")
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet {
			fmt.Fprintf(w, `[{"ID": 1, "Name": "Tom", "Age": 18}]`)
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})
	handler := loggingMiddleWare(mux)
	http.ListenAndServe(":5567", handler)

}
