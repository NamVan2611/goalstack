import type { Goal, CreateSubtaskPayload } from '../types'
import SubtaskCard from './SubtaskCard'
import SubtaskForm from './SubtaskForm'

interface Props {
  goal: Goal
  onAddSubtask: (goalId: string, payload: CreateSubtaskPayload) => void
}

export default function Timeline({ goal, onAddSubtask }: Props) {
  const percentageComplete = goal.subtasks.length
    ? Math.round(goal.subtasks.reduce((sum, item) => sum + item.progress, 0) / goal.subtasks.length)
    : 0

  return (
    <div className="rounded-3xl border border-slate-800 bg-slate-950/95 p-5 shadow-xl shadow-slate-950/20">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h3 className="text-lg font-semibold text-slate-100">{goal.title}</h3>
          <p className="mt-1 text-sm text-slate-400">Start: {new Date(goal.startDate).toLocaleDateString()}</p>
        </div>
        <div className="rounded-full bg-slate-900 px-4 py-2 text-sm text-slate-300">
          {goal.totalDuration} {goal.durationType}
        </div>
      </div>

      <div className="mt-4 rounded-2xl bg-slate-900 p-4">
        <div className="flex items-center justify-between text-sm text-slate-400">
          <span>Progress</span>
          <span>{percentageComplete}%</span>
        </div>
        <div className="mt-2 h-3 overflow-hidden rounded-full bg-slate-800">
          <div className="h-full rounded-full bg-sky-500" style={{ width: `${percentageComplete}%` }} />
        </div>
      </div>

      <div className="mt-5 space-y-3">
        {goal.subtasks.length === 0 ? (
          <p className="text-sm text-slate-400">No subtasks yet. Add tasks below.</p>
        ) : (
          goal.subtasks.map((task) => <SubtaskCard key={task.id} task={task} />)
        )}
      </div>

      <div className="mt-6 border-t border-slate-800 pt-5">
        <h4 className="text-base font-semibold text-slate-100">Add a subtask</h4>
        <p className="mt-1 text-sm text-slate-400">Subtasks are allocated in sequence based on weight.</p>
        <SubtaskForm onSubmit={(payload) => onAddSubtask(goal.id, payload)} />
      </div>
    </div>
  )
}
