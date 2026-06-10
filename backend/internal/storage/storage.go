package storage

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"goalstack/internal/models"
)

// Store interface defines the storage contract
type Store interface {
	// Goal operations
	CreateGoal(goal *models.Goal) error
	GetGoal(id string) (*models.Goal, error)
	UpdateGoal(id string, updates *models.UpdateGoalRequest) (*models.Goal, error)
	DeleteGoal(id string) error
	ListGoals() ([]models.Goal, error)

	// Subtask operations
	AddSubtask(goalID string, subtask *models.Subtask) error
	GetSubtask(goalID, subtaskID string) (*models.Subtask, error)
	UpdateSubtask(goalID, subtaskID string, updates *models.UpdateSubtaskRequest) (*models.Subtask, error)
	UpdateSubtaskProgress(subtaskID string, progress float64) (*models.Goal, error)
	CompleteSubtask(subtaskID string) (*models.Goal, error)
	ReorderSubtasks(goalID string, subtaskIDs []string) (*models.Goal, error)
	DeleteSubtask(goalID, subtaskID string) error

	// Note operations
	AddNote(goalID, subtaskID string, note *models.Note) error
	DeleteNote(goalID, subtaskID, noteID string) error

	// Link operations
	AddLink(goalID, subtaskID string, link *models.Link) error
	DeleteLink(goalID, subtaskID, linkID string) error

	// Checklist operations
	AddChecklistItem(goalID, subtaskID string, item *models.ChecklistItem) error
	UpdateChecklistItem(goalID, subtaskID, itemID string, completed bool) error
	DeleteChecklistItem(goalID, subtaskID, itemID string) error
}

// MemoryStore implements Store interface with in-memory storage
type MemoryStore struct {
	mu    sync.RWMutex
	goals map[string]*models.Goal
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		goals: make(map[string]*models.Goal),
	}
}

// CreateGoal creates a new goal
func (m *MemoryStore) CreateGoal(goal *models.Goal) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.goals[goal.ID]; exists {
		return fmt.Errorf("goal with ID %s already exists", goal.ID)
	}

	goal.CreatedAt = time.Now()
	goal.UpdatedAt = time.Now()
	goal.Subtasks = []models.Subtask{}
	m.goals[goal.ID] = goal
	return nil
}

// GetGoal retrieves a goal by ID
func (m *MemoryStore) GetGoal(id string) (*models.Goal, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	goal, exists := m.goals[id]
	if !exists {
		return nil, fmt.Errorf("goal not found")
	}

	return goal, nil
}

// UpdateGoal updates an existing goal
func (m *MemoryStore) UpdateGoal(id string, updates *models.UpdateGoalRequest) (*models.Goal, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	goal, exists := m.goals[id]
	if !exists {
		return nil, fmt.Errorf("goal not found")
	}

	if updates.Title != nil {
		goal.Title = *updates.Title
	}
	if updates.StartDate != nil {
		goal.StartDate = updates.StartDate.Time
	}
	if updates.TotalDuration != nil {
		goal.TotalDuration = *updates.TotalDuration
	}
	if updates.DurationType != nil {
		goal.DurationType = models.DurationType(*updates.DurationType)
	}

	goal.UpdatedAt = time.Now()
	return goal, nil
}

// DeleteGoal deletes a goal
func (m *MemoryStore) DeleteGoal(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.goals[id]; !exists {
		return fmt.Errorf("goal not found")
	}

	delete(m.goals, id)
	return nil
}

// ListGoals returns all goals
func (m *MemoryStore) ListGoals() ([]models.Goal, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	goals := make([]models.Goal, 0, len(m.goals))
	for _, goal := range m.goals {
		goals = append(goals, *goal)
	}

	return goals, nil
}

// AddSubtask adds a subtask to a goal
func (m *MemoryStore) AddSubtask(goalID string, subtask *models.Subtask) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	goal, exists := m.goals[goalID]
	if !exists {
		return fmt.Errorf("goal not found")
	}

	subtask.CreatedAt = time.Now()
	subtask.UpdatedAt = time.Now()
	subtask.Notes = []models.Note{}
	subtask.Links = []models.Link{}
	subtask.ChecklistItems = []models.ChecklistItem{}
	subtask.Progress = 0
	subtask.Status = models.Todo

	// Set order if not provided
	if subtask.Order == 0 {
		subtask.Order = len(goal.Subtasks) + 1
	}

	goal.Subtasks = append(goal.Subtasks, *subtask)
	goal.UpdatedAt = time.Now()

	return nil
}

// GetSubtask retrieves a subtask by ID
func (m *MemoryStore) GetSubtask(goalID, subtaskID string) (*models.Subtask, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	goal, exists := m.goals[goalID]
	if !exists {
		return nil, fmt.Errorf("goal not found")
	}

	for i, subtask := range goal.Subtasks {
		if subtask.ID == subtaskID {
			return &goal.Subtasks[i], nil
		}
	}

	return nil, fmt.Errorf("subtask not found")
}

// UpdateSubtask updates a subtask
func (m *MemoryStore) UpdateSubtask(goalID, subtaskID string, updates *models.UpdateSubtaskRequest) (*models.Subtask, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	goal, exists := m.goals[goalID]
	if !exists {
		return nil, fmt.Errorf("goal not found")
	}

	for i, subtask := range goal.Subtasks {
		if subtask.ID == subtaskID {
			if updates.Title != nil {
				goal.Subtasks[i].Title = *updates.Title
			}
			if updates.Weight != nil {
				goal.Subtasks[i].Weight = *updates.Weight
			}
			if updates.Order != nil {
				goal.Subtasks[i].Order = *updates.Order
			}

			goal.Subtasks[i].UpdatedAt = time.Now()
			goal.UpdatedAt = time.Now()
			if updates.Order != nil {
				m.normalizeSubtaskOrdersLocked(goal)
			}

			for j := range goal.Subtasks {
				if goal.Subtasks[j].ID == subtaskID {
					return &goal.Subtasks[j], nil
				}
			}

			return &goal.Subtasks[i], nil
		}
	}

	return nil, fmt.Errorf("subtask not found")
}

// UpdateSubtaskProgress updates progress for a subtask by ID.
func (m *MemoryStore) UpdateSubtaskProgress(subtaskID string, progress float64) (*models.Goal, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, goal := range m.goals {
		for i := range goal.Subtasks {
			if goal.Subtasks[i].ID == subtaskID {
				goal.Subtasks[i].Progress = progress
				goal.Subtasks[i].Status = statusForProgress(progress)
				goal.Subtasks[i].UpdatedAt = time.Now()
				goal.UpdatedAt = time.Now()
				return goal, nil
			}
		}
	}

	return nil, fmt.Errorf("subtask not found")
}

// CompleteSubtask marks a subtask complete by ID.
func (m *MemoryStore) CompleteSubtask(subtaskID string) (*models.Goal, error) {
	return m.UpdateSubtaskProgress(subtaskID, 100)
}

// ReorderSubtasks updates the order of subtasks within a goal.
func (m *MemoryStore) ReorderSubtasks(goalID string, subtaskIDs []string) (*models.Goal, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	goal, exists := m.goals[goalID]
	if !exists {
		return nil, fmt.Errorf("goal not found")
	}

	if len(subtaskIDs) != len(goal.Subtasks) {
		return nil, fmt.Errorf("subtask order must include every subtask")
	}

	positions := make(map[string]int, len(subtaskIDs))
	for i, id := range subtaskIDs {
		if _, exists := positions[id]; exists {
			return nil, fmt.Errorf("subtask order contains duplicate IDs")
		}
		positions[id] = i + 1
	}

	for i := range goal.Subtasks {
		position, exists := positions[goal.Subtasks[i].ID]
		if !exists {
			return nil, fmt.Errorf("subtask order contains unknown IDs")
		}
		goal.Subtasks[i].Order = position
		goal.Subtasks[i].UpdatedAt = time.Now()
	}

	m.normalizeSubtaskOrdersLocked(goal)
	goal.UpdatedAt = time.Now()
	return goal, nil
}

// DeleteSubtask deletes a subtask
func (m *MemoryStore) DeleteSubtask(goalID, subtaskID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	goal, exists := m.goals[goalID]
	if !exists {
		return fmt.Errorf("goal not found")
	}

	for i, subtask := range goal.Subtasks {
		if subtask.ID == subtaskID {
			goal.Subtasks = append(goal.Subtasks[:i], goal.Subtasks[i+1:]...)
			m.normalizeSubtaskOrdersLocked(goal)
			goal.UpdatedAt = time.Now()
			return nil
		}
	}

	return fmt.Errorf("subtask not found")
}

// AddNote adds a note to a subtask
func (m *MemoryStore) AddNote(goalID, subtaskID string, note *models.Note) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	goal, exists := m.goals[goalID]
	if !exists {
		return fmt.Errorf("goal not found")
	}

	for i, subtask := range goal.Subtasks {
		if subtask.ID == subtaskID {
			note.CreatedAt = time.Now()
			goal.Subtasks[i].Notes = append(goal.Subtasks[i].Notes, *note)
			goal.UpdatedAt = time.Now()
			return nil
		}
	}

	return fmt.Errorf("subtask not found")
}

// DeleteNote deletes a note from a subtask
func (m *MemoryStore) DeleteNote(goalID, subtaskID, noteID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	goal, exists := m.goals[goalID]
	if !exists {
		return fmt.Errorf("goal not found")
	}

	for i, subtask := range goal.Subtasks {
		if subtask.ID == subtaskID {
			for j, note := range subtask.Notes {
				if note.ID == noteID {
					goal.Subtasks[i].Notes = append(
						goal.Subtasks[i].Notes[:j],
						goal.Subtasks[i].Notes[j+1:]...,
					)
					goal.UpdatedAt = time.Now()
					return nil
				}
			}
		}
	}

	return fmt.Errorf("note not found")
}

// AddLink adds a link to a subtask
func (m *MemoryStore) AddLink(goalID, subtaskID string, link *models.Link) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	goal, exists := m.goals[goalID]
	if !exists {
		return fmt.Errorf("goal not found")
	}

	for i, subtask := range goal.Subtasks {
		if subtask.ID == subtaskID {
			goal.Subtasks[i].Links = append(goal.Subtasks[i].Links, *link)
			goal.UpdatedAt = time.Now()
			return nil
		}
	}

	return fmt.Errorf("subtask not found")
}

// DeleteLink deletes a link from a subtask
func (m *MemoryStore) DeleteLink(goalID, subtaskID, linkID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	goal, exists := m.goals[goalID]
	if !exists {
		return fmt.Errorf("goal not found")
	}

	for i, subtask := range goal.Subtasks {
		if subtask.ID == subtaskID {
			for j, link := range subtask.Links {
				if link.ID == linkID {
					goal.Subtasks[i].Links = append(
						goal.Subtasks[i].Links[:j],
						goal.Subtasks[i].Links[j+1:]...,
					)
					goal.UpdatedAt = time.Now()
					return nil
				}
			}
		}
	}

	return fmt.Errorf("link not found")
}

// AddChecklistItem adds a checklist item to a subtask
func (m *MemoryStore) AddChecklistItem(goalID, subtaskID string, item *models.ChecklistItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	goal, exists := m.goals[goalID]
	if !exists {
		return fmt.Errorf("goal not found")
	}

	for i, subtask := range goal.Subtasks {
		if subtask.ID == subtaskID {
			goal.Subtasks[i].ChecklistItems = append(goal.Subtasks[i].ChecklistItems, *item)
			goal.UpdatedAt = time.Now()
			return nil
		}
	}

	return fmt.Errorf("subtask not found")
}

// UpdateChecklistItem updates a checklist item
func (m *MemoryStore) UpdateChecklistItem(goalID, subtaskID, itemID string, completed bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	goal, exists := m.goals[goalID]
	if !exists {
		return fmt.Errorf("goal not found")
	}

	for i, subtask := range goal.Subtasks {
		if subtask.ID == subtaskID {
			for j, item := range subtask.ChecklistItems {
				if item.ID == itemID {
					goal.Subtasks[i].ChecklistItems[j].Completed = completed
					goal.UpdatedAt = time.Now()
					return nil
				}
			}
		}
	}

	return fmt.Errorf("checklist item not found")
}

func (m *MemoryStore) normalizeSubtaskOrdersLocked(goal *models.Goal) {
	sort.SliceStable(goal.Subtasks, func(i, j int) bool {
		return goal.Subtasks[i].Order < goal.Subtasks[j].Order
	})

	for i := range goal.Subtasks {
		goal.Subtasks[i].Order = i + 1
	}
}

func statusForProgress(progress float64) models.StatusType {
	switch {
	case progress >= 100:
		return models.Completed
	case progress > 0:
		return models.InProgress
	default:
		return models.Todo
	}
}

// DeleteChecklistItem deletes a checklist item
func (m *MemoryStore) DeleteChecklistItem(goalID, subtaskID, itemID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	goal, exists := m.goals[goalID]
	if !exists {
		return fmt.Errorf("goal not found")
	}

	for i, subtask := range goal.Subtasks {
		if subtask.ID == subtaskID {
			for j, item := range subtask.ChecklistItems {
				if item.ID == itemID {
					goal.Subtasks[i].ChecklistItems = append(
						goal.Subtasks[i].ChecklistItems[:j],
						goal.Subtasks[i].ChecklistItems[j+1:]...,
					)
					goal.UpdatedAt = time.Now()
					return nil
				}
			}
		}
	}

	return fmt.Errorf("checklist item not found")
}
