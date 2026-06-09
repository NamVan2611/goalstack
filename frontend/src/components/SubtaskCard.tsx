import type { Subtask } from '../types'

interface Props {
  task: Subtask
}

export default function SubtaskCard({ task }: Props) {
  return (
    <div className="rounded-3xl border border-slate-800 bg-slate-900 p-4">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h4 className="text-base font-semibold text-slate-100">{task.title}</h4>
          <p className="mt-1 text-sm text-slate-400">Weight: {task.weight}</p>
        </div>
        <div className="text-right text-sm text-slate-400">
          <p>Start: {new Date(task.startDate).toLocaleString()}</p>
          <p>End: {new Date(task.endDate).toLocaleString()}</p>
        </div>
      </div>

      <div className="mt-4 grid gap-3 sm:grid-cols-2">
        <div className="rounded-2xl bg-slate-950 p-3 text-sm text-slate-300">
          Allocated: {task.allocatedTime.toFixed(1)}
        </div>
        <div className="rounded-2xl bg-slate-950 p-3 text-sm text-slate-300">
          Progress: {task.progress}%
        </div>
      </div>
    </div>
  )
}
