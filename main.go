package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Todo{})

	tmpl := template.Must(template.ParseFiles("template/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err:  %v", err)
				return
			}

			todo := r.FormValue("todo")
			assignee := r.FormValue("assignee")
			deadline := r.FormValue("deadline")
			db.Create(&Todo{Title: todo, Assignee: assignee, Deadline: deadline, Done: false})
		}

		var todos []Todo
		db.Find(&todos)
		data := TodoPageData{
			PageTitle: "Employee TODO List",
			Todos:     todos,
		}
		tmpl.Execute(w, data)
	})

	http.HandleFunc("/edit/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/edit/")
		var todo Todo
		db.First(&todo, id)
		todo.Title = r.FormValue("editTodo")
		todo.Assignee = r.FormValue("editAssignee")
		todo.Deadline = r.FormValue("editDeadline")
		db.Save(&todo)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/done/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/done/")
		var todo Todo
		db.First(&todo, id)
		todo.Done = true
		db.Save(&todo)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/delete/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/delete/")
		db.Delete(&Todo{}, id)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.ListenAndServe(":8080", nil)
}

type Todo struct {
	gorm.Model
	ID       uint
	Title    string
	Assignee string
	Deadline string
	Done     bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}
