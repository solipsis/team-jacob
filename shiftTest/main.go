package main

import (
	"fmt"
	"log"

	ss "github.com/solipsis/shapeshift"
)

func main() {
	n := ss.New{
		Pair:        "eth_btc",
		ToAddress:   "16FdfRFVPUwiKAceRSqgEfn1tmB4sVUmLh",
		FromAddress: "0xcf2f204aC8D7714990912fA422874371c001217D",
	}

	resp, err := n.Shift()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
	//&{0xa6bd216e8e5f463742f37aaab169cabce601835c ETH 16FdfRFVPUwiKAceRSqgEfn1tmB4sVUmLh BTC   shapeshift {}}
}
