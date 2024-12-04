package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// 登陆后的请求处理函数
func logininHandler(w http.ResponseWriter, r *http.Request) {
	// 鉴别用户是否登录
	if !loggedIn {
		http.Error(w, "用户未登录", http.StatusUnauthorized)
		return
	}

	// 打印接收到的请求信息
	//log.Printf("Received request: %s %s\n", r.Method, r.URL)

	// 定义一个结构体，用于存储要返回给前端的数据
	type Response struct {
		Status         string `json:"status"`
		BlockchainInfo string `json:"blockchainInfo"`
		BlogInfo       string `json:"blogInfo"`
		Refresh        string `json:"refresh"`
	}

	// 处理获取博客和区块链信息请求
	if r.Method == http.MethodGet {
		// 获取区块链信息的字符串
		blockchainInfo := viewBlockchainInfo(currentUser.Username)

		// 获取博客信息的字符串
		blogInfo := viewBlogEntries(currentUser.Username)

		// 构建要返回的结构体
		responseData := Response{
			Status:         "success",
			BlockchainInfo: blockchainInfo,
			BlogInfo:       blogInfo,
			Refresh:        "True",
		}

		// 将结构体转换为 JSON 格式的字符串
		responseJSON, err := json.Marshal(responseData)
		if err != nil {
			http.Error(w, "无法序列化响应数据", http.StatusInternalServerError)
			return
		}

		// 打印要发送到前端的响应信息
		//log.Printf("Sent response (length=%d): %q\n", len(responseJSON), responseJSON)
		// 设置响应头为 application/json
		w.Header().Set("Content-Type", "application/json")
		// 发送到前端
		w.Write(responseJSON)

		// 处理写博客请求
	} else if r.Method == http.MethodPost {
		var entry struct {
			Content string `json:"content"`
		}

		err := json.NewDecoder(r.Body).Decode(&entry)
		if err != nil {
			http.Error(w, "无效的请求数据", http.StatusBadRequest)
			return
		}

		// 提交博客条目
		submitBlogEntry(entry.Content, w)

		// 获取更新后的区块链信息的字符串
		blockchainInfo := viewBlockchainInfo(currentUser.Username)

		// 获取更新后的博客信息的字符串
		blogInfo := viewBlogEntries(currentUser.Username)

		// 构建要返回的结构体
		responseData := Response{
			Status:         "success",
			BlockchainInfo: blockchainInfo,
			BlogInfo:       blogInfo,
			Refresh:        "True",
		}

		// 将结构体转换为 JSON 格式的字符串
		responseJSON, err := json.Marshal(responseData)
		if err != nil {
			http.Error(w, "无法序列化响应数据", http.StatusInternalServerError)
			return
		}

		// 打印要发送到前端的响应信息
		//log.Printf("Sent response (length=%d): %q\n", len(responseJSON), responseJSON)
		// 设置响应头为 application/json
		w.Header().Set("Content-Type", "application/json")
		// 发送到前端
		w.Write(responseJSON)
	} else {
		http.Error(w, "不允许的请求方法", http.StatusMethodNotAllowed)
	}
}

// 登出处理函数
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// 处理登出请求
	if r.Method == http.MethodPost {
		// 登出成功，将登录状态保存在服务器端
		loggedIn = false
		currentUser = User{}

		// 返回前端 JSON 数据
		responseData := map[string]string{
			"status":  "success",
			"message": "登出成功",
		}

		jsonResponse, err := json.Marshal(responseData)
		if err != nil {
			http.Error(w, "无法生成 JSON 响应", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)

		//fmt.Println("登出成功")
		//fmt.Println("后端返回给前端的响应：", string(jsonResponse))
	} else {
		http.Error(w, "不允许的请求方法", http.StatusMethodNotAllowed)
	}
}

// 登录请求处理函数
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// 处理登录请求
	if r.Method == http.MethodPost || r.Method == http.MethodGet {
		var credentials User
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(w, "无效的请求数据", http.StatusBadRequest)
			return
		}

		username := credentials.Username
		password := credentials.Password

		if validateLogin(username, password) {
			// 登录成功，将登录状态保存在服务器端
			loggedIn = true
			currentUser = User{Username: username, Password: password}

			// 加载区块链和博客信息
			loadBlockchain(username)

			// 添加CORS头部
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:80")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// 返回前端 JSON 数据
			responseData := map[string]string{
				"status":  "success",
				"message": fmt.Sprintf("登录成功，用户名：%s，密码：%s", username, password),
			}

			jsonResponse, err := json.Marshal(responseData)
			if err != nil {
				http.Error(w, "无法生成 JSON 响应", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponse)

			// 便于调试，将发送数据打印在终端
			//fmt.Printf("登录成功，用户名：%s，密码：%s\n", username, password)
			//fmt.Println("后端返回给前端的响应：", string(jsonResponse))
		} else {
			// 返回前端 JSON 数据
			responseData := map[string]string{
				"status": "error",
				"error":  "用户名或密码错误",
			}

			jsonResponse, err := json.Marshal(responseData)
			if err != nil {
				http.Error(w, "无法生成 JSON 响应", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponse)

			//fmt.Println("登录失败")
			//fmt.Println("后端返回给前端的响应：", string(jsonResponse))
		}
	} else {
		http.Error(w, "不允许的请求方法", http.StatusMethodNotAllowed)
	}
}

func validateLogin(username, password string) bool {
	// 从密码文件中验证用户
	users, err := loadUsers()
	if err != nil {
		fmt.Println("无法验证登录:", err)
		return false
	}

	for _, user := range users {
		if user.Username == username && user.Password == password {
			currentUser = user
			return true
		}
	}

	return false
}

func userExists(username string) bool {
	// 检查用户名是否存在于密码文件中
	users, err := loadUsers()
	if err != nil {
		fmt.Println("无法检查用户名:", err)
		return false
	}

	for _, user := range users {
		if user.Username == username {
			return true
		}
	}

	return false
}

func addUser(username, password string) error {
	// 添加新用户到密码文件
	users, err := loadUsers()
	if err != nil {
		return err
	}

	newUser := User{Username: username, Password: password}
	users = append(users, newUser)

	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(passwdFile, data, 0755)
	if err != nil {
		return err
	}

	return nil
}

// 处理注册请求
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var credentials User
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(w, "无效的请求数据", http.StatusBadRequest)
			return
		}

		// 打印出注册所用的用户信息
		fmt.Printf("Username: %s, Password: %s\n", credentials.Username, credentials.Password)

		mutex.Lock()
		defer mutex.Unlock()

		if userExists(credentials.Username) {
			http.Error(w, "用户名已存在，请选择登录或使用其他用户名", http.StatusBadRequest)
			return
		}

		err = addUser(credentials.Username, credentials.Password)
		if err != nil {
			http.Error(w, "注册失败", http.StatusInternalServerError)
			return
		}

		// 注册成功响应
		responseData := map[string]interface{}{
			"status":  "success",
			"message": "注册成功，正在自动登录...",
		}

		// 调用初始化函数
		initialize()

		// 自动登录的逻辑
		// 设置登录状态
		loggedIn = true
		currentUser = User{Username: credentials.Username, Password: credentials.Password}

		// 创建创世区块并写入文件
		createGenesisBlock(credentials.Username)

		// 将响应数据转换为 JSON
		responseJSON, err := json.Marshal(responseData)
		if err != nil {
			http.Error(w, "无法序列化响应数据", http.StatusInternalServerError)
			return
		}

		// 打印后端响应
		//fmt.Printf("Backend Response: %s\n", string(responseJSON))

		// 将响应发送到前端
		w.Write(responseJSON)
	} else {
		http.Error(w, "请求方法无效", http.StatusMethodNotAllowed)
	}
}

func loadUsers() ([]User, error) {
	// 从文件中读取用户信息
	data, err := os.ReadFile(passwdFile)
	if err != nil {
		// 如果文件不存在，创建一个新文件
		if os.IsNotExist(err) {
			err := os.WriteFile(passwdFile, []byte{}, 0755)
			if err != nil {
				fmt.Println("Error creating password file:", err)
				os.Exit(1)
			}
			return []User{}, nil
		}
		return nil, err
	}

	// 如果文件为空，返回空用户列表
	if len(data) == 0 {
		return []User{}, nil
	}

	var users []User
	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}
