package main

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	teamcityURL   = "http://192.168.17.132:30111/app/rest/buildQueue"
	teamcityToken = "eyJ0eXAiOiAiVENWMiJ9.TXhvNDRNRjNCc3VXVlV1TWhlRGtqZ3otRTVK.NzI0MmE1ZTMtZWEyNS00NjQ5LWExM2QtMWM0NmUwNmM0N2Qx"
)

func main() {
	r := gin.Default()

	r.POST("/webhook", func(c *gin.Context) {
		// 可选：验证 GitHub/GitLab Secret
		// githubSecret := c.GetHeader("X-Hub-Signature-256")
		// gitlabSecret := c.GetHeader("X-Gitlab-Token")

		// 原封不动读取请求体
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("Failed to read request body:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read body"})
			return
		}

		// 构建转发请求
		req, err := http.NewRequest("POST", teamcityURL, bytes.NewReader(bodyBytes))
		if err != nil {
			log.Println("Failed to create request:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}

		// 添加 Authorization Header
		req.Header.Set("Authorization", "Bearer "+teamcityToken)

		// 保留原 Content-Type
		req.Header.Set("Content-Type", c.GetHeader("Content-Type"))
		req.Header.Set("Accept", "application/json")

		// 发请求到 TeamCity
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Failed to call TeamCity:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call TeamCity"})
			return
		}
		defer resp.Body.Close()

		// 将 TeamCity 返回原封不动返回给 Webhook 调用方
		c.Status(resp.StatusCode)
		io.Copy(c.Writer, resp.Body)
	})

	log.Println("Webhook proxy listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
