package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "qapi",
	Short: "A command line tool to fetch data from an API",
	Run:   qapi,
}

var apiURL string
var flagNoColor bool

func init() {
	rootCmd.PersistentFlags().StringVarP(&apiURL, "url", "u", "", "The API URL to fetch data from")
	rootCmd.MarkPersistentFlagRequired("url")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "This flag disables color output to stdout")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func qapi(cmd *cobra.Command, args []string) {
	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if flagNoColor {
		color.NoColor = true
	}
	color.Green(string(body))
}
