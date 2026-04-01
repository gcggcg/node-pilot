package handler

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"node-pilot/internal/logger"
	"node-pilot/internal/model"
	"node-pilot/internal/repository"
	"node-pilot/internal/service"
	ws "node-pilot/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	repo    *repository.Repository
	sshPool *service.SSHPool
	wsHub   *ws.Hub
	taskSvc *service.TaskExecutor
	encKey  []byte
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewHandler(repo *repository.Repository, sshPool *service.SSHPool, wsHub *ws.Hub, taskSvc *service.TaskExecutor) *Handler {
	key := []byte("12345678901234567890123456789012") // exactly 32 bytes for AES-256
	return &Handler{
		repo:    repo,
		sshPool: sshPool,
		wsHub:   wsHub,
		taskSvc: taskSvc,
		encKey:  key,
	}
}

func (h *Handler) ListServers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	servers, total, err := h.repo.ListServersWithPagination(page, pageSize)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if servers == nil {
		servers = []*model.Server{}
	}
	c.JSON(200, gin.H{
		"data":     servers,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *Handler) GetServer(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	server, err := h.repo.GetServer(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "server not found"})
		return
	}
	c.JSON(200, server)
}

type CreateServerInput struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) CreateServer(c *gin.Context) {
	var input CreateServerInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	encrypted, err := h.encrypt(input.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "encryption failed"})
		return
	}

	server := &model.Server{
		Name:              input.Name,
		Host:              input.Host,
		Port:              input.Port,
		Username:          input.Username,
		PasswordEncrypted: encrypted,
	}

	id, err := h.repo.CreateServer(server)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"id": id})
}

func (h *Handler) UpdateServer(c *gin.Context) {
	var input CreateServerInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	server, err := h.repo.GetServer(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "server not found"})
		return
	}

	encrypted, err := h.encrypt(input.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "encryption failed"})
		return
	}

	server.Name = input.Name
	server.Host = input.Host
	server.Port = input.Port
	server.Username = input.Username
	server.PasswordEncrypted = encrypted

	if err := h.repo.UpdateServer(server); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "updated"})
}

func (h *Handler) DeleteServer(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.repo.DeleteServer(id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "deleted"})
}

func (h *Handler) DeleteServers(c *gin.Context) {
	var input struct {
		IDs []int64 `json:"ids"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.DeleteServers(input.IDs); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "deleted"})
}

func (h *Handler) TestServerConnection(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	server, err := h.repo.GetServer(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "server not found"})
		return
	}

	password, err := h.decrypt(server.PasswordEncrypted)
	if err != nil {
		c.JSON(500, gin.H{"error": "decryption failed"})
		return
	}

	if err := h.taskSvc.TestConnection(server, password); err != nil {
		h.repo.UpdateServerConnectionStatus(id, "offline")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	h.repo.UpdateServerConnectionStatus(id, "online")
	c.JSON(200, gin.H{"message": "connection successful"})
}

func (h *Handler) ListScripts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	scripts, total, err := h.repo.ListScriptsWithPagination(page, pageSize)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if scripts == nil {
		scripts = []*model.Script{}
	}
	c.JSON(200, gin.H{
		"data":     scripts,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *Handler) GetScript(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	script, err := h.repo.GetScript(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "script not found"})
		return
	}
	c.JSON(200, script)
}

type CreateScriptInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Content     string `json:"content"`
	TargetPath  string `json:"target_path"`
}

func (h *Handler) CreateScript(c *gin.Context) {
	var input CreateScriptInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	script := &model.Script{
		Name:        input.Name,
		Description: input.Description,
		Content:     input.Content,
		TargetPath:  input.TargetPath,
	}

	id, err := h.repo.CreateScript(script)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"id": id})
}

func (h *Handler) UpdateScript(c *gin.Context) {
	var input CreateScriptInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	script, err := h.repo.GetScript(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "script not found"})
		return
	}

	script.Name = input.Name
	script.Description = input.Description
	script.Content = input.Content
	script.TargetPath = input.TargetPath

	if err := h.repo.UpdateScript(script); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "updated"})
}

func (h *Handler) DeleteScript(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.repo.DeleteScript(id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "deleted"})
}

func (h *Handler) DeleteScripts(c *gin.Context) {
	var input struct {
		IDs []int64 `json:"ids"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.DeleteScripts(input.IDs); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "deleted"})
}

func (h *Handler) ListTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	tasks, total, err := h.repo.ListTasksWithPagination(page, pageSize)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if tasks == nil {
		tasks = []*model.Task{}
	}
	c.JSON(200, gin.H{
		"data":     tasks,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *Handler) GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	task, err := h.repo.GetTask(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "task not found"})
		return
	}
	taskServers, _ := h.repo.GetTaskServers(task.ID)
	c.JSON(200, gin.H{
		"task":    task,
		"servers": taskServers,
	})
}

type CreateTaskInput struct {
	ScriptID  int64   `json:"script_id"`
	Name      string  `json:"name"`
	ServerIDs []int64 `json:"server_ids"`
}

func (h *Handler) CreateTask(c *gin.Context) {
	var input CreateTaskInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	task := &model.Task{
		ScriptID: input.ScriptID,
		Name:     input.Name,
		Status:   "pending",
	}

	id, err := h.repo.CreateTask(task)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	task.ID = id

	var servers []*model.Server
	for _, sid := range input.ServerIDs {
		server, err := h.repo.GetServer(sid)
		if err == nil {
			servers = append(servers, server)
		}
	}

	go func() {
		if len(servers) > 0 {
			srv := servers[0]
			password, _ := h.decrypt(srv.PasswordEncrypted)
			script, _ := h.repo.GetScript(input.ScriptID)
			h.taskSvc.ExecuteScript(task, script, servers, password)
		}
	}()

	c.JSON(201, gin.H{"id": id})
}

func (h *Handler) CancelTask(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	// 调用 taskSvc.CancelTask 标记任务为取消状态（内存中跟踪）
	h.taskSvc.CancelTask(id)
	// 更新数据库状态为 cancelled
	if err := h.repo.CancelTask(id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "cancelled"})
}

func (h *Handler) DeleteTasks(c *gin.Context) {
	var input struct {
		IDs []int64 `json:"ids"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.DeleteTasks(input.IDs); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "deleted"})
}

func (h *Handler) GetTaskOutput(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, _ := strconv.ParseInt(taskIDStr, 10, 64)

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	clientGone := c.Request.Context().Done()
	for {
		select {
		case <-clientGone:
			return
		case <-ticker.C:
			task, err := h.repo.GetTask(taskID)
			if err != nil {
				return
			}
			taskServers, _ := h.repo.GetTaskServers(task.ID)
			data, _ := json.Marshal(gin.H{
				"task":    task,
				"servers": taskServers,
			})
			c.SSEvent("message", string(data))
			c.Writer.Flush()
			if task.Status == "completed" || task.Status == "cancelled" {
				return
			}
		}
	}
}

func (h *Handler) WebSocketHandler(c *gin.Context) {
	taskIDStr := c.Query("task_id")
	taskID, _ := strconv.ParseUint(taskIDStr, 10, 64)
	logger.Debug("[WS-HANDLER] New WebSocket connection, task_id=%d, remote=%s", taskID, c.ClientIP())

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Debug("[WS-HANDLER] Failed to upgrade connection: %v", err)
		return
	}

	client := &ws.Client{
		Hub:    h.wsHub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		TaskID: taskID,
	}

	logger.Debug("[WS-HANDLER] Registering client to hub...")
	h.wsHub.RegisterClient(client)
	logger.Debug("[WS-HANDLER] Client registered (taskID=%d), starting pumps", taskID)

	go client.WritePump()
	go client.ReadPump(taskIDStr)
}

func (h *Handler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "no file uploaded"})
		return
	}
	defer file.Close()

	localPath := filepath.Join("/tmp", header.Filename)
	out, err := os.Create(localPath)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to save file"})
		return
	}
	defer out.Close()

	io.Copy(out, file)
	c.JSON(200, gin.H{"path": localPath})
}

type DeployInput struct {
	ServerIDs  []int64 `json:"server_ids"`
	LocalPath  string  `json:"local_path"`
	RemotePath string  `json:"remote_path"`
}

func (h *Handler) DeployFile(c *gin.Context) {
	var input DeployInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	for _, sid := range input.ServerIDs {
		server, err := h.repo.GetServer(sid)
		if err != nil {
			continue
		}
		password, _ := h.decrypt(server.PasswordEncrypted)
		go h.taskSvc.DeployFile(server, password, input.LocalPath, input.RemotePath)
	}

	c.JSON(200, gin.H{"message": "deployment started"})
}

func (h *Handler) encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(h.encKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (h *Handler) decrypt(encrypted string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(h.encKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
