package tool

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func PathList(dirPath string) []FileItem {
	// 目录路径
	// 读取目录下的文件列表，返回 DirEntry 列表
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Println(err)
		return nil
	}
	items := make([]FileItem, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			log.Println("Error getting file info:", err)
			continue
		}
		//fmt.Printf("Name: %s, Size: %d bytes, IsDir: %v, ModTime: %v\n", info.Name(), info.Size(), info.IsDir(), info.ModTime())
		item := NewFileItem()
		item.Name = info.Name()
		truncatedTime := info.ModTime().Truncate(24 * time.Hour)
		if info.IsDir() {
			item.IsDir = true
			item.Time = truncatedTime
		} else {
			item.Size = info.Size()
			item.Time = truncatedTime
		}
		if item.Name == "desktop.ini" {
			continue
		}
		items = append(items, *item)
	}
	return items
}

func deferRemoveFile(delFile string) {
	var count = 60
	var del = false
	for i := 0; i < count; i++ {
		var del_err = os.Remove(delFile)
		if del_err != nil {
			fmt.Println("defer_del:", del_err)
		} else {
			del = true
			break
		}
		time.Sleep(3 * time.Second)
	}
	if del {
		var delMark = delFile + ".ifile_delete"
		os.Remove(delMark)
	}
}

// 判断目标路径是否是父目录的直接子文件夹
func isDirectSubDirectory(parentDir, targetDir string) (bool, error) {
	targetParent := filepath.Dir(targetDir)
	if parentDir == targetParent {
		return true, nil
	}
	// 否则，目标路径不是父目录的直接子目录
	return false, nil
}
