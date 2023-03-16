package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/alecthomas/chroma/quick"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	apiURL      string
	apiKey      string
	format      string
	flagNoColor bool
	rootCmd     = &cobra.Command{
		Use:   "qapi",
		Short: "A command line tool to fetch data from an API",
		Run:   qapi,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&apiURL, "url", "u", "", "The API URL to fetch data from")
	rootCmd.PersistentFlags().StringVarP(&apiKey, "key", "k", "", "This flag attaches an API Key header to the HTTP Request")
	rootCmd.PersistentFlags().StringVarP(&format, "output", "o", "json", "This flag set the output format")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "This flag disables color output to stdout")
	_ = rootCmd.MarkPersistentFlagRequired("url")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func qapi(cmd *cobra.Command, args []string) {
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Fatal("Error:", err)
		return
	}
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error: No URL Provided ", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error:", err)
		return
	}

	if flagNoColor {
		color.NoColor = true
	}

	switch format {
	case "json":
		printJSONOutput(body)
	case "yaml":
		printYAMLOutput(body)
	// case "csv":
	// TODO: Fix function printCSVOutput(body)
	default:
		fmt.Println("Error: unsupported output format")
	}
}

func printJSONOutput(body []byte) {
	var data interface{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal("Error:", err)
		return
	}

	prettyJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Fatal("Error:", err)
		return
	}

	// Colorize and print the output
	quick.Highlight(os.Stdout, string(prettyJSON), "json", "terminal256", "native")
}

func printYAMLOutput(body []byte) {
	var data interface{}
	err := yaml.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Pretty print the output
	prettyYAML, err := yaml.Marshal(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Colorize and print the output
	quick.Highlight(os.Stdout, string(prettyYAML), "yaml", "terminal256", "native")
}

// func printCSVOutput(body []byte) {
// 	var data []map[string]interface{}
// 	err := json.Unmarshal(body, &data)
// 	if err != nil {
// 		log.Fatal("Error:", err)
// 		return
// 	}

// 	// Convert data to CSV format
// 	csvOutput, err := gocsv.MarshalString(data)
// 	if err != nil {
// 		log.Fatal("Error:", err)
// 		return
// 	}

// 	// Replace "," with ";" to avoid conflicts with CSV formatting
// 	csvOutput = strings.ReplaceAll(csvOutput, ",", ";")

// 	// Colorize and print the output
// 	csvReader := csv.NewReader(strings.NewReader(csvOutput))
// 	records, err := csvReader.ReadAll()
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}

// 	// Colorize and print the output
// 	for _, record := range records {
// 		for _, field := range record {
// 			color.Blue(field)
// 		}
// 		fmt.Println()
// 	}
// }
