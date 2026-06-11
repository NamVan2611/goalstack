import { useState } from 'react'
import type { CreateSubtaskPayload } from '../types'

interface Props {
  onSubmit: (payload: CreateSubtaskPayload) => void
}

const initialState = {
  title: '',
  weight: 10,
}

export default function SubtaskForm({ onSubmit }: Props) {
  const [form, setForm] = useState(initialState)

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    if (!form.title.trim() || form.weight <= 0) {
      return
    }

    onSubmit(form)
    setForm(initialState)
  }

  return (
    <form className="space-y-3" onSubmit={handleSubmit}>
      <div>
        <label className="mb-1 block text-[11px] font-semibold text-slate-600">
          Subtask title
        </label>
        <input
          value={form.title}
          onChange={(event) =>
            setForm((current) => ({ ...current, title: event.target.value }))
          }
          placeholder="e.g. Build REST API"
          className="w-full rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm text-slate-800 outline-none transition focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/15"
          required
        />
      </div>

      <div>
        <div className="flex items-center justify-between">
          <label className="text-[11px] font-semibold text-slate-600">Weight</label>
          <span className="rounded-full bg-indigo-100 px-2 py-0.5 text-[10px] font-bold text-indigo-700">
            {form.weight}%
          </span>
        </div>
        <input
          type="range"
          min={1}
          max={100}
          value={form.weight}
          onChange={(event) =>
            setForm((current) => ({
              ...current,
              weight: Number(event.target.value),
            }))
          }
          className="mt-2 w-full cursor-pointer accent-indigo-600"
        />
      </div>

      <button
        type="submit"
        className="w-full rounded-xl bg-indigo-600 px-3 py-2 text-sm font-semibold text-white transition hover:bg-indigo-500"
      >
        Add subtask
      </button>
    </form>
  )
}
