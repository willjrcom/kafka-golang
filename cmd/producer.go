package cmd

import (
	"github.com/spf13/cobra"
	"github.com/willjrcom/kafka-golang/producer"
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
