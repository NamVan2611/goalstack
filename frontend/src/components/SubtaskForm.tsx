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
    <form className="mt-6 rounded-3xl border border-slate-800 bg-slate-950/90 p-4" onSubmit={handleSubmit}>
      <div className="flex flex-col gap-4">
        <div>
          <label className="block text-sm font-medium text-slate-300">Subtask title</label>
          <input
            value={form.title}
            onChange={(event) => setForm((current) => ({ ...current, title: event.target.value }))}
            placeholder="Write user stories"
            className="mt-2 w-full rounded-2xl border border-slate-700 bg-slate-900 px-4 py-3 text-slate-100 outline-none focus:border-sky-500 focus:ring-2 focus:ring-sky-500/20"
            required
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-slate-300">Weight</label>
          <input
            type="number"
            min={1}
            value={form.weight}
            onChange={(event) => setForm((current) => ({ ...current, weight: Number(event.target.value) }))}
            className="mt-2 w-full rounded-2xl border border-slate-700 bg-slate-900 px-4 py-3 text-slate-100 outline-none focus:border-sky-500 focus:ring-2 focus:ring-sky-500/20"
            required
          />
        </div>

        <button
          type="submit"
          className="rounded-2xl bg-sky-500 px-4 py-3 text-sm font-semibold text-slate-950 transition hover:bg-sky-400"
        >
          Add subtask
        </button>
      </div>
    </form>
  )
}
