import type { TagSummaryWithDiff } from '../api/client'

interface Props {
  tagTotals: TagSummaryWithDiff[]
  onTagDetail: (params: { tagId: number; tagName: string; year: number; month: number }) => void
  year: number
  month: number
}

function formatDiff(diff: number): string {
  if (diff > 0) return `+¥${diff.toLocaleString()}`
  if (diff < 0) return `-¥${Math.abs(diff).toLocaleString()}`
  return '±¥0'
}

function diffClassName(diff: number): string {
  if (diff > 0) return 'diff-increase'
  if (diff < 0) return 'diff-decrease'
  return 'diff-zero'
}

export function TagSummaryCard({ tagTotals, onTagDetail, year, month }: Props) {
  return (
    <div className="card">
      <h2>支出理由別</h2>
      {tagTotals.length === 0 ? (
        <p className="empty">データなし</p>
      ) : (
        <ul className="tag-list">
          {tagTotals.map((t) => (
            <li
              key={t.tag}
              className="tag-list-clickable"
              onClick={() => onTagDetail({ tagId: t.tag_id, tagName: t.tag, year, month })}
            >
              <span>{t.tag}</span>
              <span>
                &yen;{t.total.toLocaleString()}
                <span className={`diff ${diffClassName(t.diff)}`}>{formatDiff(t.diff)}</span>
              </span>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
