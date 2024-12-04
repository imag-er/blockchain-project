package main

import (
	"time"
)

// User 结构表示用户信息
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// BlogBook 结构表示博客本
type BlogBook struct {
	Entries []BlogEntry `json:"entries"`
}

// BlogEntry 结构表示博客条目
type BlogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
}

// Block 结构表示区块链中的块
type Block struct {
	Index       int       `json:"index"`
	Timestamp   time.Time `json:"timestamp"`
	BlogEntries BlogBook  `json:"blogEntries"`
	PrevHash    string    `json:"prevHash"`
	Hash        string    `json:"hash"`
	Nonce       int       `json:"nonce"`
}

// Blockchain 结构表示区块链
type Blockchain struct {
	Blocks []Block `json:"blocks"`
}
