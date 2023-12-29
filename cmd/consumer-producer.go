/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	consumerproducer "github.com/willjrcom/kafka-golang/internal/consumer_producer_segmentio"
)

// consumerProducerCmd represents the consumerProducer command
var consumerProducerCmd = &cobra.Command{
	Use:   "consumer-producer",
	Short: "Consumer and Producer kafka",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("consumerProducer called")
		consumerproducer.Main()
	},
}

func init() {
	rootCmd.AddCommand(consumerProducerCmd)
}
