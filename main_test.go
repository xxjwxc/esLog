package main

import (
	"testing"

	"github.com/xxjwxc/esLog/logger"
)

func TestMain(m *testing.M) {
	var log logger.Logger
	log.AddLog()
	log.Search()
}
