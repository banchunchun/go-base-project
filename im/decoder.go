package im

import (
	"encoding/binary"
	"fmt"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/utils"
)

type PackageDecoder struct {
	Role string
}

func (x *PackageDecoder) HandleRead(ctx netty.InboundContext, message netty.Message) {
	fmt.Println("decoder in")
	//v := utils.MustToBytes(message)
	////if len(v) < 4 {
	////	logger.Log().Infof("invalid package length")
	////	return
	////}
	////buffer := bytes.NewBuffer(v)
	//////immessage := &ImMessage{}
	//magicHeader := v[0:4]
	//crcCode := v[4:8]
	//fmt.Println("magicHeader", string(magicHeader))
	//fmt.Println("crcCode", string(crcCode))
	bufferReader := utils.MustToReader(message)
	magicHeader := make([]byte, 4)
	bufferReader.Read(magicHeader)
	fmt.Println(string(magicHeader))
	crcCode := make([]byte, 4)
	bufferReader.Read(crcCode)
	fmt.Println(string(crcCode))
	var contentLength int64
	binary.Read(bufferReader, binary.BigEndian, &contentLength)
	fmt.Println(contentLength)
	commandBytes := make([]byte, 5)
	bufferReader.Read(commandBytes)
	fmt.Println(string(commandBytes))
	contentBytes := make([]byte, contentLength-5)
	bufferReader.Read(contentBytes)
	fmt.Println(string(contentBytes))
}
