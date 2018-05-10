package main

import (
	"sort"
	"strings"

	ss "github.com/solipsis/shapeshift"
)

type Coin struct {
	Name      string
	Symbol    string
	Available bool
}

type shift struct {
	*ss.NewTransactionResponse
	receiveAddr string
}

// extract the fields we need from a shapeshift coin response object
func toCoin(sc ss.Coin) *Coin {
	return &Coin{
		Name:      sc.Name,
		Symbol:    sc.Symbol,
		Available: sc.Status == "Available",
	}
}

// initiate a new shift with Shapeshift
func newShift(pair, recAddr string) (*shift, error) {

	s := ss.New{
		//TODO; check other similar method on select screen
		Pair:      pair,
		ToAddress: recAddr,
	}
	Log.Println("Pair: ", selectScreen.activePair())

	response, err := s.Shift()
	if err != nil {
		Log.Println(err)
		panic(err)
	}
	Log.Println("received from ss ", response)

	if response.ErrorMsg() != "" {
		Log.Println(response.ErrorMsg())
		panic(response.ErrorMsg())
	}
	return &shift{response, recAddr}, nil
}

// activeCoins returns a slice of all the currently active coins on shapeshift
func activeCoins() ([]*Coin, error) {
	ssCoins, err := ss.CoinsAsList()
	active := make([]*Coin, 0)
	if err != nil {
		// Add 2 dummy coins so the scroll wheels still function
		active = append(active, &Coin{Name: "Unable to contact Shapeshift"})
		active = append(active, &Coin{Name: "Unable to contact Shapeshift"})
		return active, err
	}

	// Ignore any coins that aren't available
	for _, c := range ssCoins {
		if c.Status == "available" {
			active = append(active, toCoin(c))
		}
	}

	// Sort alphabetically
	sort.Slice(active, func(i, j int) bool {
		return strings.ToLower(active[i].Name) < strings.ToLower(active[j].Name)
	})
	return active, nil
}
