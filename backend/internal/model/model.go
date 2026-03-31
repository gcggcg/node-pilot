package model

import "time"

type Server struct {
	ID                int64     `json:"id"`
	Name              string    `json:"name"`
	Host              string    `json:"host"`
	Port              int       `json:"port"`
	Username          string    `json:"username"`
	PasswordEncrypted string    `json:"-"`
	ConnectionStatus  string    `json:"connection_status"` // online|offline|unknown
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type Script struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	TargetPath  string    `json:"target_path"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Task struct {
	ID         int64      `json:"id"`
	ScriptID   int64      `json:"script_id"`
	Name       string     `json:"name"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}

type TaskServer struct {
	ID         int64      `json:"id"`
	TaskID     int64      `json:"task_id"`
	ServerID   int64      `json:"server_id"`
	ServerName string     `json:"server_name,omitempty"`
	Status     string     `json:"status"`
	Output     string     `json:"output"`
	Error      string     `json:"error"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}

type WSMessage struct {
	Type       string    `json:"type"`
	TaskID     int64     `json:"task_id,omitempty"`
	ServerID   int64     `json:"server_id,omitempty"`
	ServerName string    `json:"server_name,omitempty"`
	Content    string    `json:"content,omitempty"`
	Status     string    `json:"status,omitempty"`
	ExitCode   int       `json:"exit_code,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
	Total      int       `json:"total,omitempty"`
	Success    int       `json:"success,omitempty"`
	Failed     int       `json:"failed,omitempty"`
}
