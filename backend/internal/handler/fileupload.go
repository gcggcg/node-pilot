package handler

import (
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

	"github.com/gin-gonic/gin"
)

type FileUploadHandler struct {
	repo      *repository.Repository
	uploadSvc *service.FileUploadService
	baseDir   string
}

func NewFileUploadHandler(repo *repository.Repository, uploadSvc *service.FileUploadService, baseDir string) *FileUploadHandler {
	return &FileUploadHandler{
		repo:      repo,
		uploadSvc: uploadSvc,
		baseDir:   baseDir,
	}
}

func (h *FileUploadHandler) ListFileUploads(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var startTime, endTime *time.Time
	if st := c.Query("start_time"); st != "" {
		if t, err := time.Parse("2006-01-02", st); err == nil {
			startTime = &t
		}
	}
	if et := c.Query("end_time"); et != "" {
		if t, err := time.Parse("2006-01-02", et); err == nil {
			t = t.Add(24*time.Hour - time.Second)
			endTime = &t
		}
	}

	uploads, total, err := h.repo.ListFileUploads(page, pageSize, status, keyword, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if uploads == nil {
		uploads = []*model.FileUpload{}
	}

	// 填充服务器信息
	type FileUploadWithServers struct {
		model.FileUpload
		Servers []*model.Server `json:"servers"`
	}

	result := make([]*FileUploadWithServers, 0, len(uploads))
	for _, fu := range uploads {
		fus, _ := h.repo.GetFileUploadServers(fu.ID)
		serverIDs := make([]int64, 0, len(fus))
		for _, f := range fus {
			serverIDs = append(serverIDs, f.ServerID)
		}
		var servers []*model.Server
		for _, sid := range serverIDs {
			if s, err := h.repo.GetServer(sid); err == nil {
				servers = append(servers, s)
			}
		}
		result = append(result, &FileUploadWithServers{
			FileUpload: *fu,
			Servers:    servers,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     result,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *FileUploadHandler) CreateFileUpload(c *gin.Context) {
	type CreateFileUploadInput struct {
		Name       string  `json:"name" binding:"required"`
		LocalPath  string  `json:"local_path" binding:"required"`
		RemotePath string  `json:"remote_path" binding:"required"`
		ServerIDs  []int64 `json:"server_ids" binding:"required"`
	}

	var input CreateFileUploadInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证本地文件存在
	fullLocalPath := filepath.Join(h.baseDir, input.LocalPath)
	if _, err := os.Stat(fullLocalPath); os.IsNotExist(err) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "本地文件不存在，请先上传文件"})
		return
	}

	fu := &model.FileUpload{
		Name:       input.Name,
		LocalPath:  input.LocalPath,
		RemotePath: input.RemotePath,
		Status:     "pending",
	}

	id, err := h.repo.CreateFileUpload(fu)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 创建服务器关联记录
	var fus []*model.FileUploadServer
	for _, sid := range input.ServerIDs {
		fus = append(fus, &model.FileUploadServer{
			FileUploadID: id,
			ServerID:     sid,
			Status:       "pending",
		})
	}
	if err := h.repo.CreateFileUploadServers(fus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *FileUploadHandler) UpdateFileUpload(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	type UpdateFileUploadInput struct {
		Name       string  `json:"name"`
		LocalPath  string  `json:"local_path"`
		RemotePath string  `json:"remote_path"`
		ServerIDs  []int64 `json:"server_ids"`
	}

	var input UpdateFileUploadInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fu, err := h.repo.GetFileUploadByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file upload not found"})
		return
	}

	if input.Name != "" {
		fu.Name = input.Name
	}
	if input.LocalPath != "" {
		fu.LocalPath = input.LocalPath
	}
	if input.RemotePath != "" {
		fu.RemotePath = input.RemotePath
	}

	if err := h.repo.UpdateFileUpload(fu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 更新服务器关联
	if len(input.ServerIDs) > 0 {
		// 删除旧的关联
		if err := h.repo.DeleteFileUploads([]int64{id}); err != nil {
			logger.Warn("Failed to delete old file upload servers: %v", err)
		}
		// 创建新的关联
		var fus []*model.FileUploadServer
		for _, sid := range input.ServerIDs {
			fus = append(fus, &model.FileUploadServer{
				FileUploadID: id,
				ServerID:     sid,
				Status:       "pending",
			})
		}
		if err := h.repo.CreateFileUploadServers(fus); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *FileUploadHandler) DeleteFileUploads(c *gin.Context) {
	var input struct {
		IDs []int64 `json:"ids"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.DeleteFileUploads(input.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *FileUploadHandler) ExecuteFileUpload(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	fu, err := h.repo.GetFileUploadByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file upload not found"})
		return
	}

	if fu.Status == "running" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "upload is already running"})
		return
	}

	// 异步执行上传
	go func() {
		if err := h.uploadSvc.ExecuteUpload(id); err != nil {
			logger.Error("[FILE-UPLOAD] Execute failed: %v", err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "upload started"})
}

func (h *FileUploadHandler) GetFileUploadResults(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	fu, err := h.repo.GetFileUploadByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file upload not found"})
		return
	}

	servers, err := h.repo.GetFileUploadServers(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type ServerResult struct {
		ServerID       int64  `json:"server_id"`
		ServerName     string `json:"server_name"`
		Status         string `json:"status"`
		ErrorMessage   string `json:"error_message,omitempty"`
		RemoteFullPath string `json:"remote_full_path,omitempty"`
	}

	results := make([]*ServerResult, 0, len(servers))
	for _, s := range servers {
		server, _ := h.repo.GetServer(s.ServerID)
		serverName := ""
		if server != nil {
			serverName = server.Name
		}
		results = append(results, &ServerResult{
			ServerID:       s.ServerID,
			ServerName:     serverName,
			Status:         s.Status,
			ErrorMessage:   s.ErrorMessage,
			RemoteFullPath: s.RemoteFullPath,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"file_upload": fu,
		"results":     results,
	})
}

// UploadFileToStorage 处理文件上传到本地存储
func (h *FileUploadHandler) UploadFileToStorage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}
	defer file.Close()

	// 创建上传目录
	if err := os.MkdirAll(h.baseDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
		return
	}

	// 生成唯一文件名
	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + filepath.Base(header.Filename)
	localPath := filepath.Join(h.baseDir, filename)

	out, err := os.Create(localPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// 返回相对路径（相对于baseDir）
	relativePath := filename
	c.JSON(http.StatusOK, gin.H{
		"path": relativePath,
		"name": header.Filename,
		"size": header.Size,
	})
}
