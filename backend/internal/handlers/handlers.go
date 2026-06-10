package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"goalstack/internal/models"
	"goalstack/internal/storage"
	"goalstack/internal/timeline"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	store storage.Store
}

// NewHandler creates a new handler
func NewHandler(store storage.Store) *Handler {
	return &Handler{
		store: store,
	}
}

// CreateGoal handles POST /goals
func (h *Handler) CreateGoal(c *gin.Context) {
	var req models.CreateGoalRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal := &models.Goal{
		ID:            uuid.New().String(),
		Title:         req.Title,
		StartDate:     req.StartDate.Time,
		TotalDuration: req.TotalDuration,
		DurationType:  req.DurationType,
	}

	if err := h.store.CreateGoal(goal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate timeline for consistency with other endpoints
	goalWithTimeline := timeline.GetTimelineWithCalculations(goal)

	c.JSON(http.StatusCreated, goalWithTimeline)
}

// GetGoal handles GET /goals/:id
func (h *Handler) GetGoal(c *gin.Context) {
	id := c.Param("id")

	goal, err := h.store.GetGoal(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Calculate timeline
	goalWithTimeline := timeline.GetTimelineWithCalculations(goal)

	c.JSON(http.StatusOK, goalWithTimeline)
}

// UpdateGoal handles PUT /goals/:id
func (h *Handler) UpdateGoal(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateGoalRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal, err := h.store.UpdateGoal(id, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Calculate timeline after update
	goalWithTimeline := timeline.GetTimelineWithCalculations(goal)

	c.JSON(http.StatusOK, goalWithTimeline)
}

// DeleteGoal handles DELETE /goals/:id
func (h *Handler) DeleteGoal(c *gin.Context) {
	id := c.Param("id")

	if err := h.store.DeleteGoal(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "goal deleted"})
}

// ListGoals handles GET /goals
func (h *Handler) ListGoals(c *gin.Context) {
	goals, err := h.store.ListGoals()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate timeline for each goal
	for i := range goals {
		calculated := timeline.GetTimelineWithCalculations(&goals[i])
		goals[i] = *calculated
	}

	if len(goals) == 0 {
		c.JSON(http.StatusOK, []models.Goal{})
		return
	}

	c.JSON(http.StatusOK, goals)
}

// AddSubtask handles POST /goals/:id/subtasks
func (h *Handler) AddSubtask(c *gin.Context) {
	goalID := c.Param("id")
	var req models.CreateSubtaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subtask := &models.Subtask{
		ID:     uuid.New().String(),
		GoalID: goalID,
		Title:  req.Title,
		Weight: req.Weight,
	}

	if err := h.store.AddSubtask(goalID, subtask); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	goal, _ := h.store.GetGoal(goalID)
	goalWithTimeline := timeline.GetTimelineWithCalculations(goal)

	c.JSON(http.StatusCreated, goalWithTimeline)
}

// UpdateSubtask handles PUT /goals/:id/subtasks/:taskId
func (h *Handler) UpdateSubtask(c *gin.Context) {
	goalID := c.Param("id")
	taskID := c.Param("taskId")
	var req models.UpdateSubtaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.store.UpdateSubtask(goalID, taskID, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	goal, _ := h.store.GetGoal(goalID)
	goalWithTimeline := timeline.GetTimelineWithCalculations(goal)

	c.JSON(http.StatusOK, goalWithTimeline)
}

// UpdateSubtaskProgress handles PATCH /subtasks/:taskId/progress
func (h *Handler) UpdateSubtaskProgress(c *gin.Context) {
	taskID := c.Param("taskId")
	var req models.UpdateSubtaskProgressRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal, err := h.store.UpdateSubtaskProgress(taskID, req.Progress)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	goalWithTimeline := timeline.GetTimelineWithCalculations(goal)

	c.JSON(http.StatusOK, goalWithTimeline)
}

// CompleteSubtask handles PATCH /subtasks/:taskId/complete
func (h *Handler) CompleteSubtask(c *gin.Context) {
	taskID := c.Param("taskId")

	goal, err := h.store.CompleteSubtask(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	goalWithTimeline := timeline.GetTimelineWithCalculations(goal)

	c.JSON(http.StatusOK, goalWithTimeline)
}

// ReorderSubtasks handles PATCH /goals/:id/subtasks/reorder
func (h *Handler) ReorderSubtasks(c *gin.Context) {
	goalID := c.Param("id")
	var req models.ReorderSubtasksRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal, err := h.store.ReorderSubtasks(goalID, req.SubtaskIDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goalWithTimeline := timeline.GetTimelineWithCalculations(goal)

	c.JSON(http.StatusOK, goalWithTimeline)
}

// DeleteSubtask handles DELETE /goals/:id/subtasks/:taskId
func (h *Handler) DeleteSubtask(c *gin.Context) {
	goalID := c.Param("id")
	taskID := c.Param("taskId")

	if err := h.store.DeleteSubtask(goalID, taskID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	goal, _ := h.store.GetGoal(goalID)
	goalWithTimeline := timeline.GetTimelineWithCalculations(goal)

	c.JSON(http.StatusOK, goalWithTimeline)
}

// AddNote handles POST /goals/:id/subtasks/:taskId/notes
func (h *Handler) AddNote(c *gin.Context) {
	goalID := c.Param("id")
	taskID := c.Param("taskId")
	var req models.AddNoteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note := &models.Note{
		ID:      uuid.New().String(),
		Content: req.Content,
	}

	if err := h.store.AddNote(goalID, taskID, note); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	goal, _ := h.store.GetGoal(goalID)
	goalWithTimeline := timeline.GetTimelineWithCalculations(goal)

	c.JSON(http.StatusCreated, goalWithTimeline)
}

// DeleteNote handles DELETE /goals/:id/subtasks/:taskId/notes/:noteId
func (h *Handler) DeleteNote(c *gin.Context) {
	goalID := c.Param("id")
	taskID := c.Param("taskId")
	noteID := c.Param("noteId")

	if err := h.store.DeleteNote(goalID, taskID, noteID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	goal, _ := h.store.GetGoal(goalID)
	goalWithTimeline := timeline.GetTimelineWithCalculations(goal)

	c.JSON(http.StatusOK, goalWithTimeline)
}

// AddLink handles POST /goals/:id/subtasks/:taskId/links
func (h *Handler) AddLink(c *gin.Context) {
	goalID := c.Param("id")
	taskID := c.Param("taskId")
	var req models.AddLinkRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	link := &models.Link{
		ID:    uuid.New().String(),
		Title: req.Title,
		URL:   req.URL,
	}

	if err := h.store.AddLink(goalID, taskID, link); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	goal, _ := h.store.GetGoal(goalID)
	goalWithTimeline := timeline.GetTimelineWithCalculations(goal)

	c.JSON(http.StatusCreated, goalWithTimeline)
}

// DeleteLink handles DELETE /goals/:id/subtasks/:taskId/links/:linkId
func (h *Handler) DeleteLink(c *gin.Context) {
	goalID := c.Param("id")
	taskID := c.Param("taskId")
	linkID := c.Param("linkId")

	if err := h.store.DeleteLink(goalID, taskID, linkID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	goal, _ := h.store.GetGoal(goalID)
	goalWithTimeline := timeline.GetTimelineWithCalculations(goal)

	c.JSON(http.StatusOK, goalWithTimeline)
}

// AddChecklistItem handles POST /goals/:id/subtasks/:taskId/checklist
func (h *Handler) AddChecklistItem(c *gin.Context) {
	goalID := c.Param("id")
	taskID := c.Param("taskId")
	var req models.AddChecklistItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := &models.ChecklistItem{
		ID:        uuid.New().String(),
		Title:     req.Title,
		Completed: false,
	}

	if err := h.store.AddChecklistItem(goalID, taskID, item); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	goal, _ := h.store.GetGoal(goalID)
	goalWithTimeline := timeline.GetTimelineWithCalculations(goal)

	c.JSON(http.StatusCreated, goalWithTimeline)
}

// UpdateChecklistItem handles PATCH /goals/:id/subtasks/:taskId/checklist/:itemId
func (h *Handler) UpdateChecklistItem(c *gin.Context) {
	goalID := c.Param("id")
	taskID := c.Param("taskId")
	itemID := c.Param("itemId")

	var req struct {
		Completed bool `json:"completed"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.store.UpdateChecklistItem(goalID, taskID, itemID, req.Completed); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	goal, _ := h.store.GetGoal(goalID)
	goalWithTimeline := timeline.GetTimelineWithCalculations(goal)

	c.JSON(http.StatusOK, goalWithTimeline)
}

// DeleteChecklistItem handles DELETE /goals/:id/subtasks/:taskId/checklist/:itemId
func (h *Handler) DeleteChecklistItem(c *gin.Context) {
	goalID := c.Param("id")
	taskID := c.Param("taskId")
	itemID := c.Param("itemId")

	if err := h.store.DeleteChecklistItem(goalID, taskID, itemID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	goal, _ := h.store.GetGoal(goalID)
	goalWithTimeline := timeline.GetTimelineWithCalculations(goal)

	c.JSON(http.StatusOK, goalWithTimeline)
}
