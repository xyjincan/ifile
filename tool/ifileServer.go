package tool

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Root  string     `json:"root"`
	Paths []RootPath `json:"paths"`
}

// 全局配置实例
// TODO 优化list -> map
var config Config

/*
加载配置信息
*/
func load() Config {
	// 打开 JSON 文件
	file, err := os.Open("./config.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return config
	}
	defer file.Close()
	// 解码 JSON 数据
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return config
	}
	fmt.Println("Paths:", config.Paths)
	return config
}

func save() error {
	file, err := os.Create("./config.json")
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // 格式化输出为易读的 JSON
	err = encoder.Encode(&config)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}
	fmt.Println("配置已保存")
	return nil
}

/*
启动服务
*/
func StartHttpd() {
	gin.SetMode(gin.ReleaseMode)
	config = load()
	r := gin.Default()
	// 配置 CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://127.0.0.1:8000", "http://127.0.0.1:8080",
			"http://localhost:8000", "http://localhost:8080",
		}, // 允许的域名
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},                  // 允许的方法, "PUT"  "DELETE",
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // 允许的请求头
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},          // 可被客户端访问的响应头
		AllowCredentials: true,                                                // 是否允许传递 Cookie
		MaxAge:           1 * time.Hour,                                       // 预检请求的缓存时间
	}))
	r.Static("/ui", "./static")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")
	// 基础目录webdav
	for index, value := range config.Paths {
		_, err := os.ReadDir(value.Path)
		if err != nil {
			fmt.Println("init_web_fs_fail:", value.Name, " fail_path:", value.Path)
			config.Paths[index].NotExist = true
			continue
		}
		fmt.Println("init_web_fs: " + value.Name)
		//循环添加 静态文件服务
		r.StaticFS(value.Id, http.Dir(value.Path))
	}
	// 基础目录信息接口
	r.GET("/fs/paths", apiPaths)
	// ifile接口
	r.GET("/fs/delete", deleteFile)
	r.GET("/fs/path", listFiles)
	r.GET("/fs/file", viewFile)

	// admin config
	r.POST("/api/admin/add_dir", addBaseDir)
	r.POST("/api/admin/remove_dir", removeBaseDir)
	r.POST("/api/admin/restart", restartApp)
	r.GET("/api/v1", configInfoV1)
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/ui/")
	})
	r.GET("/ping", func(c *gin.Context) {
		//go restart()
		c.JSON(200, gin.H{
			"msg": "pong",
		})
	})
	err := r.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
}

/*
重启当前进程
*/
func restart() {
	args := os.Args
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("无法获取可执行文件路径: %v", err)
	}
	log.Println("ifile", execPath)
	isDev := isGoRun(execPath)
	log.Println("goRun", isDev)
	if isDev {
		ppid_info_win()
		os.Exit(0)
		return
	}
	cmd := exec.Command(execPath, args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		log.Fatalf("重启失败: %v", err)
	}
	time.Sleep(3 * time.Second)
	log.Println("重启成功，退出旧进程")
	os.Exit(0)
}

func isGoRun(exePath string) bool {
	// 获取可执行文件的路径 // 检查路径是否包含临时目录的特征
	// 注意：这种方法不保证100%准确，但可以作为一种有效的启发式判断
	var goRun = false
	if strings.Contains(exePath, filepath.Join("go-build")) || strings.Contains(exePath, filepath.Join("go-run")) {
		goRun = true
	}
	if strings.Contains(exePath, filepath.Join("\\AppData\\Local\\Temp\\")) {
		goRun = true
	}
	return goRun
}

func ppid_info_win() {
	ppid := os.Getppid()
	cmd := exec.Command("wmic", "process", "where", fmt.Sprintf("(processid=%d)", ppid), "get", "name")
	output, err := cmd.Output()
	if err != nil {
		return
	}
	// 解析 wmic 的输出，它通常包含 "Name" 字段和进程名称
	lines := strings.Split(string(output), "\n")
	if len(lines) > 1 {
		parentProcessName := strings.TrimSpace(lines[1])
		if parentProcessName != "" {
			fmt.Printf("父进程: %s\n", parentProcessName)
			return
		}
	}
}
