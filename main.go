package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

var (
	db *gorm.DB
)

func initMysql() (err error) {
	dsn := "root:root@tcp(127.0.0.1:3306)/bubble?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}
func main() {
	// 连接数据库
	err := initMysql()
	if err != nil {
		panic(err.Error())
	}
	// 数据库模型绑定
	db.AutoMigrate(&Todo{})

	// 启动gin服务
	r := gin.Default()
	r.Static("/static", "static")

	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})

	// todolist

	v1Group := r.Group("v1")
	{

		// 添加事项
		v1Group.POST("/todo", func(context *gin.Context) {
			// 从前端填写待办事项 点击提交 请求会发到这里
			// 1.从请求中把数据拿出来
			var todo Todo
			context.BindJSON(&todo)
			// 2.存入数据库
			err := db.Create(&todo).Error
			if err != nil {
				context.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			} else {
				context.JSON(http.StatusOK, todo)
			}
			// 3.返回响应
		})
		// 查看事项
		v1Group.GET("/todo", func(context *gin.Context) {
			var todoList []Todo
			err := db.Find(&todoList).Error
			if err != nil {
				context.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			} else {
				context.JSON(http.StatusOK, todoList)
			}
		})
		// 修改事项
		v1Group.PUT("/todo/:id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{"error": "无效的id"})
				return
			}
			var todo Todo
			if err = db.Where("id=?", id).First(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
				return
			}
			c.BindJSON(&todo)
			if err = db.Save(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todo)
			}
		})
		// 删除事项
		v1Group.DELETE("/todo/:id", func(c *gin.Context) {
			id := c.Param("id")
			if err = db.Where("id=?", id).Delete(Todo{}).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{id: "deleted"})
			}
		})

	}

	r.Run(":8080")
}
