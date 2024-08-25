package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"edgefusion-video-push/config"
	"edgefusion-video-push/server"
	"edgefusion-video-push/service"
	"gopkg.in/yaml.v3"
)

func main() {
	group := sync.WaitGroup{}
	group.Add(1)
	defer group.Done()
	// 加载配置文件
	yamlFile, err := ioutil.ReadFile("etc/conf.yml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}
	// 解析 YAML 文件
	var cfg config.Config
	if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
		log.Fatalf("Error unmarshalling YAML data: %v", err)
	}
	queue := service.NewQueue()

	var transmit *server.TransmitInfo
	var push *server.CommandStatus
	var sysPush, userPush string
	if cfg.Push.IsCloudLive == "1" || cfg.Push.IsCloudStorage == "1" {
		path := GetRtmpPutPath(cfg)
		sysPush = path
	}
	if cfg.Push.DistributionSetting {
		switch cfg.Push.CloudLiveMode {
		case "0":
			if sysPush != "" {
				userPush = cfg.Push.InputSrc
			} else {
				sysPush = cfg.Push.InputSrc
			}
			log.Printf("视频[ 直播推流rtmp ] 启动")
		case "1":
			transmit = server.NewTransmit(cfg.Push.InputSrc)
			log.Printf("视频[ 透传转发 ] 启动: %v", transmit.RemoteAddr)
		default:
			log.Printf("错误的启动类型")
		}
	}
	if sysPush != "" && userPush != "" {
		push = server.NewPushRtmp(sysPush, userPush)
	} else if sysPush != "" {
		push = server.NewOnePushRtmp(sysPush)
	}
	lister := server.NewLister()
	//启动数据接收
	go lister.Lister(queue)
	go server.Consume(lister, queue, transmit, push)

	group.Wait()
}

func GetRtmpPutPath(cfg config.Config) string {
	//节点id
	NodeId := os.Getenv("EF_NODE_ID")
	//所属应用名称
	AppName := os.Getenv("EF_APP_NAME")
	//服务名称
	ServiceName := os.Getenv("EF_SERVICE_NAME")
	if cfg.Push.IsCloudStorage == "1" {
		//拼接带录播的直播地址
		return fmt.Sprintf("rtmp://%s:1935/%s-%s/%s?vhost=edgefusiondvr", cfg.Push.CloudAddress, NodeId, AppName, ServiceName)
	} else {
		//拼接不录播的直播地址
		return fmt.Sprintf("rtmp://%s:1935/%s-%s/%s?vhost=edgefusion", cfg.Push.CloudAddress, NodeId, AppName, ServiceName)
	}
}
