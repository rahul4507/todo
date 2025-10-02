package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rahul4507/todo/internal/todo"
)

const (
	todoFile = "todos.json"
)

func main() {
	//define flags
	interactiveFlag := flag.Bool("i", false, "Run in interactive mode")
	helpFlag := flag.Bool("h", false, "Show help Information")

	//parse flags but keep access to non-flag arguments
	flag.Parse()
	args := flag.Args()

	// Show help if the -h flag is provided
	if *helpFlag {
		printHelp()
		return
	}

	// Load Existing todos
	todoList := todo.NewList()
	if _, err := os.Stat(todoFile); err == nil {
		if err := todoList.Load(todoFile); err != nil {
			fmt.Fprintln(os.Stderr, "Error Loading todos: ", err)
			os.Exit(1)
		}
	}

	//Handle interactive mode
	if *interactiveFlag {
		runInteractive(todoList)
		return
	}

	if len(args) == 0 {
		// default action print the todo list
		fmt.Println(todoList)
		return
	}

	command := args[0]
	switch command {
	case "add":
		if len(args) < 2 {
			fmt.Println("Error : Missing todo text")
			os.Exit(1)
		}

		// join all remaining args as todo text.
		text := strings.Join(args[1:], " ")
		if err := todoList.Add(text); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		saveTodos(todoList)

		fmt.Println("Added:", text)

	case "list":
		fmt.Println(todoList)

	case "complete":
		if len(args) < 2 {
			fmt.Println("Error : Missing the item number")
			os.Exit(1)
		}

		num, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: Invalid item number: ", args[1])
			os.Exit(1)
		}

		if err := todoList.Complete(num - 1); err != nil {
			fmt.Fprintln(os.Stderr, "Error completing todo:", err)
			os.Exit(1)
		}

		saveTodos(todoList)
		fmt.Println("Marked item as completed")

	case "uncomplete":
		if len(args) < 2 {
			fmt.Println("Error: Missing the item number")
			os.Exit(1)
		}

		num, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: Invalid item number:", args[1])
			os.Exit(1)
		}

		if err := todoList.Uncomplete(num - 1); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		saveTodos(todoList)
		fmt.Println("Marked item as incomplete")

	case "delete", "remove":
		if len(args) < 2 {
			fmt.Println("Error: Missing the item number")
			os.Exit(1)
		}

		num, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: Invalid item number:", args[1])
			os.Exit(1)
		}

		if err := todoList.Delete(num - 1); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		saveTodos(todoList)
		fmt.Println("Deleted item")

	case "edit":
		if len(args) < 3 {
			fmt.Println("Error: Missing item number or new text")
			os.Exit(1)
		}

		num, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: Invalid item number:", args[1])
			os.Exit(1)
		}

		newText := strings.Join(args[2:], " ")
		if err := todoList.Edit(num-1, newText); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		saveTodos(todoList)
		fmt.Println("Updated item")

	case "clear":
		count := todoList.ClearCompleted()
		saveTodos(todoList)
		fmt.Printf("Cleared %d completed item(s)\n", count)

	case "stats":
		stats := todoList.GetStats()
		fmt.Printf("Total: %d | Pending: %d | Completed: %d\n",
			stats.Total, stats.Pending, stats.Completed)

	case "priority":
		if len(args) < 3 {
			fmt.Println("Error: Missing item number or priority level")
			fmt.Println("Usage: todo priority <n> <high|medium|low>")
			os.Exit(1)
		}

		num, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: Invalid item number:", args[1])
			os.Exit(1)
		}

		priority := todo.ParsePriority(args[2])
		if err := todoList.SetPriority(num-1, priority); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		saveTodos(todoList)
		fmt.Printf("Set priority to %s\n", priority)

	case "due":
		if len(args) < 3 {
			fmt.Println("Error: Missing item number or due date")
			fmt.Println("Usage: todo due <n> <YYYY-MM-DD>")
			os.Exit(1)
		}

		num, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: Invalid item number:", args[1])
			os.Exit(1)
		}

		dueDate, err := time.Parse("2006-01-02", args[2])
		if err != nil {
			fmt.Println("Error: Invalid date format. Use YYYY-MM-DD")
			os.Exit(1)
		}

		if err := todoList.SetDueDate(num-1, dueDate); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		saveTodos(todoList)
		fmt.Printf("Set due date to %s\n", dueDate.Format("2006-01-02"))

	case "tag":
		if len(args) < 3 {
			fmt.Println("Error: Missing item number or tag")
			fmt.Println("Usage: todo tag <n> <tag>")
			os.Exit(1)
		}

		num, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: Invalid item number:", args[1])
			os.Exit(1)
		}

		tag := args[2]
		if err := todoList.AddTag(num-1, tag); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		saveTodos(todoList)
		fmt.Printf("Added tag: %s\n", tag)

	case "untag":
		if len(args) < 3 {
			fmt.Println("Error: Missing item number or tag")
			fmt.Println("Usage: todo untag <n> <tag>")
			os.Exit(1)
		}

		num, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error: Invalid item number:", args[1])
			os.Exit(1)
		}

		tag := args[2]
		if err := todoList.RemoveTag(num-1, tag); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		saveTodos(todoList)
		fmt.Printf("Removed tag: %s\n", tag)

	case "search":
		if len(args) < 2 {
			fmt.Println("Error: Missing search query")
			os.Exit(1)
		}

		query := strings.Join(args[1:], " ")
		results := todoList.Search(query)

		if len(results) == 0 {
			fmt.Println("No items found")
		} else {
			fmt.Printf("Found %d item(s):\n", len(results))
			for i, item := range results {
				status := " "
				if item.Done {
					status = "✓"
				}
				fmt.Printf("%d. [%s] %s\n", i+1, status, item.Text)
			}
		}

	case "overdue":
		results := todoList.GetOverdue()
		if len(results) == 0 {
			fmt.Println("No overdue items")
		} else {
			fmt.Printf("Overdue items (%d):\n", len(results))
			for i, item := range results {
				dueStr := item.DueDate.Format("2006-01-02")
				fmt.Printf("%d. %s (Due: %s)\n", i+1, item.Text, dueStr)
			}
		}

	case "help":
		printHelp()

	default:
		fmt.Printf("Unknown Command: %s\n", command)
		fmt.Println("Run 'todo help' for available commands")
		os.Exit(1)
	}

}

func saveTodos(list *todo.List) {
	if err := list.Save(todoFile); err != nil {
		fmt.Fprintln(os.Stderr, "Error saving todos: ", err)
		os.Exit(1)
	}
}

func printHelp() {
	helpText := `
Todo - A powerful command line todo manager

Usage:
  todo [command] [arguments]
  todo [flags]

Commands:
  add <text>              Add a new todo item
  list                    List all todo items
  complete <n>            Mark item n as completed
  uncomplete <n>          Mark item n as incomplete
  delete <n>              Delete item n
  edit <n> <text>         Edit the text of item n
  clear                   Remove all completed items
  stats                   Show statistics

  priority <n> <level>    Set priority (high/medium/low)
  due <n> <YYYY-MM-DD>    Set due date
  tag <n> <tag>           Add a tag to item
  untag <n> <tag>         Remove a tag from item

  search <query>          Search tasks by text or tag
  overdue                 Show overdue tasks

  help                    Show this help message

Flags:
  -h                      Show this help message
  -i                      Run in interactive mode

Examples:
  todo add "Learn Go testing"
  todo list
  todo complete 2
  todo priority 1 high
  todo due 1 2025-12-31
  todo tag 1 work
  todo search "go"
  todo overdue
  todo -i

Priority Levels:
  🔴 high    - Critical/urgent tasks
  🟡 medium  - Normal priority (default)
  🟢 low     - Nice to have

Symbols:
  [✓] - Completed task
  [ ] - Pending task
  📅 - Due date
  🏷️  - Tags
`
	fmt.Println(helpText)
}

func runInteractive(list *todo.List) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n" + list.String())
		stats := list.GetStats()
		fmt.Printf("\nStats: Total: %d | Pending: %d | Completed: %d\n",
			stats.Total, stats.Pending, stats.Completed)
		fmt.Println("\nCommands: add, complete, uncomplete, delete, edit, clear, help, quit")
		fmt.Print("\n> ")

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}
		cmd := parts[0]

		switch cmd {
		case "add":
			if len(parts) < 2 {
				fmt.Println("Error: missing todo text")
				continue
			}
			text := strings.Join(parts[1:], " ")
			err := list.Add(text)
			if err != nil {
				fmt.Println(err)
				continue
			}
			saveTodos(list)
			fmt.Println("Added:", text)

		case "complete":
			if len(parts) < 2 {
				fmt.Println("Error: missing item number")
				continue
			}
			num, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Error: invalid item number")
				continue
			}
			if err := list.Complete(num - 1); err != nil {
				fmt.Println("Error:", err)
				continue
			}
			saveTodos(list)
			fmt.Println("Marked item as completed")

		case "uncomplete":
			if len(parts) < 2 {
				fmt.Println("Error: missing item number")
				continue
			}
			num, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Error: invalid item number")
				continue
			}
			if err := list.Uncomplete(num - 1); err != nil {
				fmt.Println("Error:", err)
				continue
			}
			saveTodos(list)
			fmt.Println("Marked item as incomplete")

		case "delete", "remove":
			if len(parts) < 2 {
				fmt.Println("Error: missing item number")
				continue
			}
			num, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Error: invalid item number")
				continue
			}
			if err := list.Delete(num - 1); err != nil {
				fmt.Println("Error:", err)
				continue
			}
			saveTodos(list)
			fmt.Println("Deleted item")

		case "edit":
			if len(parts) < 3 {
				fmt.Println("Error: missing item number or new text")
				continue
			}
			num, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Error: invalid item number")
				continue
			}
			newText := strings.Join(parts[2:], " ")
			if err := list.Edit(num-1, newText); err != nil {
				fmt.Println("Error:", err)
				continue
			}
			saveTodos(list)
			fmt.Println("Updated item")

		case "clear":
			count := list.ClearCompleted()
			saveTodos(list)
			fmt.Printf("Cleared %d completed item(s)\n", count)

		case "stats":
			s := list.GetStats()
			fmt.Printf("Total: %d | Pending: %d | Completed: %d\n",
				s.Total, s.Pending, s.Completed)

		case "list":
			fmt.Println(list)

		case "help":
			fmt.Println("\nAvailable commands:")
			fmt.Println("  add <text>       - Add a new todo")
			fmt.Println("  list             - List all todos")
			fmt.Println("  complete <n>     - Mark item n as completed")
			fmt.Println("  uncomplete <n>   - Mark item n as incomplete")
			fmt.Println("  delete <n>       - Delete item n")
			fmt.Println("  edit <n> <text>  - Edit the text of item n")
			fmt.Println("  clear            - Remove all completed items")
			fmt.Println("  stats            - Show statistics")
			fmt.Println("  help             - Show this help message")
			fmt.Println("  quit             - Exit the program")

		case "quit", "exit":
			return

		default:
			fmt.Println("Unknown command:", cmd, "- Type 'help' for available commands")
		}
	}
}
