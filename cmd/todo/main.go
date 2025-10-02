package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

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

	case "help":
		printHelp()

	default:
		fmt.Printf("Unknown Command: %s\n", command)
		fmt.Println("Available Commands: add")
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
Todo - A simple command line todo manager

Usage:
  todo [command] [arguments]
  todo [flags]

Commands:
  add <text>     Add a new todo item
  list           List all todo items
  complete <n>   Mark item n as completed
  help           Show this help message

Flags:
  -h             Show this help message
  -i             Run in interactive mode

Examples:
  todo add "Learn Go testing"
  todo list
  todo complete 2
  todo -i
`
	fmt.Println(helpText)
}

func runInteractive(list *todo.List) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n" + list.String())
		fmt.Println("\nCommands:")
		fmt.Println("  add <text>    - Add a new todo")
		fmt.Println("  complete <n>  - Mark item n as completed")
		fmt.Println("  help          - get help")
		fmt.Println("  quit          - Exit the program")
		fmt.Print("\n> ")

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		parts := strings.SplitN(input, " ", 2)
		cmd := parts[0]

		switch cmd {
		case "add":
			if len(parts) < 2 {
				fmt.Println("Error: missing todo text")
				continue
			}
			err := list.Add(parts[1])
			if err != nil {
				fmt.Println(err)
				continue
			}
			saveTodos(list)
			fmt.Println("Added:", parts[1])

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

		case "help":
			fmt.Println("\nAvailable commands:")
			fmt.Println("  add <text>    - Add a new todo")
			fmt.Println("  list          - List all todos")
			fmt.Println("  complete <n>  - Mark item n as completed")
			fmt.Println("  help          - Show this help message")
			fmt.Println("  quit          - Exit the program")

		case "quit", "exit":
			return

		default:
			fmt.Println("Unknown command:", cmd)
		}
	}
}
