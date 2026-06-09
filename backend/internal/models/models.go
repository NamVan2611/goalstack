package models

import "time"

// DurationType defines whether duration is in hours or days
type DurationType string

const (
	Hours DurationType = "hours"
	Days  DurationType = "days"
)

// Link represents a reference link for a subtask
type Link struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

// Note represents a note attached to a subtask
type Note struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

// ChecklistItem represents a checklist item for a subtask
type ChecklistItem struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// Subtask represents a subtask within a goal
type Subtask struct {
	ID             string          `json:"id"`
	GoalID         string          `json:"goalId"`
	Title          string          `json:"title"`
	Weight         float64         `json:"weight"`
	Order          int             `json:"order"`
	Notes          []Note          `json:"notes"`
	Links          []Link          `json:"links"`
	ChecklistItems []ChecklistItem `json:"checklistItems"`
	Progress       float64         `json:"progress"` // 0-100
	AllocatedTime  float64         `json:"allocatedTime"`
	StartDate      time.Time       `json:"startDate"`
	EndDate        time.Time       `json:"endDate"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

// Goal represents a user's goal
type Goal struct {
	ID            string       `json:"id"`
	Title         string       `json:"title"`
	StartDate     time.Time    `json:"startDate"`
	TotalDuration float64      `json:"totalDuration"`
	DurationType  DurationType `json:"durationType"` // "hours" or "days"
	Subtasks      []Subtask    `json:"subtasks"`
	CreatedAt     time.Time    `json:"createdAt"`
	UpdatedAt     time.Time    `json:"updatedAt"`
}

// CreateGoalRequest represents the request body for creating a goal
type CreateGoalRequest struct {
	Title         string       `json:"title" binding:"required"`
	StartDate     time.Time    `json:"startDate" binding:"required"`
	TotalDuration float64      `json:"totalDuration" binding:"required,gt=0"`
	DurationType  DurationType `json:"durationType" binding:"required,oneof=hours days"`
}

// UpdateGoalRequest represents the request body for updating a goal
type UpdateGoalRequest struct {
	Title         *string  `json:"title"`
	TotalDuration *float64 `json:"totalDuration" binding:"omitempty,gt=0"`
	DurationType  *string  `json:"durationType" binding:"omitempty,oneof=hours days"`
}

// CreateSubtaskRequest represents the request body for creating a subtask
type CreateSubtaskRequest struct {
	Title  string  `json:"title" binding:"required"`
	Weight float64 `json:"weight" binding:"required,gt=0"`
}

// UpdateSubtaskRequest represents the request body for updating a subtask
type UpdateSubtaskRequest struct {
	Title  *string  `json:"title"`
	Weight *float64 `json:"weight" binding:"omitempty,gt=0"`
	Order  *int     `json:"order"`
}

// AddNoteRequest represents the request body for adding a note
type AddNoteRequest struct {
	Content string `json:"content" binding:"required"`
}

// AddLinkRequest represents the request body for adding a link
type AddLinkRequest struct {
	Title string `json:"title" binding:"required"`
	URL   string `json:"url" binding:"required,url"`
}

// AddChecklistItemRequest represents the request body for adding a checklist item
type AddChecklistItemRequest struct {
	Title string `json:"title" binding:"required"`
}
