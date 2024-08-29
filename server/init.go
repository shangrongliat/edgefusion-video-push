package server

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"edgefusion-video-push/config"
)

func PushInit(cfg config.Config) (transmit, localTransmit *net.UDPAddr, push *CommandStatus) {
	var sysPush, userPush string
	var err error
	if cfg.Push.IsCloudLive || cfg.Push.IsCloudStorage {
		localTransmit, _ = NewTransmit("127.0.0.1:65525")
		path := GetRtmpPutPath(&cfg)
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
			transmit, err = NewTransmit(cfg.Push.InputSrc)
			if err != nil {
				log.Printf("视频[ 透传转发 ] 启动失败,转发地址: %s \n", cfg.Push.InputSrc)
			}
			log.Printf("视频[ 透传转发 ] 启动: %v \n", transmit)
		default:
			log.Println("错误的启动类型")
		}
	}
	log.Printf("系统直播地址:%s", sysPush)
	log.Printf("用户直播地址:%s", userPush)
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
	if cfg.Push.IsCloudLive || cfg.Push.IsCloudStorage {
		path := GetRtmpPutPath(&cfg)
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

func GetRtmpPutPath(cfg *config.Config) string {
	//节点id
	NodeId := os.Getenv("EF_NODE_ID")
	if NodeId == "" {
		NodeId = "test"
	}
	//所属应用名称
	AppName := os.Getenv("EF_APP_NAME")
	if AppName == "" {
		AppName = fmt.Sprintf("%v", time.Now().Unix())
	}
	//服务名称
	streamStr := fmt.Sprintf("%sn%sn%s", cfg.Push.With, cfg.Push.Height, cfg.Push.Fps)
	log.Printf("分辨率: with: %s hight: %s fps: %s", streamStr, cfg.Push.Height, cfg.Push.Fps)
	stream, err := base64ToHex(streamStr)
	if err != nil {
		log.Printf("stream序列化失敗", err)
	}
	if cfg.Push.IsCloudStorage {
		//拼接带录播的直播地址
		return fmt.Sprintf("rtmp://%s:1935/%s-%s/%s?vhost=edgefusiondvr", cfg.Push.CloudAddress, NodeId, AppName, stream)
	} else {
		//拼接不录播的直播地址
		return fmt.Sprintf("rtmp://%s:1935/%s-%s/%s?vhost=edgefusion", cfg.Push.CloudAddress, NodeId, AppName, stream)
	}
}

func base64ToHex(s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
