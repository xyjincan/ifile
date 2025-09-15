package tool

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// PC端管理配置接口
// 127.0.0.1网络
// TODO：禁止跨站请求
// const res = await post("/api/admin/remove_dir?id=" + tmp.id);
func apiFail(c *gin.Context) {
	c.JSON(500, gin.H{
		"msg": "fail",
	})
}

func removeBaseDir(c *gin.Context) {
	if !passAdminAPI(c) {
		apiFail(c)
		return
	}
	var base, _ = c.GetQuery("id")
	var item RootPath
	newPaths := []RootPath{}
	for _, p := range config.Paths {
		if p.Id == base {
			item = p
			item.Path = "ifile"
			// 移除元素，返回新的切片
		} else {
			newPaths = append(newPaths, p)
		}
	}
	config.Paths = newPaths
	save()
	fmt.Println("remove_dir", base)
	var res = gin.H{
		"item_id": base,
		"item":    item,
	}
	res["code"] = 200
	c.JSON(http.StatusOK, res)
}

// todo 重启后生效
func addBaseDir(c *gin.Context) {
	if !passAdminAPI(c) {
		apiFail(c)
		return
	}
	var newItem RootPath
	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	config = load()
	var same = false
	for _, value := range config.Paths {
		if newItem.Id == value.Id || newItem.Path == value.Path {
			fmt.Println("重复项！")
			same = true
			break
		}
	}
	if !same {
		_, err := os.ReadDir(newItem.Path)
		if err != nil {
			fmt.Println("new_fs_fail:", newItem.Name, " fail_path:", newItem.Path)
			newItem.NotExist = true
		}
	}
	// file_not_exists
	var res = gin.H{
		"duplicate": same,
		"item":      newItem,
	}
	if !same && !newItem.NotExist {
		if strings.TrimSpace(newItem.Name) == "" {
			filename := filepath.Base(newItem.Path)
			newItem.Name = filename
		}
		config.Paths = append(config.Paths, newItem)
		save()
		res["code"] = 200
		c.JSON(http.StatusOK, res)
	} else {
		res["code"] = 500
		c.JSON(http.StatusInternalServerError, res)
	}

}

func restartApp(c *gin.Context) {
	if !passAdminAPI(c) {
		apiFail(c)
		return
	}
	go restart()
	c.JSON(200, gin.H{
		"msg": "重启中",
	})
}
