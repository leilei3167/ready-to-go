package model

import "mime/multipart"

// Param 用于构容纳前端发起的请求参数以便于后续解析
type Param struct {
	Type  string
	Ports []int
	Ip    []string
	File  *multipart.FileHeader
}
type ScanReq struct {
	ScanType string   `json:"scanType,omitempty" binding:"required"`
	Targets  []string `json:"targets,omitempty" binding:"required"`
}

type Response struct {
	InvalidIP []string `json:"invalidIP,omitempty"`
	HasSent   uint     `json:"hasSent,omitempty"`
}
