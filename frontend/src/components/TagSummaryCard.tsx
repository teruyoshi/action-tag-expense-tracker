import type { TagSummaryWithDiff } from '../api/client'

interface Props {
  tagTotals: TagSummaryWithDiff[]
  onTagDetail: (params: { tagId: number; tagName: string; year: number; month: number }) => void
  year: number
  month: number
}

function formatDiff(diff: number): string {
  if (diff > 0) return `前月比 ▲¥${diff.toLocaleString()}`
  if (diff < 0) return `前月比 ▼¥${Math.abs(diff).toLocaleString()}`
  return '前月比 ±¥0'
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
        <table className="tag-summary-table">
          <tbody>
            {tagTotals.map((t) => (
              <tr
                key={t.tag}
                className="tag-list-clickable"
                onClick={() => onTagDetail({ tagId: t.tag_id, tagName: t.tag, year, month })}
              >
                <td>{t.tag}</td>
                <td className="amount-col">&yen;{t.total.toLocaleString()}</td>
                <td className={`diff ${diffClassName(t.diff)}`}>{formatDiff(t.diff)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  )
}
