package timeline

import (
	"sort"
	"time"

	"goalstack/internal/models"
)

// CalculateTimeline calculates start and end dates for all subtasks based on weights
func CalculateTimeline(goal *models.Goal) {
	goal.Progress = CalculateGoalProgress(goal.Subtasks)

	if len(goal.Subtasks) == 0 {
		return
	}

	// Sort subtasks by order
	sort.Slice(goal.Subtasks, func(i, j int) bool {
		return goal.Subtasks[i].Order < goal.Subtasks[j].Order
	})

	// Calculate total weight
	totalWeight := 0.0
	for _, subtask := range goal.Subtasks {
		totalWeight += subtask.Weight
	}

	if totalWeight == 0 {
		return
	}

	// Convert duration to hours for consistent calculations
	totalHours := goal.TotalDuration
	if goal.DurationType == models.Days {
		totalHours = goal.TotalDuration * 24
	}

	// Calculate start and end dates for each subtask
	currentDate := goal.StartDate

	for i := range goal.Subtasks {
		// Calculate allocated time based on weight proportion
		allocatedHours := (goal.Subtasks[i].Weight / totalWeight) * totalHours

		goal.Subtasks[i].AllocatedTime = allocatedHours
		goal.Subtasks[i].StartDate = currentDate

		// Calculate end date: add hours to start date
		goal.Subtasks[i].EndDate = currentDate.Add(time.Duration(allocatedHours) * time.Hour)

		// Next subtask starts when current one ends
		currentDate = goal.Subtasks[i].EndDate
	}
}

// CalculateGoalProgress calculates weighted progress across all subtasks.
func CalculateGoalProgress(subtasks []models.Subtask) float64 {
	totalWeight := 0.0
	weightedProgress := 0.0

	for _, subtask := range subtasks {
		totalWeight += subtask.Weight
		weightedProgress += subtask.Progress * subtask.Weight
	}

	if totalWeight == 0 {
		return 0
	}

	return weightedProgress / totalWeight
}

// GetTimelineWithCalculations returns a goal with calculated timeline
func GetTimelineWithCalculations(goal *models.Goal) *models.Goal {
	// Make a copy to avoid modifying the original
	goalCopy := *goal
	goalCopy.Subtasks = make([]models.Subtask, len(goal.Subtasks))
	copy(goalCopy.Subtasks, goal.Subtasks)

	CalculateTimeline(&goalCopy)
	return &goalCopy
}
