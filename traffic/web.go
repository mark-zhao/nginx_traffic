package traffic

import (
	"context"
	"encoding/json"
	"glog"
	"io/ioutil"
	"net/http"
	"nginx_traffic/common"
	"time"
)

const uri = "/rgw"
var RF = new(ReadFile)
var PartitionPointFileName  = config.Nginx.PartitionPointFileName

type RequestBody struct {
	BucketName 		string 		`json:"bucket_name"`
	StartTime  		string		`json:"start_time"`
	EndTime 		string  	`json:"end_time"`
}

type ResponseBody struct {
	BucketName 		string 		`json:"bucket_name"`
	TrafficSize  	interface{}		`json:"trafficsize"`
}

func ParseReq(origin []byte) (*RequestBody, error) {
	req := RequestBody{}
	err := json.Unmarshal(origin, &req)
	if err != nil {
		glog.Error("json Unmarshal failed; ", err)
		return &req, err
	}
	glog.Info("获取request body 成功 ")
	return &req, nil
}

func server(w http.ResponseWriter, r *http.Request) {
	defer glog.Flush()
	res := ResponseBody{}
	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/xml")

		originData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			glog.Error("Invalid request; ", err)
			http.Error(w, "Invalid request", 405)
			return
		}
		//glog.Info("content: ", originData)
		rawData, err := ParseReq(originData)
		if err != nil {
			glog.Error("parse request failed;")
			http.Error(w, "parse data failed;", 500)
			return
		}
		_, err = time.ParseInLocation("2006-01-02 15:04:05", rawData.StartTime, time.Local)
		if err != nil {
			glog.Error("时间格式有问题;")
			http.Error(w, "parse data failed，start time format err;", 500)
			return
		}
		_, err = time.ParseInLocation("2006-01-02 15:04:05", rawData.EndTime, time.Local)
		if err != nil {
			glog.Error("时间格式有问题;")
			http.Error(w, "parse data failed，end time format err;", 500)
			return
		}
		trafficSize, err := SelectTraffic(rawData.BucketName, rawData.StartTime, rawData.EndTime)
		glog.Info(rawData.BucketName," traffic size is ",trafficSize)
		res = ResponseBody{rawData.BucketName,trafficSize[0].Series[0].Values[0][1]}
		msg, err := json.Marshal(res)
		_, _ = w.Write(msg)

	} else {
		glog.Error("Invalid request method.", 400)
		http.Error(w, "Invalid request method", 400)
	}

}

func StartHTTP(ctx context.Context) error {
	http.HandleFunc(uri, server)
	server := &http.Server{
		Addr:         config.Nginx.Addr,
		Handler:      http.DefaultServeMux,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	//读文件
	//R := new(traffic.ReadFile)
	r := common.NewReflector(time.Second*60)
	go func() {
		resyncCh, cleanup := r.ResyncChan()
		defer func() {
			cleanup() // Call the last one written into cleanup
		}()
		for {
			//当时间没到时会堵塞在这边，时间到了定时器才会有信号继续往下执行
			select {
			case <-resyncCh:
			}
			//初始化PartitionPoint
			if RF.CheckFileIsExist(PartitionPointFileName) {
				PartitionPoint = common.GetPartition().PartitionPoint
			}
			//处理逻辑
			if err := RF.ReadFile(config.Nginx.FileName, RF.Todo);err != nil {
				glog.Error(err)
			}
			glog.Flush()
			// 清理掉当前的计时器，获取下一个同步时间定时器
			cleanup()
			resyncCh, cleanup = r.ResyncChan()
		}
	}()

	go func(s *http.Server) {
		if err := server.ListenAndServe(); err != nil {
			glog.Error("http server listen and Serve failed; message: ", err)
		}
	}(server)
	glog.Info("http server launch sucess!")

	for {
		select {
		case <-ctx.Done():
			_ = server.Close()
			return nil
		}
	}
}
