# Todo List CLI (Go)

A lightweight, portable command-line task manager built with Go. 

## Features

- **Portability**: Tasks are stored in your home directory (`~/.todo-list-go/`), so your data is available anywhere on your system.
- **Simplicity**: No complex database setup. Uses a clean CSV-based storage.
- **Fast**: Built with Go for instantaneous startup and execution.

## Installation

To install the tool globally on your system, ensure you have Go installed and run:

```bash
go install ./cmd/todolist-cli/
```

This will install the `todolist-cli` executable into your `$GOPATH/bin` (or `$HOME/go/bin`). Ensure this directory is in your system's `PATH`.

## Usage

### Commands

| Command | Shorthand | Argument | Description |
|---------|-----------|----------|-------------|
| `add` | `a` | `"<task>"` | Add a new task |
| `view` | `v` | - | List all tasks |
| `markdone` | `m` | `<id>` | Mark a task as completed |
| `delete` | `x` | `<id>` | Remove a task permanently |

### Options

- `-h, --help`: Show full help documentation
- `-l, --list`: Show summary of available commands
- `-V, --version`: Show version information

## Examples

```bash
# Add a task
todolist-cli add "Finish README"

# View all tasks
todolist-cli view

# Mark task #1 as done
todolist-cli m 1

# Delete task #2
todolist-cli x 2
```

## Development

To run the project locally without installing:

```bash
go run ./cmd/todolist-cli/main.go view
```

## License

MIT
