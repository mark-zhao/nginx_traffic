package main

import (
	"context"
	"encoding/json"
	"glog"
	"io/ioutil"
	"nginx_traffic/common"
	"nginx_traffic/traffic"
	"os"
)
func main() {
	//退出处理
	/*exitChan = make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	go exitHandle()*/
	defer  glog.Flush()

	//http 服务
	glog.Info("http服务 start!")
	ctx, cancel := context.WithCancel(context.Background())
	go common.SetSignal(cancel,exitHandle)
	if err := traffic.StartHTTP(ctx); err != nil {
		glog.Error("start http service failed; ", err)
	}
}

func exitHandle() {
	defer glog.Flush()
	glog.Info("收到退出信号")
	var RF = new(traffic.ReadFile)

	//if listener != nil {
	//	listener.Close()
	partition := common.Partition{PartitionPoint: traffic.PartitionPoint}
	fileJson, _ := json.Marshal(partition)
	if RF.CheckFileIsExist(common.Conf.Nginx.PartitionPointFileName) {
		if err := ioutil.WriteFile(common.Conf.Nginx.PartitionPointFileName, fileJson, 0644); err != nil{
			glog.Error("写入 PartitionPoint 失败", err)
			os.Exit(111)
		}else{
			glog.Info("写入PartitionPoint 到文件成功")
		}
	}else {
		_, err := os.Create(common.Conf.Nginx.PartitionPointFileName) //创建文件
		glog.Info("文件不存在,创建文件")
		if err != nil {
			glog.Info("文件创建失败。")
			//创建文件失败的原因有：
			//1、路径不存在  2、权限不足  3、打开文件数量超过上限  4、磁盘空间不足等
		}
		if err := ioutil.WriteFile(common.Conf.Nginx.PartitionPointFileName, fileJson, 0644); err != nil {
			glog.Error("写入 PartitionPoint 失败", err)
			os.Exit(111)
		} else {
			glog.Info("写入PartitionPoint 到文件成功")
		}
	}
}
