import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../api/client'
import type { TagSummaryWithDiff } from '../api/client'
import { MonthNav } from '../components/MonthNav'
import { BalanceCard } from '../components/BalanceCard'
import { TagSummaryCard } from '../components/TagSummaryCard'
import { useBalance } from '../hooks/useBalance'

export function Home() {
  const navigate = useNavigate()
  const now = new Date()
  const [year, setYear] = useState(now.getFullYear())
  const [month, setMonth] = useState(now.getMonth() + 1)
  const [total, setTotal] = useState(0)
  const [tagTotals, setTagTotals] = useState<TagSummaryWithDiff[]>([])
  const { balance, updateBalance } = useBalance()

  useEffect(() => {
    let ignore = false
    api.getMonthTotal(year, month).then((r) => {
      if (!ignore) setTotal(r.total)
    })
    api.getTagTotalsWithDiff(year, month).then((data) => {
      if (!ignore) setTagTotals(data)
    })
    return () => {
      ignore = true
    }
  }, [year, month])

  const changeMonth = (delta: number) => {
    let m = month + delta
    let y = year
    if (m < 1) {
      m = 12
      y--
    }
    if (m > 12) {
      m = 1
      y++
    }
    setYear(y)
    setMonth(m)
  }

  const handleTagDetail = (params: {
    tagId: number
    tagName: string
    year: number
    month: number
  }) => {
    navigate(
      `/tags/${params.tagId}/details?year=${params.year}&month=${params.month}&name=${encodeURIComponent(params.tagName)}`,
    )
  }

  return (
    <div>
      <h1>家計簿</h1>

      <BalanceCard
        balance={balance}
        onBalanceUpdate={async (amount) => {
          await updateBalance(amount)
          api.getMonthTotal(year, month).then((r) => setTotal(r.total))
          api.getTagTotalsWithDiff(year, month).then(setTagTotals)
        }}
      />

      <MonthNav year={year} month={month} onChangeMonth={changeMonth} />

      <div className="card">
        <h2>月支出合計</h2>
        <p className="total">&yen;{total.toLocaleString()}</p>
      </div>

      <TagSummaryCard
        tagTotals={tagTotals}
        onTagDetail={handleTagDetail}
        year={year}
        month={month}
      />

      <div className="actions">
        <button className="btn-primary" onClick={() => navigate('/tags/select')}>
          支出入力
        </button>
        <button className="btn-secondary" onClick={() => navigate('/tags/manage')}>
          支出理由管理
        </button>
      </div>
    </div>
  )
}
