package ntp

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

const ntpEpochOffset = 2208988800

type packet struct {
	Settings       uint8
	Stratum        uint8
	Poll           int8
	Precision      int8
	RootDelay      uint32
	RootDispersion uint32
	ReferenceID    uint32
	RefTimeSec     uint32
	RefTimeFrac    uint32
	OrigTimeSec    uint32
	OrigTimeFrac   uint32
	RxTimeSec      uint32
	RxTimeFrac     uint32
	TxTimeSec      uint32
	TxTimeFrac     uint32
}

// GetRemoteTime 从指定的NTP服务器获取精确时间
// 通过UDP协议与NTP服务器通信，获取当前世界标准时间
// 参数 site: NTP服务器地址，如 "pool.ntp.org"
// 返回值: 获取到的时间和可能的错误
func GetRemoteTime(site string) (time.Time, error) {
	// 建立UDP连接到NTP服务器的123端口(标准NTP端口)
	conn, err := net.Dial("udp", site+":123")
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()
	// 设置15秒超时，防止连接挂起
	if err := conn.SetDeadline(time.Now().Add(15 * time.Second)); err != nil {
		return time.Time{}, fmt.Errorf("failed to set deadline: %v", err)
	}

	// 创建NTP请求数据包
	// 0x1B (00011011) 表示:
	// - 版本号为3 (00011xxx)
	// - 模式为客户端 (xxxxx011)
	req := &packet{Settings: 0x1B}

	// 发送NTP请求到服务器
	if err := binary.Write(conn, binary.BigEndian, req); err != nil {
		return time.Time{}, fmt.Errorf("failed to set request: %v", err)
	}

	// 接收NTP服务器的响应
	rsp := &packet{}
	if err := binary.Read(conn, binary.BigEndian, rsp); err != nil {
		return time.Time{}, fmt.Errorf("failed to read server response: %v", err)
	}

	// 将NTP时间转换为UNIX时间
	// NTP时间从1900年开始计算，需要减去与UNIX纪元(1970年)的差值
	secs := float64(rsp.TxTimeSec) - ntpEpochOffset
	// 将NTP小数部分转换为纳秒
	nanos := (int64(rsp.TxTimeFrac) * 1e9) >> 32

	// 使用秒和纳秒创建Go时间对象
	showtime := time.Unix(int64(secs), nanos)

	return showtime, nil
}
