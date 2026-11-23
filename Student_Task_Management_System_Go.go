package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type StudentManager struct {
	Students map[int]Student
}

func (sm *StudentManager) AddStudent(s Student) {
	sm.Students[s.ID] = s
}

func (sm *StudentManager) GetStudent(id int) (Student, bool) {
	student, exists := sm.Students[id]
	return student, exists
}

func (sm *StudentManager) DeleteStudent(id int) {
	delete(sm.Students, id)
}

func (sm *StudentManager) SaveToFile() error {
	data, err := json.MarshalIndent(sm.Students, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("students.json", data, 0644)
}

func main() {
	manager := StudentManager{
		Students: make(map[int]Student),
	}

	manager.AddStudent(Student{
		ID: 1, Name: "Akshat", Email: "akshat@example.com",
	})

	manager.AddStudent(Student{
		ID: 2, Name: "Rahul", Email: "rahul@example.com",
	})

	student, exists := manager.GetStudent(1)

	if exists {
		fmt.Println("Student:", student)
	}

	manager.DeleteStudent(2)

	if err := manager.SaveToFile(); err != nil {
		fmt.Println("Save error:", err)
	}
}
