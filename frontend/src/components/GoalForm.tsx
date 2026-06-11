import { useState } from 'react'
import type { CreateGoalPayload } from '../types'

interface Props {
  onSubmit: (payload: CreateGoalPayload) => void
}

const initialState = {
  title: '',
  startDate: new Date().toISOString().slice(0, 10),
  totalDuration: 7,
  durationType: 'days' as const,
}

const inputClass =
  'mt-1.5 w-full rounded-xl border border-slate-200 bg-white px-3.5 py-2.5 text-sm text-slate-800 outline-none transition focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/15'

export default function GoalForm({ onSubmit }: Props) {
  const [form, setForm] = useState(initialState)

  const handleChange = (field: keyof typeof initialState, value: string) => {
    setForm((current) => ({
      ...current,
      [field]: field === 'totalDuration' ? Number(value) : value,
    }))
  }

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    onSubmit(form)
    setForm(initialState)
  }

  return (
    <form className="space-y-4" onSubmit={handleSubmit}>
      <label className="block text-xs font-semibold text-slate-600">
        Title
        <input
          value={form.title}
          onChange={(event) => handleChange('title', event.target.value)}
          className={inputClass}
          placeholder="Launch marketing sprint"
          required
        />
      </label>

      <label className="block text-xs font-semibold text-slate-600">
        Start date
        <input
          type="date"
          value={form.startDate}
          onChange={(event) => handleChange('startDate', event.target.value)}
          className={inputClass}
          required
        />
      </label>

      <div className="grid gap-4 sm:grid-cols-2">
        <label className="block text-xs font-semibold text-slate-600">
          Total duration
          <input
            type="number"
            min={1}
            value={form.totalDuration}
            onChange={(event) => handleChange('totalDuration', event.target.value)}
            className={inputClass}
            required
          />
        </label>

        <label className="block text-xs font-semibold text-slate-600">
          Unit
          <select
            value={form.durationType}
            onChange={(event) => handleChange('durationType', event.target.value)}
            className={inputClass}
          >
            <option value="days">Days</option>
            <option value="hours">Hours</option>
          </select>
        </label>
      </div>

      <button
        type="submit"
        className="w-full rounded-xl bg-indigo-600 px-4 py-2.5 text-sm font-semibold text-white transition hover:bg-indigo-500"
      >
        Create goal
      </button>
    </form>
  )
}
