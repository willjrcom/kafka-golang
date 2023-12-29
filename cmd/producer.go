package cmd

import (
	"github.com/spf13/cobra"
	producer "github.com/willjrcom/kafka-golang/internal/producer_ibm"
)

// producerCmd represents the base command when called without any subcommands
var producerCmd = &cobra.Command{
	Use:   "producer",
	Short: "Start producer kafka",
	Run: func(cmd *cobra.Command, args []string) {
		producer.Main()
	},
}

func init() {
	rootCmd.AddCommand(producerCmd)
}
