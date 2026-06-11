import type { DurationType, Goal } from '../types'

const DAY_MS = 24 * 60 * 60 * 1000

export function addDuration(
  start: Date,
  amount: number,
  type: DurationType,
): Date {
  const ms = type === 'hours' ? amount * 60 * 60 * 1000 : amount * DAY_MS
  return new Date(start.getTime() + ms)
}

export function goalEndDate(goal: Goal): Date {
  return addDuration(
    new Date(goal.startDate),
    goal.totalDuration,
    goal.durationType,
  )
}

export function sameDay(a: Date, b: Date): boolean {
  return (
    a.getFullYear() === b.getFullYear() &&
    a.getMonth() === b.getMonth() &&
    a.getDate() === b.getDate()
  )
}

export function startOfDay(date: Date): Date {
  const d = new Date(date)
  d.setHours(0, 0, 0, 0)
  return d
}

/** Returns true if `day` falls within the goal's [start, end] range (inclusive). */
export function goalCoversDay(goal: Goal, day: Date): boolean {
  const start = startOfDay(new Date(goal.startDate))
  const end = startOfDay(goalEndDate(goal))
  const target = startOfDay(day)
  return target >= start && target <= end
}

/** Builds the 6x7 grid of dates for a given month (Sun-first). */
export function buildMonthGrid(year: number, month: number): Date[] {
  const first = new Date(year, month, 1)
  const startWeekday = first.getDay()
  const gridStart = new Date(year, month, 1 - startWeekday)

  return Array.from({ length: 42 }, (_, i) => {
    return new Date(
      gridStart.getFullYear(),
      gridStart.getMonth(),
      gridStart.getDate() + i,
    )
  })
}

export const MONTH_NAMES = [
  'January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December',
]

export const WEEKDAY_LABELS = ['SUN', 'MON', 'TUE', 'WED', 'THU', 'FRI', 'SAT']

/** Stable color per goal id, used for calendar bars. */
const BAR_COLORS = [
  'bg-indigo-500',
  'bg-emerald-500',
  'bg-sky-500',
  'bg-amber-500',
  'bg-rose-500',
  'bg-violet-500',
]

export function goalColor(goalId: string): string {
  let hash = 0
  for (let i = 0; i < goalId.length; i += 1) {
    hash = (hash * 31 + goalId.charCodeAt(i)) >>> 0
  }
  return BAR_COLORS[hash % BAR_COLORS.length]
}
