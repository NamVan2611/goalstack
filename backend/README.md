# GoalStack Backend

A REST API for managing goals and subtasks with automatic timeline calculation.

## Tech Stack

- **Language**: Go 1.21
- **Framework**: Gin
- **Storage**: In-memory (PostgreSQL integration coming soon)

## Project Structure

```
backend/
├── internal/
│   ├── models/        # Data models and request/response structs
│   ├── storage/       # In-memory storage implementation
│   ├── handlers/      # HTTP request handlers
│   └── timeline/      # Timeline calculation logic
├── main.go            # Server entry point
└── go.mod             # Go module definition
```

## Getting Started

### Prerequisites

- Go 1.21 or higher

### Installation

```bash
cd backend
go mod download
go mod tidy
```

### Running the Server

```bash
go run main.go
```

The server will start on `http://localhost:8080`

### Health Check

```bash
curl http://localhost:8080/health
```

## API Endpoints

### Goals

- `POST /goals` - Create a new goal
- `GET /goals` - List all goals
- `GET /goals/:id` - Get a specific goal with timeline
- `PUT /goals/:id` - Update a goal
- `DELETE /goals/:id` - Delete a goal

### Subtasks

- `POST /goals/:id/subtasks` - Add a subtask
- `PUT /goals/:id/subtasks/:taskId` - Update a subtask
- `DELETE /goals/:id/subtasks/:taskId` - Delete a subtask

### Notes

- `POST /goals/:id/subtasks/:taskId/notes` - Add a note
- `DELETE /goals/:id/subtasks/:taskId/notes/:noteId` - Delete a note

### Links

- `POST /goals/:id/subtasks/:taskId/links` - Add a link
- `DELETE /goals/:id/subtasks/:taskId/links/:linkId` - Delete a link

### Checklist Items

- `POST /goals/:id/subtasks/:taskId/checklist` - Add a checklist item
- `PATCH /goals/:id/subtasks/:taskId/checklist/:itemId` - Update a checklist item
- `DELETE /goals/:id/subtasks/:taskId/checklist/:itemId` - Delete a checklist item

## Timeline Calculation

The system automatically calculates start and end dates for each subtask based on:

1. **Total Goal Duration**: Specified in hours or days
2. **Subtask Weights**: Proportional allocation of time
3. **Sequential Execution**: Each subtask starts when the previous one ends

Example:
- Goal duration: 100 hours
- Task A weight: 20 → allocated 20 hours
- Task B weight: 30 → allocated 30 hours
- Task C weight: 50 → allocated 50 hours

## Development

### Project Layout

- All models are in `internal/models/`
- Storage layer is abstracted in `internal/storage/`
- All HTTP handlers are in `internal/handlers/`
- Timeline calculations are in `internal/timeline/`

### Next Steps

- [ ] Integrate PostgreSQL with GORM
- [ ] Add JWT authentication
- [ ] Add input validation middleware
- [ ] Add logging
- [ ] Write unit tests
