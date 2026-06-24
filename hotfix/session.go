package hotfix

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type LogMessage struct {
	Time    string
	Message string
}

type PipelineSession struct {
	ID       string
	Cancel   context.CancelFunc
	Messages []LogMessage
	Status   string
	mu       sync.Mutex
}

const maxLogMessages = 500

func (s *PipelineSession) AddMsg(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Messages = append(s.Messages, LogMessage{
		Time:    time.Now().Format("15:04:05"),
		Message: msg,
	})
	if len(s.Messages) > maxLogMessages {
		s.Messages = s.Messages[len(s.Messages)-maxLogMessages:]
	}
}

func (s *PipelineSession) SetStatus(status string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Status = status
}

func (s *PipelineSession) GetStatus() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Status
}

func (s *PipelineSession) GetMessages(since int) ([]LogMessage, int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if since >= len(s.Messages) {
		return nil, len(s.Messages)
	}
	return s.Messages[since:], len(s.Messages)
}

var (
	sessions   = map[string]*PipelineSession{}
	sessionsMu sync.Mutex
	sessionSeq int

	pipelineMu sync.Mutex
)

// StartSession creates a new pipeline session with a cancelable context.
// Returns an error if a pipeline is already running.
func StartSession(desc string) (*PipelineSession, error) {
	pipelineMu.Lock()
	defer pipelineMu.Unlock()

	for _, s := range sessions {
		if s.GetStatus() == "running" {
			return nil, fmt.Errorf("pipeline %q is already running", s.ID)
		}
	}

	sessionSeq++
	_, cancel := context.WithCancel(context.Background())
	session := &PipelineSession{
		ID:     fmt.Sprintf("hotfix-%d", sessionSeq),
		Cancel: cancel,
		Status: "running",
	}
	session.AddMsg(fmt.Sprintf("Pipeline %q started: %s", session.ID, desc))

	sessionsMu.Lock()
	sessions[session.ID] = session
	sessionsMu.Unlock()

	return session, nil
}

// GetSession retrieves a pipeline session by ID.
func GetSession(id string) *PipelineSession {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	return sessions[id]
}

// CancelSession cancels a running pipeline.
func CancelSession(id string) error {
	sessionsMu.Lock()
	session := sessions[id]
	sessionsMu.Unlock()
	if session == nil {
		return fmt.Errorf("session %q not found", id)
	}
	session.Cancel()
	session.SetStatus("cancelled")
	session.AddMsg("Pipeline cancelled by user.")
	return nil
}

// FinishSession marks a session as done or error and cleans up after a delay.
func FinishSession(session *PipelineSession, err error) {
	if err != nil {
		session.SetStatus("error")
		session.AddMsg(fmt.Sprintf("Pipeline failed: %s", err.Error()))
	} else {
		session.SetStatus("done")
		session.AddMsg("All Is OK!")
	}

	go func() {
		time.Sleep(10 * time.Minute)
		sessionsMu.Lock()
		delete(sessions, session.ID)
		sessionsMu.Unlock()
	}()
}

