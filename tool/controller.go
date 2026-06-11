package tool

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// 功能服务接口
// 本地局域网服务
func configInfoV1(c *gin.Context) {
	c.JSON(200, gin.H{
		"view": "/",
	})
}

// 删除指定文件
func deleteFile(c *gin.Context) {
	fmt.Println("delete")
	var base = c.Query("base")
	var file = c.DefaultQuery("file", "")

	var root = "" //
	var check_base = false
	for _, value := range config.Paths {
		if base == value.Id {
			root = value.Path
			check_base = true
			break
		}
	}
	if !check_base {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "error": "Base ID not found"})
		return
	}
	_root := filepath.Clean(root)
	file = filepath.Clean(file)
	targetPath := filepath.Join(_root, file)
	cleanPath := filepath.Clean(targetPath)
	// 相对根目录的子目录
	relPath, err := filepath.Rel(_root, cleanPath) // 禁止目录穿越、逃逸
	if err != nil || strings.HasPrefix(relPath, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_list_fail"})
		return
	}

	// 禁止链接
	realPath, err := filepath.EvalSymlinks(cleanPath)
	if err != nil || !strings.EqualFold(realPath, cleanPath) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_path"})
		return
	}
	fmt.Println("root", root)
	fmt.Println("file", file)
	fmt.Println("targetPath", targetPath)
	fmt.Println("cleanPath", cleanPath)
	var delFile = cleanPath
	fmt.Println("delete: " + delFile)
	var del_err = os.Remove(delFile)
	if del_err != nil {
		fmt.Println("删除文件失败:", del_err)
		fileInfo, _ := os.Stat(delFile)
		if !fileInfo.IsDir() {
			var delMark = delFile + ".ifile_delete"
			del, _ := os.Create(delMark)
			defer del.Close()
			go deferRemoveFile(delFile)
		}
		// 标记删除文件
		c.JSON(500, gin.H{
			"base":  base,
			"view":  "/",
			"error": "delete_fail",
			"file":  file,
		})
		return
	}
	c.JSON(200, gin.H{
		"base": base,
		"view": "/",
		"file": file,
	})
}

func apiPaths(c *gin.Context) {
	newPaths := make([]RootPath, 0, len(config.Paths))
	for _, tmp := range config.Paths {
		tmp.Path = "ifile"
		newPaths = append(newPaths, tmp)
	}
	c.JSON(200, gin.H{
		"list": newPaths,
	})
}

// 文件目录接口，支持递归访问子目录，禁止访问文件
func listFiles(c *gin.Context) {
	var base = c.Query("base")
	var parent = c.DefaultQuery("view", "")

	var root = "" //
	var check_base = false
	for _, value := range config.Paths {
		if base == value.Id {
			root = value.Path
			check_base = true
			break
		}
	}
	if !check_base {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "error": "Base ID not found"})
		return
	}
	_root := filepath.Clean(root)
	parent = filepath.Clean(parent)
	targetPath := filepath.Join(_root, parent)
	cleanPath := filepath.Clean(targetPath)
	// 相对根目录的子目录
	relPath, err := filepath.Rel(_root, cleanPath) // 禁止目录穿越、逃逸
	if err != nil || strings.HasPrefix(relPath, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_list_fail"})
		return
	}

	// 禁止链接
	realPath, err := filepath.EvalSymlinks(cleanPath)
	if err != nil || !strings.EqualFold(realPath, cleanPath) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_path"})
		return
	}

	// 禁止文件
	fileInfo, err := os.Stat(realPath)
	if err != nil || !fileInfo.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "access_denied_file"})
		return
	}

	// 相对根目录的子目录
	relPath = filepath.Join(string(filepath.Separator), relPath)
	if relPath != "/" && relPath != "\\" {
		relPath = relPath + string(filepath.Separator)
	}
	var res = gin.H{
		"base": base,
		"view": strings.ReplaceAll(relPath, "\\", "/"),
	}
	if check_base {
		var items = PathList(cleanPath)
		res["list"] = items
		res["code"] = 200
	} else {
		res["code"] = 404
	}
	c.Header("Cache-Control", "no-store, no-cache")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.JSON(200, res)
}
