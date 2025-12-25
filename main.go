package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// 定义喝水记录结构
type WaterRecord struct {
	ID     int       `json:"id"`
	Amount int       `json:"amount"`
	Time   time.Time `json:"time"`
}

// 全局数据库连接
var db *sql.DB

// 初始化数据库
func initDB() error {
	var err error
	// 打开数据库连接
	db, err = sql.Open("sqlite3", "./water_records.db")
	if err != nil {
		return err
	}

	// 创建表
	createTableSQL := `CREATE TABLE IF NOT EXISTS water_records (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"amount" INTEGER NOT NULL,
		"time" DATETIME NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	return err
}

func main() {
	// 初始化数据库
	err := initDB()
	if err != nil {
		fmt.Println("数据库初始化失败:", err)
		return
	}
	defer db.Close()

	// 创建Gin路由器
	router := gin.Default()

	// 设置静态文件目录
	router.Static("/static", "./static")

	// 加载HTML模板
	router.LoadHTMLGlob("templates/*")

	// 主页路由
	router.GET("/", func(c *gin.Context) {
		// 计算今天的喝水总量
		total := getTodayTotal()

		// 渲染HTML模板
		c.HTML(http.StatusOK, "index.html", gin.H{
			"total":   total,
			"records": getTodayRecords(),
		})
	})

	// 添加喝水记录的路由
	router.POST("/add", func(c *gin.Context) {
		// 获取表单数据
		amountStr := c.PostForm("amount")
		amount, err := strconv.Atoi(amountStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的喝水量"})
			return
		}

		// 添加到数据库
		_, err = db.Exec("INSERT INTO water_records (amount, time) VALUES (?, ?)", amount, time.Now())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "添加记录失败"})
			return
		}

		// 重定向到主页
		c.Redirect(http.StatusFound, "/")
	})

	// 创建HTML模板
	createHTMLTemplate()

	// 启动服务器
	fmt.Println("服务器运行在 http://localhost:8080")
	router.Run(":8080")
}

// 获取今天的喝水记录
func getTodayRecords() []WaterRecord {
	var records []WaterRecord
	// 查询今天的所有记录
	todayStart := time.Now().Format("2006-01-02") + " 00:00:00"
	rows, err := db.Query("SELECT id, amount, time FROM water_records WHERE time >= ? ORDER BY time DESC", todayStart)
	if err != nil {
		fmt.Println("查询记录失败:", err)
		return records
	}
	defer rows.Close()

	// 遍历结果
	for rows.Next() {
		var record WaterRecord
		var timeStr string

		err := rows.Scan(&record.ID, &record.Amount, &timeStr)
		if err != nil {
			fmt.Println("扫描记录失败:", err)
			continue
		}

		// 解析时间
		record.Time, err = time.Parse("2006-01-02 15:04:05", timeStr)
		if err != nil {
			fmt.Println("解析时间失败:", err)
			continue
		}

		records = append(records, record)
	}

	return records
}

// 获取今天的喝水总量
func getTodayTotal() int {
	total := 0
	// 查询今天的总喝水量
	todayStart := time.Now().Format("2006-01-02") + " 00:00:00"
	err := db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM water_records WHERE time >= ?", todayStart).Scan(&total)
	if err != nil {
		fmt.Println("查询总量失败:", err)
		return 0
	}

	return total
}

// 创建HTML模板文件
func createHTMLTemplate() {
	// 这里我们直接在路由中使用字符串模板，而不是创建文件
}
