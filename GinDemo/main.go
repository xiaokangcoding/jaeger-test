package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"log"
	"net/http"
	"github.com/thinkerou/favicon"
)

// 请求体的数据结构
type Login struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// 错误处理中间件
func ErrorHandling() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": c.Errors})
			c.Abort()
		}
	}
}
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var Users = []User{
	{ID: "1", Name: "Tom"},
	{ID: "2", Name: "Jerry"},
}

func main() {
	router := gin.Default()
	router.Use(favicon.New("./images.ico"))
	// 设置跨域处理
	router.Use(cors.Default())
	// 使用错误处理中间件
	router.Use(ErrorHandling())

	//加载静态文件
	router.LoadHTMLGlob("templates/*")


	router.GET("/index", func(c *gin.Context) {
		// 使用HTML模板文件
		// 第一个参数是HTTP状态码
		// 第二个参数是模板文件的名字（我们在LoadHTMLGlob调用中已经加载了）
		// 第三个参数是你想在模板文件中使用的数据
		c.HTML(200, "index.html",gin.H{"msg":"你好"} )
	})

	router.POST("/user/add", func(context *gin.Context) {
		userName := context.PostForm("username")
		password := context.PostForm("password")

		context.JSON(http.StatusOK,gin.H{
			"username":userName,
			"password" :password,
		})
	})

	//http://localhost:8080/login?userid=1&name=carl
	router.GET("/login", func(context *gin.Context) {
		userid := context.Query("userid")
		name := context.Query("name")

		context.JSON(http.StatusOK,gin.H{
			"userid" : userid,
			"name" : name,
		})
	})
	var mp map[string]string
	router.POST("/login", func(context *gin.Context) {
		data , _ := context.GetRawData()
		log.Println(data)
		json.Unmarshal(data,&mp)
		log.Println(mp)
		context.JSON(http.StatusOK,mp)
	})

	// 设置路由分组
	api := router.Group("/api")
	{
		api.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Hello Gin!"})
		})

		api.POST("/login", func(c *gin.Context) {
			var json Login
			if err := c.ShouldBindJSON(&json); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			// 检查用户名密码
			if json.User != "user" || json.Password != "password" {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
		})
	}


	// RestFul Api
	// 获取所有用户 http://localhost:8080/users
	router.GET("/users", func(c *gin.Context) {
		c.JSON(200, gin.H{"users": Users})
	})

	// 获取单个用户 http://localhost:8080/users/1
	router.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, user := range Users {
			if user.ID == id {
				c.JSON(200, gin.H{"user": user})
				return
			}
		}
		c.Status(404)
	})

	// 创建新用户 http://localhost:8080/users  body:  {
	//     "id": "3",
	//      "name": "Carl"
	//}
	router.POST("/users", func(c *gin.Context) {
		var newUser User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		Users = append(Users, newUser)
		c.JSON(200, newUser)
	})

	// 更新用户 http://localhost:8080/users/3   body:  {
	//     "id": "5",
	//      "name": "jk"
	//}
	router.PUT("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var updateUser User
		if err := c.ShouldBindJSON(&updateUser); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		for i, user := range Users {
			if user.ID == id {
				Users[i] = updateUser
				c.JSON(200, updateUser)
				return
			}
		}
		c.Status(404)
	})

	// 删除用户
	router.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		for i, user := range Users {
			if user.ID == id {
				Users = append(Users[:i], Users[i+1:]...)
				c.Status(200)
				return
			}
		}
		c.Status(404)
	})

	//
	router.Run(":8080")
}
