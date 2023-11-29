package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

// Global Variables -------------------------------- //
const API_KEY = os.Getenv("THEODDSAPIKEYY")

// System gv --------------------------------------- //
var (
	log     = log.Println
	error   = log.Println
	warning = log.Println
	success = log.Println
)

// Main -------------------------------------------- //
func main() {
	log("BetArbit v1.0.0")
	warning("Arbitrage gambling has a variety of factors that can negatively affect profit. Do not use this cli to place real-world bets.")
	log("Starting...")

	demo := false
	demoFile := "./test_data.json"
	bet := 100
	sport := "upcoming"
	region := "au"
	verbose := false

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--demo":
			demo = true
		case "--demoFile":
			if i+1 < len(args) {
				demoFile = args[i+1]
				i++
			}
		case "--bet":
			if i+1 < len(args) {
				bet, _ = strconv.Atoi(args[i+1])
				i++
			}
		case "--sport":
			if i+1 < len(args) {
				sport = args[i+1]
				i++
			}
		case "--region":
			if i+1 < len(args) {
				region = args[i+1]
				i++
			}
		case "--verbose":
			verbose = true
		}
	}

	if demo {
		if demoFile == "" {
			log("!! When in demo mode please provide a demoFile using --demoFile or do not provide the argument to use the default.")
			exit(true)
		}
	}

	DEMO := demo
	DEMO_FILE := demoFile
	BET := bet
	SPORT := sport
	REGION := region
	VERBOSE := verbose

	if VERBOSE {
		log(`
        Arguments
        demo: ` + strconv.FormatBool(DEMO) + `
        demoFile: ` + DEMO_FILE + `
        bet: ` + strconv.Itoa(BET) + `
        sport: ` + SPORT + `
        region: ` + REGION + `
    `)
	}

	if VERBOSE {
		log("Collecting Input Data")
	}

	var data []map[string]interface{}

	if DEMO {
		if VERBOSE {
			log("     L_____ Demo mode enabled")
		}
		file, err := ioutil.ReadFile(DEMO_FILE)
		if err != nil {
			log(err)
			exit(true)
		}
		err = json.Unmarshal(file, &data)
		if err != nil {
			log("!! Failed to collect data: Could not get demoFile: " + DEMO_FILE)
			exit(true)
		}
	} else {
		if VERBOSE {
			log("     L_____ Production mode enabled")
		}
		resp, err := http.Get("https://api.the-odds-api.com/v3/odds/?apiKey=" + API_KEY + "&sport=" + SPORT + "&region=" + REGION + "&mkt=h2h")
		if err != nil {
			reqError := ""
			if resp == nil {
				reqError = "Could not reach server."
			} else {
				reqError = err.Error()
			}
			log("!! Failed to collect data: " + reqError)
			exit(true)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log("!! Failed to collect data: " + err.Error())
			exit(true)
		}
		err = json.Unmarshal(body, &data)
		if err != nil {
			log("!! Failed to collect data: " + err.Error())
			exit(true)
		}
	}

	if data == nil {
		log("!! An unknown error occurred while collecting data")
		exit(true)
	}

	if VERBOSE {
		log(success("Successfully retrieved data object"))
		log("Processing data")
	}

	counts := map[string]int{
		"inProfit":    0,
		"totalProfit": 0,
	}

	for i := 0; i < len(data); i++ {
		match := data[i]

		log("Checking for arbitrage on " + match["sport_key"].(string) + ", " + match["teams"].(string))

		oddsMatrix := make([][]map[string]interface{}, len(match["sites"].([]interface{})))
		for j := 0; j < len(match["sites"].([]interface{})); j++ {
			site := match["sites"].([]interface{})[j].(map[string]interface{})
			for k := 0; k < len(site["odds"].(map[string]interface{})["h2h"].([]interface{})); k++ {
				if len(oddsMatrix[k]) == 0 {
					oddsMatrix[k] = []map[string]interface{}{}
				}
				oddsMatrix[k] = append(oddsMatrix[k], map[string]interface{}{
					"site": site["site_key"].(string),
					"odd":  site["odds"].(map[string]interface{})["h2h"].([]interface{})[k].(float64),
				})
			}
		}

		if VERBOSE {
			log("     L_____ Successfully organized odds")
		}

		highestOdds := make([]map[string]interface{}, len(oddsMatrix))
		for j := 0; j < len(oddsMatrix); j++ {
			odds := oddsMatrix[j]
			sort.Slice(odds, func(a, b int) bool {
				return odds[a]["odd"].(float64) > odds[b]["odd"].(float64)
			})
			highestOdds[j] = odds[0]
		}

		if VERBOSE {
			log("     L_____ Found Highest odds across all sites")
		}

		arbitrage := 0.0
		for j := 0; j < len(highestOdds); j++ {
			arbitrage += 1 / highestOdds[j]["odd"].(float64)
		}

		if VERBOSE {
			log("     L_____ Calculating arbitrage values")
		}

		arbitrage *= 100

		if arbitrage < 100 {
			counts["inProfit"]++

			log("     L__________ Profitable arbitrage found at " + strconv.Itoa(int(arbitrage)) + ", calculating ideal wagers")

			wagers := []map[string]interface{}{}
			for j := 0; j < len(highestOdds); j++ {
				bet := 1.0
				for k := 0; k < len(highestOdds); k++ {
					if k != j {
						odd := highestOdds[k]["odd"].(float64)
						oddForOutcome := highestOdds[j]["odd"].(float64) / odd
						bet += oddForOutcome
					}
				}
				wager := BET / bet
				profit := (wager * highestOdds[j]["odd"].(float64)) - float64(BET)
				wagers = append(wagers, map[string]interface{}{
					"site":   highestOdds[j]["site"].(string),
					"wager":  fmt.Sprintf("%.2f", wager),
					"odd":    highestOdds[j]["odd"].(float64),
					"profit": fmt.Sprintf("%.2f", profit),
				})
			}

			log("     L_______________ Found ideal wagers")

			for j := 0; j < len(wagers); j++ {
				log("     L_______________ Selection " + strconv.Itoa(j) + "(" + strconv.FormatFloat(wagers[j]["odd"].(float64), 'f', -1, 64) + ") on " + wagers[j]["site"].(string) + " with $" + wagers[j]["wager"].(string))
			}

			log("     L_______________________ Profit if win: $" + wagers[0]["profit"].(string))

			counts["totalProfit"] += int(wagers[0]["profit"].(float64))
		} else {
			log("     L__________ No profitable arbitrage found.\n")
		}
	}

	if counts["inProfit"] > 0 {
		log("Successfully found a total of " + strconv.Itoa(counts["inProfit"]) + " possible arbitrage bets with a total potential profit of " + strconv.Itoa(counts["totalProfit"]))
	} else {
		log("Could not find any arbitrage bets for the given data.")
	}

	exit()
}

func exit(error bool) {
	if error {
		log("better luck next time!")
	}
	os.Exit(1)
}
