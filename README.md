# TODO App

A powerful command-line TODO application written in Go to help you manage your tasks efficiently with advanced features like priorities, due dates, tags, and search.

[![CI](https://github.com/rahul4507/todo/actions/workflows/ci.yml/badge.svg)](https://github.com/rahul4507/todo/actions/workflows/ci.yml)

## Features

✨ **Core Features**
- Add, edit, and delete tasks
- Mark tasks as completed/incomplete
- Auto-sort: completed tasks move to bottom
- Persistent storage (JSON)

🎯 **Advanced Features**
- **Priority Levels**: High 🔴, Medium 🟡, Low 🟢
- **Due Dates**: Set deadlines with overdue detection
- **Tags**: Organize tasks with custom tags
- **Search**: Find tasks by text or tags
- **Statistics**: Track completion rates
- **Interactive Mode**: Full-featured TUI

## Installation

### From Source

1. Clone the repository:
    ```sh
    git clone https://github.com/rahul4507/todo.git
    cd TODO-APP
    ```

2. Build the application:
    ```sh
    go build -o todo cmd/todo/main.go
    ```

3. (Optional) Move to PATH:
    ```sh
    sudo mv todo /usr/local/bin/
    ```

## Usage

### Basic Commands

```sh
# Add a task
./todo add "Buy groceries"

# List all tasks
./todo list

# Complete a task
./todo complete 1

# Delete a task
./todo delete 2

# Edit a task
./todo edit 1 "Buy groceries and cook dinner"
```

### Priority Management

```sh
# Set priority (high, medium, low)
./todo priority 1 high
./todo priority 2 low

# Priority indicators:
# 🔴 high    - Critical/urgent tasks
# 🟡 medium  - Normal priority (default)
# 🟢 low     - Nice to have
```

### Due Dates

```sh
# Set due date (YYYY-MM-DD format)
./todo due 1 2025-12-31

# View overdue tasks
./todo overdue
```

### Tags

```sh
# Add tags to organize tasks
./todo tag 1 work
./todo tag 1 urgent

# Remove a tag
./todo untag 1 urgent
```

### Search & Filter

```sh
# Search tasks by text or tags
./todo search "groceries"
./todo search "work"

# View statistics
./todo stats
```

### Batch Operations

```sh
# Clear all completed tasks
./todo clear

# Uncomplete a task (reopen)
./todo uncomplete 3
```

### Interactive Mode

```sh
# Launch interactive mode
./todo -i
```

In interactive mode, you get a full-featured interface with:
- Real-time task list display
- Live statistics
- All commands available
- Auto-refresh after each action

## Examples

```sh
# Create a high-priority work task with deadline
./todo add "Finish project presentation"
./todo priority 1 high
./todo due 1 2025-10-15
./todo tag 1 work
./todo tag 1 urgent

# Search for work tasks
./todo search work

# View the formatted list
./todo list
# Output:
# 1. [ ] 🔴 Finish project presentation 📅 2025-10-15 🏷️  work, urgent
```

## CI/CD

This project uses GitHub Actions for continuous integration:

- ✅ Automated testing on every PR
- 🔍 Code linting with golangci-lint
- 🏗️ Build verification
- 📊 Test coverage reports

See `.github/workflows/ci.yml` for details.

## Development

### Run Tests

```sh
go test ./... -v
```

### Test Coverage

```sh
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Project Structure

```
TODO-APP/
├── cmd/
│   └── todo/
│       └── main.go          # CLI entry point
├── internal/
│   └── todo/
│       ├── todo.go          # Core logic
│       └── todo_test.go     # Unit tests
├── .github/
│   └── workflows/
│       └── ci.yml           # CI/CD pipeline
└── todos.json               # Data storage
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

All PRs will automatically run tests via GitHub Actions.

## License

This project is licensed under the MIT License.
