package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
)

type Task struct {
	ID int
	Description string
	Status string
}

const fileName = "tasks.json"

func loadTasks() ([]Task, error){
	if _, err := os.Stat(fileName); os.IsNotExist(err){
		return []Task{}, nil
	}

	data, err := os.ReadFile(fileName)

	if err != nil {
		return nil, err
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func saveTasks(tasks []Task) error{
	data, err := json.MarshalIndent(tasks, "", "	")

	if err != nil {
		return err
	}

	return os.WriteFile(fileName, data, 0644)
}

func main() {
	// Your code here
	if len(os.Args) < 2 {
		fmt.Println("Penggunaan: task-cli <command> [arguments]")
		fmt.Println("Command tersedia: add, list, mark-done, delete")
		return
	}

	command := os.Args[1]
	tasks, err := loadTasks()

	if err != nil {
		fmt.Println("Error saat membaca data: ", err)
		return
	}

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Harap masukkan deskripsi tugas.")
			fmt.Println("Contoh: task-cli add \"Beli susu\"")
			return
		}
		
		var randNumber int = rand.IntN(1000)
		desc := os.Args[2]
		newTask := Task{
			ID: randNumber,
			Description: desc,
			Status: "todo",
		}

		tasks = append(tasks, newTask)

		if err := saveTasks(tasks); err != nil {
			fmt.Println("Gagal menyimpan tugas: ", err)
			return
		}

		fmt.Printf("Tugas berhasil ditambahkan (ID: %d)\n", newTask.ID)

	case "list":
		if len(tasks) == 0 {
			fmt.Println("Belum ada tugas tersimpan")
			return
		}

		fmt.Println("\n--- Daftar tugas ---")
		
		for _, task := range tasks {
			fmt.Printf("[%d] %s (Status: %s)\n", task.ID, task.Description, task.Status)
		}
	
	case "mark-done":
		if len(os.Args) < 3 {
			fmt.Println("Harap masukkan ID tugas")
			fmt.Println("Contoh: task-cli mark-done 100")
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("ID harus berupa angka")
			return
		}

		found := false 

		for i := range tasks {
			if tasks[i].ID == id {
				tasks[i].Status = "done"
				found = true
				break
			}
		}

		if !found {
			fmt.Printf("Tugas dengan ID %d tidak ditemukan. \n", id)
			return
		}

		saveTasks(tasks)
		fmt.Printf("Tugas dengan ID %d telah ditandai sebagai selesai \n", id)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Harap masukkan ID tugas yang ingin dihapus")
			fmt.Println("Contoh: task-cli delete 100")
			return
		}

		id, err := strconv.Atoi(os.Args[2])

		if err != nil {
			fmt.Println("ID harus berupa angka")
			return
		}

		updatedTasks := []Task{}
		
		found := false

		for _, task := range tasks {
			if task.ID == id {
				found = true
				continue
			}
			updatedTasks = append(updatedTasks, task)
		}

		if !found {
			fmt.Println("ID tugas tidak ditemukan")
			return
		}

		saveTasks(updatedTasks)
		fmt.Printf("Tugas dengan ID %d berhasil dihapus. \n", id)
	
	default:
		fmt.Println("Perintah tidak dikenali: ", command)
	}
}