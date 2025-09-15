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
	var base, _ = c.GetQuery("base")
	var file, _ = c.GetQuery("file")
	var path = ""
	for _, value := range config.Paths {
		if base == value.Id {
			path = value.Path
		}
	}
	var delFile = path + "\\" + file
	fmt.Println("delete: " + delFile)
	var del_err = os.Remove(delFile)
	if del_err != nil {
		fmt.Println("删除文件失败:", del_err)
		var delMark = delFile + ".ifile_delete"
		del, _ := os.Create(delMark)
		defer del.Close()
		go deferRemoveFile(delFile)
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

func listFiles(c *gin.Context) {
	var base = c.Query("base")
	var parent = c.DefaultQuery("view", "")
	// parent to abs path
	var root = "" //
	var check_base = false
	for _, value := range config.Paths {
		if base == value.Id {
			root = value.Path
			check_base = true
		}
	}
	_root := filepath.Clean(root)
	targetPath := filepath.Join(_root, parent)
	cleanPath := filepath.Clean(targetPath)
	relPath, err := filepath.Rel(_root, targetPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path specified"})
		return
	}
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
