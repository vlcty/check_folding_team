package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Donor struct {
	Points             uint   `json:"credit"`
	Name               string `json:"name"`
	CompletedWorkUnits uint   `json:"wus"`
}

type TeamResult struct {
	Points             uint    `json:"credit"`
	CompletedWorkUnits uint    `json:"wus"`
	ActiveCPUs         uint    `json:"active_50"`
	Donors             []Donor `json:"donors"`
}

func main() {
	teamID := flag.Uint("team", 1, "Team ID")
	flag.Parse()

	for {
		if teamResult := FetchStats(*teamID); teamResult != nil {
			fmt.Printf("OK - Team points: %d Work units: %d Active CPUs (last 50d): %d|team_points=%d team_work_units=%d team_active_cpus=%d",
				teamResult.Points, teamResult.CompletedWorkUnits, teamResult.ActiveCPUs,
				teamResult.Points, teamResult.CompletedWorkUnits, teamResult.ActiveCPUs)

			for _, donor := range teamResult.Donors {
				fmt.Printf(" 'donor_%s_points'=%d", donor.Name, donor.Points)
				fmt.Printf(" 'donor_%s_wus'=%d", donor.Name, donor.CompletedWorkUnits)
			}

			os.Exit(0)
		}
	}
}

func FetchStats(teamID uint) *TeamResult {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://stats.foldingathome.org/api/team/%d", teamID), nil)
	request.Header.Set("Referer", fmt.Sprintf("https://stats.foldingathome.org/team/%d", teamID))
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:74.0) Gecko/20100101 Firefox/74.0")
	request.Header.Set("Connection", "close")

	client := &http.Client{}

	response, responseError := client.Do(request)

	if responseError != nil {
		ExitCritical("Response error: " + responseError.Error())
		return nil
	} else {
		defer response.Body.Close()

		switch response.StatusCode {
		case http.StatusOK:
			teamResult := &TeamResult{}
			decodeError := json.NewDecoder(response.Body).Decode(teamResult)

			if decodeError != nil {
				ExitCritical("Decode error: " + decodeError.Error())
				return nil
			} else {
				return teamResult
			}

		case http.StatusGatewayTimeout, http.StatusBadGateway:
			time.Sleep(time.Second * 5)
			return nil

		default:
			ExitWarning("Didn't get any result. HTTP code was " + response.Status)
			return nil
		}
	}
}

func ExitCritical(message string) {
	fmt.Printf("CRITICAL - %s\n", message)
	os.Exit(2)
}

func ExitWarning(message string) {
	fmt.Printf("WARNING - %s\n", message)
	os.Exit(1)
}
