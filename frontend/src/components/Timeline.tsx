import type { Goal, CreateSubtaskPayload } from '../types'
import SubtaskCard from './SubtaskCard'
import SubtaskForm from './SubtaskForm'

interface Props {
  goal: Goal
  onAddSubtask: (
    goalId: string,
    payload: CreateSubtaskPayload,
  ) => void
  onProgressChange: (
    subtaskId: string,
    progress: number,
  ) => void
  onCompleteSubtask: (subtaskId: string) => void
  onReorderSubtasks: (
    goalId: string,
    subtaskIds: string[],
  ) => void
}

export default function Timeline({
  goal,
  onAddSubtask,
  onProgressChange,
  onCompleteSubtask,
  onReorderSubtasks,
}: Props) {
  const percentageComplete = Math.round(goal.progress ?? 0)

  const moveSubtask = (index: number, direction: -1 | 1) => {
    const nextIndex = index + direction

    if (
      nextIndex < 0 ||
      nextIndex >= goal.subtasks.length
    ) {
      return
    }

    const nextOrder = goal.subtasks.map((task) => task.id)
    const [moved] = nextOrder.splice(index, 1)
    nextOrder.splice(nextIndex, 0, moved)
    onReorderSubtasks(goal.id, nextOrder)
  }

  return (
    <div className="rounded-3xl border border-slate-800 bg-slate-950/95 p-6 shadow-2xl shadow-slate-950/30">

      {/* Hero */}

      <div className="rounded-3xl bg-gradient-to-r from-sky-500/20 to-cyan-500/10 p-6">
        <div className="flex flex-col gap-5 md:flex-row md:items-start md:justify-between">

          <div>
            <p className="text-xs uppercase tracking-[0.3em] text-sky-400">
              Goal
            </p>

            <h2 className="mt-2 text-3xl font-bold text-white">
              {goal.title}
            </h2>

            <div className="mt-4 flex flex-wrap gap-3 text-sm text-slate-300">
              <span>
                Start:{' '}
                {new Date(
                  goal.startDate,
                ).toLocaleDateString()}
              </span>

              <span>
                Duration: {goal.totalDuration}{' '}
                {goal.durationType}
              </span>
            </div>
          </div>

          <div className="rounded-2xl bg-slate-900/80 px-6 py-5 text-center">
            <p className="text-3xl font-bold text-sky-400">
              {percentageComplete}%
            </p>

            <p className="mt-1 text-xs uppercase tracking-wider text-slate-400">
              Complete
            </p>
          </div>
        </div>
      </div>

      {/* Statistics */}

      <div className="mt-6 grid gap-4 sm:grid-cols-3">

        <div className="rounded-2xl bg-slate-900 p-4">
          <p className="text-xs uppercase tracking-wider text-slate-400">
            Tasks
          </p>

          <p className="mt-2 text-2xl font-bold text-white">
            {goal.subtasks.length}
          </p>
        </div>

        <div className="rounded-2xl bg-slate-900 p-4">
          <p className="text-xs uppercase tracking-wider text-slate-400">
            Duration
          </p>

          <p className="mt-2 text-2xl font-bold text-white">
            {goal.totalDuration}
          </p>
        </div>

        <div className="rounded-2xl bg-slate-900 p-4">
          <p className="text-xs uppercase tracking-wider text-slate-400">
            Progress
          </p>

          <p className="mt-2 text-2xl font-bold text-sky-400">
            {percentageComplete}%
          </p>
        </div>

      </div>

      {/* Progress */}

      <div className="mt-6 rounded-3xl border border-slate-800 bg-slate-900 p-5">

        <div className="flex items-center justify-between">
          <span className="text-sm text-slate-400">
            Overall Progress
          </span>

          <span className="text-sm font-semibold text-white">
            {percentageComplete}%
          </span>
        </div>

        <div className="mt-3 h-4 overflow-hidden rounded-full bg-slate-800">

          <div
            className="h-full rounded-full bg-gradient-to-r from-sky-500 to-cyan-400 transition-all duration-500"
            style={{
              width: `${percentageComplete}%`,
            }}
          />

        </div>

      </div>

      {/* Timeline */}

      <div className="mt-8">

        <h3 className="text-lg font-semibold text-white">
          Timeline
        </h3>

        {goal.subtasks.length === 0 ? (
          <div className="mt-4 rounded-2xl border border-dashed border-slate-700 p-6 text-center text-slate-400">
            No subtasks yet.
          </div>
        ) : (
          <div className="relative mt-6 space-y-5">

            <div className="absolute bottom-0 left-6 top-0 w-px bg-slate-800" />

            {goal.subtasks.map((task, index) => (
              <div
                key={task.id}
                className="relative pl-14"
              >
                <div className="absolute left-3 top-6 h-6 w-6 rounded-full border-4 border-slate-950 bg-sky-500" />

                <SubtaskCard
                  task={task}
                  canMoveUp={index > 0}
                  canMoveDown={
                    index < goal.subtasks.length - 1
                  }
                  onMoveUp={() =>
                    moveSubtask(index, -1)
                  }
                  onMoveDown={() =>
                    moveSubtask(index, 1)
                  }
                  onProgressChange={
                    onProgressChange
                  }
                  onComplete={
                    onCompleteSubtask
                  }
                />
              </div>
            ))}

          </div>
        )}

      </div>

      {/* Add Subtask */}

      <div className="mt-8 rounded-3xl border border-dashed border-slate-700 bg-slate-900/60 p-6">

        <h3 className="text-lg font-semibold text-white">
          Add New Subtask
        </h3>

        <p className="mt-2 text-sm text-slate-400">
          Timeline will automatically recalculate
          based on task weights.
        </p>

        <SubtaskForm
          onSubmit={(payload) =>
            onAddSubtask(goal.id, payload)
          }
        />

      </div>

    </div>
  )
}
