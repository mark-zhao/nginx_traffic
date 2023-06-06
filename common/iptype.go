package common
import (
	"glog"
	"net"
)
//判断ip是否为公网ip
func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

//string转换为net.IP
func StrToIP(ip string) (net.IP,error)  {
	IP, _, err := net.ParseCIDR(ip+"/24")
	if err != nil{
		glog.Info("IP地址解析失败",err)
		return nil, err
	}
	return  IP, nil
}