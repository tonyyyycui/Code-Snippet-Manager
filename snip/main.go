package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

type Snippet struct {
	Name     string   `json:"name"`
	Language string   `json:"language"`
	Tags     []string `json:"tags"`
	Content  string   `json:"content"`
}

// const snippetFile = "snippets.json"

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found, relying on system env")
	}
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

func addSnippet(cmd *cobra.Command, args []string) {
	name, _ := cmd.Flags().GetString("name")
	language, _ := cmd.Flags().GetString("language")
	userTags, _ := cmd.Flags().GetString("tags")

	content, err := openEditor()
	if err != nil {
		log.Fatal(err)
	}

	// Generate GPT tags
	gptTags, err := GenerateTagsFromContent(content)
	if err != nil {
		fmt.Println("⚠️  Could not generate GPT tags, using user-provided tags.")
		fmt.Printf("Error details: %v\n", err)
		gptTags = []string{}
	}

	finalTags := gptTags
	if userTags != "" {
		finalTags = append(finalTags, strings.Split(userTags, ",")...)
	}

	db, err := initDB()
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	defer db.Close()

	_, err = db.Exec(
		`INSERT INTO snippets (name, language, tags, content) VALUES ($1, $2, $3, $4)`,
		name, language, strings.Join(finalTags, ","), content,
	)
	if err != nil {
		log.Fatal("Error inserting snippet:", err)
	}

	fmt.Println("✅ Snippet added with GPT tags:", finalTags)
}

func listSnippets(cmd *cobra.Command, args []string) {
	db, err := initDB()
	if err != nil {
		fmt.Println("DB connection error:", err)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, language, tags FROM snippets ORDER BY id")
	if err != nil {
		fmt.Println("Error querying snippets:", err)
		return
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var id int
		var name, language, tags string
		if err := rows.Scan(&id, &name, &language, &tags); err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		fmt.Printf("[%d] %s (%s) - Tags: %s\n", id, name, language, tags)
		found = true
	}

	if !found {
		fmt.Println("No snippets found.")
	}
}

func searchSnippets(cmd *cobra.Command, args []string) {
	query, _ := cmd.Flags().GetString("query")
	if query == "" {
		fmt.Println("Please provide a search query with -q")
		return
	}

	db, err := initDB()
	if err != nil {
		fmt.Println("DB connection error:", err)
		return
	}
	defer db.Close()

	// Search name, content, or tags for the query
	rows, err := db.Query(
		`SELECT id, name, language, tags FROM snippets
         WHERE LOWER(name) LIKE $1 OR LOWER(content) LIKE $1 OR LOWER(tags) LIKE $1
         ORDER BY id`, "%"+strings.ToLower(query)+"%")
	if err != nil {
		fmt.Println("Error querying snippets:", err)
		return
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var id int
		var name, language, tags string
		if err := rows.Scan(&id, &name, &language, &tags); err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		fmt.Printf("[%d] %s (%s) - Tags: %s\n", id, name, language, tags)
		found = true
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

func openEditor() (string, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano" // fallback if EDITOR not set
	}

	tmpFile, err := os.CreateTemp("", "snippet-*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	// Open the editor
	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	// Read edited content
	contentBytes, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", err
	}

	return string(contentBytes), nil
}

func initDB() (*sql.DB, error) {
	connStr := "dbname=snippets sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS snippets (
            id SERIAL PRIMARY KEY,
            name TEXT UNIQUE NOT NULL,
            language TEXT,
            tags TEXT,
            content TEXT
        )
    `)
	if err != nil {
		return nil, err
	}
	return db, nil
}
