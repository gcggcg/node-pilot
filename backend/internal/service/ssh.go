package service

import (
	"fmt"
	"sync"
	"time"

	"node-pilot/internal/logger"

	"golang.org/x/crypto/ssh"
)

type SSHPool struct {
	mu       sync.Mutex
	clients  map[int64]*ssh.Client
	maxConns int
	debug    bool
}

func NewSSHPool() *SSHPool {
	return &SSHPool{
		clients:  make(map[int64]*ssh.Client),
		maxConns: 50,
	}
}

func (p *SSHPool) GetClient(serverID int64, host string, port int, username, password string) (*ssh.Client, error) {
	if p.debug {
		logger.Debug("[SSH-POOL] 获取SSH客户端: serverID=%d, host=%s:%d, user=%s", serverID, host, port, username)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// Check for existing connection
	if client, ok := p.clients[serverID]; ok {
		_, _, err := client.SendRequest("keepalive@golang.org", true, nil)
		if err == nil {
			if p.debug {
				logger.Debug("[SSH-POOL] 复用现有连接: serverID=%d", serverID)
			}
			return client, nil
		}
		if p.debug {
			logger.Debug("[SSH-POOL] 现有连接已失效，重新建立: serverID=%d, err=%v", serverID, err)
		}
		client.Close()
		delete(p.clients, serverID)
	}

	// Check pool size limit
	if len(p.clients) >= p.maxConns {
		if p.debug {
			logger.Debug("[SSH-POOL] 连接池已满，关闭一个旧连接")
		}
		for id, client := range p.clients {
			client.Close()
			delete(p.clients, id)
			break
		}
	}

	// Create new connection
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	if p.debug {
		logger.Debug("[SSH-POOL] 建立新SSH连接: %s", addr)
	}
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		if p.debug {
			logger.Error("[SSH-POOL] SSH连接失败: %s, err=%v", addr, err)
		}
		return nil, err
	}

	if p.debug {
		logger.Debug("[SSH-POOL] SSH连接成功: serverID=%d, addr=%s", serverID, addr)
	}
	p.clients[serverID] = client
	return client, nil
}

func (p *SSHPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for id, client := range p.clients {
		client.Close()
		delete(p.clients, id)
	}
}

func (p *SSHPool) CloseServer(serverID int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if client, ok := p.clients[serverID]; ok {
		client.Close()
		delete(p.clients, serverID)
	}
}
