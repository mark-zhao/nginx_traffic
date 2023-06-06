package traffic

import (
	"glog"
	"io/ioutil"
	"nginx_traffic/common"
	"os"
	"strconv"
	"strings"
)
type ReadFile  struct{
}
type F interface {
	ReadLine(fileName string, handler func([]string)) error
	Todo(line []string)
	CheckFileIsExist(filename string) bool
}
var PartitionPoint = 0
func (R *ReadFile) ReadFile(fileName string, handler func([]string)) error {
	if R.CheckFileIsExist(fileName) { //如果文件存在
		glog.Info("文件存在")
		lines, err := ioutil.ReadFile(fileName)
		if err != nil {
			glog.Error("read file err")
			return err
		} else {
			contents := string(lines)
			lines := strings.Split(contents, "\n")
			if len(lines) >PartitionPoint {
				handler(lines[PartitionPoint:])
			}else {
				handler(lines)
			}
			handler(lines[PartitionPoint:])
			PartitionPoint = len(lines)
		}
	}else {
		glog.Info("文件不存在")
	}
	return nil
}

func (R *ReadFile) Todo(lines []string) {
	var RF []*Fields
	b := new(Fields)
	conn := connInflux()
	defer conn.Close()
	for _, line := range lines {
		line := strings.Split(line, " ")
		if len(line[0]) == 0 {
			continue
		}
		IP, err := common.StrToIP(line[0])
		if err != nil {
			glog.Error("ip 转换失败")
		} else {
			if Res := common.IsPublicIP(IP); Res {
				bucket := strings.Split(line[6], "/")
				size, err := strconv.ParseInt(line[8], 10, 64)
				if err != nil {
					glog.Error("size string 转换int64失败")
				}
				timestamp := strings.Split(line[3], "[")
				b = &Fields{BucketName: bucket[1], size: size, timestamp: timestamp[1]}
				RF = append(RF, b)
				glog.Info(bucket[1], size, timestamp[1])
			}
		}
	}
	InsertInfluxdb(RF, conn)
}

func (R *ReadFile) CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
