package im

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/go-netty/go-netty"
)

type PackageEncoder struct {
	Name string
}

func (x *PackageEncoder) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	fmt.Println("encoder in")
	var immessage *ImMessage
	v := message.([]byte)
	json.Unmarshal(v, &immessage)
	buffer := bytes.NewBuffer(make([]byte, 0, len(v)))
	//写入头
	buffer.Write([]byte(immessage.Header.Header))
	buffer.Write([]byte(immessage.Header.Crc))
	//以大端的方式写入协议体的长度
	binary.Write(buffer, binary.BigEndian, immessage.Header.Length)
	buffer.Write([]byte(immessage.Command))
	buffer.Write([]byte(immessage.Content))
	//固定消息体分割
	buffer.Write([]byte(DelimiterMagicHeader))
	ctx.HandleWrite(buffer.Bytes())
}
