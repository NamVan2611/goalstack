import { useState } from 'react'
import { AnimatePresence, motion } from 'framer-motion'
import type { Subtask } from '../types'

interface Props {
  task: Subtask
  canMoveUp: boolean
  canMoveDown: boolean
  onMoveUp: () => void
  onMoveDown: () => void
  onProgressChange: (
    subtaskId: string,
    progress: number,
  ) => void
  onComplete: (subtaskId: string) => void
}

export default function SubtaskCard({
  task,
  canMoveUp,
  canMoveDown,
  onMoveUp,
  onMoveDown,
  onProgressChange,
  onComplete,
}: Props) {
  const [expanded, setExpanded] = useState(false)
  const isCompleted = task.status === 'completed'

  return (
    <div
      className={[
        'rounded-3xl border p-5 transition hover:border-sky-500/40',
        isCompleted
          ? 'border-emerald-500/40 bg-emerald-950/30'
          : 'border-slate-800 bg-slate-900',
      ].join(' ')}
    >
      <div className="flex items-start justify-between gap-4">
        <div>
          <h4 className="text-lg font-semibold text-white">
            {task.title}
          </h4>

          <div className="mt-3 flex flex-wrap gap-2">
            <span className="rounded-full bg-sky-500/20 px-3 py-1 text-xs font-medium text-sky-400">
              Weight {task.weight}
            </span>

            <span className="rounded-full bg-emerald-500/20 px-3 py-1 text-xs font-medium text-emerald-400">
              {task.allocatedTime.toFixed(1)} Hours
            </span>

            <span className="rounded-full bg-slate-800 px-3 py-1 text-xs font-medium capitalize text-slate-300">
              {task.status.replace('_', ' ')}
            </span>
          </div>
        </div>

        <div className="flex shrink-0 flex-wrap justify-end gap-2">
          <button
            type="button"
            onClick={onMoveUp}
            disabled={!canMoveUp}
            aria-label={`Move ${task.title} up`}
            title="Move up"
            className="h-9 w-9 rounded-xl border border-slate-700 text-sm text-slate-300 transition hover:border-sky-500 hover:text-white disabled:cursor-not-allowed disabled:opacity-40"
          >
            ^
          </button>

          <button
            type="button"
            onClick={onMoveDown}
            disabled={!canMoveDown}
            aria-label={`Move ${task.title} down`}
            title="Move down"
            className="h-9 w-9 rounded-xl border border-slate-700 text-sm text-slate-300 transition hover:border-sky-500 hover:text-white disabled:cursor-not-allowed disabled:opacity-40"
          >
            v
          </button>

          <button
            type="button"
            onClick={() => setExpanded(!expanded)}
            className="rounded-xl border border-slate-700 px-3 py-2 text-sm text-slate-300 transition hover:border-sky-500 hover:text-white"
          >
            {expanded ? 'Hide' : 'Details'}
          </button>
        </div>
      </div>

      <div className="mt-5">
        <div className="flex items-center justify-between text-sm">
          <span className="text-slate-400">
            Progress
          </span>

          <span className="font-medium text-white">
            {task.progress}%
          </span>
        </div>

        <div className="mt-2 h-3 overflow-hidden rounded-full bg-slate-800">
          <div
            className={[
              'h-full rounded-full transition-all duration-500',
              isCompleted
                ? 'bg-emerald-400'
                : 'bg-gradient-to-r from-sky-500 to-cyan-400',
            ].join(' ')}
            style={{
              width: `${task.progress}%`,
            }}
          />
        </div>
      </div>

      <div className="mt-4 grid gap-3 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-center">
        <input
          type="range"
          min={0}
          max={100}
          value={task.progress}
          onChange={(event) =>
            onProgressChange(
              task.id,
              Number(event.target.value),
            )
          }
          aria-label={`Set progress for ${task.title}`}
          className="w-full cursor-pointer"
        />

        <button
          type="button"
          onClick={() => onComplete(task.id)}
          disabled={isCompleted}
          className="rounded-xl bg-emerald-400 px-4 py-2 text-sm font-semibold text-emerald-950 transition hover:bg-emerald-300 disabled:cursor-not-allowed disabled:opacity-60"
        >
          Complete
        </button>
      </div>

      <div className="mt-5 grid gap-3 md:grid-cols-2">
        <div className="rounded-2xl bg-slate-950 p-4">
          <p className="text-xs uppercase tracking-wider text-slate-500">
            Start
          </p>

          <p className="mt-2 text-sm text-slate-200">
            {new Date(
              task.startDate,
            ).toLocaleString()}
          </p>
        </div>

        <div className="rounded-2xl bg-slate-950 p-4">
          <p className="text-xs uppercase tracking-wider text-slate-500">
            End
          </p>

          <p className="mt-2 text-sm text-slate-200">
            {new Date(
              task.endDate,
            ).toLocaleString()}
          </p>
        </div>
      </div>

      <AnimatePresence initial={false}>
        {expanded && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.22, ease: 'easeOut' }}
            className="overflow-hidden"
          >
            <div className="mt-6 space-y-4 border-t border-slate-800 pt-6">
          <div className="rounded-2xl bg-slate-950 p-4">
            <h5 className="text-sm font-semibold text-white">
              Notes
            </h5>

            {task.notes.length === 0 ? (
              <p className="mt-2 text-sm text-slate-500">
                No notes
              </p>
            ) : (
              <div className="mt-2 space-y-2">
                {task.notes.map((note) => (
                  <div
                    key={note.id}
                    className="rounded-xl bg-slate-900 p-3 text-sm text-slate-300"
                  >
                    {note.content}
                  </div>
                ))}
              </div>
            )}
          </div>

          <div className="rounded-2xl bg-slate-950 p-4">
            <h5 className="text-sm font-semibold text-white">
              Links
            </h5>

            {task.links.length === 0 ? (
              <p className="mt-2 text-sm text-slate-500">
                No links
              </p>
            ) : (
              <div className="mt-2 space-y-2">
                {task.links.map((link) => (
                  <a
                    key={link.id}
                    href={link.url}
                    target="_blank"
                    rel="noreferrer"
                    className="block rounded-xl bg-slate-900 p-3 text-sm text-sky-400 hover:bg-slate-800"
                  >
                    {link.title}
                  </a>
                ))}
              </div>
            )}
          </div>

          <div className="rounded-2xl bg-slate-950 p-4">
            <h5 className="text-sm font-semibold text-white">
              Checklist
            </h5>

            {task.checklistItems.length === 0 ? (
              <p className="mt-2 text-sm text-slate-500">
                No checklist items
              </p>
            ) : (
              <div className="mt-2 space-y-2">
                {task.checklistItems.map((item) => (
                  <div
                    key={item.id}
                    className="flex items-center gap-3 rounded-xl bg-slate-900 p-3 text-sm text-slate-300"
                  >
                    <span>
                      {item.completed ? '[x]' : '[ ]'}
                    </span>

                    <span>{item.title}</span>
                  </div>
                ))}
              </div>
            )}
          </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}
