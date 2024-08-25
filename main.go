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
	lister := server.NewLister()
	//go server.RTSPServer(&group, videosteam)
	var transmit *server.TransmitInfo
	if cfg.Push.IsCloudLive == "1" || cfg.Push.IsCloudStorage == "1" {
		path := GetRtmpPutPath(cfg)
		sysPush := server.NewPushRtmp(path)
	}
	if cfg.Push.DistributionSetting {
		switch cfg.Push.CloudLiveMode {
		case "0":
			transmit = server.NewTransmit("")
			log.Fatalf("视频[ 直播推流rtmp ] 启动")
		case "1":
			transmit = server.NewTransmit(cfg.Push.InputSrc)
			log.Fatalf("视频[ 透传转发 ] 启动")
		}
	}
	//启动数据接收
	go lister.Lister(&group, queue)

	//// 判断是否开启视频多路分发
	//if cfg.DistributionSetting {
	//	fmt.Println("distribution setting")
	//	go func() {
	//		switch cfg.CloudLiveMode {
	//		case "0":
	//			server.Transmit(&group, cfg.InputSrc)
	//			log.Fatalf("视频[ 透传转发 ] 启动")
	//		case "1":
	//			server.PushRtmp(&group, cfg.InputSrc)
	//			log.Fatalf("视频[ 直播推流rtmp ] 启动")
	//		}
	//	}()
	//}
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
