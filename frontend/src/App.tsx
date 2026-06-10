import { useEffect, useState } from 'react'
import {
  addSubtask,
  completeSubtask,
  createGoal,
  fetchGoals,
  reorderSubtasks,
  updateSubtaskProgress,
} from './services/api'
import type {
  Goal,
  CreateGoalPayload,
  CreateSubtaskPayload,
} from './types'

import GoalForm from './components/GoalForm'
import Timeline from './components/Timeline'

function App() {
  const [goals, setGoals] = useState<Goal[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchGoals()
      .then(setGoals)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

  const handleCreateGoal = async (
    payload: CreateGoalPayload,
  ) => {
    setError(null)

    try {
      const created = await createGoal(payload)

      setGoals((current) => [
        created,
        ...current,
      ])
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : 'Failed to create goal',
      )
    }
  }

  const handleAddSubtask = async (
    goalId: string,
    payload: CreateSubtaskPayload,
  ) => {
    setError(null)

    try {
      const updatedGoal =
        await addSubtask(goalId, payload)

      setGoals((current) =>
        current.map((goal) =>
          goal.id === goalId
            ? updatedGoal
            : goal,
        ),
      )
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : 'Failed to add subtask',
      )
    }
  }

  const replaceGoal = (updatedGoal: Goal) => {
    setGoals((current) =>
      current.map((goal) =>
        goal.id === updatedGoal.id
          ? updatedGoal
          : goal,
      ),
    )
  }

  const handleProgressChange = async (
    subtaskId: string,
    progress: number,
  ) => {
    setError(null)

    try {
      replaceGoal(
        await updateSubtaskProgress(subtaskId, progress),
      )
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : 'Failed to update progress',
      )
    }
  }

  const handleCompleteSubtask = async (
    subtaskId: string,
  ) => {
    setError(null)

    try {
      replaceGoal(await completeSubtask(subtaskId))
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : 'Failed to complete subtask',
      )
    }
  }

  const handleReorderSubtasks = async (
    goalId: string,
    subtaskIds: string[],
  ) => {
    setError(null)

    try {
      replaceGoal(await reorderSubtasks(goalId, subtaskIds))
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : 'Failed to reorder subtasks',
      )
    }
  }

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100">

      {/* Background */}

      <div className="absolute inset-0 -z-10 bg-[radial-gradient(circle_at_top_right,#0ea5e920,transparent_30%),radial-gradient(circle_at_bottom_left,#06b6d420,transparent_30%)]" />

      <div className="mx-auto max-w-7xl px-6 py-10">

        {/* Header */}

        <header className="mb-10">

          <p className="text-sm uppercase tracking-[0.35em] text-sky-400">
            GoalStack
          </p>

          <h1 className="mt-3 text-5xl font-bold">
            Build Goals.
            <br />
            Track Progress.
          </h1>

          <p className="mt-5 max-w-2xl text-slate-400">
            Create weighted goals,
            automatically generate timelines,
            and manage project execution
            with an intuitive planner.
          </p>

        </header>

        {/* Dashboard */}

        <div className="mb-8 grid gap-5 md:grid-cols-3">

          <div className="rounded-3xl border border-slate-800 bg-slate-900 p-5">

            <p className="text-xs uppercase tracking-wider text-slate-500">
              Total Goals
            </p>

            <p className="mt-3 text-4xl font-bold text-white">
              {goals.length}
            </p>

          </div>

          <div className="rounded-3xl border border-slate-800 bg-slate-900 p-5">

            <p className="text-xs uppercase tracking-wider text-slate-500">
              Total Tasks
            </p>

            <p className="mt-3 text-4xl font-bold text-white">
              {goals.reduce(
                (sum, goal) =>
                  sum +
                  goal.subtasks.length,
                0,
              )}
            </p>

          </div>

          <div className="rounded-3xl border border-slate-800 bg-slate-900 p-5">

            <p className="text-xs uppercase tracking-wider text-slate-500">
              Planner Status
            </p>

            <p className="mt-3 text-2xl font-bold text-sky-400">
              Active
            </p>

          </div>

        </div>

        {/* Main Layout */}

        <div className="grid gap-8 xl:grid-cols-[380px_minmax(0,1fr)]">

          {/* Left */}

          <aside>

            <div className="sticky top-6 rounded-3xl border border-slate-800 bg-slate-900/90 p-6">

              <h2 className="text-2xl font-bold">
                New Goal
              </h2>

              <p className="mt-2 text-sm text-slate-400">
                Create a weighted goal
                and start planning.
              </p>

              <GoalForm
                onSubmit={
                  handleCreateGoal
                }
              />

              {error && (
                <div className="mt-5 rounded-2xl border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-300">
                  {error}
                </div>
              )}

            </div>

          </aside>

          {/* Right */}

          <main>

            {loading ? (
              <div className="rounded-3xl border border-slate-800 bg-slate-900 p-10 text-center text-slate-400">
                Loading...
              </div>
            ) : goals.length === 0 ? (
              <div className="rounded-3xl border border-dashed border-slate-700 bg-slate-900/50 p-10 text-center">

                <div className="text-6xl font-bold text-sky-400">
                  +
                </div>

                <h3 className="mt-4 text-2xl font-bold">
                  No goals yet
                </h3>

                <p className="mt-2 text-slate-400">
                  Create your first
                  goal to begin.
                </p>

              </div>
            ) : (
              <div className="space-y-8">

                {goals.map((goal) => (
                  <Timeline
                    key={goal.id}
                    goal={goal}
                    onAddSubtask={
                      handleAddSubtask
                    }
                    onProgressChange={
                      handleProgressChange
                    }
                    onCompleteSubtask={
                      handleCompleteSubtask
                    }
                    onReorderSubtasks={
                      handleReorderSubtasks
                    }
                  />
                ))}

              </div>
            )}

          </main>

        </div>

      </div>
    </div>
  )
}

export default App
