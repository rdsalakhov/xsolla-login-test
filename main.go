package main

import (
	"github.com/rdsalakhov/xsolla-login-test/server"
)

func main(){
	server := server.NewServer()
	server.Start(":5000")
}



