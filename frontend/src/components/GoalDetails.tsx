import { useState } from 'react'
import type { CreateSubtaskPayload, Goal, Subtask } from '../types'
import { goalEndDate } from '../utils/date'
import SubtaskForm from './SubtaskForm'

interface Props {
  goal: Goal | null
  onAddSubtask: (goalId: string, payload: CreateSubtaskPayload) => void
  onDeleteSubtask: (goalId: string, subtaskId: string) => void
  onProgressChange: (subtaskId: string, progress: number) => void
  onCompleteSubtask: (subtaskId: string) => void
  onAddNote: (goalId: string, subtaskId: string, content: string) => void
  onDeleteNote: (goalId: string, subtaskId: string, noteId: string) => void
  onAddLink: (
    goalId: string,
    subtaskId: string,
    payload: { title: string; url: string },
  ) => void
  onDeleteLink: (goalId: string, subtaskId: string, linkId: string) => void
  onAddChecklistItem: (goalId: string, subtaskId: string, title: string) => void
  onToggleChecklistItem: (
    goalId: string,
    subtaskId: string,
    itemId: string,
    completed: boolean,
  ) => void
  onDeleteChecklistItem: (
    goalId: string,
    subtaskId: string,
    itemId: string,
  ) => void
}

function statusBadge(status: Subtask['status']) {
  switch (status) {
    case 'completed':
      return { label: 'COMPLETED', className: 'bg-emerald-100 text-emerald-700' }
    case 'in_progress':
      return { label: 'IN PROGRESS', className: 'bg-indigo-100 text-indigo-700' }
    default:
      return { label: 'TODO', className: 'bg-slate-100 text-slate-500' }
  }
}

const sectionLabel =
  'text-[10px] font-bold uppercase tracking-wider text-slate-400'
const miniInput =
  'w-full rounded-lg border border-slate-200 bg-white px-2.5 py-1.5 text-xs text-slate-800 outline-none transition focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/15'
const miniBtn =
  'shrink-0 rounded-lg bg-indigo-600 px-2.5 py-1.5 text-xs font-semibold text-white transition hover:bg-indigo-500 disabled:cursor-not-allowed disabled:opacity-50'

function SubtaskBlock({
  goalId,
  task,
  onProgressChange,
  onComplete,
  onDelete,
  onAddNote,
  onDeleteNote,
  onAddLink,
  onDeleteLink,
  onAddChecklistItem,
  onToggleChecklistItem,
  onDeleteChecklistItem,
}: {
  goalId: string
  task: Subtask
  onProgressChange: Props['onProgressChange']
  onComplete: Props['onCompleteSubtask']
  onDelete: (goalId: string, subtaskId: string) => void
  onAddNote: Props['onAddNote']
  onDeleteNote: Props['onDeleteNote']
  onAddLink: Props['onAddLink']
  onDeleteLink: Props['onDeleteLink']
  onAddChecklistItem: Props['onAddChecklistItem']
  onToggleChecklistItem: Props['onToggleChecklistItem']
  onDeleteChecklistItem: Props['onDeleteChecklistItem']
}) {
  const badge = statusBadge(task.status)
  const isCompleted = task.status === 'completed'
  const [collapsed, setCollapsed] = useState(true)
  const [noteText, setNoteText] = useState('')
  const [linkTitle, setLinkTitle] = useState('')
  const [linkUrl, setLinkUrl] = useState('')
  const [checkText, setCheckText] = useState('')

  const submitNote = () => {
    const content = noteText.trim()
    if (!content) return
    onAddNote(goalId, task.id, content)
    setNoteText('')
  }

  const submitLink = () => {
    const title = linkTitle.trim()
    const url = linkUrl.trim()
    if (!title || !url) return
    onAddLink(goalId, task.id, { title, url })
    setLinkTitle('')
    setLinkUrl('')
  }

  const submitCheck = () => {
    const title = checkText.trim()
    if (!title) return
    onAddChecklistItem(goalId, task.id, title)
    setCheckText('')
  }

  return (
    <div className="rounded-2xl border border-slate-200 bg-white p-4">
      {/* Header row */}
      <div className="flex items-start justify-between gap-3">
        <div className="flex items-center gap-2">
          <button
            type="button"
            onClick={() => setCollapsed(!collapsed)}
            className="flex h-6 w-6 items-center justify-center rounded-md text-slate-400 hover:bg-slate-100"
            aria-label={collapsed ? 'Expand subtask' : 'Collapse subtask'}
          >
            {collapsed ? '▶' : '▼'}
          </button>
          <h4 className="text-sm font-bold text-slate-900">{task.title}</h4>
        </div>

        <div className="flex shrink-0 items-center gap-2">
          <span
            className={`rounded-full px-2.5 py-0.5 text-[10px] font-bold tracking-wide ${badge.className}`}
          >
            {badge.label}
          </span>
          <button
            type="button"
            onClick={() => onDelete(goalId, task.id)}
            aria-label="Delete subtask"
            className="flex h-6 w-6 items-center justify-center rounded-md text-slate-300 transition hover:bg-red-50 hover:text-red-500"
          >
            🗑
          </button>
        </div>
      </div>

      {/* Collapsible body */}
      {!collapsed && (
        <>
          {/* Progress */}
          <div className="mt-3">
            <div className="flex items-center justify-between text-[11px] text-slate-500">
              <span>Progress</span>
              <span className="font-semibold text-slate-700">{task.progress}%</span>
            </div>
            <div className="mt-1.5 h-1.5 overflow-hidden rounded-full bg-slate-100">
              <div
                className={`h-full rounded-full ${isCompleted ? 'bg-emerald-500' : 'bg-indigo-500'}`}
                style={{ width: `${task.progress}%` }}
              />
            </div>
            <div className="mt-2 flex items-center gap-2">
              <input
                type="range"
                min={0}
                max={100}
                value={task.progress}
                onChange={(e) => onProgressChange(task.id, Number(e.target.value))}
                aria-label={`Set progress for ${task.title}`}
                className="h-1.5 w-full cursor-pointer accent-indigo-600"
              />
              <button
                type="button"
                onClick={() => onComplete(task.id)}
                disabled={isCompleted}
                className="shrink-0 rounded-lg bg-emerald-500 px-2.5 py-1 text-[11px] font-semibold text-white transition hover:bg-emerald-400 disabled:cursor-not-allowed disabled:opacity-50"
              >
                Done
              </button>
            </div>
          </div>

          {/* Notes */}
          <div className="mt-4">
            <p className={sectionLabel}>Notes</p>
            <div className="mt-1.5 space-y-1.5">
              {task.notes.map((note) => (
                <div
                  key={note.id}
                  className="group flex items-start justify-between gap-2 rounded-lg bg-slate-50 px-2.5 py-1.5 text-xs text-slate-600"
                >
                  <span className="whitespace-pre-wrap break-words">{note.content}</span>
                  <button
                    type="button"
                    onClick={() => onDeleteNote(goalId, task.id, note.id)}
                    aria-label="Delete note"
                    className="shrink-0 text-slate-300 transition hover:text-red-500"
                  >
                    ✕
                  </button>
                </div>
              ))}
            </div>
            <div className="mt-1.5 flex gap-2">
              <input
                value={noteText}
                onChange={(e) => setNoteText(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && submitNote()}
                placeholder="Add a note"
                className={miniInput}
              />
              <button
                type="button"
                onClick={submitNote}
                disabled={!noteText.trim()}
                className={miniBtn}
              >
                Add
              </button>
            </div>
          </div>

          {/* Links */}
          <div className="mt-4">
            <p className={sectionLabel}>Links</p>
            <div className="mt-1.5 space-y-1">
              {task.links.map((link) => (
                <div
                  key={link.id}
                  className="group flex items-center justify-between gap-2"
                >
                  <a
                    href={link.url}
                    target="_blank"
                    rel="noreferrer"
                    className="flex items-center gap-1.5 truncate text-xs font-medium text-indigo-600 hover:underline"
                  >
                    <span>🔗</span>
                    {link.title}
                  </a>
                  <button
                    type="button"
                    onClick={() => onDeleteLink(goalId, task.id, link.id)}
                    aria-label="Delete link"
                    className="shrink-0 text-slate-300 transition hover:text-red-500"
                  >
                    ✕
                  </button>
                </div>
              ))}
            </div>
            <div className="mt-1.5 space-y-1.5">
              <input
                value={linkTitle}
                onChange={(e) => setLinkTitle(e.target.value)}
                placeholder="Link title"
                className={miniInput}
              />
              <div className="flex gap-2">
                <input
                  value={linkUrl}
                  onChange={(e) => setLinkUrl(e.target.value)}
                  onKeyDown={(e) => e.key === 'Enter' && submitLink()}
                  placeholder="https://..."
                  className={miniInput}
                />
                <button
                  type="button"
                  onClick={submitLink}
                  disabled={!linkTitle.trim() || !linkUrl.trim()}
                  className={miniBtn}
                >
                  Add
                </button>
              </div>
            </div>
          </div>

          {/* Checklist */}
          <div className="mt-4">
            <p className={sectionLabel}>Checklist</p>
            <div className="mt-1.5 space-y-1">
              {task.checklistItems.map((item) => (
                <div
                  key={item.id}
                  className="group flex items-center justify-between gap-2 text-xs text-slate-600"
                >
                  <label className="flex flex-1 cursor-pointer items-center gap-2">
                    <input
                      type="checkbox"
                      checked={item.completed}
                      onChange={(e) =>
                        onToggleChecklistItem(goalId, task.id, item.id, e.target.checked)
                      }
                      className="h-3.5 w-3.5 cursor-pointer accent-emerald-500"
                    />
                    <span className={item.completed ? 'text-slate-400 line-through' : ''}>
                      {item.title}
                    </span>
                  </label>
                  <button
                    type="button"
                    onClick={() => onDeleteChecklistItem(goalId, task.id, item.id)}
                    aria-label="Delete checklist item"
                    className="shrink-0 text-slate-300 transition hover:text-red-500"
                  >
                    ✕
                  </button>
                </div>
              ))}
            </div>
            <div className="mt-1.5 flex gap-2">
              <input
                value={checkText}
                onChange={(e) => setCheckText(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && submitCheck()}
                placeholder="Add checklist item"
                className={miniInput}
              />
              <button
                type="button"
                onClick={submitCheck}
                disabled={!checkText.trim()}
                className={miniBtn}
              >
                Add
              </button>
            </div>
          </div>
        </>
      )}
    </div>
  )
}

export default function GoalDetails({
  goal,
  onAddSubtask,
  onDeleteSubtask,
  onProgressChange,
  onCompleteSubtask,
  onAddNote,
  onDeleteNote,
  onAddLink,
  onDeleteLink,
  onAddChecklistItem,
  onToggleChecklistItem,
  onDeleteChecklistItem,
}: Props) {
  const [adding, setAdding] = useState(false)

  if (!goal) {
    return (
      <aside className="flex h-full min-h-[400px] flex-col items-center justify-center rounded-3xl border border-slate-200 bg-white p-6 text-center shadow-sm">
        <div className="text-4xl">🎯</div>
        <h3 className="mt-3 text-base font-bold text-slate-800">No goal selected</h3>
        <p className="mt-1 text-sm text-slate-400">
          Pick a goal from the calendar to see its details.
        </p>
      </aside>
    )
  }

  const progress = Math.round(goal.progress ?? 0)
  const start = new Date(goal.startDate)
  const end = goalEndDate(goal)
  const fmt = (d: Date) =>
    d.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })

  return (
    <aside className="flex h-full flex-col rounded-3xl border border-slate-200 bg-white shadow-sm">
      {/* Header */}
      <div className="border-b border-slate-100 p-6">
        <div className="flex items-center justify-between">
          <p className="text-[11px] font-bold uppercase tracking-wider text-slate-400">
            Goal Details
          </p>
        </div>
        <h2 className="mt-2 text-xl font-bold text-slate-900">{goal.title}</h2>

        <div className="mt-4 flex items-center justify-between">
          <span className="flex items-center gap-1.5 text-xs text-slate-500">
            📅 {fmt(start)} – {fmt(end)}
          </span>
          <span className="text-2xl font-bold text-indigo-600">{progress}%</span>
        </div>
        <div className="mt-2 h-2 overflow-hidden rounded-full bg-slate-100">
          <div
            className="h-full rounded-full bg-indigo-600 transition-all duration-500"
            style={{ width: `${progress}%` }}
          />
        </div>
      </div>

      {/* Subtasks */}
      <div className="scroll-thin flex-1 space-y-3 overflow-y-auto p-6">
        <div className="flex items-center justify-between">
          <p className="text-[11px] font-bold uppercase tracking-wider text-slate-400">
            Subtasks
          </p>
          <button
            type="button"
            onClick={() => setAdding((v) => !v)}
            className="rounded-lg border border-slate-200 px-2.5 py-1 text-[11px] font-semibold text-slate-600 transition hover:bg-slate-50"
          >
            {adding ? 'Cancel' : '+ Add'}
          </button>
        </div>

        {adding && (
          <div className="rounded-2xl border border-dashed border-slate-300 bg-slate-50 p-3">
            <SubtaskForm
              onSubmit={(payload) => {
                onAddSubtask(goal.id, payload)
                setAdding(false)
              }}
            />
          </div>
        )}

        {goal.subtasks.length === 0 && !adding ? (
          <p className="rounded-2xl border border-dashed border-slate-200 p-6 text-center text-sm text-slate-400">
            No subtasks yet.
          </p>
        ) : (
          goal.subtasks.map((task) => (
            <SubtaskBlock
              key={task.id}
              goalId={goal.id}
              task={task}
              onProgressChange={onProgressChange}
              onComplete={onCompleteSubtask}
              onDelete={onDeleteSubtask}
              onAddNote={onAddNote}
              onDeleteNote={onDeleteNote}
              onAddLink={onAddLink}
              onDeleteLink={onDeleteLink}
              onAddChecklistItem={onAddChecklistItem}
              onToggleChecklistItem={onToggleChecklistItem}
              onDeleteChecklistItem={onDeleteChecklistItem}
            />
          ))
        )}
      </div>
    </aside>
  )
}