# GoalStack — Goal Planning App

GoalStack helps users plan time-bound goals by splitting a goal into weighted subtasks and generating a timeline with per-subtask allocations.

## Core concept
 - Create a Goal with: title, description, start date, duration (in days or hours).
 - Add Subtasks to a Goal. Each subtask has: title, weight (positive number), optional min/max time, notes, links, and attachments.
 - The system normalizes weights and allocates the goal's total duration proportionally to each subtask.

## Allocation algorithm (simple)
1. Convert goal duration to a base unit (hours or minutes) depending on user preference.
2. Sum weights of all subtasks: W = sum(w_i).
3. For each subtask i: allocated_i = total_duration * (w_i / W).
4. If min/max constraints exist, clamp allocated_i and re-distribute remaining time among unconstrained subtasks.

Edge cases: if all weights are zero, distribute uniformly; if total duration is shorter than sum of mins, mark conflict and surface to user.

## Timeline generation
 - Sort subtasks by user-defined order or priority.
 - Starting from goal.start_date, assign every subtask a start and end timestamp using its allocated duration.
 - Support day/hour granularity and respect business hours (optional setting).

Example timeline entry (JSON):

```
{
	"goal_id": "g1",
	"subtask_id": "s1",
	"title": "Research",
	"start": "2026-06-10T09:00:00Z",
	"end": "2026-06-12T17:00:00Z",
	"allocated_hours": 24,
	"notes": "Read papers",
	"links": ["https://example.com"]
}
```

## UI
 - Goal editor: create/edit goal start date and duration units.
 - Subtask list: add/edit weight, min/max time, order, attachments.
 - Timeline view: horizontal or vertical timeline showing each subtask as an expandable card. Card shows title, allocated time, start/end, and can be expanded to show notes, links, and attachments. Cards are draggable to reorder (which updates timeline).

Interactions:
 - Click card to expand and show a detail pane with an editor for notes and link attachments.
 - Attachments stored as links or file metadata; actual upload handled separately.
 - Conflict UI: highlight subtasks whose constraints cannot be satisfied and provide suggestions (reduce mins, extend duration, or reweight).

## Data model (suggested)

Goal:
 - id, title, description, start_date, duration_value, duration_unit (days/hours), timezone, settings

Subtask:
 - id, goal_id, title, weight, min_duration, max_duration, order, notes, links[], attachment_meta[]

## Backend APIs (examples)
 - POST /goals — create goal
 - GET /goals/:id — fetch goal with subtasks and timeline
 - POST /goals/:id/subtasks — add subtask
 - PUT /goals/:id/subtasks/:sid — update subtask (weight/order/notes)
 - POST /goals/:id/generate-timeline — compute allocations and return timeline

## Implementation notes
 - Keep allocation logic deterministic and testable. Write unit tests for edge cases (zero weights, constraints, rounding errors).
 - Consider progressive enhancement: compute timeline client-side for immediate feedback, validate on server.
 - Persist attachments separately (S3 or equivalent) and store metadata in subtask records.

## TODO / Roadmap
 - Business-hours-aware scheduling
 - Support recurring tasks and milestones
 - Sync with calendar providers (Google/Outlook)
 - Timeline zoom and filtering by tag/status


