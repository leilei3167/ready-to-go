package main

import (
	"testing"
)

var s []string

func BenchmarkReadIPs(b *testing.B) {
	a := []string{}
	for i := 0; i < b.N; i++ {
		a, _ = ReadIPs("./golbalips.txt")
	}
	//	log.Printf("len:%d,cap:%d,size:%d", len(a), cap(a), unsafe.Sizeof(a))
	s = a
}
