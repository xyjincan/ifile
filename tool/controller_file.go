package tool

import (
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// 提供文件下载或在线查看
func viewFile(c *gin.Context) {
	var base = c.Query("base")
	var parent = c.DefaultQuery("view", "")
	var root = "" //
	var check_base = false
	for _, value := range config.Paths {
		if base == value.Id {
			root = value.Path
			check_base = true
		}
	}
	if !check_base {
		var res = gin.H{
			"base": base,
			"view": parent,
			"code": 404,
		}
		c.JSON(404, res)
		return
	}
	targetPath := filepath.Join(root, parent)
	cleanPath := filepath.Clean(targetPath)
	if !strings.HasPrefix(cleanPath, root) {
		c.JSON(400, gin.H{"error": "Invalid path specified"})
		return
	}
	view_file(c, cleanPath)
}

// 提供视频 video/*
// 支持http range 断点续传
// 尽快关闭文件，消除文件锁
// 采取措施，提高视频文件随机读取速度
// 你可以扩展一些函数
func view_file(c *gin.Context, cleanPath string) {
	// 打开文件 (O_RDONLY 避免文件锁)
	f, err := os.OpenFile(cleanPath, os.O_RDONLY, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "open_file_error"})
		return
	}
	defer f.Close()
	// 获取文件信息
	stat, err := f.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 设置 MIME 类型 (video/*)
	contentType := mime.TypeByExtension(filepath.Ext(cleanPath))
	fmt.Println("contentType", contentType)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	c.Header("Cache-Control", "no-store, no-cache")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Header("Content-Type", contentType) // 你可以根据扩展名判断，比如 .webm .avi
	c.Header("Accept-Ranges", "bytes")
	filename := stat.Name()
	safeFilename := url.PathEscape(filename)
	c.Header("Content-Disposition", fmt.Sprintf(`inline; filename*=UTF-8''%s`, safeFilename))
	// 处理 Range 请求
	http.ServeContent(c.Writer, c.Request, filename, stat.ModTime(), f)
}
