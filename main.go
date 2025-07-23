package main

import (
	"agency/entities/agency"
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	command := flag.String("command", "no-command", "command to run")
	// region := flag.String("region", "no-region", "region of the data you want")

	flag.Parse()

	checkCommand(*command)
}

func checkCommand(command string) {
	switch command {
	case "list":
	case "get":
	case "create":
		createNewAgency()
	case "edit":
	case "status":

	}
}

func createNewAgency() {
	scanner := bufio.NewScanner(os.Stdin)
	agencyValues := map[string]any{
		"name":              "",
		"address":           "",
		"phoneNumber":       "",
		"membershipDate":    "",
		"numberOfEmployees": 0,
		"regionName":        "",
	}
	for k := range agencyValues {
	startOfGettingNumber:
		fmt.Println("please enter the agency", k)
		scanner.Scan()
		if k == "numberOfEmployees" {
			if v, err := strconv.Atoi(strings.TrimSpace(scanner.Text())); err != nil || v < 1 {
				fmt.Println("Invalid number !")
				goto startOfGettingNumber
			} else {
				agencyValues[k] = uint64(v)
				continue
			}
		}
		agencyValues[k] = strings.TrimSpace(scanner.Text())
	}

	newAgency := agency.New(
		agencyValues["name"].(string),
		agencyValues["address"].(string),
		agencyValues["phoneNumber"].(string),
		agencyValues["membershipDate"].(string),
		agencyValues["regionName"].(string),
		agencyValues["numberOfEmployees"].(uint64),
	)
	if err := saveAgency(newAgency); err != nil {
		log.Fatalln("An error occurred while creating your agency:", err)
	}
	fmt.Println()
	fmt.Println("New agency Created!")
	fmt.Println()
	fmt.Println(newAgency.String())
}

func saveAgency(a *agency.Agency) error {
	err := writeToCSVFile("agency.csv", agency.ToCSV(a))
	if err != nil {
		return err
	}
	return nil
}

func writeToCSVFile(filePath string, record []string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	if err := writer.Write(record); err != nil {
		return err
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}
	return nil
}
