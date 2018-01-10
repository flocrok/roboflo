package cex

import (
    "github.com/jhautefeuille/go-cexio-api"
    "fmt"
)

func getticker() {
     cexapi := cexio.CexKey{
         Username:"your_username", 
         Api_key:"your_api_key", 
    	 Api_secret:"your_api_secret"}

    // Public
    fmt.Printf("Ticker => %s\n", cexapi.Ticker("USD"))
    //fmt.Printf("Order Book => %s\n", cexapi.OrderBook("GHS/BTC"))
    //fmt.Printf("Trade History => %s\n", cexapi.TradeHistory("GHS/BTC"))

    // Private
    // fmt.Printf("Balance => %s\n", cexapi.Balance())
    // fmt.Printf("Open Orders => %s\n", cexapi.OpenOrders("GHS/BTC"))

    // Trading orders
    //fmt.Printf("Place Order => %s\n", cexapi.PlaceOrder("buy", "0.001", "0.017", "GHS/BTC"))
    //fmt.Printf("Cancel Order => %s\n", cexapi.CancelOrder("477571539"))

    // Workers 
    // fmt.Printf("Hashrate => %s\n", cexapi.Hashrate())
    // fmt.Printf("Workers => %s\n", cexapi.Workers())
}