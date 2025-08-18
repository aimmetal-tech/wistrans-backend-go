package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

// DB 数据库连接实例
var DB *sql.DB

// InitDB 初始化数据库连接
func InitDB() error {
	// 先从系统环境变量获取数据库URL
	databaseURL := os.Getenv("DATABASE_URL")
	
	// 如果环境变量中没有，则从.env文件获取
	if databaseURL == "" {
		err := godotenv.Load()
		if err != nil {
			log.Printf("未找到.env文件: %v", err)
		}
		databaseURL = os.Getenv("DATABASE_URL")
	}
	
	if databaseURL == "" {
		return fmt.Errorf("未配置数据库URL")
	}
	
	// 如果数据库URL中没有sslmode参数，则添加禁用SSL的参数
	if !strings.Contains(databaseURL, "sslmode=") {
		if strings.Contains(databaseURL, "?") {
			databaseURL += "&sslmode=disable"
		} else {
			databaseURL += "?sslmode=disable"
		}
	}
	
	var err error
	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("无法连接数据库: %v", err)
	}
	
	// 测试连接
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}
	
	log.Println("数据库连接成功")
	
	// 检查并创建表
	err = createTables()
	if err != nil {
		return fmt.Errorf("创建表失败: %v", err)
	}
	
	return nil
}

// createTables 检查并创建必要的表
func createTables() error {
	// 创建 conversations 表
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS conversations (
			id TEXT PRIMARY KEY,
			title TEXT,
			model TEXT,
			service TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("创建 conversations 表失败: %v", err)
	}
	
	// 创建 messages 表
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id SERIAL PRIMARY KEY,
			conversation_id TEXT REFERENCES conversations(id) ON DELETE CASCADE,
			role TEXT,
			content TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("创建 messages 表失败: %v", err)
	}
	
	log.Println("数据库表检查完成")
	return nil
}