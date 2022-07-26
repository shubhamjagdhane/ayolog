package main

import (
	"bytes"
	"fmt"
	"testing"

	"bitbucket.org/shubhamjagdhane_ayoconnect/mylog"
)

func BenchmarkDebug(b *testing.B) {
	log := mylog.New()
	var buf bytes.Buffer
	msg := "testing"
	log.Out = &buf

	for n := 1; n <= 20; n += 1 {
		name := fmt.Sprintf("%v", n)
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				buf.Reset()
				log.Debug(msg)
			}
		})
	}
}
