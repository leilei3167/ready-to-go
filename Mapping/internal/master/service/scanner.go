package service

import "sync"

type scanner interface {
}

//每一个IP都是一个任务,且量特别大 应该考虑对象复用

var TcpScannerPool = &sync.Pool{
	New: func() interface{} {
		return new(TcpScanner)
	},
}

type TcpScanner struct {
	IP    string
	Ports []int
}

func (t *TcpScanner) Reset() {
	t.IP = ""
}
