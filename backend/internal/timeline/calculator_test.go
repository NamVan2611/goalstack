package timeline

import (
	"testing"

	"goalstack/internal/models"
)

func TestCalculateGoalProgressUsesSubtaskWeights(t *testing.T) {
	progress := CalculateGoalProgress([]models.Subtask{
		{Weight: 40, Progress: 25},
		{Weight: 60, Progress: 100},
	})

	if progress != 70 {
		t.Fatalf("expected weighted progress 70, got %v", progress)
	}
}

func TestCalculateGoalProgressWithNoWeight(t *testing.T) {
	progress := CalculateGoalProgress([]models.Subtask{
		{Weight: 0, Progress: 100},
	})

	if progress != 0 {
		t.Fatalf("expected progress 0 with no weight, got %v", progress)
	}
}
