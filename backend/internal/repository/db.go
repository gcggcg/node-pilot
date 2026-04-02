package repository

import (
	"database/sql"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"node-pilot/internal/model"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db *sql.DB
}

func NewDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := initSchema(db); err != nil {
		return nil, err
	}
	if err := migrateSchema(db); err != nil {
		return nil, err
	}

	// Initialize root user
	repo := &Repository{db: db}
	if err := repo.InitRootUser(); err != nil {
		return nil, err
	}

	return db, nil
}

func initSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS servers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		host TEXT NOT NULL,
		port INTEGER DEFAULT 22,
		username TEXT NOT NULL,
		password_encrypted TEXT,
		connection_status TEXT DEFAULT 'unknown',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS scripts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		content TEXT NOT NULL,
		target_path TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		script_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		status TEXT DEFAULT 'pending',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		started_at DATETIME,
		finished_at DATETIME,
		FOREIGN KEY (script_id) REFERENCES scripts(id)
	);

	CREATE TABLE IF NOT EXISTS task_servers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id INTEGER NOT NULL,
		server_id INTEGER NOT NULL,
		status TEXT DEFAULT 'pending',
		output TEXT DEFAULT '',
		error TEXT DEFAULT '',
		started_at DATETIME,
		finished_at DATETIME,
		FOREIGN KEY (task_id) REFERENCES tasks(id),
		FOREIGN KEY (server_id) REFERENCES servers(id)
	);

	CREATE INDEX IF NOT EXISTS idx_task_servers_task ON task_servers(task_id);
	CREATE INDEX IF NOT EXISTS idx_task_servers_server ON task_servers(server_id);
	CREATE INDEX IF NOT EXISTS idx_task_servers_status ON task_servers(status);

	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		email TEXT DEFAULT '',
		phone TEXT DEFAULT '',
		role TEXT DEFAULT 'ROLE_USER',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS file_uploads (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		local_path TEXT NOT NULL,
		remote_path TEXT NOT NULL,
		status TEXT DEFAULT 'pending',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS file_upload_servers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_upload_id INTEGER NOT NULL,
		server_id INTEGER NOT NULL,
		server_name TEXT DEFAULT '',
		status TEXT DEFAULT 'pending',
		error_message TEXT DEFAULT '',
		file_name TEXT NOT NULL,
		remote_full_path TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (file_upload_id) REFERENCES file_uploads(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_file_upload_servers_upload ON file_upload_servers(file_upload_id);
	CREATE INDEX IF NOT EXISTS idx_file_upload_servers_server ON file_upload_servers(server_id);
	`
	_, err := db.Exec(schema)
	return err
}

func migrateSchema(db *sql.DB) error {
	// Add connection_status column if not exists (for existing databases)
	_, err := db.Exec(`ALTER TABLE servers ADD COLUMN connection_status TEXT DEFAULT 'unknown'`)
	if err != nil && err.Error() != "duplicate column name: connection_status" {
		return err
	}
	// Set default status for existing servers that have NULL
	_, err = db.Exec(`UPDATE servers SET connection_status = 'unknown' WHERE connection_status IS NULL`)
	if err != nil {
		return err
	}

	// Add users table migration (for existing databases without users table)
	_, err = db.Exec(`ALTER TABLE users ADD COLUMN email TEXT DEFAULT ''`)
	if err != nil && err.Error() != "duplicate column name: email" {
		// Table might not exist, ignore error
	}
	_, err = db.Exec(`ALTER TABLE users ADD COLUMN phone TEXT DEFAULT ''`)
	if err != nil && err.Error() != "duplicate column name: phone" {
		// Table might not exist, ignore error
	}
	_, err = db.Exec(`ALTER TABLE users ADD COLUMN role TEXT DEFAULT 'ROLE_USER'`)
	if err != nil && err.Error() != "duplicate column name: role" {
		// Table might not exist, ignore error
	}

	return nil
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListServers() ([]*model.Server, error) {
	rows, err := r.db.Query(`SELECT id, name, host, port, username, password_encrypted, connection_status, created_at, updated_at FROM servers ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []*model.Server
	for rows.Next() {
		s := &model.Server{}
		if err := rows.Scan(&s.ID, &s.Name, &s.Host, &s.Port, &s.Username, &s.PasswordEncrypted, &s.ConnectionStatus, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		servers = append(servers, s)
	}
	return servers, nil
}

func (r *Repository) ListServersWithPagination(page, pageSize int) ([]*model.Server, int64, error) {
	offset := (page - 1) * pageSize

	var total int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM servers").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(
		"SELECT id, name, host, port, username, password_encrypted, connection_status, created_at, updated_at FROM servers ORDER BY id DESC LIMIT ? OFFSET ?",
		pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var servers []*model.Server
	for rows.Next() {
		s := &model.Server{}
		if err := rows.Scan(&s.ID, &s.Name, &s.Host, &s.Port, &s.Username, &s.PasswordEncrypted, &s.ConnectionStatus, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, err
		}
		servers = append(servers, s)
	}
	return servers, total, nil
}

func (r *Repository) GetServer(id int64) (*model.Server, error) {
	s := &model.Server{}
	err := r.db.QueryRow(`SELECT id, name, host, port, username, password_encrypted, connection_status, created_at, updated_at FROM servers WHERE id = ?`, id).
		Scan(&s.ID, &s.Name, &s.Host, &s.Port, &s.Username, &s.PasswordEncrypted, &s.ConnectionStatus, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *Repository) CreateServer(s *model.Server) (int64, error) {
	result, err := r.db.Exec(`INSERT INTO servers (name, host, port, username, password_encrypted, connection_status) VALUES (?, ?, ?, ?, ?, 'unknown')`,
		s.Name, s.Host, s.Port, s.Username, s.PasswordEncrypted)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *Repository) UpdateServer(s *model.Server) error {
	_, err := r.db.Exec(`UPDATE servers SET name = ?, host = ?, port = ?, username = ?, password_encrypted = ?, connection_status = ?, updated_at = ? WHERE id = ?`,
		s.Name, s.Host, s.Port, s.Username, s.PasswordEncrypted, s.ConnectionStatus, time.Now(), s.ID)
	return err
}

func (r *Repository) UpdateServerConnectionStatus(id int64, status string) error {
	_, err := r.db.Exec(`UPDATE servers SET connection_status = ?, updated_at = ? WHERE id = ?`,
		status, time.Now(), id)
	return err
}

func (r *Repository) DeleteServer(id int64) error {
	_, err := r.db.Exec(`DELETE FROM servers WHERE id = ?`, id)
	return err
}

func (r *Repository) DeleteServers(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	query := `DELETE FROM servers WHERE id IN (` + strings.Join(placeholders, ",") + `)`
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *Repository) ListScripts() ([]*model.Script, error) {
	rows, err := r.db.Query(`SELECT id, name, description, content, target_path, created_at, updated_at FROM scripts ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scripts []*model.Script
	for rows.Next() {
		s := &model.Script{}
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.Content, &s.TargetPath, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		scripts = append(scripts, s)
	}
	return scripts, nil
}

func (r *Repository) ListScriptsWithPagination(page, pageSize int) ([]*model.Script, int64, error) {
	offset := (page - 1) * pageSize

	var total int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM scripts").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(
		"SELECT id, name, description, content, target_path, created_at, updated_at FROM scripts ORDER BY id DESC LIMIT ? OFFSET ?",
		pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var scripts []*model.Script
	for rows.Next() {
		s := &model.Script{}
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.Content, &s.TargetPath, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, err
		}
		scripts = append(scripts, s)
	}
	return scripts, total, nil
}

func (r *Repository) GetScript(id int64) (*model.Script, error) {
	s := &model.Script{}
	err := r.db.QueryRow(`SELECT id, name, description, content, target_path, created_at, updated_at FROM scripts WHERE id = ?`, id).
		Scan(&s.ID, &s.Name, &s.Description, &s.Content, &s.TargetPath, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *Repository) CreateScript(s *model.Script) (int64, error) {
	result, err := r.db.Exec(`INSERT INTO scripts (name, description, content, target_path) VALUES (?, ?, ?, ?)`,
		s.Name, s.Description, s.Content, s.TargetPath)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *Repository) UpdateScript(s *model.Script) error {
	_, err := r.db.Exec(`UPDATE scripts SET name = ?, description = ?, content = ?, target_path = ?, updated_at = ? WHERE id = ?`,
		s.Name, s.Description, s.Content, s.TargetPath, time.Now(), s.ID)
	return err
}

func (r *Repository) DeleteScript(id int64) error {
	_, err := r.db.Exec(`DELETE FROM scripts WHERE id = ?`, id)
	return err
}

func (r *Repository) DeleteScripts(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	query := `DELETE FROM scripts WHERE id IN (` + strings.Join(placeholders, ",") + `)`
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *Repository) ListTasks() ([]*model.Task, error) {
	rows, err := r.db.Query(`SELECT id, script_id, name, status, created_at, started_at, finished_at FROM tasks ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*model.Task
	for rows.Next() {
		t := &model.Task{}
		if err := rows.Scan(&t.ID, &t.ScriptID, &t.Name, &t.Status, &t.CreatedAt, &t.StartedAt, &t.FinishedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *Repository) ListTasksWithPagination(page, pageSize int) ([]*model.Task, int64, error) {
	offset := (page - 1) * pageSize

	var total int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(
		"SELECT id, script_id, name, status, created_at, started_at, finished_at FROM tasks ORDER BY id DESC LIMIT ? OFFSET ?",
		pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []*model.Task
	for rows.Next() {
		t := &model.Task{}
		if err := rows.Scan(&t.ID, &t.ScriptID, &t.Name, &t.Status, &t.CreatedAt, &t.StartedAt, &t.FinishedAt); err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, t)
	}
	return tasks, total, nil
}

func (r *Repository) GetTask(id int64) (*model.Task, error) {
	t := &model.Task{}
	err := r.db.QueryRow(`SELECT id, script_id, name, status, created_at, started_at, finished_at FROM tasks WHERE id = ?`, id).
		Scan(&t.ID, &t.ScriptID, &t.Name, &t.Status, &t.CreatedAt, &t.StartedAt, &t.FinishedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *Repository) CreateTask(t *model.Task) (int64, error) {
	result, err := r.db.Exec(`INSERT INTO tasks (script_id, name, status) VALUES (?, ?, ?)`,
		t.ScriptID, t.Name, t.Status)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *Repository) UpdateTaskStatus(id int64, status string, started, finished *time.Time) error {
	_, err := r.db.Exec(`UPDATE tasks SET status = ?, started_at = ?, finished_at = ? WHERE id = ?`,
		status, started, finished, id)
	return err
}

func (r *Repository) CancelTask(id int64) error {
	_, err := r.db.Exec(`UPDATE tasks SET status = 'cancelled', finished_at = ? WHERE id = ?`, time.Now(), id)
	return err
}

func (r *Repository) CreateTaskServer(ts *model.TaskServer) (int64, error) {
	result, err := r.db.Exec(`INSERT INTO task_servers (task_id, server_id, status) VALUES (?, ?, ?)`,
		ts.TaskID, ts.ServerID, ts.Status)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *Repository) GetTaskServers(taskID int64) ([]*model.TaskServer, error) {
	rows, err := r.db.Query(`
		SELECT ts.id, ts.task_id, ts.server_id, s.name, ts.status, ts.output, ts.error, ts.started_at, ts.finished_at 
		FROM task_servers ts 
		JOIN servers s ON ts.server_id = s.id 
		WHERE ts.task_id = ?`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.TaskServer
	for rows.Next() {
		ts := &model.TaskServer{}
		if err := rows.Scan(&ts.ID, &ts.TaskID, &ts.ServerID, &ts.ServerName, &ts.Status, &ts.Output, &ts.Error, &ts.StartedAt, &ts.FinishedAt); err != nil {
			return nil, err
		}
		results = append(results, ts)
	}
	return results, nil
}

func (r *Repository) UpdateTaskServerStatus(id int64, status, output, errMsg string, started, finished *time.Time) error {
	_, err := r.db.Exec(`UPDATE task_servers SET status = ?, output = ?, error = ?, started_at = ?, finished_at = ? WHERE id = ?`,
		status, output, errMsg, started, finished, id)
	return err
}

func (r *Repository) DeleteTasks(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	query := `DELETE FROM tasks WHERE id IN (` + strings.Join(placeholders, ",") + `)`

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 先删除关联的 task_servers
	_, err = tx.Exec(`DELETE FROM task_servers WHERE task_id IN (`+strings.Join(placeholders, ",")+`)`, args...)
	if err != nil {
		return err
	}

	// 再删除 tasks
	_, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// User Repository Methods

func (r *Repository) CreateUser(user *model.User) (int64, error) {
	result, err := r.db.Exec(
		`INSERT INTO users (username, password_hash, email, phone, role) VALUES (?, ?, ?, ?, ?)`,
		user.Username, user.PasswordHash, user.Email, user.Phone, user.Role)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *Repository) GetUserByUsername(username string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(
		`SELECT id, username, password_hash, email, phone, role, created_at, updated_at FROM users WHERE username = ?`,
		username).
		Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.Phone, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) GetUserByID(id int64) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(
		`SELECT id, username, password_hash, email, phone, role, created_at, updated_at FROM users WHERE id = ?`,
		id).
		Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.Phone, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) UpdateUser(user *model.User) error {
	_, err := r.db.Exec(
		`UPDATE users SET username = ?, password_hash = ?, email = ?, phone = ?, role = ?, updated_at = ? WHERE id = ?`,
		user.Username, user.PasswordHash, user.Email, user.Phone, user.Role, time.Now(), user.ID)
	return err
}

func (r *Repository) ListUsers(page, pageSize int, keyword string) ([]*model.User, int64, error) {
	offset := (page - 1) * pageSize

	var total int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, username, password_hash, email, phone, role, created_at, updated_at FROM users`
	var args []interface{}

	if keyword != "" {
		query += ` WHERE username LIKE ?`
		keyword = "%" + keyword + "%"
		args = append(args, keyword)
	}

	query += ` ORDER BY id DESC LIMIT ? OFFSET ?`
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.Phone, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}
	return users, total, nil
}

func (r *Repository) DeleteUsers(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	query := `DELETE FROM users WHERE id IN (` + strings.Join(placeholders, ",") + `)`
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *Repository) InitRootUser() error {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'root'").Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		// password: root, use bcrypt hash
		hash, err := bcrypt.GenerateFromPassword([]byte("root"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		_, err = r.db.Exec("INSERT INTO users (username, password_hash, role) VALUES ('root', ?, 'ROLE_ADMIN')", string(hash))
		return err
	}
	return nil
}
