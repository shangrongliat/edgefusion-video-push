package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"edgefusion-video-push/config"
	"edgefusion-video-push/server"
	"edgefusion-video-push/service"
	"gopkg.in/yaml.v3"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() / 2)
	group := sync.WaitGroup{}
	initLog(false)
	// 设置 log 包的日志输出
	group.Add(1)
	defer group.Done()
	// 加载配置文件
	//yamlFile, err := ioutil.ReadFile("./config.yml")
	yamlFile, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}
	// 解析 YAML 文件
	var cfg config.Config
	if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
		log.Fatalf("Error unmarshalling YAML data: %v", err)
	}
	queue := service.NewQueue()

	lister := server.NewLister()
	time.Sleep(1 * time.Millisecond)
	forward := server.NewForward()
	//启动数据接收
	go lister.Lister(queue, forward)

	transmit, localTransmit, push := server.PushInit(cfg)
	forward.SetTransmitAddr(transmit, 1)
	forward.SetTransmitAddr(localTransmit, 2)

	go server.Consume2(push, queue, cfg)

	group.Wait()
}

func initLog(terminal bool) {
	// 构建日志文件的完整路径
	logFilePath := filepath.Join("/etc/edgefusion/video/push/", "logs", "app.log")
	//logFilePath := filepath.Join("D:\\go-project\\edgefusion\\edgefusion-video-push", "logs", "app.log")
	// 创建文件夹 "logs" 如果它不存在
	err := os.MkdirAll(filepath.Dir(logFilePath), 0755)
	if err != nil {
		log.Fatalf("Error creating logs folder: %v", err)
	}
	// 打开一个文件用于写入日志
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	// 设置 log 包的日志输出
	log.SetOutput(logFile)
	if terminal {
		// 创建一个 io.MultiWriter 实例，它允许我们将日志输出到多个地方
		multiWriter := io.MultiWriter(os.Stdout, logFile)

		// 设置 log 包的日志输出
		log.SetOutput(multiWriter)
	}
}
