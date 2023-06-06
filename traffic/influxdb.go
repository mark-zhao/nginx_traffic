package traffic

import (
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"glog"
	"nginx_traffic/common"
	"time"
)
var config = common.Conf
type Fields struct {
	BucketName string
	size       int64
	timestamp  string
}

func InsertInfluxdb(Fs []*Fields, conn client.Client) {
	defer glog.Flush()
	MyDB := config.DBInfo.MyDB
	MyMeasurement := config.DBInfo.MyMeasurement
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  MyDB,
		Precision: "s",
	})
	if err != nil {
		glog.Error("error", err)
	}
	//数据录入
	for _, v := range Fs {
		fields := map[string]interface{}{
			"size": v.size,
		}
		//the_time, err := time.Parse("2006-01-02 15:04:05", v.timestamp)
		the_time, err := time.ParseInLocation("2006-01-02 15:04:05", v.timestamp, time.Local)
		if err != nil {
			glog.Error(err)
		}
		Tags := map[string]string{"BucketName": v.BucketName}
		//a := the_time.Unix()
		glog.Info("time:", the_time, "MyMeasurement:", MyMeasurement, "tags:", Tags, "fields:", fields)
		//pt, err := client.NewPoint(MyMeasurement, Tags, fields, time.Unix(a, 0))
		//pt, err := client.NewPoint(MyMeasurement, Tags, fields, the_time.Add(8*h))
		pt, err := client.NewPoint(MyMeasurement, Tags, fields, the_time	)
		bp.AddPoint(pt)
		if err := conn.Write(bp); err != nil {
			glog.Error("error", err)
		}
	}
}

//查询bucket 某段时间内的流量
func SelectTraffic(bucketName string, startTime string, endTime string) ([]client.Result, error){
	cli := connInflux()
	defer cli.Close()
	qs := fmt.Sprintf("SELECT SUM(size) FROM traffic WHERE BucketName = '%s' AND time >= '%s' AND time <= '%s'", bucketName, startTime, endTime)
	fmt.Println(qs)
	todo := client.Query{
		Command: qs,
		Database: "objectStorage",
	}
	if Res, err := cli.Query(todo);err == nil {
		if Res.Error() != nil {
			fmt.Println("查询语句有问题", Res.Error())
			return nil, Res.Error()
		}
		return Res.Results, nil
	}else {
		return nil, err
	}
}

//建立influxdb链接
func connInflux() client.Client {
	defer glog.Flush()
	username := config.DBInfo.Username
	password := config.DBInfo.Password
	DBip := config.DBInfo.DBip
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://" + DBip,
		Username: username,
		Password: password,
	})
	if err != nil {
		glog.Error("error", err)
	}
	return cli
}
