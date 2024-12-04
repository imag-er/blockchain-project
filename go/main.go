package main

import (
	"fmt"
	"net/http"
	"os"
)

// 在main函数之前添加一个初始化函数
func initialize() {
	// 创建博客文件夹
	err := os.MkdirAll(blogFolder, 0755)
	if err != nil {
		fmt.Println("Error creating blog folder:", err)
		os.Exit(1)
	}

	// 创建区块链文件夹
	err = os.MkdirAll(blockchainFolder, 0755)
	if err != nil {
		fmt.Println("Error creating blockchain folder:", err)
		os.Exit(1)
	}
}

func main() {

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/loginin", logininHandler)
	http.HandleFunc("/validateBlockchain", validateBlockchainHandler)

	fmt.Println("后端服务器启动在82端口")
	http.ListenAndServe("0.0.0.0:82", nil)
}
