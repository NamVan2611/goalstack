package service

import (
	"time"

	"github.com/google/uuid"

	"goalstack/internal/models"
	"goalstack/internal/repository"
	"goalstack/internal/timeline"
)

type GoalService struct {
	repo repository.GoalRepository
}

func NewGoalService(repo repository.GoalRepository) *GoalService {
	return &GoalService{repo: repo}
}

func (s *GoalService) CreateGoal(req *models.CreateGoalRequest) (*models.Goal, error) {
	goal := &models.Goal{
		ID:            uuid.New().String(),
		Title:         req.Title,
		StartDate:     req.StartDate.Time,
		TotalDuration: req.TotalDuration,
		DurationType:  req.DurationType,
		Subtasks:      []models.Subtask{},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.CreateGoal(goal); err != nil {
		return nil, err
	}

	return timeline.GetTimelineWithCalculations(goal), nil
}

func (s *GoalService) GetGoal(id string) (*models.Goal, error) {
	goal, err := s.repo.GetGoal(id)
	if err != nil {
		return nil, err
	}
	return timeline.GetTimelineWithCalculations(goal), nil
}

func (s *GoalService) UpdateGoal(id string, req *models.UpdateGoalRequest) (*models.Goal, error) {
	goal, err := s.repo.UpdateGoal(id, req)
	if err != nil {
		return nil, err
	}
	return timeline.GetTimelineWithCalculations(goal), nil
}

func (s *GoalService) DeleteGoal(id string) error {
	return s.repo.DeleteGoal(id)
}

func (s *GoalService) ListGoals() ([]models.Goal, error) {
	goals, err := s.repo.ListGoals()
	if err != nil {
		return nil, err
	}

	for i := range goals {
		calculated := timeline.GetTimelineWithCalculations(&goals[i])
		goals[i] = *calculated
	}

	if goals == nil {
		return []models.Goal{}, nil
	}
	return goals, nil
}

func (s *GoalService) AddSubtask(goalID string, req *models.CreateSubtaskRequest) (*models.Goal, error) {
	subtask := &models.Subtask{
		ID:             uuid.New().String(),
		GoalID:         goalID,
		Title:          req.Title,
		Weight:         req.Weight,
		Notes:          []models.Note{},
		Links:          []models.Link{},
		ChecklistItems: []models.ChecklistItem{},
		Progress:       0,
		Status:         models.Todo,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.AddSubtask(goalID, subtask); err != nil {
		return nil, err
	}

	return s.GetGoal(goalID)
}

func (s *GoalService) UpdateSubtask(goalID, taskID string, req *models.UpdateSubtaskRequest) (*models.Goal, error) {
	if _, err := s.repo.UpdateSubtask(goalID, taskID, req); err != nil {
		return nil, err
	}
	return s.GetGoal(goalID)
}

func (s *GoalService) UpdateSubtaskProgress(taskID string, progress float64) (*models.Goal, error) {
	goal, err := s.repo.UpdateSubtaskProgress(taskID, progress)
	if err != nil {
		return nil, err
	}
	return timeline.GetTimelineWithCalculations(goal), nil
}

func (s *GoalService) CompleteSubtask(taskID string) (*models.Goal, error) {
	return s.UpdateSubtaskProgress(taskID, 100)
}

func (s *GoalService) ReorderSubtasks(goalID string, subtaskIDs []string) (*models.Goal, error) {
	goal, err := s.repo.ReorderSubtasks(goalID, subtaskIDs)
	if err != nil {
		return nil, err
	}
	return timeline.GetTimelineWithCalculations(goal), nil
}

func (s *GoalService) DeleteSubtask(goalID, taskID string) (*models.Goal, error) {
	if err := s.repo.DeleteSubtask(goalID, taskID); err != nil {
		return nil, err
	}
	return s.GetGoal(goalID)
}

func (s *GoalService) AddNote(goalID, taskID string, req *models.AddNoteRequest) (*models.Goal, error) {
	note := &models.Note{
		ID:        uuid.New().String(),
		SubtaskID: taskID,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	if err := s.repo.AddNote(goalID, taskID, note); err != nil {
		return nil, err
	}
	return s.GetGoal(goalID)
}

func (s *GoalService) DeleteNote(goalID, taskID, noteID string) (*models.Goal, error) {
	if err := s.repo.DeleteNote(goalID, taskID, noteID); err != nil {
		return nil, err
	}
	return s.GetGoal(goalID)
}

func (s *GoalService) AddLink(goalID, taskID string, req *models.AddLinkRequest) (*models.Goal, error) {
	link := &models.Link{
		ID:        uuid.New().String(),
		SubtaskID: taskID,
		Title:     req.Title,
		URL:       req.URL,
	}

	if err := s.repo.AddLink(goalID, taskID, link); err != nil {
		return nil, err
	}
	return s.GetGoal(goalID)
}

func (s *GoalService) DeleteLink(goalID, taskID, linkID string) (*models.Goal, error) {
	if err := s.repo.DeleteLink(goalID, taskID, linkID); err != nil {
		return nil, err
	}
	return s.GetGoal(goalID)
}

func (s *GoalService) AddChecklistItem(goalID, taskID string, req *models.AddChecklistItemRequest) (*models.Goal, error) {
	item := &models.ChecklistItem{
		ID:        uuid.New().String(),
		SubtaskID: taskID,
		Title:     req.Title,
		Completed: false,
	}

	if err := s.repo.AddChecklistItem(goalID, taskID, item); err != nil {
		return nil, err
	}
	return s.GetGoal(goalID)
}

func (s *GoalService) UpdateChecklistItem(goalID, taskID, itemID string, completed bool) (*models.Goal, error) {
	if err := s.repo.UpdateChecklistItem(goalID, taskID, itemID, completed); err != nil {
		return nil, err
	}
	return s.GetGoal(goalID)
}

func (s *GoalService) DeleteChecklistItem(goalID, taskID, itemID string) (*models.Goal, error) {
	if err := s.repo.DeleteChecklistItem(goalID, taskID, itemID); err != nil {
		return nil, err
	}
	return s.GetGoal(goalID)
}
