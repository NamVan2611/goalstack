import { useEffect, useState } from 'react'
import { addSubtask, createGoal, fetchGoals } from './services/api'
import type { Goal, CreateGoalPayload, CreateSubtaskPayload } from './types'
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

  const handleCreateGoal = async (payload: CreateGoalPayload) => {
    setError(null)

    try {
      const created = await createGoal(payload)
      setGoals((current) => [created, ...current])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create goal')
    }
  }

  const handleAddSubtask = async (goalId: string, payload: CreateSubtaskPayload) => {
    setError(null)

    try {
      const updatedGoal = await addSubtask(goalId, payload)
      setGoals((current) => current.map((goal) => (goal.id === goalId ? updatedGoal : goal)))
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to add subtask')
    }
  }

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 px-4 py-8 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-6xl">
        <header className="mb-8 flex flex-col gap-4">
          <div>
            <p className="text-sm uppercase tracking-[0.3em] text-sky-400">GoalStack</p>
            <h1 className="text-3xl font-semibold sm:text-4xl">Weighted goal planning</h1>
          </div>
          <p className="max-w-2xl text-slate-300">
            Create a goal, assign subtasks by weight, and see the computed timeline for sequential execution.
          </p>
        </header>

        <section className="grid gap-8 lg:grid-cols-[360px_minmax(0,1fr)]">
          <div className="rounded-3xl border border-slate-800 bg-slate-900/90 p-6 shadow-lg shadow-slate-950/20">
            <h2 className="text-xl font-semibold">New goal</h2>
            <p className="mt-2 text-sm text-slate-400">Define title, start date, duration, and unit.</p>
            <GoalForm onSubmit={handleCreateGoal} />
            {error ? <p className="mt-4 rounded-lg bg-rose-500/10 px-4 py-3 text-sm text-rose-300">{error}</p> : null}
          </div>

          <div className="space-y-6">
            <div className="rounded-3xl border border-slate-800 bg-slate-900/90 p-6 shadow-lg shadow-slate-950/20">
              <div className="flex items-start justify-between gap-4">
                <div>
                  <h2 className="text-xl font-semibold">Goals</h2>
                  <p className="mt-1 text-sm text-slate-400">View goals and their timeline summaries.</p>
                </div>
              </div>
              {loading ? (
                <p className="mt-4 text-slate-400">Loading goals…</p>
              ) : goals.length === 0 ? (
                <p className="mt-4 text-slate-400">No goals yet. Add one to begin planning.</p>
              ) : (
                <div className="mt-6 space-y-4">
                  {goals.map((goal) => (
                    <Timeline key={goal.id} goal={goal} onAddSubtask={handleAddSubtask} />
                  ))}
                </div>
              )}
            </div>
          </div>
        </section>
      </div>
    </div>
  )
}

export default App
