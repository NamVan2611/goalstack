import type { Goal } from '../types'
import {
  buildMonthGrid,
  goalColor,
  goalCoversDay,
  MONTH_NAMES,
  sameDay,
  startOfDay,
  WEEKDAY_LABELS,
} from '../utils/date'

interface Props {
  year: number
  month: number
  goals: Goal[]
  selectedGoalId: string | null
  onPrevMonth: () => void
  onNextMonth: () => void
  onToday: () => void
  onSelectGoal: (goalId: string) => void
}

export default function Calendar({
  year,
  month,
  goals,
  selectedGoalId,
  onPrevMonth,
  onNextMonth,
  onToday,
  onSelectGoal,
}: Props) {
  const grid = buildMonthGrid(year, month)
  const today = startOfDay(new Date())

  return (
    <section className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
      {/* Calendar header */}
      <div className="mb-6 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <h2 className="text-xl font-bold text-slate-900">
            {MONTH_NAMES[month]} {year}
          </h2>
          <button
            type="button"
            onClick={onToday}
            className="rounded-lg border border-slate-200 px-3 py-1.5 text-xs font-semibold text-slate-600 transition hover:bg-slate-50"
          >
            Today
          </button>
        </div>

        <div className="flex items-center gap-1">
          <button
            type="button"
            onClick={onPrevMonth}
            aria-label="Previous month"
            className="flex h-9 w-9 items-center justify-center rounded-lg border border-slate-200 text-slate-500 transition hover:bg-slate-50 hover:text-slate-900"
          >
            ‹
          </button>
          <button
            type="button"
            onClick={onNextMonth}
            aria-label="Next month"
            className="flex h-9 w-9 items-center justify-center rounded-lg border border-slate-200 text-slate-500 transition hover:bg-slate-50 hover:text-slate-900"
          >
            ›
          </button>
        </div>
      </div>

      {/* Weekday labels */}
      <div className="mb-2 grid grid-cols-7 gap-2">
        {WEEKDAY_LABELS.map((label) => (
          <div
            key={label}
            className="text-center text-xs font-semibold tracking-wider text-slate-400"
          >
            {label}
          </div>
        ))}
      </div>

      {/* Day grid */}
      <div className="grid grid-cols-7 gap-2">
        {grid.map((day) => {
          const isCurrentMonth = day.getMonth() === month
          const isToday = sameDay(day, today)
          const dayGoals = goals.filter((goal) => goalCoversDay(goal, day))

          return (
            <div
              key={day.toISOString()}
              className={[
                'min-h-[96px] rounded-xl border p-2 transition',
                isCurrentMonth
                  ? 'border-slate-200 bg-white'
                  : 'border-transparent bg-slate-50/60',
              ].join(' ')}
            >
              <div
                className={[
                  'mb-1 flex h-6 w-6 items-center justify-center rounded-full text-xs font-semibold',
                  isToday
                    ? 'bg-indigo-600 text-white'
                    : isCurrentMonth
                    ? 'text-slate-700'
                    : 'text-slate-300',
                ].join(' ')}
              >
                {day.getDate()}
              </div>

              <div className="space-y-1">
                {dayGoals.slice(0, 3).map((goal) => {
                  const isStart = sameDay(new Date(goal.startDate), day)
                  const active = goal.id === selectedGoalId
                  return (
                    <button
                      key={goal.id}
                      type="button"
                      onClick={() => onSelectGoal(goal.id)}
                      title={goal.title}
                      className={[
                        'flex w-full items-center gap-1 truncate rounded-md px-1.5 py-1 text-left text-[11px] font-medium text-white transition',
                        goalColor(goal.id),
                        active ? 'ring-2 ring-slate-900 ring-offset-1' : 'opacity-90 hover:opacity-100',
                      ].join(' ')}
                    >
                      {isStart && <span className="text-[9px]">●</span>}
                      <span className="truncate">{goal.title}</span>
                    </button>
                  )
                })}
                {dayGoals.length > 3 && (
                  <p className="px-1 text-[10px] font-medium text-slate-400">
                    +{dayGoals.length - 3} more
                  </p>
                )}
              </div>
            </div>
          )
        })}
      </div>
    </section>
  )
}
