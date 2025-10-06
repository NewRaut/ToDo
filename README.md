Todo List GUI
A simple and elegant desktop todo list application built with Go and Fyne GUI toolkit.

Features
✅ Add, delete, and mark tasks as completed

🎯 Priority levels (Low, Medium, High)

🔄 Sort tasks by priority

💾 Automatic saving to JSON file

🖼️ Clean and intuitive GUI interface

📁 Load/Save functionality via menu

Installation
Prerequisites
Go 1.16 or higher

Fyne GUI toolkit

Build from source
bash
git clone <repository-url>
cd <project-directory>
go build -o todo-app
./todo-app
Usage
Adding a task: Click "Добавить задачу" (Add Task), enter description and select priority

Completing tasks: Check the checkbox next to any task

Sorting: Click "Сортировать по приоритету" (Sort by Priority) to organize tasks

Deleting: Click "Удалить выбранную" (Delete Selected) and enter task ID

File operations: Use the File menu to save/load task lists

Data Storage
Tasks are automatically saved to tasks.json in the application directory. The file uses JSON format for easy readability and manual editing if needed.

License
MIT License
