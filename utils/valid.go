package utils

import "regexp"

var rtmpURLRegexp = regexp.MustCompile(`^rtmp://([^/:]+)(:(\d+))?/([^/]+)(/.*)?$`)

// IsRTMPURLValid 检查RTMP推流地址格式是否正确
func IsRTMPURLValid(url string) bool {
	return rtmpURLRegexp.MatchString(url)
}
