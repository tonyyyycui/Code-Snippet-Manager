
# Code Snippet Manager

A simple tool for managing, storing, and organizing code snippets. Built with Go, this project aims to help developers quickly save, search, and retrieve useful code fragments.

## Features

- Store code snippets in a JSON file (with plans to migrate to SQL)
- Easily add, search, and manage snippets via CLI
- Future plans for OpenAI-powered auto-tagging and analysis
- Environment manager for OpenAI API keys

## Getting Started

### Prerequisites

- Go 1.18+
- (Optional) OpenAI API key for advanced features

### Installation

```bash
git clone https://github.com/tonyyyycui/Code-Snippet-Manager.git
cd Code-Snippet-Manager/snip
go build -o snip main.go
```

### Usage

```bash
./snip add "Title" "Your code snippet here" --tags="go,example"
./snip list
./snip search "keyword"
```

## Roadmap

- [ ] OpenAI API Key Environment Manager
- [ ] OpenAI Auto-Tagging and Analysis
- [ ] Migrate storage from JSON to SQL
- [X] Add-only file content management
- [ ] Add PostgreSQL auto-table management

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

## License

MIT