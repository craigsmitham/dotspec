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
	"slices"
	"strconv"
	"strings"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
)

//go:embed list.html
var markup string

var listPage = template.Must(template.New("list").Parse(markup))

type task struct {
	ID          int64
	Title       string
	WaitingFor  string
	FollowUp    string
	FollowUpNow bool
}

type list struct {
	State string
	Title string
	Tasks []task
}

type app struct{ db *sql.DB }

// A pending reply belongs to one unfinished task. Its presence distinguishes
// waiting commitments from errands the person can act on now.
const schema = `CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    completed INTEGER NOT NULL DEFAULT 0 CHECK (completed IN (0, 1))
);
CREATE TABLE IF NOT EXISTS pending_replies (
    task_id INTEGER PRIMARY KEY REFERENCES tasks(id),
    waiting_for TEXT NOT NULL,
    follow_up TEXT NOT NULL
);
CREATE VIEW IF NOT EXISTS task_list AS
    SELECT id, title, waiting_for, follow_up,
        CASE WHEN completed THEN 'completed'
             WHEN task_id IS NOT NULL THEN 'waiting'
             ELSE 'open' END AS state
    FROM tasks LEFT JOIN pending_replies ON task_id = id`

var listTitles = map[string]string{"open": "Open", "waiting": "Waiting", "completed": "Completed"}

var allowedStates = map[string][]string{
	"wait":             {"open"},
	"resume":           {"waiting"},
	"update-follow-up": {"waiting"},
	"complete":         {"open", "waiting"},
	"reopen":           {"completed"},
	"rename":           {"open", "waiting", "completed"},
	"delete":           {"open", "waiting", "completed"},
}

func validFollowUpDate(day string) bool {
	year, _, _ := strings.Cut(day, "-")
	if len(year) < 4 || strings.Trim(year, "0123456789") != "" || strings.Trim(year, "0") == "" {
		return false
	}
	// HTML dates allow longer years. Gregorian leap years repeat every 400
	// years, so the last four digits suffice for the month and day check.
	_, err := time.Parse(time.DateOnly, day[len(year)-4:])
	return err == nil
}

func onOrBefore(day, other string) bool {
	// The month and day have fixed widths. Compare the significant year
	// lengths first so 10000 sorts after 9999.
	day = strings.TrimLeft(day, "0")
	other = strings.TrimLeft(other, "0")
	if len(day) != len(other) {
		return len(day) < len(other)
	}
	return day <= other
}

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
	_, err = db.Exec(schema)
	if err != nil {
		return err
	}
	a := app{db}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", a.show)
	mux.HandleFunc("GET /waiting", a.show)
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
	state := strings.TrimPrefix(r.URL.Path, "/")
	if state == "" {
		state = "open"
	}
	view := list{State: state, Title: listTitles[state]}
	today := time.Now().Format(time.DateOnly)
	rows, err := a.db.QueryContext(r.Context(), `SELECT id, title, coalesce(waiting_for, ''), coalesce(follow_up, '')
        FROM task_list WHERE state = ? ORDER BY id`, state)
	if err != nil {
		storageError(w, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var item task
		if err := rows.Scan(&item.ID, &item.Title, &item.WaitingFor, &item.FollowUp); err != nil {
			storageError(w, err)
			return
		}
		item.FollowUpNow = state == "waiting" && onOrBefore(item.FollowUp, today)
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
		if err != nil || id < 1 || allowedStates[action] == nil {
			http.NotFound(w, r)
			return
		}
	}
	ctx := r.Context()
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		storageError(w, err)
		return
	}
	defer tx.Rollback()
	state := "open"
	if action != "" {
		err := tx.QueryRowContext(ctx, "SELECT state FROM task_list WHERE id = ?", id).Scan(&state)
		if err == sql.ErrNoRows {
			http.Error(w, "This entry is no longer on your list.", http.StatusNotFound)
			return
		}
		if err != nil {
			storageError(w, err)
			return
		}
		if state != r.PostForm.Get("view") || !slices.Contains(allowedStates[action], state) {
			http.Error(w, "This action is no longer available. Return to the list and refresh it.", http.StatusConflict)
			return
		}
	}
	title := strings.TrimSpace(r.PostForm.Get("title"))
	if (action == "" || action == "rename") && title == "" {
		http.Error(w, "Enter a task name.", http.StatusBadRequest)
		return
	}
	waitingFor := strings.TrimSpace(r.PostForm.Get("waiting_for"))
	followUp := r.PostForm.Get("follow_up")
	if action == "wait" || action == "update-follow-up" {
		if waitingFor == "" || !validFollowUpDate(followUp) {
			http.Error(w, "Enter the expected reply and a valid follow-up date.", http.StatusBadRequest)
			return
		}
	}
	// Leaving Waiting discards the old reply details in the same commit.
	if state == "waiting" && (action == "complete" || action == "resume" || action == "delete") {
		if _, err := tx.ExecContext(ctx, "DELETE FROM pending_replies WHERE task_id = ?", id); err != nil {
			storageError(w, err)
			return
		}
	}
	switch action {
	case "":
		_, err = tx.ExecContext(ctx, "INSERT INTO tasks (title) VALUES (?)", title)
	case "rename":
		_, err = tx.ExecContext(ctx, "UPDATE tasks SET title = ? WHERE id = ?", title, id)
	case "complete", "reopen":
		_, err = tx.ExecContext(ctx, "UPDATE tasks SET completed = ? WHERE id = ?", action == "complete", id)
	case "wait", "update-follow-up":
		_, err = tx.ExecContext(ctx, `INSERT INTO pending_replies (task_id, waiting_for, follow_up)
            VALUES (?, ?, ?) ON CONFLICT(task_id) DO UPDATE
            SET waiting_for = excluded.waiting_for, follow_up = excluded.follow_up`, id, waitingFor, followUp)
	case "resume":
		// Removing the pending reply returns the unfinished task to Open.
	case "delete":
		_, err = tx.ExecContext(ctx, "DELETE FROM tasks WHERE id = ?", id)
	}
	if err != nil {
		storageError(w, err)
		return
	}
	if err := tx.Commit(); err != nil {
		storageError(w, err)
		return
	}
	dest := "/"
	if state != "open" {
		dest += state
	}
	http.Redirect(w, r, dest, http.StatusSeeOther)
}
