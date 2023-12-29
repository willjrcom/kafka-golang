/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	consumer "github.com/willjrcom/kafka-golang/internal/consumer_ibm"
)

// consumerCmd represents the consumer command
var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Start consumer kafka",
	Run: func(cmd *cobra.Command, args []string) {
		consumer.Main()
	},
}

func init() {
	rootCmd.AddCommand(consumerCmd)
}
