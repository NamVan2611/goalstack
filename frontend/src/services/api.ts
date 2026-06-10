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
