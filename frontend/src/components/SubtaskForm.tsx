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

  const handleSubmit = (
    event: React.FormEvent<HTMLFormElement>,
  ) => {
    event.preventDefault()

    if (!form.title.trim() || form.weight <= 0) {
      return
    }

    onSubmit(form)
    setForm(initialState)
  }

  return (
    <form
      className="mt-6"
      onSubmit={handleSubmit}
    >
      <div className="space-y-5">

        {/* Title */}

        <div>
          <label className="mb-2 block text-sm font-medium text-slate-300">
            Subtask Title
          </label>

          <input
            value={form.title}
            onChange={(event) =>
              setForm((current) => ({
                ...current,
                title: event.target.value,
              }))
            }
            placeholder="Example: Build REST API"
            className="w-full rounded-2xl border border-slate-700 bg-slate-900 px-4 py-3 text-slate-100 outline-none transition focus:border-sky-500 focus:ring-2 focus:ring-sky-500/20"
            required
          />
        </div>

        {/* Weight */}

        <div>
          <div className="flex items-center justify-between">

            <label className="text-sm font-medium text-slate-300">
              Weight
            </label>

            <span className="rounded-full bg-sky-500/20 px-3 py-1 text-xs font-medium text-sky-400">
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
                weight: Number(
                  event.target.value,
                ),
              }))
            }
            className="mt-4 w-full cursor-pointer"
          />

          <input
            type="number"
            min={1}
            max={100}
            value={form.weight}
            onChange={(event) =>
              setForm((current) => ({
                ...current,
                weight: Number(
                  event.target.value,
                ),
              }))
            }
            className="mt-3 w-full rounded-2xl border border-slate-700 bg-slate-900 px-4 py-3 text-slate-100 outline-none transition focus:border-sky-500 focus:ring-2 focus:ring-sky-500/20"
            required
          />

          <p className="mt-2 text-xs text-slate-500">
            Timeline duration will be calculated
            automatically based on weight.
          </p>
        </div>

        {/* Submit */}

        <button
          type="submit"
          className="flex w-full items-center justify-center gap-2 rounded-2xl bg-gradient-to-r from-sky-500 to-cyan-400 px-4 py-3 text-sm font-semibold text-slate-950 transition hover:scale-[1.02]"
        >
          <span className="text-lg">+</span>
          Add Subtask
        </button>

      </div>
    </form>
  )
}
