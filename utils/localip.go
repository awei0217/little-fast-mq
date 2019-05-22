package utils

import (
	"fmt"
	"net"
	"os"
)

func GetLocalIp() string {
	address, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println("get local ip error：", err)
		os.Exit(1)
	}
	for _, address := range address {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	fmt.Println("get local ip is null")
	return ""
}
