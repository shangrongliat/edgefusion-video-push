package decode

import (
	"errors"
	"fmt"
)

type CustomRTP struct {
	Header    uint8
	Mark      uint8  // 1b  *
	Seq       uint16 // 16b **
	Timestamp uint32 // 32b **** samples
	Ssrc      uint32 // 32b **** Synchronization source
	Payload   []byte //
}

func Decode(data []byte) (c *CustomRTP, err error) {
	fmt.Println(data)
	if len(data) > 12 {
		//c.Header = data[0]
		//c.Mark = data[1]
		//c.Seq = binary.BigEndian.Uint16(data[2:4])
		//c.Timestamp = binary.BigEndian.Uint32(data[4:8])
		//c.Ssrc = binary.BigEndian.Uint32(data[8:12])
		c.Payload = data[12:]
	} else {
		err = errors.New("数据长度异常")
	}
	return
}
