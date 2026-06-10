package repository

import (
	"errors"
	"fmt"
	"time"

	"goalstack/internal/models"

	"gorm.io/gorm"
)

// GoalRepository defines persistence operations for goals and their children.
type GoalRepository interface {
	CreateGoal(goal *models.Goal) error
	GetGoal(id string) (*models.Goal, error)
	UpdateGoal(id string, updates *models.UpdateGoalRequest) (*models.Goal, error)
	DeleteGoal(id string) error
	ListGoals() ([]models.Goal, error)

	AddSubtask(goalID string, subtask *models.Subtask) error
	GetSubtask(goalID, subtaskID string) (*models.Subtask, error)
	UpdateSubtask(goalID, subtaskID string, updates *models.UpdateSubtaskRequest) (*models.Subtask, error)
	UpdateSubtaskProgress(subtaskID string, progress float64) (*models.Goal, error)
	ReorderSubtasks(goalID string, subtaskIDs []string) (*models.Goal, error)
	DeleteSubtask(goalID, subtaskID string) error

	AddNote(goalID, subtaskID string, note *models.Note) error
	DeleteNote(goalID, subtaskID, noteID string) error

	AddLink(goalID, subtaskID string, link *models.Link) error
	DeleteLink(goalID, subtaskID, linkID string) error

	AddChecklistItem(goalID, subtaskID string, item *models.ChecklistItem) error
	UpdateChecklistItem(goalID, subtaskID, itemID string, completed bool) error
	DeleteChecklistItem(goalID, subtaskID, itemID string) error
}

type GormGoalRepository struct {
	db *gorm.DB
}

func NewGormGoalRepository(db *gorm.DB) *GormGoalRepository {
	return &GormGoalRepository{db: db}
}

func (r *GormGoalRepository) CreateGoal(goal *models.Goal) error {
	return r.db.Create(goal).Error
}

func (r *GormGoalRepository) GetGoal(id string) (*models.Goal, error) {
	var goal models.Goal
	err := r.withGoalPreloads(r.db).
		First(&goal, "id = ?", id).
		Error
	if err != nil {
		return nil, notFound(err, "goal not found")
	}
	return &goal, nil
}

func (r *GormGoalRepository) UpdateGoal(id string, updates *models.UpdateGoalRequest) (*models.Goal, error) {
	changes := map[string]any{"updated_at": time.Now()}

	if updates.Title != nil {
		changes["title"] = *updates.Title
	}
	if updates.StartDate != nil {
		changes["start_date"] = updates.StartDate.Time
	}
	if updates.TotalDuration != nil {
		changes["total_duration"] = *updates.TotalDuration
	}
	if updates.DurationType != nil {
		changes["duration_type"] = models.DurationType(*updates.DurationType)
	}

	result := r.db.Model(&models.Goal{}).Where("id = ?", id).Updates(changes)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("goal not found")
	}

	return r.GetGoal(id)
}

func (r *GormGoalRepository) DeleteGoal(id string) error {
	result := r.db.Delete(&models.Goal{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("goal not found")
	}
	return nil
}

func (r *GormGoalRepository) ListGoals() ([]models.Goal, error) {
	var goals []models.Goal
	err := r.withGoalPreloads(r.db).
		Order("created_at DESC").
		Find(&goals).
		Error
	return goals, err
}

func (r *GormGoalRepository) AddSubtask(goalID string, subtask *models.Subtask) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := ensureGoalExists(tx, goalID); err != nil {
			return err
		}

		if subtask.Order == 0 {
			var count int64
			if err := tx.Model(&models.Subtask{}).Where("goal_id = ?", goalID).Count(&count).Error; err != nil {
				return err
			}
			subtask.Order = int(count) + 1
		}

		if err := tx.Create(subtask).Error; err != nil {
			return err
		}

		return touchGoal(tx, goalID)
	})
}

func (r *GormGoalRepository) GetSubtask(goalID, subtaskID string) (*models.Subtask, error) {
	var subtask models.Subtask
	err := r.withSubtaskPreloads(r.db).
		First(&subtask, "id = ? AND goal_id = ?", subtaskID, goalID).
		Error
	if err != nil {
		return nil, notFound(err, "subtask not found")
	}
	return &subtask, nil
}

func (r *GormGoalRepository) UpdateSubtask(goalID, subtaskID string, updates *models.UpdateSubtaskRequest) (*models.Subtask, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := ensureGoalExists(tx, goalID); err != nil {
			return err
		}

		changes := map[string]any{"updated_at": time.Now()}
		if updates.Title != nil {
			changes["title"] = *updates.Title
		}
		if updates.Weight != nil {
			changes["weight"] = *updates.Weight
		}
		if updates.Order != nil {
			changes["sort_order"] = *updates.Order
		}

		result := tx.Model(&models.Subtask{}).
			Where("id = ? AND goal_id = ?", subtaskID, goalID).
			Updates(changes)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("subtask not found")
		}

		if updates.Order != nil {
			if err := normalizeSubtaskOrders(tx, goalID); err != nil {
				return err
			}
		}

		return touchGoal(tx, goalID)
	})
	if err != nil {
		return nil, err
	}

	return r.GetSubtask(goalID, subtaskID)
}

func (r *GormGoalRepository) UpdateSubtaskProgress(subtaskID string, progress float64) (*models.Goal, error) {
	var goalID string
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var subtask models.Subtask
		if err := tx.First(&subtask, "id = ?", subtaskID).Error; err != nil {
			return notFound(err, "subtask not found")
		}

		goalID = subtask.GoalID
		if err := tx.Model(&subtask).Updates(map[string]any{
			"progress":   progress,
			"status":     statusForProgress(progress),
			"updated_at": time.Now(),
		}).Error; err != nil {
			return err
		}

		return touchGoal(tx, goalID)
	})
	if err != nil {
		return nil, err
	}

	return r.GetGoal(goalID)
}

func (r *GormGoalRepository) ReorderSubtasks(goalID string, subtaskIDs []string) (*models.Goal, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var subtasks []models.Subtask
		if err := tx.Where("goal_id = ?", goalID).Find(&subtasks).Error; err != nil {
			return err
		}
		if len(subtasks) == 0 {
			if err := ensureGoalExists(tx, goalID); err != nil {
				return err
			}
		}
		if len(subtaskIDs) != len(subtasks) {
			return fmt.Errorf("subtask order must include every subtask")
		}

		existing := make(map[string]struct{}, len(subtasks))
		for _, subtask := range subtasks {
			existing[subtask.ID] = struct{}{}
		}

		seen := make(map[string]struct{}, len(subtaskIDs))
		for index, id := range subtaskIDs {
			if _, ok := seen[id]; ok {
				return fmt.Errorf("subtask order contains duplicate IDs")
			}
			if _, ok := existing[id]; !ok {
				return fmt.Errorf("subtask order contains unknown IDs")
			}
			seen[id] = struct{}{}

			if err := tx.Model(&models.Subtask{}).
				Where("id = ? AND goal_id = ?", id, goalID).
				Updates(map[string]any{
					"sort_order": index + 1,
					"updated_at": time.Now(),
				}).Error; err != nil {
				return err
			}
		}

		return touchGoal(tx, goalID)
	})
	if err != nil {
		return nil, err
	}

	return r.GetGoal(goalID)
}

func (r *GormGoalRepository) DeleteSubtask(goalID, subtaskID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&models.Subtask{}, "id = ? AND goal_id = ?", subtaskID, goalID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("subtask not found")
		}
		if err := normalizeSubtaskOrders(tx, goalID); err != nil {
			return err
		}
		return touchGoal(tx, goalID)
	})
}

func (r *GormGoalRepository) AddNote(goalID, subtaskID string, note *models.Note) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := ensureSubtaskExists(tx, goalID, subtaskID); err != nil {
			return err
		}
		note.SubtaskID = subtaskID
		if err := tx.Create(note).Error; err != nil {
			return err
		}
		return touchGoal(tx, goalID)
	})
}

func (r *GormGoalRepository) DeleteNote(goalID, subtaskID, noteID string) error {
	return r.deleteChild(goalID, subtaskID, noteID, &models.Note{}, "note not found")
}

func (r *GormGoalRepository) AddLink(goalID, subtaskID string, link *models.Link) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := ensureSubtaskExists(tx, goalID, subtaskID); err != nil {
			return err
		}
		link.SubtaskID = subtaskID
		if err := tx.Create(link).Error; err != nil {
			return err
		}
		return touchGoal(tx, goalID)
	})
}

func (r *GormGoalRepository) DeleteLink(goalID, subtaskID, linkID string) error {
	return r.deleteChild(goalID, subtaskID, linkID, &models.Link{}, "link not found")
}

func (r *GormGoalRepository) AddChecklistItem(goalID, subtaskID string, item *models.ChecklistItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := ensureSubtaskExists(tx, goalID, subtaskID); err != nil {
			return err
		}
		item.SubtaskID = subtaskID
		if err := tx.Create(item).Error; err != nil {
			return err
		}
		return touchGoal(tx, goalID)
	})
}

func (r *GormGoalRepository) UpdateChecklistItem(goalID, subtaskID, itemID string, completed bool) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := ensureSubtaskExists(tx, goalID, subtaskID); err != nil {
			return err
		}
		result := tx.Model(&models.ChecklistItem{}).
			Where("id = ? AND subtask_id = ?", itemID, subtaskID).
			Update("completed", completed)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("checklist item not found")
		}
		return touchGoal(tx, goalID)
	})
}

func (r *GormGoalRepository) DeleteChecklistItem(goalID, subtaskID, itemID string) error {
	return r.deleteChild(goalID, subtaskID, itemID, &models.ChecklistItem{}, "checklist item not found")
}

func (r *GormGoalRepository) deleteChild(goalID, subtaskID, id string, model any, notFoundMessage string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := ensureSubtaskExists(tx, goalID, subtaskID); err != nil {
			return err
		}
		result := tx.Delete(model, "id = ? AND subtask_id = ?", id, subtaskID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New(notFoundMessage)
		}
		return touchGoal(tx, goalID)
	})
}

func (r *GormGoalRepository) withGoalPreloads(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Subtasks", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Subtasks.Notes", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("Subtasks.Links").
		Preload("Subtasks.ChecklistItems")
}

func (r *GormGoalRepository) withSubtaskPreloads(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Notes", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("Links").
		Preload("ChecklistItems")
}

func ensureGoalExists(db *gorm.DB, goalID string) error {
	var count int64
	if err := db.Model(&models.Goal{}).Where("id = ?", goalID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("goal not found")
	}
	return nil
}

func ensureSubtaskExists(db *gorm.DB, goalID, subtaskID string) error {
	var count int64
	if err := db.Model(&models.Subtask{}).
		Where("id = ? AND goal_id = ?", subtaskID, goalID).
		Count(&count).
		Error; err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("subtask not found")
	}
	return nil
}

func normalizeSubtaskOrders(db *gorm.DB, goalID string) error {
	var subtasks []models.Subtask
	if err := db.Where("goal_id = ?", goalID).Order("sort_order ASC").Find(&subtasks).Error; err != nil {
		return err
	}

	for i, subtask := range subtasks {
		if err := db.Model(&models.Subtask{}).
			Where("id = ?", subtask.ID).
			Update("sort_order", i+1).
			Error; err != nil {
			return err
		}
	}

	return nil
}

func touchGoal(db *gorm.DB, goalID string) error {
	return db.Model(&models.Goal{}).
		Where("id = ?", goalID).
		Update("updated_at", time.Now()).
		Error
}

func notFound(err error, message string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(message)
	}
	return err
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
