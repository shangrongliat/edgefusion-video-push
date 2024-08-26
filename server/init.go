package server

import (
	"edgefusion-video-push/config"
	"fmt"
	"log"
	"net"
	"os"
)

func PushInit(cfg config.Config) (transmit, localTransmit *net.UDPAddr, push *CommandStatus) {
	var sysPush, userPush string
	var err error
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
			localTransmit, _ = NewTransmit("127.0.0.1:65525")
			log.Printf("视频[ 直播推流rtmp ] 启动,系统推流地址：%s ,用户推流地址: %s \n", sysPush, userPush)
		case "1":
			transmit, err := NewTransmit(cfg.Push.InputSrc)
			if err != nil {
				log.Printf("视频[ 透传转发 ] 启动失败,转发地址: %s \n", cfg.Push.InputSrc)
			}
			log.Printf("视频[ 透传转发 ] 启动: %v \n", transmit)
		default:
			log.Println("错误的启动类型")
		}
	}
	if sysPush != "" && userPush != "" {
		push, err = NewPushRtmp(sysPush, userPush)
		if err != nil {
			log.Printf("[ 双路 ]直播推流初始化失败. %v", err)
		}
	} else if sysPush != "" {
		push, err = NewOnePushRtmp(sysPush)
		if err != nil {
			log.Printf("[ 单路 ]直播推流初始化失败. %v", err)
		}
	}
	return
}

func RetryPush(cfg config.Config) (push *CommandStatus) {
	var sysPush, userPush string
	var err error
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
			log.Printf("视频[ 直播推流rtmp ] 重试,系统推流地址：%s ,用户推流地址: %s \n", sysPush, userPush)
		default:
			log.Println("错误的启动类型")
		}
	}
	if sysPush != "" && userPush != "" {
		push, err = NewPushRtmp(sysPush, userPush)
		if err != nil {
			log.Printf("[ 双路 ]直播推流重试初始化失败. %v", err)
		}
	} else if sysPush != "" {
		push, err = NewOnePushRtmp(sysPush)
		if err != nil {
			log.Printf("[ 单路 ]直播推流重试初始化失败. %v", err)
		}
	}
	return
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
