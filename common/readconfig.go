package common

import (
	"encoding/json"
	"flag"
	"glog"
	"io/ioutil"
)

var Conf *Config
var configFile  string
func init() {
	flag.StringVar(&configFile, "conf", "/root/zzc/config.json", "define config file ")
	flag.Parse()
	Conf = Getconfig()
	glog.Info(Conf)
}

//配置文件结构
type DB struct {
	DBip          string `json:"DBip"`
	MyDB          string `json:"MyDB"`
	MyMeasurement string `json:"MyMeasurement"`
	Username      string `json:"username"`
	Password      string `json:"password"`
}
type Nginx struct {
	FileName      			string `json:"FileName"`
	PartitionPointFileName 	string `json:"PartitionPointFileName"`
	Addr 					string `json:"Addr"`
}
type Config struct {
	DBInfo DB `json:"db"`
	Nginx  Nginx `json:"nginx"`
}

type Partition struct {
	PartitionPoint int `json:"PartitionPoint"`
}


// json读取
type JsonStruct struct {
}

func Getconfig() *Config {
	JsonParse := NewJsonStruct()
	v := new(Config)
	JsonParse.Load(configFile, &v)
	return v
}

func GetPartition() *Partition {
	JsonParse := NewJsonStruct()
	v := new(Partition)
	JsonParse.Load(Conf.Nginx.PartitionPointFileName, &v)
	return v
}

//读json 文件
func (jst *JsonStruct) Load(filename string, v interface{}) {
	defer glog.Flush()
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		glog.Info("error:", err)
		return
	}
	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		glog.Info("error:", err)
		return
	}
}
func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}
