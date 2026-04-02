package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/sftp"
	"node-pilot/internal/logger"
	"node-pilot/internal/model"
	"node-pilot/internal/repository"
)

type FileUploadService struct {
	repo    *repository.Repository
	sshPool *SSHPool
	baseDir string
}

func NewFileUploadService(repo *repository.Repository, sshPool *SSHPool, baseDir string) *FileUploadService {
	return &FileUploadService{
		repo:    repo,
		sshPool: sshPool,
		baseDir: baseDir,
	}
}

func (s *FileUploadService) UploadFile(serverID int64, localFilePath, remotePath string) (string, error) {
	server, err := s.repo.GetServer(serverID)
	if err != nil {
		return "", fmt.Errorf("failed to get server: %w", err)
	}

	client, err := s.sshPool.GetClient(server.ID, server.Host, server.Port, server.Username, "")
	if err != nil {
		return "", fmt.Errorf("SSH connection failed: %w", err)
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return "", fmt.Errorf("SFTP client creation failed: %w", err)
	}
	defer sftpClient.Close()

	fullLocalPath := filepath.Join(s.baseDir, localFilePath)

	if _, err := os.Stat(fullLocalPath); os.IsNotExist(err) {
		return "", fmt.Errorf("local file not found: %s", fullLocalPath)
	}

	srcFile, err := os.Open(fullLocalPath)
	if err != nil {
		return "", fmt.Errorf("failed to open local file: %w", err)
	}
	defer srcFile.Close()

	remoteFullPath := filepath.Join(remotePath, filepath.Base(localFilePath))

	err = sftpClient.MkdirAll(filepath.Dir(remoteFullPath))
	if err != nil {
		return "", fmt.Errorf("failed to create remote directory: %w", err)
	}

	dstFile, err := sftpClient.Create(remoteFullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create remote file: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	return remoteFullPath, nil
}

func (s *FileUploadService) ExecuteUpload(fileUploadID int64) error {
	fu, err := s.repo.GetFileUploadByID(fileUploadID)
	if err != nil {
		return fmt.Errorf("failed to get file upload: %w", err)
	}

	servers, err := s.repo.GetFileUploadServers(fileUploadID)
	if err != nil {
		return fmt.Errorf("failed to get file upload servers: %w", err)
	}

	if len(servers) == 0 {
		return fmt.Errorf("no servers configured for this upload")
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	allSuccess := true

	for _, serverRecord := range servers {
		wg.Add(1)
		go func(srv *model.FileUploadServer) {
			defer wg.Done()

			fullLocalPath := filepath.Join(s.baseDir, fu.LocalPath)
			if _, err := os.Stat(fullLocalPath); os.IsNotExist(err) {
				s.repo.UpdateFileUploadServerStatus(srv.ID, "failed", "本地文件不存在，请重新上传")
				mu.Lock()
				allSuccess = false
				mu.Unlock()
				return
			}

			server, err := s.repo.GetServer(srv.ServerID)
			if err != nil {
				s.repo.UpdateFileUploadServerStatus(srv.ID, "failed", fmt.Sprintf("获取服务器信息失败: %v", err))
				mu.Lock()
				allSuccess = false
				mu.Unlock()
				return
			}

			client, err := s.sshPool.GetClient(server.ID, server.Host, server.Port, server.Username, "")
			if err != nil {
				s.repo.UpdateFileUploadServerStatus(srv.ID, "failed", fmt.Sprintf("SSH连接失败: %v", err))
				mu.Lock()
				allSuccess = false
				mu.Unlock()
				return
			}
			defer client.Close()

			sftpClient, err := sftp.NewClient(client)
			if err != nil {
				s.repo.UpdateFileUploadServerStatus(srv.ID, "failed", fmt.Sprintf("SFTP客户端创建失败: %v", err))
				mu.Lock()
				allSuccess = false
				mu.Unlock()
				return
			}
			defer sftpClient.Close()

			srcFile, err := os.Open(fullLocalPath)
			if err != nil {
				s.repo.UpdateFileUploadServerStatus(srv.ID, "failed", fmt.Sprintf("打开本地文件失败: %v", err))
				mu.Lock()
				allSuccess = false
				mu.Unlock()
				return
			}
			defer srcFile.Close()

			remoteFullPath := filepath.Join(fu.RemotePath, filepath.Base(fu.LocalPath))
			if srv.RemoteFullPath != "" {
				remoteFullPath = srv.RemoteFullPath
			}

			err = sftpClient.MkdirAll(filepath.Dir(remoteFullPath))
			if err != nil {
				s.repo.UpdateFileUploadServerStatus(srv.ID, "failed", fmt.Sprintf("创建远程目录失败: %v", err))
				mu.Lock()
				allSuccess = false
				mu.Unlock()
				return
			}

			dstFile, err := sftpClient.Create(remoteFullPath)
			if err != nil {
				s.repo.UpdateFileUploadServerStatus(srv.ID, "failed", fmt.Sprintf("创建远程文件失败: %v", err))
				mu.Lock()
				allSuccess = false
				mu.Unlock()
				return
			}
			defer dstFile.Close()

			_, err = io.Copy(dstFile, srcFile)
			if err != nil {
				s.repo.UpdateFileUploadServerStatus(srv.ID, "failed", fmt.Sprintf("文件传输失败: %v", err))
				mu.Lock()
				allSuccess = false
				mu.Unlock()
				return
			}

			s.repo.UpdateFileUploadServerStatus(srv.ID, "success", "")
			logger.Info("[FILE-UPLOAD] Successfully uploaded to server %s: %s", server.Host, remoteFullPath)
		}(serverRecord)
	}

	wg.Wait()

	if allSuccess {
		fu.Status = "success"
	} else {
		fu.Status = "failed"
	}
	s.repo.UpdateFileUpload(fu)

	return nil
}
