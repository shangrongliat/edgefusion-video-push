package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"edgefusion-video-push/config"
	"edgefusion-video-push/utils"
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
			addr := domainResolver(cfg.Push.InputSrc)
			transmit, err = NewTransmit(fmt.Sprintf("%s:65506", addr))
			if err != nil {
				log.Printf("视频[ 透传转发 ] 启动失败,转发地址: %s:65506 \n", addr)
			}
			log.Printf("视频[ 透传转发 ] 启动: %v \n", transmit)
		default:
			log.Println("错误的启动类型")
		}
	}
	log.Printf("系统直播地址:%s", sysPush)
	log.Printf("用户直播地址:%s", userPush)
	if sysPush != "" && userPush != "" && utils.IsRTMPURLValid(userPush) {
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
	// 默认直播标记为0
	var liveFlag = "0"
	if cfg.Push.IsCloudLive {
		// 如果开启直播则标记为1
		liveFlag = "1"
	}
	//服务名称
	streamStr := fmt.Sprintf("%s@%s@%s@%s@%s", cfg.Video.VencWith, cfg.Video.VencHeight, cfg.Video.VencFps, cfg.Video.SrcType, liveFlag)
	log.Printf("分辨率: with: %s hight: %s fps: %s liveFlag: %v", streamStr, cfg.Push.Height, cfg.Push.Fps, liveFlag)
	stream, err := base64ToHex(streamStr)
	if err != nil {
		log.Printf("stream序列化失敗", err)
	}
	if cfg.Push.IsCloudStorage {
		//拼接带录播的直播地址
		return fmt.Sprintf("rtmp://%s:1935/%s@%s/%s?vhost=edgefusiondvr", cfg.Push.CloudAddress, NodeId, AppName, stream)
	} else {
		//拼接不录播的直播地址
		return fmt.Sprintf("rtmp://%s:1935/%s@%s/%s?vhost=edgefusion", cfg.Push.CloudAddress, NodeId, AppName, stream)
	}
}

func base64ToHex(s string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(s)), nil
}

func domainResolver(domain string) string {
	// Consul DNS 服务器地址和端口
	consulDNSServer := "127.0.0.1:8600"
	// 要解析的域名
	//domain := "montage01.service.consul"
	// 创建一个自定义的 DNS 解析器
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 30 * time.Second,
			}
			return d.Dial("udp", consulDNSServer)
		},
	}
	// 解析域名
	addrs, err := resolver.LookupHost(context.Background(), domain)
	if err != nil {
		fmt.Printf("Error resolving %s: %v\n", domain, err)
		return ""
	}
	// 打印解析结果
	if len(addrs) == 0 {
		return ""
	}
	return addrs[0]
}
