// Command todo serves one person's task list on the loopback interface.
package main

import (
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
)

//go:embed list.html
var markup string

var listPage = template.Must(template.New("list").Parse(markup))

type task struct {
	ID    int64
	Title string
}

type list struct {
	Completed bool
	Tasks     []task
}

type app struct{ db *sql.DB }

func main() {
	file := flag.String("db", "todo.sqlite", "SQLite database file")
	port := flag.Int("port", 8080, "loopback port (0 selects an available port)")
	flag.Parse()
	if err := serve(*file, *port); err != nil {
		log.Fatal(err)
	}
}

func serve(file string, port int) error {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return err
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		completed INTEGER NOT NULL DEFAULT 0 CHECK (completed IN (0, 1))
	)`)
	if err != nil {
		return err
	}
	a := app{db}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", a.show)
	mux.HandleFunc("GET /completed", a.show)
	mux.HandleFunc("POST /tasks", a.change)
	mux.HandleFunc("POST /tasks/{id}/{action}", a.change)
	listener, err := net.Listen("tcp4", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return err
	}
	host := listener.Addr().String()
	protected := http.NewCrossOriginProtection().Handler(mux)
	server := &http.Server{ReadHeaderTimeout: 5 * time.Second, Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != host {
			http.Error(w, "Use the local address printed at startup.", http.StatusForbidden)
			return
		}
		protected.ServeHTTP(w, r)
	})}
	fmt.Printf("http://%s\n", host)
	return server.Serve(listener)
}

func (a app) show(w http.ResponseWriter, r *http.Request) {
	view := list{Completed: r.URL.Path == "/completed"}
	rows, err := a.db.QueryContext(r.Context(), "SELECT id, title FROM tasks WHERE completed = ? ORDER BY id", view.Completed)
	if err != nil {
		storageError(w, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var item task
		if err := rows.Scan(&item.ID, &item.Title); err != nil {
			storageError(w, err)
			return
		}
		view.Tasks = append(view.Tasks, item)
	}
	if err := rows.Err(); err != nil {
		storageError(w, err)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	if err := listPage.Execute(w, view); err != nil {
		log.Print(err)
	}
}

func storageError(w http.ResponseWriter, err error) {
	log.Print(err)
	http.Error(w, "Could not access the task list.", http.StatusInternalServerError)
}

func (a app) change(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Could not read the form.", http.StatusBadRequest)
		return
	}
	action := r.PathValue("action")
	var id int64
	if action != "" {
		var err error
		id, err = strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil || id < 1 {
			http.NotFound(w, r)
			return
		}
	}
	title := strings.TrimSpace(r.PostForm.Get("title"))
	if (action == "" || action == "rename") && title == "" {
		http.Error(w, "Enter a task name.", http.StatusBadRequest)
		return
	}
	var result sql.Result
	var err error
	ctx := r.Context()
	switch action {
	case "":
		result, err = a.db.ExecContext(ctx, "INSERT INTO tasks (title) VALUES (?)", title)
	case "rename":
		result, err = a.db.ExecContext(ctx, "UPDATE tasks SET title = ? WHERE id = ?", title, id)
	case "complete", "reopen":
		result, err = a.db.ExecContext(ctx, "UPDATE tasks SET completed = ? WHERE id = ?", action == "complete", id)
	case "delete":
		result, err = a.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = ?", id)
	default:
		http.NotFound(w, r)
		return
	}
	if err != nil {
		storageError(w, err)
		return
	}
	changed, err := result.RowsAffected()
	if err != nil {
		storageError(w, err)
		return
	}
	if changed == 0 {
		http.NotFound(w, r)
		return
	}
	dest := "/"
	if r.PostForm.Get("view") == "completed" {
		dest = "/completed"
	}
	http.Redirect(w, r, dest, http.StatusSeeOther)
}
