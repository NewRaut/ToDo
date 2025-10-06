package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Task — структура для одной задачи (ОБНОВЛЕНО: добавлено Priority)
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	Priority    string `json:"priority"` // "Low", "Medium", "High"
}

// TodoList — список задач
type TodoList struct {
	Tasks []Task `json:"tasks"`
}

// LoadTasks загружает задачи из файла (с фиксом для пустого файла)
func (tl *TodoList) LoadTasks(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			tl.Tasks = []Task{}
			return nil
		}
		return err
	}
	if len(data) == 0 {
		tl.Tasks = []Task{}
		return nil
	}
	return json.Unmarshal(data, tl)
}

// SaveTasks сохраняет задачи в файл
func (tl *TodoList) SaveTasks(filename string) error {
	data, err := json.MarshalIndent(tl, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

// AddTask добавляет новую задачу (ОБНОВЛЕНО: с приоритетом)
func (tl *TodoList) AddTask(description, priority string) {
	id := len(tl.Tasks) + 1
	task := Task{ID: id, Description: description, Completed: false, Priority: priority}
	tl.Tasks = append(tl.Tasks, task)
}

// ToggleTask переключает статус задачи
func (tl *TodoList) ToggleTask(id int) bool {
	for i := range tl.Tasks {
		if tl.Tasks[i].ID == id {
			tl.Tasks[i].Completed = !tl.Tasks[i].Completed
			return true
		}
	}
	return false
}

// DeleteTask удаляет задачу
func (tl *TodoList) DeleteTask(id int) bool {
	for i := range tl.Tasks {
		if tl.Tasks[i].ID == id {
			tl.Tasks = append(tl.Tasks[:i], tl.Tasks[i+1:]...)
			return true
		}
	}
	return false
}

// SortByPriority сортирует задачи по приоритету (НОВОЕ: High -> Medium -> Low)
func (tl *TodoList) SortByPriority() {
	priorityOrder := map[string]int{"High": 3, "Medium": 2, "Low": 1}
	sort.Slice(tl.Tasks, func(i, j int) bool {
		// Сначала по приоритету (высокий выше), потом по ID
		if tl.Tasks[i].Priority != tl.Tasks[j].Priority {
			return priorityOrder[tl.Tasks[i].Priority] > priorityOrder[tl.Tasks[j].Priority]
		}
		return tl.Tasks[i].ID < tl.Tasks[j].ID
	})
}

// GetTasksContainer возвращает контейнер с чекбоксами для задач (ОБНОВЛЕНО: показывает приоритет)
func (tl *TodoList) GetTasksContainer(app *TodoApp) fyne.CanvasObject {
	if len(tl.Tasks) == 0 {
		return widget.NewLabel("Нет задач. Добавьте первую!")
	}

	vbox := container.NewVBox()
	for _, task := range tl.Tasks {
		status := " "
		if task.Completed {
			status = "✓"
		}
		check := widget.NewCheck(fmt.Sprintf("[%s] %s %d: %s", task.Priority, status, task.ID, task.Description), func(checked bool) {
			// При клике на чекбокс переключаем статус
			tl.ToggleTask(task.ID)
			app.refreshTasks() // Перестраиваем список для визуального обновления
			app.saveTasks()    // Автосохранение
		})
		check.Checked = task.Completed // Устанавливаем начальный статус
		vbox.Add(check)
	}
	return vbox
}

// GUI-структура для приложения (ОБНОВЛЕНО: добавлена кнопка сортировки)
type TodoApp struct {
	window    fyne.Window
	todo      *TodoList
	tasksCont fyne.CanvasObject
	filename  string
	content   fyne.CanvasObject
	addBtn    *widget.Button
	deleteBtn *widget.Button
	sortBtn   *widget.Button // НОВОЕ: кнопка сортировки
}

func (a *TodoApp) refreshTasks() {
	// Перестраиваем список задач
	a.tasksCont = a.todo.GetTasksContainer(a)
	// Обновляем основной контент (ОБНОВЛЕНО: добавлена кнопка сортировки)
	a.content = container.NewVScroll(container.NewVBox(a.tasksCont, container.NewHBox(a.addBtn, a.sortBtn, a.deleteBtn)))
	a.window.SetContent(a.content)
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Todo List GUI")
	myWindow.Resize(fyne.NewSize(500, 600)) // Увеличили размер для большего списка

	todo := &TodoList{}
	filename := "tasks.json"
	if err := todo.LoadTasks(filename); err != nil {
		dialog.ShowError(err, myWindow)
	}

	app := &TodoApp{
		window:   myWindow,
		todo:     todo,
		filename: filename,
	}

	// Кнопки (ИСПРАВЛЕНО: Callback перед parent)
	app.addBtn = widget.NewButton("Добавить задачу", func() {
		// Выпадающий список приоритетов
		priorities := []string{"Low", "Medium", "High"}
		prioritySelect := widget.NewSelect(priorities, nil)
		prioritySelect.SetSelected("Medium") // По умолчанию Medium

		descEntry := widget.NewEntry()
		descEntry.SetPlaceHolder("Введите описание задачи...")

		// ИСПРАВЛЕНО: Callback перед parent (myWindow)
		dialog.ShowCustomConfirm("Новая задача", "Добавить", "Отмена", container.NewVBox(
			widget.NewLabel("Описание:"),
			descEntry,
			widget.NewLabel("Приоритет:"),
			prioritySelect,
		), func(ok bool) {
			if ok && descEntry.Text != "" {
				app.todo.AddTask(descEntry.Text, prioritySelect.Selected)
				app.refreshTasks() // Обновляем список
				app.saveTasks()
			}
		}, myWindow)
	})

	app.sortBtn = widget.NewButton("Сортировать по приоритету", func() {
		app.todo.SortByPriority()
		app.refreshTasks() // Обновляем список
		app.saveTasks()
	})

	app.deleteBtn = widget.NewButton("Удалить выбранную", func() {
		if len(app.todo.Tasks) == 0 {
			dialog.ShowInformation("Инфо", "Нет задач для удаления", myWindow)
			return
		}
		dialog.ShowEntryDialog("Удалить задачу", "ID:", func(idStr string) {
			id, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil {
				dialog.ShowError(fmt.Errorf("Неверный ID: %v", err), myWindow)
				return
			}
			if app.todo.DeleteTask(id) {
				app.refreshTasks() // Обновляем список
				app.saveTasks()
				dialog.ShowInformation("Успех", "Задача удалена!", myWindow)
			} else {
				dialog.ShowError(fmt.Errorf("Задача с ID %d не найдена", id), myWindow)
			}
		}, myWindow)
	})

	// Инициализируем список задач
	app.refreshTasks()

	// Главное меню
	menu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Сохранить", func() { app.saveTasks() }),
			fyne.NewMenuItem("Загрузить", func() {
				dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
					if err != nil || reader == nil {
						return
					}
					defer reader.Close()
					data, _ := ioutil.ReadAll(reader)
					if len(data) == 0 {
						app.todo.Tasks = []Task{}
					} else {
						json.Unmarshal(data, app.todo)
					}
					app.refreshTasks()
				}, myWindow)
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Выход", func() { myApp.Quit() }),
		),
		fyne.NewMenu("Help",
			fyne.NewMenuItem("О программе", func() {
				dialog.ShowInformation("Todo List", "GUI версия на Fyne + Go с приоритетами", myWindow)
			}),
		),
	)

	myWindow.SetMainMenu(menu)
	myWindow.ShowAndRun()
}

func (a *TodoApp) saveTasks() {
	if err := a.todo.SaveTasks(a.filename); err != nil {
		dialog.ShowError(err, a.window)
	}
}
