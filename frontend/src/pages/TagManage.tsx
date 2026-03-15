import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../api/client'
import type { ActionTag } from '../api/client'

export function TagManage() {
  const navigate = useNavigate()
  const [tags, setTags] = useState<ActionTag[]>([])
  const [newName, setNewName] = useState('')
  const [editingId, setEditingId] = useState<number | null>(null)
  const [editName, setEditName] = useState('')

  const load = () => api.getTags().then(setTags)

  useEffect(() => {
    load()
  }, [])

  const handleCreate = async () => {
    if (!newName.trim()) return
    await api.createTag(newName.trim())
    setNewName('')
    load()
  }

  const handleUpdate = async (id: number) => {
    if (!editName.trim()) return
    await api.updateTag(id, editName.trim())
    setEditingId(null)
    load()
  }

  const handleDelete = async (id: number) => {
    if (!confirm('削除しますか？')) return
    await api.deleteTag(id)
    load()
  }

  return (
    <div>
      <button className="btn-back" onClick={() => navigate('/')}>
        &larr; 戻る
      </button>
      <h1>支出理由管理</h1>

      <div className="card">
        <h2>新規支出理由</h2>
        <div className="input-row">
          <input
            type="text"
            placeholder="支出理由名"
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
        <h2>支出理由一覧</h2>
        {tags.length === 0 ? (
          <p className="empty">支出理由がありません</p>
        ) : (
          <ul className="manage-list">
            {tags.map((tag) => (
              <li key={tag.id}>
                {editingId === tag.id ? (
                  <div className="input-row">
                    <input
                      type="text"
                      value={editName}
                      onChange={(e) => setEditName(e.target.value)}
                      onKeyDown={(e) => e.key === 'Enter' && handleUpdate(tag.id)}
                    />
                    <button className="btn-primary" onClick={() => handleUpdate(tag.id)}>
                      保存
                    </button>
                    <button className="btn-secondary" onClick={() => setEditingId(null)}>
                      取消
                    </button>
                  </div>
                ) : (
                  <div className="input-row">
                    <span>{tag.name}</span>
                    <button
                      className="btn-secondary"
                      onClick={() => {
                        setEditingId(tag.id)
                        setEditName(tag.name)
                      }}
                    >
                      編集
                    </button>
                    <button className="btn-danger" onClick={() => handleDelete(tag.id)}>
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
