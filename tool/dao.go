package tool

import "time"

type RootPath struct {
	Id       string `json:"id"`        //导航ID
	Name     string `json:"name"`      //
	Path     string `json:"path"`      //物理路径
	NotExist bool   `json:"not_exist"` //
}

type FileItem struct {
	Name  string    `json:"name"`  //
	Size  int64     `json:"size"`  //
	IsDir bool      `json:"isDir"` //
	Path  string    `json:"path"`  //物理路径
	Time  time.Time `json:"time"`
}

func NewRootPath() *RootPath {
	return &RootPath{}
}

func NewFileItem() *FileItem {
	return &FileItem{}
}
