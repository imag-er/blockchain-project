package main

import (
	"sync"
)

var (
	currentUser User
	blockchain  Blockchain
)

const (
	passwdFile       = "./passwd"
	blogFolder       = "./blogs/"     // 存放博客文件的文件夹
	blockchainFolder = "./blockchains/" // 存放区块链文件的文件夹
)

// 登录状态
var loggedIn bool

// 博客是否被篡改的标志
var blogTampered bool

var mutex sync.Mutex
