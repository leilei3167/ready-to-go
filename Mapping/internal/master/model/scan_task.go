// Package model 提供用于master和用户端消息交互的model以及构建作为kafka消息的结构体
package model

import (
	"encoding/json"
)

// ScanTaskPool produce消息时在消费端会出现ip重复的情况
/*var ScanTaskPool = &sync.Pool{
	New: func() interface{} {
		return new(ScanTask)
	},
}

func (s *ScanTask) Reset() {
	s.IP = ""
	s.Ports = nil
	s.encoded = nil
	s.err = nil
}*/

// ScanTask 实现sarama.Encoder接口,可以作为kafka的Value进行发送
type ScanTask struct {
	IP    string `json:"IP,omitempty"`
	Ports []int  `json:"ports,omitempty"`

	encoded []byte
	err     error
}

func (s *ScanTask) ensureEncoded() {
	if s.encoded == nil {
		s.encoded, s.err = json.Marshal(s)
	}
}

func (s *ScanTask) Encode() ([]byte, error) {
	s.ensureEncoded()
	return s.encoded, s.err
}

func (s *ScanTask) Length() int {
	return len(s.encoded)
}
