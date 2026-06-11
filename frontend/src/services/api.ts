import type { CreateGoalPayload, CreateSubtaskPayload, Goal } from '../types'

const BASE_URL = 'http://localhost:8080'

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${BASE_URL}${path}`, {
    headers: {
      'Content-Type': 'application/json',
    },
    ...options,
  })

  if (!response.ok) {
    const errorText = await response.text()
    throw new Error(errorText || 'API request failed')
  }

  return response.json()
}

export function fetchGoals() {
  return request<Goal[]>('/goals')
}

export function createGoal(payload: CreateGoalPayload) {
  const startDate = payload.startDate.includes('T') ? payload.startDate : `${payload.startDate}T00:00:00Z`

  return request<Goal>('/goals', {
    method: 'POST',
    body: JSON.stringify({ ...payload, startDate }),
  })
}

export function addSubtask(goalId: string, payload: CreateSubtaskPayload) {
  return request<Goal>(`/goals/${goalId}/subtasks`, {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function updateSubtaskProgress(subtaskId: string, progress: number) {
  return request<Goal>(`/subtasks/${subtaskId}/progress`, {
    method: 'PATCH',
    body: JSON.stringify({ progress }),
  })
}

export function completeSubtask(subtaskId: string) {
  return request<Goal>(`/subtasks/${subtaskId}/complete`, {
    method: 'PATCH',
  })
}

export function reorderSubtasks(goalId: string, subtaskIds: string[]) {
  return request<Goal>(`/goals/${goalId}/subtasks/reorder`, {
    method: 'PATCH',
    body: JSON.stringify({ subtaskIds }),
  })
}

// Notes
export function addNote(goalId: string, subtaskId: string, content: string) {
  return request<Goal>(`/goals/${goalId}/subtasks/${subtaskId}/notes`, {
    method: 'POST',
    body: JSON.stringify({ content }),
  })
}

export function deleteNote(goalId: string, subtaskId: string, noteId: string) {
  return request<Goal>(`/goals/${goalId}/subtasks/${subtaskId}/notes/${noteId}`, {
    method: 'DELETE',
  })
}

// Links
export function addLink(
  goalId: string,
  subtaskId: string,
  payload: { title: string; url: string },
) {
  return request<Goal>(`/goals/${goalId}/subtasks/${subtaskId}/links`, {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function deleteLink(goalId: string, subtaskId: string, linkId: string) {
  return request<Goal>(`/goals/${goalId}/subtasks/${subtaskId}/links/${linkId}`, {
    method: 'DELETE',
  })
}

// Checklist
export function addChecklistItem(goalId: string, subtaskId: string, title: string) {
  return request<Goal>(`/goals/${goalId}/subtasks/${subtaskId}/checklist`, {
    method: 'POST',
    body: JSON.stringify({ title }),
  })
}

export function updateChecklistItem(
  goalId: string,
  subtaskId: string,
  itemId: string,
  completed: boolean,
) {
  return request<Goal>(
    `/goals/${goalId}/subtasks/${subtaskId}/checklist/${itemId}`,
    {
      method: 'PATCH',
      body: JSON.stringify({ completed }),
    },
  )
}

export function deleteChecklistItem(
  goalId: string,
  subtaskId: string,
  itemId: string,
) {
  return request<Goal>(
    `/goals/${goalId}/subtasks/${subtaskId}/checklist/${itemId}`,
    {
      method: 'DELETE',
    },
  )
}
