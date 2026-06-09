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
    <form className="mt-6 space-y-4" onSubmit={handleSubmit}>
      <label className="block text-sm font-medium text-slate-300">
        Title
        <input
          value={form.title}
          onChange={(event) => handleChange('title', event.target.value)}
          className="mt-2 w-full rounded-2xl border border-slate-700 bg-slate-950 px-4 py-3 text-slate-100 outline-none focus:border-sky-500 focus:ring-2 focus:ring-sky-500/20"
          placeholder="Launch marketing sprint"
          required
        />
      </label>

      <label className="block text-sm font-medium text-slate-300">
        Start date
        <input
          type="date"
          value={form.startDate}
          onChange={(event) => handleChange('startDate', event.target.value)}
          className="mt-2 w-full rounded-2xl border border-slate-700 bg-slate-950 px-4 py-3 text-slate-100 outline-none focus:border-sky-500 focus:ring-2 focus:ring-sky-500/20"
          required
        />
      </label>

      <div className="grid gap-4 sm:grid-cols-2">
        <label className="block text-sm font-medium text-slate-300">
          Total duration
          <input
            type="number"
            min={1}
            value={form.totalDuration}
            onChange={(event) => handleChange('totalDuration', event.target.value)}
            className="mt-2 w-full rounded-2xl border border-slate-700 bg-slate-950 px-4 py-3 text-slate-100 outline-none focus:border-sky-500 focus:ring-2 focus:ring-sky-500/20"
            required
          />
        </label>

        <label className="block text-sm font-medium text-slate-300">
          Unit
          <select
            value={form.durationType}
            onChange={(event) => handleChange('durationType', event.target.value)}
            className="mt-2 w-full rounded-2xl border border-slate-700 bg-slate-950 px-4 py-3 text-slate-100 outline-none focus:border-sky-500 focus:ring-2 focus:ring-sky-500/20"
          >
            <option value="days">Days</option>
            <option value="hours">Hours</option>
          </select>
        </label>
      </div>

      <button
        type="submit"
        className="w-full rounded-2xl bg-sky-500 px-4 py-3 text-base font-semibold text-slate-950 transition hover:bg-sky-400"
      >
        Create goal
      </button>
    </form>
  )
}
