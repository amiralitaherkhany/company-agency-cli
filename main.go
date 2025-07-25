package main

import (
	"agency/entities/agency"
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	command := flag.String("command", "no-command", "command to run")
	region := flag.String("region", "no-region", "region of the data you want")

	flag.Parse()

	checkCommand(*command, *region)
}

func checkCommand(command string, region string) {
	if region == "no-region" {
		log.Fatalln("region is required!")
	}
	if command == "no-command" {
		log.Fatalln("command is required!")
	}

	switch command {
	case "list":
		listAllAgenciesByRegion(region)
	case "get":
		getExistingAgency(region)
	case "create":
		createNewAgency(region)
	case "edit":
	case "status":
		statusExistingAgency(region)
	default:
		log.Fatalln("command not found!")
	}
}

func statusExistingAgency(region string) {
	file, err := os.OpenFile("agency.csv", os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	agencies, err := reader.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}

	numberOfAgencies, numberOfEmployees := 0, 0
	for _, v := range agencies {
		if v[6] == region {
			numberOfAgencies++
			num, _ := strconv.Atoi(v[5])
			numberOfEmployees += num
		}
	}

	if numberOfAgencies == 0 {
		fmt.Printf("\nThere is no Agency in the %s region!\n", region)
		return
	}

	fmt.Println()
	fmt.Printf("There is %d Agencies and %d Employees in the %s region!\n", numberOfAgencies, numberOfEmployees, region)

}

func listAllAgenciesByRegion(region string) {
	file, err := os.OpenFile("agency.csv", os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	agencies, err := reader.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}

	result := []agency.Agency{}
	for _, v := range agencies {
		if v[6] == region {
			result = append(result, *agency.FromCSV(v))
		}
	}

	if len(result) == 0 {
		fmt.Printf("\nThere is no Agency in the %s region!\n", region)
		return
	}
	fmt.Println()
	for _, v := range result {
		fmt.Println()
		fmt.Println(v)
	}
}

func getExistingAgency(region string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("please enter the agency UUID")
	scanner.Scan()
	ID := strings.TrimSpace(scanner.Text())
	a, err := getAgencyByIDAndRegion(ID, region)
	if err != nil {
		log.Fatalln("Error in getting agency:", err)
	}
	fmt.Println()
	fmt.Println("Your Agency Found!")
	fmt.Println()
	fmt.Println(a)
	fmt.Println()
}

func createNewAgency(region string) {
	scanner := bufio.NewScanner(os.Stdin)
	agencyValues := map[string]any{
		"name":              "",
		"address":           "",
		"phoneNumber":       "",
		"membershipDate":    "",
		"numberOfEmployees": 0,
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
		region,
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
	err := writeToCsvFile("agency.csv", agency.ToCSV(a))
	if err != nil {
		return err
	}
	return nil
}
func getAgencyByIDAndRegion(ID string, region string) (*agency.Agency, error) {
	a, err := readFromCsvFile("agency.csv", ID, region)
	if err != nil {
		return nil, err
	}
	return agency.FromCSV(a), nil
}

func writeToCsvFile(filePath string, record []string) error {
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

func readFromCsvFile(filePath string, id string, region string) ([]string, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			if len(record) > 0 && record[0] == id && record[6] == region {
				return record, nil
			}
			return nil, errors.New("key doesn't exist!")
		}
		if err != nil {
			return nil, err
		}
		if len(record) > 0 && record[0] == id && record[6] == region {
			return record, nil
		}
	}

}
