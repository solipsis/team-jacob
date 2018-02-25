package main

/*
func main() {
	resp, err := http.Get("https://shapeshift.io/getcoins")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var v interface{}
	json.Unmarshal(body, &v)
	data := v.(map[string]interface{})

	coins := make([]Coin, 0)
	for k := range data {

		coins = append(coins, Coin{k})
	}

	sort.Slice(coins, func(i, j int) bool { return coins[i].Symbol < coins[j].Symbol })

	for _, coin := range coins {

		fmt.Printf("%s struct {\nCoin\n} `json:\"%s\"`\n", coin.Symbol, coin.Symbol)
	}
	/*
		for _, v := range data {
			var coin Coin
			err = json.Unmarshal(v.([]byte), &coin)
			if err != nil {
				log.Fatal(err)
			}
			coins = append(coins, coin)
			fmt.Println(coin)
		}
		/*
			var pretty bytes.Buffer
			json.Indent(&pretty, body, "", "    ")
			fmt.Println(string(pretty.Bytes()))
			err = json.Unmarshal(body, &coins)
			fmt.Println(coins)
			if err != nil {
				log.Fatal(err)
			}
			for _, c := range coins {
				fmt.Println(c)
			}
			//var pretty bytes.Buffer
			//json.Indent(&pretty, body, "", "    ")
			//fmt.Println(string(pretty.Bytes()))
}

func (c Coin) String() string {
	return fmt.Sprintf("%s struct {\nCoin\n} `json:\"%s\"`\n", c.Symbol, c.Symbol)
}
*/
