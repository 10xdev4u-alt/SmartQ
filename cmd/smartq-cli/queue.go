package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

var queueCmd = &cobra.Command{
	Use:   "queue",
	Short: "Manage queues",
	Long:  `Commands for creating, listing, and managing queues.`,
}

var queueCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new queue",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		queueName := args[0]
		createQueue(queueName)
	},
}

var queueListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all queues",
	Run: func(cmd *cobra.Command, args []string) {
		listQueues()
	},
}

func init() {
	queueCmd.AddCommand(queueCreateCmd)
	queueCmd.AddCommand(queueListCmd) // Add the new command
	rootCmd.AddCommand(queueCmd)
}

func createQueue(name string) {
	const apiBaseURL = "http://localhost:8080/api/v1"
	
	requestBody, err := json.Marshal(map[string]string{"name": name})
	if err != nil {
		fmt.Println("Error marshalling request body:", err)
		return
	}

	resp, err := http.Post(apiBaseURL+"/queues", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating queue:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Failed to create queue. Status: %s, Body: %s\n", resp.Status, string(body))
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding response body:", err)
		return
	}

	fmt.Println("Successfully created queue:")
	fmt.Printf("  ID:   %s\n", result["id"])
	fmt.Printf("  Name: %s\n", result["name"])
}

func listQueues() {
	const apiBaseURL = "http://localhost:8080/api/v1"

	resp, err := http.Get(apiBaseURL + "/queues")
	if err != nil {
		fmt.Println("Error listing queues:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Failed to list queues. Status: %s, Body: %s\n", resp.Status, string(body))
		return
	}

	var queues []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&queues); err != nil {
		fmt.Println("Error decoding response body:", err)
		return
	}

	fmt.Println("Available Queues:")
	for _, q := range queues {
		fmt.Printf("  - ID: %s, Name: %s\n", q["id"], q["name"])
	}
}
