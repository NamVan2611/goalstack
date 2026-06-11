import type { Goal } from '../types'
import { goalEndDate, startOfDay } from '../utils/date'

interface Props {
  goals: Goal[]
}

export default function StatsPanel({ goals }: Props) {
  const allSubtasks = goals.flatMap((g) => g.subtasks)
  const completed = allSubtasks.filter((s) => s.status === 'completed').length
  const inProgress = allSubtasks.filter((s) => s.status === 'in_progress').length
  const todo = allSubtasks.filter((s) => s.status === 'todo').length
  const total = allSubtasks.length || 1

  const distribution = [
    { label: 'Completed', value: completed, color: 'bg-emerald-500' },
    { label: 'In Progress', value: inProgress, color: 'bg-indigo-500' },
    { label: 'To Do', value: todo, color: 'bg-slate-300' },
  ]

  // Next milestone = nearest upcoming goal end date
  const today = startOfDay(new Date())
  const upcoming = goals
    .map((g) => ({ goal: g, end: goalEndDate(g) }))
    .filter((x) => startOfDay(x.end) >= today)
    .sort((a, b) => a.end.getTime() - b.end.getTime())[0]

  const daysLeft = upcoming
    ? Math.ceil(
        (startOfDay(upcoming.end).getTime() - today.getTime()) /
          (24 * 60 * 60 * 1000),
      )
    : null

  return (
    <div className="grid gap-5 md:grid-cols-2">
      {/* Task distribution */}
      <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
        <h3 className="text-sm font-bold text-slate-900">Task Distribution</h3>

        <div className="mt-4 flex h-3 overflow-hidden rounded-full bg-slate-100">
          {distribution.map(
            (d) =>
              d.value > 0 && (
                <div
                  key={d.label}
                  className={d.color}
                  style={{ width: `${(d.value / total) * 100}%` }}
                  title={`${d.label}: ${d.value}`}
                />
              ),
          )}
        </div>

        <div className="mt-4 grid grid-cols-3 gap-3">
          {distribution.map((d) => (
            <div key={d.label}>
              <div className="flex items-center gap-1.5">
                <span className={`h-2.5 w-2.5 rounded-full ${d.color}`} />
                <span className="text-[11px] text-slate-500">{d.label}</span>
              </div>
              <p className="mt-1 text-lg font-bold text-slate-900">{d.value}</p>
            </div>
          ))}
        </div>
      </div>

      {/* Next milestone */}
      <div className="rounded-3xl border border-slate-200 bg-gradient-to-br from-indigo-600 to-violet-600 p-6 text-white shadow-sm">
        <h3 className="text-sm font-bold">Next Milestone</h3>

        {upcoming ? (
          <div className="mt-4">
            <p className="text-lg font-bold">{upcoming.goal.title}</p>
            <p className="mt-1 text-sm text-indigo-100">
              Due{' '}
              {upcoming.end.toLocaleDateString(undefined, {
                month: 'long',
                day: 'numeric',
              })}
            </p>
            <div className="mt-5 flex items-end gap-2">
              <span className="text-4xl font-extrabold">{daysLeft}</span>
              <span className="mb-1 text-sm text-indigo-100">
                {daysLeft === 1 ? 'day left' : 'days left'}
              </span>
            </div>
          </div>
        ) : (
          <p className="mt-4 text-sm text-indigo-100">
            No upcoming milestones. Create a goal to get started.
          </p>
        )}
      </div>
    </div>
  )
}
