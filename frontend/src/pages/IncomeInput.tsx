import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useIncomeCategories } from '../hooks/useIncomeCategories'
import { api } from '../api/client'

export function IncomeInput() {
  const navigate = useNavigate()
  const { categories, loading } = useIncomeCategories()

  const today = new Date().toISOString().split('T')[0]
  const [categoryId, setCategoryId] = useState<number | ''>('')
  const [date, setDate] = useState(today)
  const [description, setDescription] = useState('')
  const [amount, setAmount] = useState('')
  const [saving, setSaving] = useState(false)

  const handleSave = async () => {
    if (!categoryId || Number(amount) <= 0) return

    setSaving(true)
    try {
      await api.createIncome(Number(categoryId), date, description, Number(amount))
      navigate('/')
    } catch (e) {
      alert('保存に失敗しました: ' + (e as Error).message)
    } finally {
      setSaving(false)
    }
  }

  if (loading) return <div>読み込み中...</div>

  return (
    <div>
      <button className="btn-back" onClick={() => navigate('/')}>
        &larr; 戻る
      </button>
      <h1>収入入力</h1>

      <div className="card">
        <label>
          カテゴリ
          <select
            value={categoryId}
            onChange={(e) => setCategoryId(e.target.value ? Number(e.target.value) : '')}
          >
            <option value="">選択してください</option>
            {categories.map((cat) => (
              <option key={cat.id} value={cat.id}>
                {cat.name}
              </option>
            ))}
          </select>
        </label>
        <label>
          日付
          <input type="date" value={date} onChange={(e) => setDate(e.target.value)} />
        </label>
        <label>
          説明
          <input
            type="text"
            placeholder="説明（任意）"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
          />
        </label>
        <label>
          金額
          <input
            type="number"
            placeholder="金額"
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
          />
        </label>
      </div>

      <button className="btn-primary" onClick={handleSave} disabled={saving || !categoryId}>
        {saving ? '保存中...' : '保存'}
      </button>
    </div>
  )
}
