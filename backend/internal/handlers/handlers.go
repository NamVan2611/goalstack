package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"goalstack/internal/models"
	"goalstack/internal/service"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	service *service.GoalService
}

// NewHandler creates a new handler
func NewHandler(service *service.GoalService) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateGoal handles POST /goals
func (h *Handler) CreateGoal(c *gin.Context) {
	var req models.CreateGoalRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal, err := h.service.CreateGoal(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, goal)
}

// GetGoal handles GET /goals/:id
func (h *Handler) GetGoal(c *gin.Context) {
	id := c.Param("id")

	goal, err := h.service.GetGoal(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goal)
}

// UpdateGoal handles PUT /goals/:id
func (h *Handler) UpdateGoal(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateGoalRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal, err := h.service.UpdateGoal(id, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goal)
}

// DeleteGoal handles DELETE /goals/:id
func (h *Handler) DeleteGoal(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteGoal(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "goal deleted"})
}

// ListGoals handles GET /goals
func (h *Handler) ListGoals(c *gin.Context) {
	goals, err := h.service.ListGoals()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	goal, err := h.service.AddSubtask(goalID, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, goal)
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

	goal, err := h.service.UpdateSubtask(goalID, taskID, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goal)
}

// UpdateSubtaskProgress handles PATCH /subtasks/:taskId/progress
func (h *Handler) UpdateSubtaskProgress(c *gin.Context) {
	taskID := c.Param("taskId")
	var req models.UpdateSubtaskProgressRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal, err := h.service.UpdateSubtaskProgress(taskID, req.Progress)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goal)
}

// CompleteSubtask handles PATCH /subtasks/:taskId/complete
func (h *Handler) CompleteSubtask(c *gin.Context) {
	taskID := c.Param("taskId")

	goal, err := h.service.CompleteSubtask(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goal)
}

// ReorderSubtasks handles PATCH /goals/:id/subtasks/reorder
func (h *Handler) ReorderSubtasks(c *gin.Context) {
	goalID := c.Param("id")
	var req models.ReorderSubtasksRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal, err := h.service.ReorderSubtasks(goalID, req.SubtaskIDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goal)
}

// DeleteSubtask handles DELETE /goals/:id/subtasks/:taskId
func (h *Handler) DeleteSubtask(c *gin.Context) {
	goalID := c.Param("id")
	taskID := c.Param("taskId")

	goal, err := h.service.DeleteSubtask(goalID, taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goal)
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

	goal, err := h.service.AddNote(goalID, taskID, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, goal)
}

// DeleteNote handles DELETE /goals/:id/subtasks/:taskId/notes/:noteId
func (h *Handler) DeleteNote(c *gin.Context) {
	goalID := c.Param("id")
	taskID := c.Param("taskId")
	noteID := c.Param("noteId")

	goal, err := h.service.DeleteNote(goalID, taskID, noteID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goal)
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

	goal, err := h.service.AddLink(goalID, taskID, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, goal)
}

// DeleteLink handles DELETE /goals/:id/subtasks/:taskId/links/:linkId
func (h *Handler) DeleteLink(c *gin.Context) {
	goalID := c.Param("id")
	taskID := c.Param("taskId")
	linkID := c.Param("linkId")

	goal, err := h.service.DeleteLink(goalID, taskID, linkID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goal)
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

	goal, err := h.service.AddChecklistItem(goalID, taskID, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, goal)
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

	goal, err := h.service.UpdateChecklistItem(goalID, taskID, itemID, req.Completed)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goal)
}

// DeleteChecklistItem handles DELETE /goals/:id/subtasks/:taskId/checklist/:itemId
func (h *Handler) DeleteChecklistItem(c *gin.Context) {
	goalID := c.Param("id")
	taskID := c.Param("taskId")
	itemID := c.Param("itemId")

	goal, err := h.service.DeleteChecklistItem(goalID, taskID, itemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goal)
}
