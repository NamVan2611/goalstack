package storage

import (
	"fmt"
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
			return &goal.Subtasks[i], nil
		}
	}

	return nil, fmt.Errorf("subtask not found")
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
