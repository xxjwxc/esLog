package main

import (
	"testing"

	"github.com/xie1xiao1jun/esLog/logger"
)

func TestMain(m *testing.M) {
	var log logger.Logger
	log.AddLog()
	log.Search()
}
