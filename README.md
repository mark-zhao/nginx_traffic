```
/root/zzc/config.json
{
  "db":{
    "DBip"          : "172.29.253.41:8086",
    "MyDB"          : "objectStorage",
    "MyMeasurement" : "traffic",
    "username"      : "csp",
    "password"      : "influx"
  },
  "nginx": {
    "FileName": "/var/log/nginx/access.log",
    "PartitionPointFileName": "/root/zzc/PartitionPoint.json",
    "Addr": "0.0.0.0:1234"
  }
}

//traffic 请求格式
 curl -H "Content-Type:application/json" -H "Data_Type:msg" -X POST --data '{"bucket_name": "default", "start_time": "2006-01-02 15:04:05","end_time":"2006-01-02 15:04:05"}'  http://127.0.0.1:1234/rgw

# /etc/systemd/system/traffic.service
[Unit]
Description=ceph rgw traffic  service
After=network-online.target firewalld.service
Wants=network-online.target

[Service]
Type=simple
# the default is not to use systemd for cgroups because the delegate issues still
# exists and systemd currently does not support the cgroup feature set required
ExecStart=/root/zzc/traffic -conf=/root/zzc/config.json
ExecStop=/bin/kill -s HUP $MAINPID
Restart=always

[Install]
WantedBy=multi-user.target
```