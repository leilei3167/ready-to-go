package db

import (
	"time"
)

/*type ScanResult struct {
	Uptime time.Time
	IP     string
	Ports  []*PortInfo
}

type PortInfo struct {
	Port int
	Type string //后续可维护一个map,包含常见端口对应的服务,显示在此
}*/

type ScanResult struct { //写入数据库的model
	IP        string    `json:"ip" binding:"ip"`
	IsAlive   bool      `json:"isalive"`
	OpenPorts []int     `json:"openports"`
	Uptime    time.Time `json:"uptime"`
}
