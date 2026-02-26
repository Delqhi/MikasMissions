package internal

import (
	"encoding/json"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Store struct {
	mu          sync.Mutex
	workflows   map[string]WorkflowTemplate
	runs        map[string]WorkflowRun
	runLogs     map[string][]WorkflowRunLog
	modelConfig map[string]ModelProfile
}

func NewStore() *Store {
	return &Store{
		workflows: map[string]WorkflowTemplate{},
		runs:      map[string]WorkflowRun{},
		runLogs:   map[string][]WorkflowRunLog{},
		modelConfig: map[string]ModelProfile{
			"nim-default": {
				ID:           "nim-default",
				Provider:     "nvidia_nim",
				BaseURL:      "http://localhost:9000",
				ModelID:      "nim-video-v1",
				TimeoutMS:    15000,
				MaxRetries:   2,
				SafetyPreset: "kids_strict",
			},
		},
	}
}

func (s *Store) ListWorkflows() ([]WorkflowTemplate, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]WorkflowTemplate, 0, len(s.workflows))
	for _, workflow := range s.workflows {
		result = append(result, workflow)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result, nil
}

func (s *Store) CreateWorkflow(workflow WorkflowTemplate, _ string) (WorkflowTemplate, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	workflow.ID = uuid.NewString()
	workflow.Version = 1
	s.workflows[workflow.ID] = workflow
	return workflow, nil
}

func (s *Store) UpdateWorkflow(workflow WorkflowTemplate, _ string) (WorkflowTemplate, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.workflows[workflow.ID]
	if !ok {
		return WorkflowTemplate{}, false, nil
	}
	workflow.Version = existing.Version + 1
	s.workflows[workflow.ID] = workflow
	return workflow, true, nil
}

func (s *Store) DeleteWorkflow(workflowID, _ string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.workflows[workflowID]; !ok {
		return false, nil
	}
	delete(s.workflows, workflowID)
	return true, nil
}

func (s *Store) FindWorkflow(workflowID string) (WorkflowTemplate, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	workflow, ok := s.workflows[workflowID]
	return workflow, ok, nil
}

func (s *Store) CreateRun(run WorkflowRun, _ string) (WorkflowRun, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run.ID = uuid.NewString()
	run.Status = "requested"
	if len(run.InputPayload) == 0 {
		run.InputPayload = json.RawMessage(`{}`)
	}
	s.runs[run.ID] = run
	return run, nil
}

func (s *Store) FindRun(runID string) (WorkflowRun, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.runs[runID]
	return run, ok, nil
}

func (s *Store) ListRunLogs(runID string) ([]WorkflowRunLog, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entries := s.runLogs[runID]
	cloned := make([]WorkflowRunLog, 0, len(entries))
	cloned = append(cloned, entries...)
	return cloned, nil
}

func (s *Store) AppendRunLog(logEntry WorkflowRunLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if logEntry.EventTime == "" {
		logEntry.EventTime = time.Now().UTC().Format(time.RFC3339)
	}
	s.runLogs[logEntry.RunID] = append(s.runLogs[logEntry.RunID], logEntry)
	return nil
}

func (s *Store) SetRunStatus(runID, status, lastError string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.runs[runID]
	if !ok {
		return false, nil
	}
	run.Status = status
	run.LastError = lastError
	s.runs[runID] = run
	return true, nil
}

func (s *Store) GetModelProfile(modelProfileID string) (ModelProfile, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	profile, ok := s.modelConfig[modelProfileID]
	return profile, ok, nil
}

func (s *Store) PutModelProfile(profile ModelProfile, _ string) (ModelProfile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.modelConfig[profile.ID] = profile
	return profile, nil
}
