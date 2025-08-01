/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/migrations"
	"github.com/spf13/cobra"
)

// deldefrecCmd represents the deldefrec command
var deldefrecCmd = &cobra.Command{
	Use:   "deldefrec",
	Short: "Delete default records from the database",
	Long:  `This command deletes default records from users, secrets, events, and event_participants tables.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running...")
		err := migrations.Migrate()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		fmt.Println("Default records deleted successfully")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(deldefrecCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deldefrecCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deldefrecCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
