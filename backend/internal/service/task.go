package service

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"node-pilot/internal/logger"
	"node-pilot/internal/model"
	"node-pilot/internal/repository"
	"node-pilot/internal/websocket"

	"github.com/pkg/sftp"
)

type TaskExecutor struct {
	repo           *repository.Repository
	sshPool        *SSHPool
	wsHub          *websocket.Hub
	debug          bool
	cancelledTasks sync.Map
}

func NewTaskExecutor(repo *repository.Repository, sshPool *SSHPool, wsHub *websocket.Hub, debug bool) *TaskExecutor {
	return &TaskExecutor{
		repo:    repo,
		sshPool: sshPool,
		wsHub:   wsHub,
		debug:   debug,
	}
}

func (e *TaskExecutor) CancelTask(taskID int64) {
	e.cancelledTasks.Store(taskID, true)
}

func (e *TaskExecutor) IsTaskCancelled(taskID int64) bool {
	_, cancelled := e.cancelledTasks.Load(taskID)
	return cancelled
}

// ExecuteTask 手动执行指定任务
func (e *TaskExecutor) ExecuteTask(taskID int64) error {
	task, err := e.repo.GetTask(taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// 只能执行 pending 状态的任务
	if task.Status != "pending" {
		return fmt.Errorf("任务不是pending状态，无法执行")
	}

	taskServers, err := e.repo.GetTaskServers(task.ID)
	if err != nil {
		return fmt.Errorf("failed to get task servers: %w", err)
	}

	if len(taskServers) == 0 {
		return fmt.Errorf("任务没有关联服务器")
	}

	// 获取服务器信息用于执行
	var servers []*model.Server
	for _, ts := range taskServers {
		server, err := e.repo.GetServer(ts.ServerID)
		if err == nil {
			servers = append(servers, server)
		}
	}

	if len(servers) == 0 {
		return fmt.Errorf("没有可用的服务器")
	}

	// 获取第一个服务器的密码（假设所有服务器密码相同）
	password, err := e.decryptPassword(servers[0].PasswordEncrypted)
	if err != nil {
		return fmt.Errorf("failed to decrypt password: %w", err)
	}

	script, err := e.repo.GetScript(task.ScriptID)
	if err != nil {
		return fmt.Errorf("failed to get script: %w", err)
	}

	// 构建 serverID -> taskServerID 映射，用于更新已有记录
	serverTaskServerMap := make(map[int64]int64)
	for _, ts := range taskServers {
		serverTaskServerMap[ts.ServerID] = ts.ID
	}

	// 在 goroutine 中执行
	go e.ExecuteScript(task, script, servers, password, serverTaskServerMap)

	return nil
}

// decryptPassword 解密服务器密码
func (e *TaskExecutor) decryptPassword(encrypted string) (string, error) {
	if encrypted == "" {
		return "", nil
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	key := []byte("12345678901234567890123456789012") // 32 bytes for AES-256
	block, err := aes.NewCipher(key)
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

func (e *TaskExecutor) TestConnection(server *model.Server, password string) error {
	client, err := e.sshPool.GetClient(server.ID, server.Host, server.Port, server.Username, password)
	if err != nil {
		return fmt.Errorf("SSH connection failed: %w", err)
	}
	defer client.Close()
	return nil
}

func (e *TaskExecutor) ExecuteScript(task *model.Task, script *model.Script, servers []*model.Server, password string, serverTaskServerMap map[int64]int64) {
	if e.debug {
		logger.Debug("[TASK-%d] ========== 任务开始执行 ==========", task.ID)
		logger.Debug("[TASK-%d] 脚本名称: %s", task.ID, script.Name)
		logger.Debug("[TASK-%d] 脚本内容:\n%s", task.ID, script.Content)
		logger.Debug("[TASK-%d] 目标路径: %s", task.ID, script.TargetPath)
		logger.Debug("[TASK-%d] 目标服务器数量: %d", task.ID, len(servers))
		for i, s := range servers {
			logger.Debug("[TASK-%d]   服务器[%d]: %s (%s:%d)", task.ID, i+1, s.Name, s.Host, s.Port)
		}
	}

	now := time.Now()
	task.Status = "running"
	task.StartedAt = &now
	e.repo.UpdateTaskStatus(task.ID, "running", &now, nil)

	e.wsHub.BroadcastToTask(&model.WSMessage{
		Type:      "task_start",
		TaskID:    task.ID,
		Timestamp: time.Now(),
	}, uint64(task.ID))

	batchSize := 10
	success := 0
	failed := 0
	wasCancelled := false

	for i := 0; i < len(servers); i += batchSize {
		// Check if task was cancelled
		if e.IsTaskCancelled(task.ID) {
			wasCancelled = true
			if e.debug {
				logger.Debug("[TASK-%d] 任务被取消，停止执行", task.ID)
			}
			break
		}

		end := i + batchSize
		if end > len(servers) {
			end = len(servers)
		}
		batch := servers[i:end]
		if e.debug {
			logger.Debug("[TASK-%d] 执行批次 %d-%d (共%d台服务器)", task.ID, i+1, end, len(servers))
		}
		e.executeBatch(task, script, batch, password, i, len(servers), serverTaskServerMap)
		time.Sleep(500 * time.Millisecond)
	}

	// Clean up cancelled task tracking
	if wasCancelled {
		e.cancelledTasks.Delete(task.ID)
	}

	taskServers, _ := e.repo.GetTaskServers(task.ID)
	for _, ts := range taskServers {
		if ts.Status == "success" {
			success++
		} else {
			failed++
		}
	}

	finished := time.Now()

	// If was cancelled, preserve cancelled status
	if wasCancelled {
		task.Status = "cancelled"
		e.repo.UpdateTaskStatus(task.ID, "cancelled", task.StartedAt, &finished)
	} else if failed > 0 {
		task.Status = "failed"
		e.repo.UpdateTaskStatus(task.ID, "failed", task.StartedAt, &finished)
	} else {
		task.Status = "completed"
		e.repo.UpdateTaskStatus(task.ID, "completed", task.StartedAt, &finished)
	}
	task.FinishedAt = &finished

	if e.debug {
		logger.Debug("[TASK-%d] ========== 任务执行完成 ==========", task.ID)
		logger.Debug("[TASK-%d] 总服务器数: %d, 成功: %d, 失败: %d, 取消: %v", task.ID, len(servers), success, failed, wasCancelled)
	}

	e.wsHub.BroadcastToTask(&model.WSMessage{
		Type:      "task_done",
		TaskID:    task.ID,
		Status:    task.Status,
		Total:     len(servers),
		Success:   success,
		Failed:    failed,
		Timestamp: time.Now(),
	}, uint64(task.ID))
}

func (e *TaskExecutor) executeBatch(task *model.Task, script *model.Script, servers []*model.Server, password string, offset, total int, serverTaskServerMap map[int64]int64) {
	var wg sync.WaitGroup
	for i, srv := range servers {
		wg.Add(1)
		go func(srv *model.Server, idx int) {
			defer wg.Done()
			e.executeOnServer(task, script, srv, password, offset+idx, total, serverTaskServerMap)
		}(srv, i)
	}
	wg.Wait()
}

func (e *TaskExecutor) executeOnServer(task *model.Task, script *model.Script, srv *model.Server, password string, idx, total int, serverTaskServerMap map[int64]int64) {
	started := time.Now()

	if e.debug {
		logger.Debug("[TASK-%d][SERVER-%d] ========== 开始在服务器 %s 上执行 ==========", task.ID, srv.ID, srv.Name)
		logger.Debug("[TASK-%d][SERVER-%d] 主机: %s:%d", task.ID, srv.ID, srv.Host, srv.Port)
		logger.Debug("[TASK-%d][SERVER-%d] 用户: %s", task.ID, srv.ID, srv.Username)
		logger.Debug("[TASK-%d][SERVER-%d] 目标路径: %s", task.ID, srv.ID, script.TargetPath)
	}

	// 查找已有的 task_servers 记录ID，如果不存在则创建新记录
	var tsID int64
	if existingID, ok := serverTaskServerMap[srv.ID]; ok {
		tsID = existingID
		// 更新已有记录的状态为 running
		e.repo.UpdateTaskServerByIDs(task.ID, srv.ID, "running", "", "", &started, nil)
	} else {
		// 如果没有找到已有记录（理论上不应该发生），创建新记录
		ts := &model.TaskServer{
			TaskID:   task.ID,
			ServerID: srv.ID,
			Status:   "running",
		}
		tsID, _ = e.repo.CreateTaskServer(ts)
	}

	e.wsHub.BroadcastToTask(&model.WSMessage{
		Type:       "server_start",
		TaskID:     task.ID,
		ServerID:   srv.ID,
		ServerName: srv.Name,
		Timestamp:  time.Now(),
	}, uint64(task.ID))

	if e.debug {
		logger.Debug("[TASK-%d][SERVER-%d] Step 1: 建立SSH连接...", task.ID, srv.ID)
	}
	client, err := e.sshPool.GetClient(srv.ID, srv.Host, srv.Port, srv.Username, password)
	if err != nil {
		if e.debug {
			logger.Error("[TASK-%d][SERVER-%d] SSH连接失败: %v", task.ID, srv.ID, err)
		}
		finished := time.Now()
		e.repo.UpdateTaskServerStatus(tsID, "failed", "", err.Error(), &started, &finished)
		e.wsHub.BroadcastToTask(&model.WSMessage{
			Type:       "server_done",
			TaskID:     task.ID,
			ServerID:   srv.ID,
			ServerName: srv.Name,
			Status:     "failed",
			Content:    err.Error(),
			Timestamp:  time.Now(),
		}, uint64(task.ID))
		return
	}
	if e.debug {
		logger.Debug("[TASK-%d][SERVER-%d] SSH连接成功!", task.ID, srv.ID)
	}
	defer client.Close()

	targetDir := filepath.Dir(script.TargetPath)
	targetFile := script.TargetPath

	// Step 1: Create target directory if not exists
	if e.debug {
		logger.Debug("[TASK-%d][SERVER-%d] Step 2: 创建目录 %s", task.ID, srv.ID, targetDir)
	}
	mkdirCmd := fmt.Sprintf("mkdir -p %s && echo OK", targetDir)
	session, err := client.NewSession()
	if err != nil {
		if e.debug {
			logger.Error("[TASK-%d][SERVER-%d] 创建SSH会话失败: %v", task.ID, srv.ID, err)
		}
		finished := time.Now()
		e.repo.UpdateTaskServerStatus(tsID, "failed", "", err.Error(), &started, &finished)
		e.wsHub.BroadcastToTask(&model.WSMessage{
			Type:       "server_done",
			TaskID:     task.ID,
			ServerID:   srv.ID,
			ServerName: srv.Name,
			Status:     "failed",
			Content:    err.Error(),
			Timestamp:  time.Now(),
		}, uint64(task.ID))
		return
	}
	out, err := session.CombinedOutput(mkdirCmd)
	session.Close()
	if err != nil {
		errMsg := fmt.Sprintf("mkdir failed: %s, output: %s", err.Error(), string(out))
		if e.debug {
			logger.Error("[TASK-%d][SERVER-%d] 创建目录失败: %s", task.ID, srv.ID, errMsg)
		}
		finished := time.Now()
		e.repo.UpdateTaskServerStatus(tsID, "failed", "", errMsg, &started, &finished)
		e.wsHub.BroadcastToTask(&model.WSMessage{
			Type:       "server_done",
			TaskID:     task.ID,
			ServerID:   srv.ID,
			ServerName: srv.Name,
			Status:     "failed",
			Content:    errMsg,
			Timestamp:  time.Now(),
		}, uint64(task.ID))
		return
	}
	if e.debug {
		logger.Debug("[TASK-%d][SERVER-%d] 目录创建成功: %s", task.ID, srv.ID, targetDir)
	}

	// Step 2: Write script content to TargetPath via SFTP
	if e.debug {
		logger.Debug("[TASK-%d][SERVER-%d] Step 3: 通过SFTP写入脚本到 %s", task.ID, srv.ID, targetFile)
		logger.Debug("[TASK-%d][SERVER-%d] 脚本大小: %d bytes", task.ID, srv.ID, len(script.Content))
	}
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		if e.debug {
			logger.Error("[TASK-%d][SERVER-%d] 创建SFTP客户端失败: %v", task.ID, srv.ID, err)
		}
		finished := time.Now()
		e.repo.UpdateTaskServerStatus(tsID, "failed", "", err.Error(), &started, &finished)
		e.wsHub.BroadcastToTask(&model.WSMessage{
			Type:       "server_done",
			TaskID:     task.ID,
			ServerID:   srv.ID,
			ServerName: srv.Name,
			Status:     "failed",
			Content:    err.Error(),
			Timestamp:  time.Now(),
		}, uint64(task.ID))
		return
	}
	f, err := sftpClient.Create(targetFile)
	if err != nil {
		sftpClient.Close()
		errMsg := fmt.Sprintf("sftp create file failed: %s, target: %s", err.Error(), targetFile)
		if e.debug {
			logger.Error("[TASK-%d][SERVER-%d] SFTP创建文件失败: %s", task.ID, srv.ID, errMsg)
		}
		finished := time.Now()
		e.repo.UpdateTaskServerStatus(tsID, "failed", "", errMsg, &started, &finished)
		e.wsHub.BroadcastToTask(&model.WSMessage{
			Type:       "server_done",
			TaskID:     task.ID,
			ServerID:   srv.ID,
			ServerName: srv.Name,
			Status:     "failed",
			Content:    errMsg,
			Timestamp:  time.Now(),
		}, uint64(task.ID))
		return
	}
	_, err = f.Write([]byte(script.Content))
	f.Close()
	sftpClient.Close()
	if err != nil {
		errMsg := fmt.Sprintf("sftp write failed: %s", err.Error())
		if e.debug {
			logger.Error("[TASK-%d][SERVER-%d] SFTP写入文件失败: %s", task.ID, srv.ID, errMsg)
		}
		finished := time.Now()
		e.repo.UpdateTaskServerStatus(tsID, "failed", "", errMsg, &started, &finished)
		e.wsHub.BroadcastToTask(&model.WSMessage{
			Type:       "server_done",
			TaskID:     task.ID,
			ServerID:   srv.ID,
			ServerName: srv.Name,
			Status:     "failed",
			Content:    errMsg,
			Timestamp:  time.Now(),
		}, uint64(task.ID))
		return
	}
	if e.debug {
		logger.Debug("[TASK-%d][SERVER-%d] SFTP写入成功!", task.ID, srv.ID)
	}

	// Step 3: Set executable permission
	if e.debug {
		logger.Debug("[TASK-%d][SERVER-%d] Step 4: 设置执行权限 chmod +x %s", task.ID, srv.ID, targetFile)
	}
	chmodCmd := fmt.Sprintf("chmod +x %s && echo OK", targetFile)
	session, err = client.NewSession()
	if err != nil {
		if e.debug {
			logger.Error("[TASK-%d][SERVER-%d] 创建chmod会话失败: %v", task.ID, srv.ID, err)
		}
		finished := time.Now()
		e.repo.UpdateTaskServerStatus(tsID, "failed", "", err.Error(), &started, &finished)
		e.wsHub.BroadcastToTask(&model.WSMessage{
			Type:       "server_done",
			TaskID:     task.ID,
			ServerID:   srv.ID,
			ServerName: srv.Name,
			Status:     "failed",
			Content:    err.Error(),
			Timestamp:  time.Now(),
		}, uint64(task.ID))
		return
	}
	out, err = session.CombinedOutput(chmodCmd)
	session.Close()
	if err != nil {
		errMsg := fmt.Sprintf("chmod failed: %s, output: %s", err.Error(), string(out))
		if e.debug {
			logger.Error("[TASK-%d][SERVER-%d] chmod失败: %s", task.ID, srv.ID, errMsg)
		}
		finished := time.Now()
		e.repo.UpdateTaskServerStatus(tsID, "failed", "", errMsg, &started, &finished)
		e.wsHub.BroadcastToTask(&model.WSMessage{
			Type:       "server_done",
			TaskID:     task.ID,
			ServerID:   srv.ID,
			ServerName: srv.Name,
			Status:     "failed",
			Content:    errMsg,
			Timestamp:  time.Now(),
		}, uint64(task.ID))
		return
	}
	if e.debug {
		logger.Debug("[TASK-%d][SERVER-%d] chmod成功!", task.ID, srv.ID)
	}

	// Step 4: Execute the script file
	if e.debug {
		logger.Debug("[TASK-%d][SERVER-%d] Step 5: 执行脚本 /bin/bash %s", task.ID, srv.ID, targetFile)
	}
	session, err = client.NewSession()
	if err != nil {
		if e.debug {
			logger.Error("[TASK-%d][SERVER-%d] 创建执行会话失败: %v", task.ID, srv.ID, err)
		}
		finished := time.Now()
		e.repo.UpdateTaskServerStatus(tsID, "failed", "", err.Error(), &started, &finished)
		e.wsHub.BroadcastToTask(&model.WSMessage{
			Type:       "server_done",
			TaskID:     task.ID,
			ServerID:   srv.ID,
			ServerName: srv.Name,
			Status:     "failed",
			Content:    err.Error(),
			Timestamp:  time.Now(),
		}, uint64(task.ID))
		return
	}
	defer session.Close()

	execCmd := fmt.Sprintf("/bin/bash %s", targetFile)

	// Use CombinedOutput to reliably capture stdout and stderr
	cmdOut, err := session.CombinedOutput(execCmd)

	// Apply 100-line limit to output for database storage
	output := limitLines(string(cmdOut), 100)

	logger.Debug("----------------------执行命令: %s, [TASK-%d][SERVER-%d] 执行结果: %s ====================", execCmd, task.ID, srv.ID, output)

	// Broadcast final output
	if len(output) > 0 {
		e.wsHub.BroadcastToTask(&model.WSMessage{
			Type:       "output",
			TaskID:     task.ID,
			ServerID:   srv.ID,
			ServerName: srv.Name,
			Content:    output,
			Timestamp:  time.Now(),
		}, uint64(task.ID))
	}

	if err != nil {
		if e.debug {
			logger.Error("[TASK-%d][SERVER-%d] 脚本执行失败: %v", task.ID, srv.ID, err)
		}
		finished := time.Now()

		// 构建详细的错误信息，包含退出码和实际输出
		errMsg := buildDetailedErrorMsg(err, execCmd, output)

		e.repo.UpdateTaskServerStatus(tsID, "failed", output, errMsg, &started, &finished)
		e.wsHub.BroadcastToTask(&model.WSMessage{
			Type:       "server_done",
			TaskID:     task.ID,
			ServerID:   srv.ID,
			ServerName: srv.Name,
			Status:     "failed",
			Content:    output,
			ExitCode:   1,
			Timestamp:  time.Now(),
		}, uint64(task.ID))
		return
	}

	finished := time.Now()
	if e.debug {
		logger.Debug("[TASK-%d][SERVER-%d] ========== 服务器执行成功 ==========", task.ID, srv.ID)
		logger.Debug("[TASK-%d][SERVER-%d] 执行输出 (最新50行):\n%s", task.ID, srv.ID, output)
	}
	e.repo.UpdateTaskServerStatus(tsID, "success", output, "", &started, &finished)
	e.wsHub.BroadcastToTask(&model.WSMessage{
		Type:       "server_done",
		TaskID:     task.ID,
		ServerID:   srv.ID,
		ServerName: srv.Name,
		Status:     "success",
		Content:    output,
		ExitCode:   0,
		Timestamp:  time.Now(),
	}, uint64(task.ID))
}

func (e *TaskExecutor) UploadFile(server *model.Server, password, localPath, remotePath string) error {
	client, err := e.sshPool.GetClient(server.ID, server.Host, server.Port, server.Username, password)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdout = io.Discard
	session.Stderr = io.Discard

	err = session.Shell()
	if err != nil {
		return err
	}

	return nil
}

type outputWriter struct {
	serverID   int64
	serverName string
	taskID     int64
	wsHub      *websocket.Hub
}

func (w *outputWriter) Write(p []byte) (n int, err error) {
	w.wsHub.BroadcastToTask(&model.WSMessage{
		Type:       "output",
		TaskID:     w.taskID,
		ServerID:   w.serverID,
		ServerName: w.serverName,
		Content:    string(p),
		Timestamp:  time.Now(),
	}, uint64(w.taskID))
	return len(p), nil
}

// streamingWriter writes output in real-time via WebSocket and keeps latest N lines
type streamingWriter struct {
	serverID   int64
	serverName string
	taskID     int64
	wsHub      *websocket.Hub
	maxLines   int
	lines      []string
}

func (w *streamingWriter) Write(p []byte) (n int, err error) {
	// Broadcast to WebSocket immediately for real-time display
	w.wsHub.BroadcastToTask(&model.WSMessage{
		Type:       "output",
		TaskID:     w.taskID,
		ServerID:   w.serverID,
		ServerName: w.serverName,
		Content:    string(p),
		Timestamp:  time.Now(),
	}, uint64(w.taskID))

	// Also accumulate lines for database storage (latest maxLines only)
	text := string(p)
	for _, line := range splitLines(text) {
		if line != "" {
			w.lines = append(w.lines, line)
		}
	}

	// Keep only the latest maxLines
	if len(w.lines) > w.maxLines {
		w.lines = w.lines[len(w.lines)-w.maxLines:]
	}

	return len(p), nil
}

func splitLines(s string) []string {
	var lines []string
	for _, line := range strings.Split(s, "\n") {
		lines = append(lines, line)
	}
	return lines
}

func (w *streamingWriter) GetOutput() string {
	return strings.Join(w.lines, "\n")
}

// limitLines returns only the latest n lines from the input
func limitLines(s string, n int) string {
	allLines := strings.Split(s, "\n")
	if len(allLines) <= n {
		return s
	}
	return strings.Join(allLines[len(allLines)-n:], "\n")
}

// buildDetailedErrorMsg 构建详细的错误信息，帮助排查问题
// 包含：退出码、错误原因、命令输出（stdout/stderr）
func buildDetailedErrorMsg(err error, cmd string, output string) string {
	var sb strings.Builder
	sb.WriteString("脚本执行失败\n")
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	sb.WriteString("命令: ")
	sb.WriteString(cmd)
	sb.WriteString("\n\n")

	// 检查是否是被信号终止的（如 Ctrl+C）
	if strings.Contains(err.Error(), "signal") {
		sb.WriteString("原因: 进程被信号终止\n")
	} else if strings.Contains(err.Error(), "exited") {
		// 尝试提取退出码
		sb.WriteString("原因: 进程异常退出\n")
	}

	// 添加错误详情
	sb.WriteString("错误: ")
	sb.WriteString(err.Error())
	sb.WriteString("\n\n")

	// 添加命令输出（如果有）
	if output != "" {
		sb.WriteString("命令输出 (stdout/stderr):\n")
		sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		sb.WriteString(output)
		sb.WriteString("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	}

	return sb.String()
}

func (e *TaskExecutor) DeployFile(server *model.Server, password, content, remotePath string) error {
	client, err := e.sshPool.GetClient(server.ID, server.Host, server.Port, server.Username, password)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdout = io.Discard
	session.Stderr = io.Discard

	err = session.Shell()
	if err != nil {
		return err
	}

	return nil
}
