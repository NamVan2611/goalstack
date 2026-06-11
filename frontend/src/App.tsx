import { useEffect, useMemo, useState } from 'react'
import {
  addChecklistItem,
  addLink,
  addNote,
  addSubtask,
  completeSubtask,
  createGoal,
  deleteChecklistItem,
  deleteLink,
  deleteNote,
  fetchGoals,
  updateChecklistItem,
  updateSubtaskProgress,
} from './services/api'
import type { Goal, CreateGoalPayload, CreateSubtaskPayload } from './types'

import GoalForm from './components/GoalForm'
import Calendar from './components/Calendar'
import GoalDetails from './components/GoalDetails'
import StatsPanel from './components/StatsPanel'

function App() {
  const [goals, setGoals] = useState<Goal[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [showGoalForm, setShowGoalForm] = useState(false)
  const [selectedGoalId, setSelectedGoalId] = useState<string | null>(null)

  const now = new Date()
  const [year, setYear] = useState(now.getFullYear())
  const [month, setMonth] = useState(now.getMonth())

  useEffect(() => {
    fetchGoals()
      .then((data) => {
        setGoals(data)
        if (data.length > 0) setSelectedGoalId(data[0].id)
      })
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

  const selectedGoal = useMemo(
    () => goals.find((g) => g.id === selectedGoalId) ?? null,
    [goals, selectedGoalId],
  )

  const totalTasks = goals.reduce((sum, g) => sum + g.subtasks.length, 0)
  const avgProgress = goals.length
    ? Math.round(goals.reduce((s, g) => s + (g.progress ?? 0), 0) / goals.length)
    : 0

  const replaceGoal = (updatedGoal: Goal) => {
    setGoals((current) =>
      current.map((goal) => (goal.id === updatedGoal.id ? updatedGoal : goal)),
    )
  }

  const runAction = async (action: () => Promise<Goal>, fallback: string) => {
    setError(null)
    try {
      replaceGoal(await action())
    } catch (err) {
      setError(err instanceof Error ? err.message : fallback)
    }
  }

  const handleCreateGoal = async (payload: CreateGoalPayload) => {
    setError(null)
    try {
      const created = await createGoal(payload)
      setGoals((current) => [created, ...current])
      setSelectedGoalId(created.id)
      setShowGoalForm(false)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create goal')
    }
  }

  const handleAddSubtask = (goalId: string, payload: CreateSubtaskPayload) =>
    runAction(() => addSubtask(goalId, payload), 'Failed to add subtask')

  const handleProgressChange = (subtaskId: string, progress: number) =>
    runAction(
      () => updateSubtaskProgress(subtaskId, progress),
      'Failed to update progress',
    )

  const handleCompleteSubtask = (subtaskId: string) =>
    runAction(() => completeSubtask(subtaskId), 'Failed to complete subtask')

  const handleAddNote = (goalId: string, subtaskId: string, content: string) =>
    runAction(() => addNote(goalId, subtaskId, content), 'Failed to add note')

  const handleDeleteNote = (goalId: string, subtaskId: string, noteId: string) =>
    runAction(() => deleteNote(goalId, subtaskId, noteId), 'Failed to delete note')

  const handleAddLink = (
    goalId: string,
    subtaskId: string,
    payload: { title: string; url: string },
  ) => runAction(() => addLink(goalId, subtaskId, payload), 'Failed to add link')

  const handleDeleteLink = (goalId: string, subtaskId: string, linkId: string) =>
    runAction(() => deleteLink(goalId, subtaskId, linkId), 'Failed to delete link')

  const handleAddChecklistItem = (
    goalId: string,
    subtaskId: string,
    title: string,
  ) =>
    runAction(
      () => addChecklistItem(goalId, subtaskId, title),
      'Failed to add checklist item',
    )

  const handleToggleChecklistItem = (
    goalId: string,
    subtaskId: string,
    itemId: string,
    completed: boolean,
  ) =>
    runAction(
      () => updateChecklistItem(goalId, subtaskId, itemId, completed),
      'Failed to update checklist item',
    )

  const handleDeleteChecklistItem = (
    goalId: string,
    subtaskId: string,
    itemId: string,
  ) =>
    runAction(
      () => deleteChecklistItem(goalId, subtaskId, itemId),
      'Failed to delete checklist item',
    )

  const goPrevMonth = () => {
    setMonth((m) => {
      if (m === 0) {
        setYear((y) => y - 1)
        return 11
      }
      return m - 1
    })
  }

  const goNextMonth = () => {
    setMonth((m) => {
      if (m === 11) {
        setYear((y) => y + 1)
        return 0
      }
      return m + 1
    })
  }

  const goToday = () => {
    const today = new Date()
    setYear(today.getFullYear())
    setMonth(today.getMonth())
  }

  return (
    <div className="min-h-screen bg-slate-100">
      {/* Top navigation */}
      <header className="sticky top-0 z-20 border-b border-slate-200 bg-white/80 backdrop-blur">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
          <div className="flex items-center gap-2.5">
            <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-indigo-600 text-sm font-bold text-white">
              GS
            </div>
            <span className="text-lg font-bold text-slate-900">GoalStack</span>
          </div>

          <div className="flex items-center gap-3">
            <div className="hidden items-center gap-4 text-sm sm:flex">
              <span className="text-slate-500">
                <span className="font-bold text-slate-900">{goals.length}</span> goals
              </span>
              <span className="text-slate-500">
                <span className="font-bold text-slate-900">{totalTasks}</span> tasks
              </span>
              <span className="text-slate-500">
                <span className="font-bold text-indigo-600">{avgProgress}%</span> avg
              </span>
            </div>
            <button
              type="button"
              onClick={() => setShowGoalForm((v) => !v)}
              className="rounded-xl bg-indigo-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-indigo-500"
            >
              + New Goal
            </button>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-7xl px-6 py-8">
        {error && (
          <div className="mb-5 rounded-2xl border border-red-200 bg-red-50 p-4 text-sm text-red-600">
            {error}
          </div>
        )}

        {showGoalForm && (
          <div className="mb-6 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <div className="mb-4 flex items-center justify-between">
              <h2 className="text-lg font-bold text-slate-900">New Goal</h2>
              <button
                type="button"
                onClick={() => setShowGoalForm(false)}
                className="text-sm text-slate-400 hover:text-slate-700"
              >
                Close
              </button>
            </div>
            <GoalForm onSubmit={handleCreateGoal} />
          </div>
        )}

        {loading ? (
          <div className="rounded-3xl border border-slate-200 bg-white p-12 text-center text-slate-400 shadow-sm">
            Loading...
          </div>
        ) : (
          <div className="space-y-6">
            {/* Calendar + details */}
            <div className="grid gap-6 xl:grid-cols-[minmax(0,1fr)_380px]">
              <Calendar
                year={year}
                month={month}
                goals={goals}
                selectedGoalId={selectedGoalId}
                onPrevMonth={goPrevMonth}
                onNextMonth={goNextMonth}
                onToday={goToday}
                onSelectGoal={setSelectedGoalId}
              />

              <GoalDetails
                  goal={selectedGoal}
                  onAddSubtask={handleAddSubtask}
                  onProgressChange={handleProgressChange}
                  onCompleteSubtask={handleCompleteSubtask}
                  onAddNote={handleAddNote}
                  onDeleteNote={handleDeleteNote}
                  onAddLink={handleAddLink}
                  onDeleteLink={handleDeleteLink}
                  onAddChecklistItem={handleAddChecklistItem}
                  onToggleChecklistItem={handleToggleChecklistItem}
                  onDeleteChecklistItem={handleDeleteChecklistItem} onDeleteSubtask={function (goalId: string, subtaskId: string): void {
                    throw new Error('Function not implemented.')
                  } }              />
            </div>

            {/* Stats */}
            <StatsPanel goals={goals} />
          </div>
        )}
      </main>
    </div>
  )
}

export default App
