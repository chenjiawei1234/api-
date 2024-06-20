package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
)

func ServerRun(port string) {
	app := gin.Default()
	app.POST("/api/test", click)
	app.GET("/api/report", getReport)
	app.GET("/download/:reportname", downReport)
	app.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "你使用了错误的请求方法或者请求路径",
		})
	})
	err := app.Run(port)
	if err != nil {
		log.Print(err)
	}
}

const filePath string = `.\testsuit\`

func click(c *gin.Context) {

	fVal := c.PostForm("f")
	newfVal := fmt.Sprintf("%s%s", filePath, fVal)
	cVal := c.PostForm("c")
	if fVal == "" || cVal == "" {
		c.JSON(400, gin.H{
			"error": "未指定参数",
		})
		return
	} else if checkc(cVal) {
		c.JSON(400, gin.H{
			"error": "参数c类型错误",
		})
		return
	} else {
		dirname := "./testsuit" // 目录路径

		files, err := os.ReadDir(dirname)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "配置文件不存在",
			})
			return
		}

		for _, file := range files {
			if !file.IsDir() {
				if file.Name() == fVal {
					cmd := exec.Command(`.\runner.exe`, "-f", newfVal, "-c", cVal)
					outputVal, err := cmd.Output()
					if err != nil {
						exitError, ok := err.(*exec.ExitError)
						if ok {
							stderr := string(exitError.Stderr)
							c.JSON(500, gin.H{
								"stderr": stderr,
							})
							return
						}
					}
					c.JSON(200, struct {
						Body string `json:"body"`
					}{Body: string(outputVal)})
					return
				}

			}

		}
	}
	c.JSON(400, gin.H{
		"error": "输入配置文件名错误",
	})
}

func checkc(val string) bool {
	_, err := strconv.ParseBool(val)
	if err != nil {
		return true
	} else {
		return false
	}
	return false
}

// 查询所有报告
func getReport(c *gin.Context) {
	dirname := "./report" // 目录路径

	files, err := os.ReadDir(dirname)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "读取日志文件失败",
		})
		return
	}
	if len(files) == 0 {
		c.JSON(400, gin.H{
			"body": nil,
		})

	} else {
		var reportName []string
		for _, file := range files {
			if !file.IsDir() {
				reportName = append(reportName, file.Name())
			}
		}
		c.JSON(400, reportName)
		return
	}
}

// 下载指定报告
func downReport(c *gin.Context) {
	reporName := c.Param("reportname")
	reporPath := path.Join("./report", reporName)
	// 检查文件是否存在，如果不存在则返回错误
	if _, err := os.Stat(reporPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"status": "日志文件不存在"})
		return
	}

	// 设置响应头，让浏览器处理文件下载
	c.Header("Content-Disposition", "attachment; filename="+reporName)
	c.Header("Content-Type", "application/octet-stream")
	c.File(reporPath)
}
