package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type endPoint struct {
	IPAddr string
	Port   int
}

type config struct {
	RX1 endPoint
	TX1 endPoint
	RX2 endPoint
	TX2 endPoint
}

func listen(ep endPoint) net.PacketConn {
	conn, err := net.ListenPacket("udp", fmt.Sprintf("%s:%d", ep.IPAddr, ep.Port))
	if err != nil {
		panic(err)
	}
	return conn
}

func dial(ep endPoint) net.Conn {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", ep.IPAddr, ep.Port))
	if err != nil {
		panic(err)
	}
	return conn
}

func forward(name string, rx endPoint, tx endPoint) {
	if rx.Port == 0 || tx.Port == 0 {
		return
	}
	src := listen(rx)
	defer src.Close()

	dst := dial(tx)
	defer dst.Close()

	buffer := make([]byte, 1500)
	fmt.Printf("%s start %v -> %v\n", name, rx, tx)
	for {
		length, _, _ := src.ReadFrom(buffer)
		fmt.Printf("%v <- %v\n%s", dst.RemoteAddr(), src.LocalAddr(), hex.Dump(buffer[:length]))
		dst.Write(buffer[:length])
	}
}

func main() {
	var config config
	_, err := toml.DecodeFile("fwudp.toml", &config)
	if err != nil {
		panic(err)
	}

	go forward("forward 1", config.RX1, config.TX1)
	time.Sleep(100 * time.Millisecond)
	go forward("forward 2", config.RX2, config.TX2)

	scan := bufio.NewScanner(os.Stdin)
	scan.Scan()
	fmt.Println("終了します")
}
