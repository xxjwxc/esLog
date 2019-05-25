package main

import "github.com/xie1xiao1jun/esLog/logger"

func TestMain() {
	var log logger.Logger
	log.AddLog()
	log.Search()
}
