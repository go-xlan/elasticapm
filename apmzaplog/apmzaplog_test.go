package apmzaplog

import (
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestNewLog(t *testing.T) {
	log := NewLog()
	a, b, c := 100, 200, 300
	log.Debugf("message a=%v b=%v c=%v", a, b, c)
	log.Errorf("message a=%v b=%v c=%v", a, b, c)
	log.Warningf("message a=%v b=%v c=%v", a, b, c)
}
