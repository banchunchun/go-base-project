package im

import (
	"github.com/go-netty/go-netty/transport"
	"github.com/go-netty/go-netty/transport/tcp"
	"time"

	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec/frame"
)

func Connect(host string, port int32) {
	pipelineInitializer := func(channel netty.Channel) {
		channel.Pipeline().AddLast(frame.DelimiterCodec(10240, "$$", true)).
			AddLast(netty.ReadIdleHandler(time.Second), netty.WriteIdleHandler(time.Second)).
			AddLast(&PackageDecoder{}).AddLast(&PackageEncoder{"server"})
	}

	//tcpOptions := &tcp.Options{
	//	Timeout:         time.Second * 3,
	//	KeepAlive:       true,
	//	KeepAlivePeriod: time.Second * 5,
	//	Linger:          0,
	//	NoDelay:         true,
	//	SockBuf:         1024,
	//}

	bs := netty.NewBootstrap(
		netty.WithChildInitializer(pipelineInitializer),
		netty.WithClientInitializer(pipelineInitializer),
		netty.WithTransport(tcp.New()),
	)

	ch, err := bs.Connect("tcp://127.0.0.1:5060", transport.WithAttachment("go-netty"))
	if err != nil {
		panic(err)
	}
	login := &LoginRequest{
		UserId:   "12345677777",
		NickName: "banchun",
	}
	message := BuildByteMessage("LOGIN", login)
	ch.Write(message)

	time.Sleep(time.Minute * 100)

	bs.Shutdown()
}
