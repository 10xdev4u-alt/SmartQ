package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

var ticketCmd = &cobra.Command{
	Use:   "ticket",
	Short: "Manage tickets",
	Long:  `Commands for creating, listing, and managing tickets.`,
}

var ticketCreateCmd = &cobra.Command{
	Use:   "create [queueId] [customerName] [customerPhone] [priority]",
	Short: "Create a new ticket",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		queueID := args[0]
		customerName := args[1]
		customerPhone := args[2]
		priority, err := strconv.Atoi(args[3])
		if err != nil {
			fmt.Println("Error: priority must be an integer")
			return
		}
		createTicket(queueID, customerName, customerPhone, priority)
	},
}

var ticketCallCmd = &cobra.Command{
	Use:   "call [ticketId]",
	Short: "Call a ticket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ticketID := args[0]
		updateTicketStatus(ticketID, "call")
	},
}

var ticketServeCmd = &cobra.Command{
	Use:   "serve [ticketId]",
	Short: "Serve a ticket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ticketID := args[0]
		updateTicketStatus(ticketID, "serve")
	},
}

var ticketCancelCmd = &cobra.Command{
	Use:   "cancel [ticketId]",
	Short: "Cancel a ticket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ticketID := args[0]
		updateTicketStatus(ticketID, "cancel")
	},
}

func init() {
	ticketCmd.AddCommand(ticketCreateCmd)
	ticketCmd.AddCommand(ticketCallCmd)
	ticketCmd.AddCommand(ticketServeCmd)
	ticketCmd.AddCommand(ticketCancelCmd)
	rootCmd.AddCommand(ticketCmd)
}

func createTicket(queueID, customerName, customerPhone string, priority int) {
	const apiBaseURL = "http://localhost:8080/api/v1"

	requestBody, err := json.Marshal(map[string]interface{}{
		"customer_name":  customerName,
		"customer_phone": customerPhone,
		"priority":       priority,
	})
	if err != nil {
		fmt.Println("Error marshalling request body:", err)
		return
	}

	resp, err := http.Post(apiBaseURL+"/queues/"+queueID+"/tickets", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating ticket:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Failed to create ticket. Status: %s, Body: %s\n", resp.Status, string(body))
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding response body:", err)
		return
	}

	fmt.Println("Successfully created ticket:")
	for k, v := range result {
		fmt.Printf("  %s: %v\n", k, v)
	}
}

func updateTicketStatus(ticketID, status string) {
	const apiBaseURL = "http://localhost:8080/api/v1"

	resp, err := http.Post(apiBaseURL+"/tickets/"+ticketID+"/"+status, "application/json", nil)
	if err != nil {
		fmt.Printf("Error updating ticket status to %s: %v\n", status, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Failed to update ticket status to %s. Status: %s, Body: %s\n", status, resp.Status, string(body))
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding response body:", err)
		return
	}

	fmt.Printf("Successfully updated ticket status to %s:\n", status)
	for k, v := range result {
		fmt.Printf("  %s: %v\n", k, v)
	}
}
