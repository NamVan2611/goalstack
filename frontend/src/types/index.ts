export type DurationType = 'hours' | 'days'

export interface Link {
  id: string
  title: string
  url: string
}

export interface Note {
  id: string
  content: string
  createdAt: string
}

export interface ChecklistItem {
  id: string
  title: string
  completed: boolean
}

export interface Subtask {
  id: string
  goalId: string
  title: string
  weight: number
  order: number
  notes: Note[]
  links: Link[]
  checklistItems: ChecklistItem[]
  progress: number
  allocatedTime: number
  startDate: string
  endDate: string
}

export interface Goal {
  id: string
  title: string
  startDate: string
  totalDuration: number
  durationType: DurationType
  subtasks: Subtask[]
}

export interface CreateSubtaskPayload {
  title: string
  weight: number
}

export interface CreateGoalPayload {
  title: string
  startDate: string
  totalDuration: number
  durationType: DurationType
}
