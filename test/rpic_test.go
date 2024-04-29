package test

import (
	"testing"

	"github.com/parnurzeal/gorequest"
)

func BenchmarkRegularGet(b *testing.B) {
	target := "http://192.168.31.128/rpic/get"
	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				gorequest.New().Get(target).End()
			}
		},
	)
}
