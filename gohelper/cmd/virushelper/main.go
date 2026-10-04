package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"gohelper/internal/api"
)

func main() {
	action := flag.String("action", "", "auth me spin case")
	initData := flag.String("init", "", "initData value")
	refCode := flag.String("ref", "", "referral code")
	token := flag.String("token", "", "auth token")
	caseID := flag.Int("id", 2, "case identifier")
	flag.Parse()

	switch *action {
	case "auth":
		if *initData == "" {
			fmt.Println("ERR|init data value is empty")
			return
		}
		if *refCode == "" {
			fmt.Println("ERR|referral code value is empty")
			return
		}
		client := api.NewClient("")
		tokenValue, err := client.Auth(*initData, *refCode)
		if err != nil {
			fmt.Println("ERR|auth request failed reason", err.Error())
			return
		}
		fmt.Println("OK|token=" + tokenValue)
	case "me":
		if *token == "" {
			fmt.Println("ERR|token value is empty")
			return
		}
		client := api.NewClient(*token)
		me, err := client.Me()
		if err != nil {
			fmt.Println("ERR|profile request failed reason", err.Error())
			return
		}
		balance := strconv.FormatFloat(me.Balance, 'f', -1, 64)
		stars := strconv.FormatFloat(me.StarsBalance, 'f', -1, 64)
		fmt.Print("OK")
		fmt.Print("|firstname=" + me.FirstName)
		fmt.Print("|balance=" + balance)
		fmt.Print("|stars=" + stars)
		if me.NextFreeSpin != "" {
			fmt.Print("|nextspin=" + me.NextFreeSpin)
		}
		if me.NextCaseFreeSpin != "" {
			fmt.Print("|nextcase=" + me.NextCaseFreeSpin)
		}
		fmt.Println()
	case "spin":
		if *token == "" {
			fmt.Println("ERR|token value is empty")
			return
		}
		client := api.NewClient(*token)
		spin, err := client.Spin()
		if err != nil {
			fmt.Println("ERR|spin request failed reason", err.Error())
			return
		}
		success := "false"
		if spin.Success {
			success = "true"
		}
		fmt.Print("OK")
		fmt.Print("|success=" + success)
		fmt.Print("|prize=" + spin.Prize.Name)
		if spin.StoryReward != nil {
			fmt.Print("|story=" + strconv.Itoa(*spin.StoryReward))
		}
		fmt.Println()
	case "case":
		if *token == "" {
			fmt.Println("ERR|token value is empty")
			return
		}
		client := api.NewClient(*token)
		result, err := client.OpenCase(*caseID)
		if err != nil {
			fmt.Println("ERR|case request failed reason", err.Error())
			return
		}
		success := "false"
		if result.Success {
			success = "true"
		}
		fmt.Print("OK")
		fmt.Print("|success=" + success)
		fmt.Print("|prize=" + result.Prize.Name)
		fmt.Println()
	default:
		fmt.Println("ERR|action value is not recognized")
		os.Exit(1)
	}
}
