package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type Snippet struct {
	Name     string   `json:"name"`
	Language string   `json:"language"`
	Tags     []string `json:"tags"`
	Content  string   `json:"content"`
}

const snippetFile = "snippets.json"

func main() {
	rootCmd := &cobra.Command{Use: "snip"}

	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new snippet",
		Run:   addSnippet,
	}
	addCmd.Flags().StringP("name", "n", "", "Snippet name")
	addCmd.Flags().StringP("language", "l", "", "Programming language")
	addCmd.Flags().StringP("tags", "t", "", "Comma-separated tags")
	addCmd.Flags().StringP("content", "c", "", "Snippet content")
	rootCmd.AddCommand(addCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all snippets",
		Run:   listSnippets,
	}
	rootCmd.AddCommand(listCmd)

	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search snippets by keyword",
		Run:   searchSnippets,
	}
	searchCmd.Flags().StringP("query", "q", "", "Search query")
	rootCmd.AddCommand(searchCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func loadSnippets() ([]Snippet, error) {
	if _, err := os.Stat(snippetFile); os.IsNotExist(err) {
		return []Snippet{}, nil
	}
	data, err := ioutil.ReadFile(snippetFile)
	if err != nil {
		return nil, err
	}
	var snippets []Snippet
	err = json.Unmarshal(data, &snippets)
	if err != nil {
		return nil, err
	}
	return snippets, nil
}

func saveSnippets(snippets []Snippet) error {
	data, err := json.MarshalIndent(snippets, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(snippetFile, data, 0644)
}

func addSnippet(cmd *cobra.Command, args []string) {
	name, _ := cmd.Flags().GetString("name")
	language, _ := cmd.Flags().GetString("language")
	tagsStr, _ := cmd.Flags().GetString("tags")
	content, _ := cmd.Flags().GetString("content")

	if name == "" || content == "" {
		fmt.Println("Name and content are required")
		return
	}

	tags := []string{}
	if tagsStr != "" {
		for _, t := range strings.Split(tagsStr, ",") {
			tags = append(tags, strings.TrimSpace(t))
		}
	}

	snippets, err := loadSnippets()
	if err != nil {
		fmt.Println("Error loading snippets:", err)
		return
	}

	// Prevent duplicates by name
	for _, s := range snippets {
		if s.Name == name {
			fmt.Println("Snippet with this name already exists.")
			return
		}
	}

	newSnippet := Snippet{
		Name:     name,
		Language: language,
		Tags:     tags,
		Content:  content,
	}
	snippets = append(snippets, newSnippet)

	if err := saveSnippets(snippets); err != nil {
		fmt.Println("Error saving snippet:", err)
		return
	}

	fmt.Printf("Snippet '%s' added successfully.\n", name)
}

func listSnippets(cmd *cobra.Command, args []string) {
	snippets, err := loadSnippets()
	if err != nil {
		fmt.Println("Error loading snippets:", err)
		return
	}

	if len(snippets) == 0 {
		fmt.Println("No snippets found.")
		return
	}

	for i, s := range snippets {
		fmt.Printf("[%d] %s (%s) - Tags: %v\n", i+1, s.Name, s.Language, s.Tags)
	}
}

func searchSnippets(cmd *cobra.Command, args []string) {
	query, _ := cmd.Flags().GetString("query")
	if query == "" {
		fmt.Println("Please provide a search query with -q")
		return
	}

	snippets, err := loadSnippets()
	if err != nil {
		fmt.Println("Error loading snippets:", err)
		return
	}

	found := false
	for i, s := range snippets {
		if strings.Contains(strings.ToLower(s.Name), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(s.Content), strings.ToLower(query)) ||
			containsTag(s.Tags, query) {
			fmt.Printf("[%d] %s (%s) - Tags: %v\n", i+1, s.Name, s.Language, s.Tags)
			found = true
		}
	}

	if !found {
		fmt.Println("No snippets matched your query.")
	}
}

func containsTag(tags []string, query string) bool {
	for _, t := range tags {
		if strings.Contains(strings.ToLower(t), strings.ToLower(query)) {
			return true
		}
	}
	return false
}
