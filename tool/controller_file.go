package tool

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
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
	_root := filepath.Clean(root)
	targetPath := filepath.Join(_root, parent)
	cleanPath := filepath.Clean(targetPath)

	relPath, err := filepath.Rel(_root, cleanPath) // 禁止穿越、逃逸
	if err != nil || strings.HasPrefix(relPath, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_file"})
		return
	}
	realPath, err := filepath.EvalSymlinks(cleanPath) // 禁止链接
	if err != nil || !strings.EqualFold(realPath, cleanPath) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_file"})
		return
	}
	fileInfo, err := os.Stat(cleanPath)
	if err != nil || fileInfo.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "access_denied_file"})
		return
	}

	view_file(c, cleanPath)
}

// 支持http range 断点续传
// 尽快关闭文件，消除文件锁
// 采取措施，提高视频文件随机读取速度
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
	// fmt.Println("contentType", contentType)
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
	etag := fmt.Sprintf(`"%x-%x"`, stat.ModTime().Unix(), stat.Size())
	c.Header("ETag", etag)
	c.Header("Last-Modified", stat.ModTime().UTC().Format(http.TimeFormat))

	// 默认返回整个文件
	// 解析 Range: bytes=start-end，完整兼容各种Range格式
	fileSize := stat.Size()
	status := http.StatusOK

	rangeHeader := c.GetHeader("Range")
	start, end, isPartial := parseRange(rangeHeader, fileSize)
	if isPartial {
		status = http.StatusPartialContent
		c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	}
	if rangeHeader != "" {
		log.Printf("view_file_range: %s", rangeHeader)
	} else {
		log.Printf("view_file_no_range")
	}
	// 设置响应头
	contentLength := end - start + 1
	c.Header("Content-Length", strconv.FormatInt(contentLength, 10))
	c.Status(status)
	// 返回指定分片
	sr := io.NewSectionReader(f, start, contentLength)
	io.Copy(c.Writer, sr)
}

func parseRange(rangeHeader string, fileSize int64) (int64, int64, bool) {
	if rangeHeader == "" || !strings.HasPrefix(rangeHeader, "bytes=") {
		return 0, fileSize - 1, false
	}
	// 去掉 "bytes="
	r := strings.TrimPrefix(rangeHeader, "bytes=")
	parts := strings.Split(r, "-")
	if len(parts) != 2 {
		return 0, fileSize - 1, false
	}
	startStr := strings.TrimSpace(parts[0])
	endStr := strings.TrimSpace(parts[1])

	var start, end int64
	if startStr == "" {
		// 格式: -500 (获取最后 500 个字节)
		val, err := strconv.ParseInt(endStr, 10, 64)
		if err != nil {
			return 0, fileSize - 1, false
		}
		start = fileSize - val
		end = fileSize - 1
	} else {
		// 格式: 500- 或 500-999
		s, err := strconv.ParseInt(startStr, 10, 64)
		if err != nil {
			return 0, fileSize - 1, false
		}
		start = s
		if endStr == "" {
			// 格式: 500- (从 500 到文件结尾)
			end = fileSize - 1
		} else {
			// 格式: 500-999
			e, err := strconv.ParseInt(endStr, 10, 64)
			if err != nil {
				return 0, fileSize - 1, false
			}
			end = e
		}
	}
	// 边界修正
	if start < 0 {
		start = 0
	}
	if end >= fileSize {
		end = fileSize - 1
	}
	if start > end {
		return 0, fileSize - 1, false
	}
	return start, end, true
}
