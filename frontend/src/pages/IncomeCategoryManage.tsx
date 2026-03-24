import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useIncomeCategories } from '../hooks/useIncomeCategories'

export function IncomeCategoryManage() {
  const navigate = useNavigate()
  const { categories, createCategory, updateCategory, deleteCategory } = useIncomeCategories()
  const [newName, setNewName] = useState('')
  const [editingId, setEditingId] = useState<number | null>(null)
  const [editName, setEditName] = useState('')

  const handleCreate = async () => {
    if (!newName.trim()) return
    await createCategory(newName.trim())
    setNewName('')
  }

  const handleUpdate = async (id: number) => {
    if (!editName.trim()) return
    await updateCategory(id, editName.trim())
    setEditingId(null)
  }

  const handleDelete = async (id: number) => {
    if (!confirm('削除しますか？')) return
    await deleteCategory(id)
  }

  return (
    <div>
      <button className="btn-back" onClick={() => navigate('/')}>
        &larr; 戻る
      </button>
      <h1>収入カテゴリ管理</h1>

      <div className="card">
        <h2>新規カテゴリ</h2>
        <div className="input-row">
          <input
            type="text"
            placeholder="カテゴリ名"
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleCreate()}
          />
          <button className="btn-primary" onClick={handleCreate}>
            追加
          </button>
        </div>
      </div>

      <div className="card">
        <h2>カテゴリ一覧</h2>
        {categories.length === 0 ? (
          <p className="empty">カテゴリがありません</p>
        ) : (
          <ul className="manage-list">
            {categories.map((cat) => (
              <li key={cat.id}>
                {editingId === cat.id ? (
                  <div className="input-row">
                    <input
                      type="text"
                      value={editName}
                      onChange={(e) => setEditName(e.target.value)}
                      onKeyDown={(e) => e.key === 'Enter' && handleUpdate(cat.id)}
                    />
                    <button className="btn-primary" onClick={() => handleUpdate(cat.id)}>
                      保存
                    </button>
                    <button className="btn-secondary" onClick={() => setEditingId(null)}>
                      取消
                    </button>
                  </div>
                ) : (
                  <div className="input-row">
                    <span>{cat.name}</span>
                    <button
                      className="btn-secondary"
                      onClick={() => {
                        setEditingId(cat.id)
                        setEditName(cat.name)
                      }}
                    >
                      編集
                    </button>
                    <button className="btn-danger" onClick={() => handleDelete(cat.id)}>
                      削除
                    </button>
                  </div>
                )}
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  )
}
