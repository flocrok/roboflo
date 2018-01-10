package main

import (
	"bufio"
	"time"
	"fmt"
	"strings"
	"github.com/bitfinexcom/bitfinex-api-go/v1"

)

const sleeptime time.Duration = 15
var Client *bitfinex.Client
var LogStoch = 1
var wstoch *bufio.Writer

var monnaie string = "XRPUSD"
var crypto string

func main() {

	crypto = strings.ToLower(string([]rune(monnaie)[0:3]))
	fmt.Println(crypto)

}