# API Test Examples

This file contains examples for testing the GoalStack API endpoints.

## Base URL
```
http://localhost:8080
```

## Health Check
```bash
curl http://localhost:8080/health
```

## Create a Goal

```bash
curl -X POST http://localhost:8080/goals \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Build GoalStack App",
    "startDate": "2026-06-10T00:00:00Z",
    "totalDuration": 100,
    "durationType": "hours"
  }'
```

Response:
```json
{
  "id": "uuid-string",
  "title": "Build GoalStack App",
  "startDate": "2026-06-10T00:00:00Z",
  "totalDuration": 100,
  "durationType": "hours",
  "subtasks": [],
  "createdAt": "2026-06-09T12:00:00Z",
  "updatedAt": "2026-06-09T12:00:00Z"
}
```

## Get All Goals

```bash
curl http://localhost:8080/goals
```

## Get a Specific Goal with Timeline

```bash
curl http://localhost:8080/goals/{goalId}
```

Response includes calculated timeline for each subtask.

## Add a Subtask

```bash
curl -X POST http://localhost:8080/goals/{goalId}/subtasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Backend Development",
    "weight": 40
  }'
```

The response will include the updated goal with all subtasks and their calculated timeline.

## Add More Subtasks

```bash
# Task 2: Frontend
curl -X POST http://localhost:8080/goals/{goalId}/subtasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Frontend Development",
    "weight": 50
  }'

# Task 3: Testing
curl -X POST http://localhost:8080/goals/{goalId}/subtasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Testing & QA",
    "weight": 10
  }'
```

After adding these subtasks, the timeline will show:
- Backend: 40/100 × 100 = 40 hours
- Frontend: 50/100 × 100 = 50 hours
- Testing: 10/100 × 100 = 10 hours

## Add a Note to a Subtask

```bash
curl -X POST http://localhost:8080/goals/{goalId}/subtasks/{subtaskId}/notes \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Need to review API design documentation"
  }'
```

## Add a Link to a Subtask

```bash
curl -X POST http://localhost:8080/goals/{goalId}/subtasks/{subtaskId}/links \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Gin Documentation",
    "url": "https://gin-gonic.com"
  }'
```

## Add a Checklist Item

```bash
curl -X POST http://localhost:8080/goals/{goalId}/subtasks/{subtaskId}/checklist \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Review code"
  }'
```

## Update a Checklist Item

```bash
curl -X PATCH http://localhost:8080/goals/{goalId}/subtasks/{subtaskId}/checklist/{itemId} \
  -H "Content-Type: application/json" \
  -d '{
    "completed": true
  }'
```

## Update a Subtask Weight

```bash
curl -X PUT http://localhost:8080/goals/{goalId}/subtasks/{subtaskId} \
  -H "Content-Type: application/json" \
  -d '{
    "weight": 45
  }'
```

The timeline will automatically recalculate with the new weight.

## Update Goal Duration

```bash
curl -X PUT http://localhost:8080/goals/{goalId} \
  -H "Content-Type: application/json" \
  -d '{
    "totalDuration": 120
  }'
```

The timeline will automatically recalculate with the new duration.

## Delete Operations

```bash
# Delete a goal
curl -X DELETE http://localhost:8080/goals/{goalId}

# Delete a subtask
curl -X DELETE http://localhost:8080/goals/{goalId}/subtasks/{subtaskId}

# Delete a note
curl -X DELETE http://localhost:8080/goals/{goalId}/subtasks/{subtaskId}/notes/{noteId}

# Delete a link
curl -X DELETE http://localhost:8080/goals/{goalId}/subtasks/{subtaskId}/links/{linkId}

# Delete a checklist item
curl -X DELETE http://localhost:8080/goals/{goalId}/subtasks/{subtaskId}/checklist/{itemId}
```

## Timeline Calculation Example

Goal: 100 hours starting 2026-06-10 00:00:00

Subtasks:
- Backend (weight 40): 40 hours
- Frontend (weight 50): 50 hours
- Testing (weight 10): 10 hours

Expected Timeline:
```
Backend:
  Start: 2026-06-10 00:00:00
  End:   2026-06-11 16:00:00
  Duration: 40 hours

Frontend:
  Start: 2026-06-11 16:00:00
  End:   2026-06-13 18:00:00
  Duration: 50 hours

Testing:
  Start: 2026-06-13 18:00:00
  End:   2026-06-14 04:00:00
  Duration: 10 hours
```

## Using Postman or Insomnia

You can import these examples into Postman or Insomnia for easier testing:

1. Create a new collection
2. Add requests using the examples above
3. Replace `{goalId}` and `{subtaskId}` with actual IDs from previous responses
4. Test different weight combinations and durations

## Notes

- All timestamps are in ISO 8601 format with UTC timezone
- IDs are generated as UUIDs (v4)
- The API automatically calculates timeline whenever subtasks or goal duration changes
- All CRUD operations return the updated goal with the recalculated timeline
