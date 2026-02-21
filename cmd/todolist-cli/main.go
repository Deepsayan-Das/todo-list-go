package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Deepsayan-Das/todo-list-go/internal/handlers"
)

const appName = "todo"
const appVersion = "1.0.0"

func printHelp() {
	fmt.Printf(`%s - A simple command-line task manager

USAGE:
  %s [command] [arguments]

DESCRIPTION:
  todo is a lightweight CLI tool to manage your personal tasks from the terminal.
  Tasks are persisted in a local CSV file. You can add, view, complete, and
  delete tasks without leaving your command line.

COMMANDS:
  add, a   <description>    Add a new task with the given description
  view, v                   List all current tasks with their status
  markdone, m  <id>         Mark a task as completed by its ID
  delete, x    <id>         Delete a task permanently by its ID

OPTIONS:
  -h, --help                Show this help message and exit
  -l, --list                Show a short summary of all available commands
  -V, --version             Show version information and exit

EXAMPLES:
  %s add "Buy groceries"
  %s a "Finish the Go project"
  %s view
  %s markdone 2
  %s m 3
  %s delete 1
  %s x 4

NOTES:
  - Task IDs are auto-assigned starting from 1
  - Completed tasks are shown with a ✓ checkmark when using 'view'
  - Deleting a task is permanent and cannot be undone

VERSION:
  %s v%s

`, appName, appName,
		appName, appName, appName,
		appName, appName, appName, appName,
		appName, appVersion)
}

func printList() {
	fmt.Printf(`%s — Available Commands

  %-20s %s
  %-20s %s
  %-20s %s
  %-20s %s
  %-20s %s
  %-20s %s
  %-20s %s

Run '%s --help' for full documentation.
`, appName,
		"add, a <desc>", "Add a new task",
		"view, v", "List all tasks",
		"markdone, m <id>", "Mark a task as done",
		"delete, x <id>", "Delete a task",
		"-h, --help", "Show full help manual",
		"-l, --list", "Show this command list",
		"-V, --version", "Show version info",
		appName)
}

func main() {
	fmt.Printf("%s v%s — Your Personal Task Manager\n", appName, appVersion)
	fmt.Println("Type -h or --help for full documentation")
	fmt.Println()

	if len(os.Args) < 2 {
		fmt.Println("No command provided. Run with -h for help.")
		return
	}

	cmd := os.Args[1]
	switch cmd {
	case "-h", "--help":
		printHelp()

	case "-l", "--list":
		printList()

	case "-V", "--version":
		fmt.Printf("%s version %s\n", appName, appVersion)

	case "add", "a":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Error: 'add' requires a task description\n\nUsage:\n  %s add \"<description>\"\n", appName)
			os.Exit(1)
		}
		taskDesc := os.Args[2]
		tasks := handlers.AddTask(taskDesc)
		if tasks == nil {
			os.Exit(1)
		}
		fmt.Printf("✓ Task added (ID: %d): %s\n", len(tasks), taskDesc)

	case "view", "v":
		handlers.ViewTasks()

	case "delete", "x":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Error: 'delete' requires a task ID\n\nUsage:\n  %s delete <id>\n", appName)
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid task ID '%s' — must be a number\n", os.Args[2])
			os.Exit(1)
		}
		handlers.DeleteTask(id)

	case "markdone", "m":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Error: 'markdone' requires a task ID\n\nUsage:\n  %s markdone <id>\n", appName)
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid task ID '%s' — must be a number\n", os.Args[2])
			os.Exit(1)
		}
		handlers.MarkDone(id)

	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command '%s'\n\nRun '%s --list' to see available commands.\n", cmd, appName)
		os.Exit(1)
	}
}
