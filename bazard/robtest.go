package main

import (
	"bufio"
	"time"
	"github.com/bitfinexcom/bitfinex-api-go/v1"
	"../utils"
	"../trc"
)

const sleeptime time.Duration = 15
var Client *bitfinex.Client
var LogStoch = 1
var wstoch *bufio.Writer

var monnaie string = "XRPUSD"

func main() {

	var ret int = 0

	ret=utils.Init(monnaie,LogStoch,4)
	if ret == 1 {return}

	trc.Trace.Printf("fin : (%d) \n",ret)

	utils.RecupMonnaies()



}