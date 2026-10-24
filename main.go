package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

const dataFile = "expenses.json"

type Expense struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
}

func loadExpenses() ([]Expense, error) {
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return []Expense{}, nil
	}

	fileData, err := os.ReadFile(dataFile)
	if err != nil {
		return nil, err
	}

	var expenses []Expense
	err = json.Unmarshal(fileData, &expenses)
	if err != nil {
		return nil, err
	}

	return expenses, nil
}

func saveExpenses(expenses []Expense) error {
	fileData, err := json.MarshalIndent(expenses, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(dataFile, fileData, 0644)
}

func addExpense(description string, amount float64) {
	if amount <= 0 {
		fmt.Println("Error: Amount harus lebih dari 0")
		os.Exit(1)
	}

	expenses, _ := loadExpenses()

	newID := 1
	if len(expenses) > 0 {
		newID = expenses[len(expenses)-1].ID + 1
	}

	expense := Expense{
		ID:          newID,
		Date:        time.Now(),
		Description: description,
		Amount:      amount,
	}

	expenses = append(expenses, expense)
	saveExpenses(expenses)

	fmt.Printf("Expense added successfully (ID: %d)\n", newID)
}

func listExpenses() {
	expenses, _ := loadExpenses()

	fmt.Printf("%-5s %-12s %-20s %s\n", "ID", "Date", "Description", "Amount")
	fmt.Println("-------------------------------------------------------")
	for _, e := range expenses {
		dateStr := e.Date.Format("2006-01-02")
		fmt.Printf("%-5d %-12s %-20s $%.2f\n", e.ID, dateStr, e.Description, e.Amount)
	}
}

func deleteExpense(id int) {
	expenses, _ := loadExpenses()

	indexToDelete := -1
	for i, e := range expenses {
		if e.ID == id {
			indexToDelete = i
			break
		}
	}

	if indexToDelete == -1 {
		fmt.Printf("Error: Expense dengan ID %d tidak ditemukan\n", id)
		os.Exit(1)
	}

	expenses = append(expenses[:indexToDelete], expenses[indexToDelete+1:]...)
	saveExpenses(expenses)

	fmt.Println("Expense deleted successfully")
}

func summaryExpenses(month int) {
	expenses, _ := loadExpenses()
	var total float64

	currentYear := time.Now().Year()

	for _, e := range expenses {
		if month == 0 || (int(e.Date.Month()) == month && e.Date.Year() == currentYear) {
			total += e.Amount
		}
	}

	if month == 0 {
		fmt.Printf("Total expenses: $%.2f\n", total)
	} else {
		monthName := time.Month(month).String()
		fmt.Printf("Total Expense for %s: $%.2f\n", monthName, total)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected 'add', 'list', 'delete', or 'summary' subcommands")
		os.Exit(1)
	}

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addDesc := addCmd.String("description", "", "Expense description")
	addAmount := addCmd.Float64("amount", 0.0, "Expense amount")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteID := deleteCmd.Int("id", 0, "Expense ID to delete")

	summaryCmd := flag.NewFlagSet("summary", flag.ExitOnError)
	summaryMonth := summaryCmd.Int("month", 0, "Month number (1-12) for summary")

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		if *addDesc == "" || *addAmount == 0 {
			fmt.Println("Usage: expense-tracker add --description \"<desc>\" --amount <amount>")
			os.Exit(1)
		}

		addExpense(*addDesc, *addAmount)

	case "list":
		listExpenses()

	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if *deleteID == 0 {
			fmt.Println("Usage: expense-tracker delete --id <id>")
			os.Exit(1)
		}

		deleteExpense(*deleteID)

	case "summary":
		summaryCmd.Parse(os.Args[2:])
		if *summaryMonth < 0 || *summaryMonth > 12 {
			fmt.Println("Error: Month harus antara 1 dan 12")
			os.Exit(1)
		}

		summaryExpenses(*summaryMonth)

	default:
		fmt.Println("Expected 'add', 'list', 'delete', or 'summary' subcommands")
		os.Exit(1)
	}
}
