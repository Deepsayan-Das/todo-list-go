package handlers

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Deepsayan-Das/todo-list-go/internal/types"
)

func getStoragePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join("storage", "storage.csv")
	}
	return filepath.Join(home, ".todo-list-go", "storage.csv")
}

func LoadTasks() []types.Task {
	path := getStoragePath()

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []types.Task{}
		}
		slog.Error("Error opening storage file", "error", err)
		return []types.Task{}
	}
	defer file.Close()

	info, err := file.Stat()
	if err == nil && info.Size() == 0 {
		return []types.Task{}
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		slog.Error("Error reading CSV", "error", err)
		return []types.Task{}
	}

	var tasks []types.Task

	for i, record := range records {
		// Skip header
		if i == 0 {
			continue
		}

		if len(record) < 3 {
			slog.Warn("Skipping invalid CSV row", "row", i+1)
			continue
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			slog.Warn("Skipping row with invalid ID", "row", i+1, "value", record[0])
			continue
		}

		tasks = append(tasks, types.Task{
			ID:         id,
			Desc:       record[1],
			CurrStatus: types.Status(record[2]),
		})
	}

	return tasks
}

func SaveTasks(tasks []types.Task) {
	path := getStoragePath()
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		slog.Error("Error creating storage directory", "error", err)
		return
	}

	file, err := os.Create(path)
	if err != nil {
		slog.Error("Error creating storage file", "error", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"id", "text", "done"}); err != nil {
		slog.Error("Error writing header", "error", err)
		return
	}

	for _, task := range tasks {
		err := writer.Write([]string{
			strconv.Itoa(task.ID),
			task.Desc,
			string(task.CurrStatus),
		})
		if err != nil {
			slog.Error("Error writing task to CSV", "error", err)
		}
	}
}

func AddTask(taskDesc string) []types.Task {
	tasks := LoadTasks()

	// Find max ID
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}

	newTask := types.Task{
		ID:         maxID + 1,
		Desc:       taskDesc,
		CurrStatus: types.Pending,
	}

	tasks = append(tasks, newTask)
	SaveTasks(tasks)

	slog.Info("Task added", "id", newTask.ID)
	return tasks
}

func ViewTasks() {
	tasks := LoadTasks()

	if len(tasks) == 0 {
		fmt.Println("No tasks yet. Use 'add' to create one.")
		return
	}

	for _, task := range tasks {
		status := " "
		if task.CurrStatus == types.Completed {
			status = "✓"
		}
		fmt.Printf("[%s] %d. %s\n", status, task.ID, task.Desc)
	}
}

func MarkDone(id int) {
	tasks := LoadTasks()

	for i, task := range tasks {
		if task.ID == id {
			tasks[i].CurrStatus = types.Completed
			SaveTasks(tasks)
			slog.Info("Task marked as done", "id", id)
			return
		}
	}

	slog.Error("Task not found", "id", id)
}

func DeleteTask(id int) {
	tasks := LoadTasks()

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			SaveTasks(tasks)
			slog.Info("Task deleted", "id", id)
			return
		}
	}

	slog.Error("Task not found", "id", id)
}
