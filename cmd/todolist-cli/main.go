package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Deepsayan-Das/todo-list-go/internal/handlers"
	"github.com/charmbracelet/lipgloss"
)

const appName = "todo"
const appVersion = "1.1.0"

var (
	// Styles
	primaryColor   = lipgloss.Color("#7D56F4")
	secondaryColor = lipgloss.Color("#04B575")
	errorColor     = lipgloss.Color("#FF5F87")
	grayColor      = lipgloss.Color("#767676")

	grayStyle = lipgloss.NewStyle().Foreground(grayColor)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(primaryColor).
			Padding(0, 1).
			Bold(true).
			MarginBottom(1)

	headerStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			MarginTop(1)

	descStyle = lipgloss.NewStyle().
			Foreground(grayColor).
			Italic(true)

	cmdStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)

	argStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))

	successStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true).
			PaddingLeft(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			PaddingLeft(1)

	infoStyle = lipgloss.NewStyle().
			Foreground(grayColor).
			PaddingLeft(1)
)

func printHelp() {
	title := titleStyle.Render(fmt.Sprintf("%s v%s", appName, appVersion))
	fmt.Println(title)
	fmt.Println(descStyle.Render("A lightweight terminal task manager for focused developers."))

	fmt.Println(headerStyle.Render("USAGE:"))
	fmt.Printf("  %s %s %s\n", cmdStyle.Render(appName), argStyle.Render("[command]"), grayStyle.Render("[arguments]"))

	fmt.Println(headerStyle.Render("COMMANDS:"))
	commands := [][]string{
		{"add, a", "<description>", "Add a new task"},
		{"view, v", "", "List all current tasks"},
		{"markdone, m", "<id>", "Mark a task as completed"},
		{"delete, x", "<id>", "Delete a task permanently"},
	}

	for _, c := range commands {
		cmdPart := fmt.Sprintf("  %-15s", cmdStyle.Render(c[0]))
		argPart := fmt.Sprintf("%-15s", argStyle.Render(c[1]))
		fmt.Printf("%s %s %s\n", cmdPart, argPart, grayStyle.Render(c[2]))
	}

	fmt.Println(headerStyle.Render("OPTIONS:"))
	options := [][]string{
		{"-h, --help", "Show this help message"},
		{"-l, --list", "Show command summary"},
		{"-V, --version", "Show version information"},
	}

	for _, o := range options {
		fmt.Printf("  %-30s %s\n", cmdStyle.Render(o[0]), grayStyle.Render(o[1]))
	}

	fmt.Println(headerStyle.Render("EXAMPLES:"))
	fmt.Printf("  %s add %s\n", appName, grayStyle.Render("\"Buy groceries\""))
	fmt.Printf("  %s view\n", appName)
	fmt.Printf("  %s markdone 2\n", appName)

	fmt.Println()
}

func printList() {
	title := headerStyle.Render(fmt.Sprintf("%s — Available Commands", appName))
	fmt.Println(title)
	fmt.Println()

	commands := [][]string{
		{"add, a <desc>", "Add a new task"},
		{"view, v", "List all tasks"},
		{"markdone, m <id>", "Mark a task as done"},
		{"delete, x <id>", "Delete a task"},
		{"-h, --help", "Show full help"},
	}

	for _, c := range commands {
		fmt.Printf("  %-25s %s\n", cmdStyle.Render(c[0]), grayStyle.Render(c[1]))
	}

	fmt.Println(infoStyle.Render(fmt.Sprintf("\nRun '%s --help' for full documentation.", appName)))
}

func main() {
	// Skip the default welcome message if commands are provided to keep it clean
	if len(os.Args) < 2 {
		fmt.Println(titleStyle.Render(fmt.Sprintf("%s v%s", appName, appVersion)))
		fmt.Println(infoStyle.Render("Your Personal Task Manager"))
		fmt.Println(infoStyle.Render("Type -h or --help for instructions\n"))
		fmt.Println(errorStyle.Render("No command provided."))
		return
	}

	cmd := os.Args[1]
	switch cmd {
	case "-h", "--help":
		printHelp()

	case "-l", "--list":
		printList()

	case "-V", "--version":
		fmt.Printf("%s version %s\n", cmdStyle.Render(appName), argStyle.Render(appVersion))

	case "add", "a":
		if len(os.Args) < 3 {
			fmt.Println(errorStyle.Render("Error: 'add' requires a task description"))
			fmt.Printf("\nUsage:\n  %s add %s\n", cmdStyle.Render(appName), argStyle.Render("\"<description>\""))
			os.Exit(1)
		}
		taskDesc := os.Args[2]
		tasks := handlers.AddTask(taskDesc)
		if tasks == nil {
			os.Exit(1)
		}
		fmt.Println(successStyle.Render(fmt.Sprintf("✓ Task added (ID: %d): %s", len(tasks), taskDesc)))

	case "view", "v":
		handlers.ViewTasks()

	case "delete", "x":
		if len(os.Args) < 3 {
			fmt.Println(errorStyle.Render("Error: 'delete' requires a task ID"))
			fmt.Printf("\nUsage:\n  %s delete %s\n", cmdStyle.Render(appName), argStyle.Render("<id>"))
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(errorStyle.Render(fmt.Sprintf("Error: Invalid task ID '%s' — must be a number", os.Args[2])))
			os.Exit(1)
		}

		// Fetch task to show name in confirmation
		tasks := handlers.LoadTasks()

		found := false
		taskDesc := ""
		for _, t := range tasks {
			if t.ID == id {
				taskDesc = t.Desc
				found = true
				break
			}
		}

		if !found {
			fmt.Println(errorStyle.Render(fmt.Sprintf("Error: Task %d not found", id)))
			os.Exit(1)
		}

		// Confirmation prompt
		prompt := lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render("?")
		message := fmt.Sprintf(" Are you sure you want to delete '%s'?", argStyle.Render(taskDesc))
		fmt.Printf("%s%s %s ", prompt, message, grayStyle.Render("(y/N)"))

		var response string
		fmt.Scanln(&response)

		if response == "y" || response == "Y" {
			handlers.DeleteTask(id)
		} else {
			fmt.Println(infoStyle.Render("Deletion cancelled."))
		}

	case "markdone", "m":
		if len(os.Args) < 3 {
			fmt.Println(errorStyle.Render("Error: 'markdone' requires a task ID"))
			fmt.Printf("\nUsage:\n  %s markdone %s\n", cmdStyle.Render(appName), argStyle.Render("<id>"))
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(errorStyle.Render(fmt.Sprintf("Error: Invalid task ID '%s' — must be a number", os.Args[2])))
			os.Exit(1)
		}
		handlers.MarkDone(id)

	default:
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error: unknown command '%s'", cmd)))
		fmt.Printf("\nRun '%s --list' to see available commands.\n", cmdStyle.Render(appName))
		os.Exit(1)
	}
}
