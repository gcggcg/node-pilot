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
	ScriptID   int64      `json:"script_id"`  // Deprecated: kept for backward compatibility, use ScriptIDs instead
	ScriptIDs  string     `json:"script_ids"` // New field: comma-separated script IDs for batch execution (e.g., "1,2,3")
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
	// Batch script execution fields
	ScriptIndex  int    `json:"script_index,omitempty"`  // Current script index (1-based)
	TotalScripts int    `json:"total_scripts,omitempty"` // Total number of scripts
	ScriptPath   string `json:"script_path,omitempty"`   // Script file path
	ScriptName   string `json:"script_name,omitempty"`   // Script name
}

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // bcrypt hash, never expose to client
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Role         string    `json:"role"` // ROLE_ADMIN or ROLE_USER
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type FileUpload struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`        // 配置名称
	LocalPath  string    `json:"local_path"`  // ./data/files/ 下的相对路径
	RemotePath string    `json:"remote_path"` // 远程目标路径，必须以/开头
	Status     string    `json:"status"`      // pending|success|failed
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type FileUploadServer struct {
	ID             int64     `json:"id"`
	FileUploadID   int64     `json:"file_upload_id"`
	ServerID       int64     `json:"server_id"`
	ServerName     string    `json:"server_name,omitempty"`
	Status         string    `json:"status"` // pending|success|failed
	ErrorMessage   string    `json:"error_message,omitempty"`
	FileName       string    `json:"file_name"`        // 文件名
	RemoteFullPath string    `json:"remote_full_path"` // 完整远程路径
	CreatedAt      time.Time `json:"created_at"`
}
