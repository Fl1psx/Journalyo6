package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Student struct {
	Name   string
	Grades []int
}

func (s Student) Average() float64 {
	if len(s.Grades) == 0 {
		return 0
	}
	sum := 0
	for _, g := range s.Grades {
		sum += g
	}
	return float64(sum) / float64(len(s.Grades))
}

func (s Student) String() string {
	status := "✓"
	if s.Average() < 3.0 {
		status = "⚠"
	}
	return fmt.Sprintf("%s %s: оценки %v (средний: %.2f)", status, s.Name, s.Grades, s.Average())
}

type Journal struct {
	students map[string]Student
}

func NewJournal() *Journal {
	return &Journal{students: make(map[string]Student)}
}

func (j *Journal) AddStudent(name string, grades []int) {
	j.students[name] = Student{Name: name, Grades: grades}
}

func (j *Journal) RemoveStudent(name string) bool {
	if _, exists := j.students[name]; exists {
		delete(j.students, name)
		return true
	}
	return false
}

func (j *Journal) ShowAll() {
	if len(j.students) == 0 {
		fmt.Println("В журнале нет студентов")
		return
	}
	
	// Создаем список студентов для сортировки
	studentList := make([]Student, 0, len(j.students))
	for _, s := range j.students {
		studentList = append(studentList, s)
	}
	
	// Сортируем по среднему баллу (по убыванию)
	sort.Slice(studentList, func(i, j int) bool {
		return studentList[i].Average() > studentList[j].Average()
	})
	
	fmt.Printf("\n=== ВСЕ СТУДЕНТЫ (%d) ===\n", len(studentList))
	for i, s := range studentList {
		fmt.Printf("%d. %s\n", i+1, s)
	}
}

func (j *Journal) FilterByAverage(threshold float64, below bool) {
	found := false
	
	// Создаем список для сортировки
	studentList := make([]Student, 0, len(j.students))
	for _, s := range j.students {
		if below && s.Average() < threshold || !below && s.Average() >= threshold {
			studentList = append(studentList, s)
			found = true
		}
	}
	
	if !found {
		if below {
			fmt.Printf("Студентов со средним баллом ниже %.2f не найдено\n", threshold)
		} else {
			fmt.Printf("Студентов со средним баллом выше %.2f не найдено\n", threshold)
		}
		return
	}
	
	// Сортируем по среднему баллу
	sort.Slice(studentList, func(i, j int) bool {
		return studentList[i].Average() > studentList[j].Average()
	})
	
	condition := "ниже"
	if !below {
		condition = "выше или равно"
	}
	fmt.Printf("\n=== СТУДЕНТЫ СО СРЕДНИМ БАЛЛОМ %s %.2f (%d) ===\n", condition, threshold, len(studentList))
	for i, s := range studentList {
		fmt.Printf("%d. %s\n", i+1, s)
	}
}

func (j *Journal) ShowStatistics() {
	if len(j.students) == 0 {
		fmt.Println("Нет данных для статистики")
		return
	}
	
	totalStudents := len(j.students)
	var totalAverage, minAverage, maxAverage float64
	minAverage = 5.0
	maxAverage = 0.0
	
	gradeCount := make(map[int]int)
	
	for _, s := range j.students {
		avg := s.Average()
		totalAverage += avg
		
		if avg < minAverage {
			minAverage = avg
		}
		if avg > maxAverage {
			maxAverage = avg
		}
		
		for _, grade := range s.Grades {
			gradeCount[grade]++
		}
	}
	
	totalAverage /= float64(totalStudents)
	
	fmt.Println("\n=== СТАТИСТИКА ===")
	fmt.Printf("Всего студентов: %d\n", totalStudents)
	fmt.Printf("Средний балл по группе: %.2f\n", totalAverage)
	fmt.Printf("Лучший средний балл: %.2f\n", maxAverage)
	fmt.Printf("Худший средний балл: %.2f\n", minAverage)
	
	fmt.Println("\nРаспределение оценок:")
	for i := 5; i >= 1; i-- {
		count := gradeCount[i]
		percentage := float64(count) / float64(totalStudents*5) * 100 // примерное распределение
		fmt.Printf("  %d: %d оценок (%.1f%%)\n", i, count, percentage)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	journal := NewJournal()

	fmt.Println("=== ЖУРНАЛ УСПЕВАЕМОСТИ СТУДЕНТОВ ===")
	
	for {
		fmt.Println("\n1. Добавить студента")
		fmt.Println("2. Удалить студента")
		fmt.Println("3. Показать всех студентов")
		fmt.Println("4. Студенты с низким средним баллом (< порога)")
		fmt.Println("5. Студенты с высоким средним баллом (>= порога)")
		fmt.Println("6. Статистика успеваемости")
		fmt.Println("7. Выход")
		fmt.Print("Выберите действие: ")

		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			fmt.Print("Введите ФИО студента: ")
			scanner.Scan()
			name := strings.TrimSpace(scanner.Text())
			
			if name == "" {
				fmt.Println("Ошибка: ФИО не может быть пустым!")
				continue
			}
			
			if _, exists := journal.students[name]; exists {
				fmt.Printf("Студент '%s' уже существует!\n", name)
				continue
			}

			fmt.Print("Введите оценки через пробел (1-5): ")
			scanner.Scan()
			gradesInput := scanner.Text()

			var grades []int
			valid := true
			parts := strings.Fields(gradesInput)
			
			for _, p := range parts {
				if g, err := strconv.Atoi(p); err == nil && g >= 1 && g <= 5 {
					grades = append(grades, g)
				} else {
					fmt.Printf("Ошибка: оценка '%s' недопустима! Допустимы оценки от 1 до 5.\n", p)
					valid = false
					break
				}
			}
			
			if !valid {
				continue
			}
			
			if len(grades) == 0 {
				fmt.Println("Ошибка: необходимо ввести хотя бы одну оценку!")
				continue
			}

			journal.AddStudent(name, grades)
			fmt.Printf("✅ Студент '%s' успешно добавлен с оценками: %v\n", name, grades)

		case "2":
			fmt.Print("Введите ФИО студента для удаления: ")
			scanner.Scan()
			name := strings.TrimSpace(scanner.Text())
			
			if journal.RemoveStudent(name) {
				fmt.Printf("✅ Студент '%s' удален\n", name)
			} else {
				fmt.Printf("❌ Студент '%s' не найден\n", name)
			}

		case "3":
			journal.ShowAll()

		case "4":
			fmt.Print("Введите пороговый средний балл: ")
			scanner.Scan()
			thresholdStr := scanner.Text()

			threshold, err := strconv.ParseFloat(thresholdStr, 64)
			if err != nil || threshold < 1 || threshold > 5 {
				fmt.Println("Ошибка: введите число от 1 до 5!")
				continue
			}

			journal.FilterByAverage(threshold, true)

		case "5":
			fmt.Print("Введите пороговый средний балл: ")
			scanner.Scan()
			thresholdStr := scanner.Text()

			threshold, err := strconv.ParseFloat(thresholdStr, 64)
			if err != nil || threshold < 1 || threshold > 5 {
				fmt.Println("Ошибка: введите число от 1 до 5!")
				continue
			}

			journal.FilterByAverage(threshold, false)

		case "6":
			journal.ShowStatistics()

		case "7":
			fmt.Println("До свидания!")
			return

		default:
			fmt.Println("❌ Неверный выбор! Пожалуйста, выберите от 1 до 7")
		}
	}
}