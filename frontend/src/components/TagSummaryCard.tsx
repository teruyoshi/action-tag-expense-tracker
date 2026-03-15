import type { TagSummary } from '../api/client'

interface Props {
  tagTotals: TagSummary[]
  onTagDetail: (params: { tagId: number; tagName: string; year: number; month: number }) => void
  year: number
  month: number
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
              <span>&yen;{t.total.toLocaleString()}</span>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
