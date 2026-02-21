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
		// Fallback to local if home is unavailable
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
		abs, _ := filepath.Abs(path)
		wd, _ := os.Getwd()
		slog.Error("Error opening storage file", "path", path, "abs", abs, "wd", wd, "error", err)
		return nil
	}
	defer file.Close()

	// Check if file is empty
	info, err := file.Stat()
	if err == nil && info.Size() == 0 {
		return []types.Task{}
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		slog.Error("Error reading CSV", "error", err)
		return nil
	}

	var tasks []types.Task
	for i, record := range records {
		// Skip header row and invalid records
		if i == 0 || len(record) < 3 {
			continue
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			slog.Error("Error parsing ID in CSV", "row", i+1, "id_value", record[0], "error", err)
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
		slog.Error("Error creating directory", "dir", dir, "error", err)
		return
	}

	file, err := os.Create(path)
	if err != nil {
		abs, _ := filepath.Abs(path)
		wd, _ := os.Getwd()
		slog.Error("Error creating file", "path", path, "abs", abs, "wd", wd, "error", err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"id", "text", "done"})

	for _, task := range tasks {
		writer.Write([]string{strconv.Itoa(task.ID), task.Desc, string(task.CurrStatus)})
	}
}

func AddTask(taskDesc string) []types.Task {
	tasks := LoadTasks()
	if tasks == nil {
		slog.Error("Cannot add task: failed to load existing tasks")
		return nil
	}
	newTask := types.Task{
		ID:         len(tasks) + 1,
		Desc:       taskDesc,
		CurrStatus: types.Pending,
	}
	tasks = append(tasks, newTask)
	SaveTasks(tasks)
	return tasks
}

func ViewTasks() {
	tasks := LoadTasks()
	if tasks == nil {
		slog.Error("Cannot view tasks: failed to load task list")
		return
	}
	if len(tasks) == 0 {
		fmt.Println("No tasks yet. Use 'add' to create one.")
		return
	}
	for _, task := range tasks {
		sts := " "
		if task.CurrStatus == types.Completed {
			sts = "✓"
		}
		fmt.Printf("[%s] %d. %s\n", sts, task.ID, task.Desc)
	}
}

func MarkDone(id int) {
	tasks := LoadTasks()
	if tasks == nil {
		slog.Error("Cannot mark task: failed to load task list")
		return
	}
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].CurrStatus = types.Completed
			slog.Info("Task marked as done", "id", id)
			SaveTasks(tasks)
			return
		}
	}
	slog.Error("Task not found", "id", id)
}
func DeleteTask(id int) {
	tasks := LoadTasks()
	if tasks == nil {
		slog.Error("Cannot delete task: failed to load task list")
		return
	}
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
