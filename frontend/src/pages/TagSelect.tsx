import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../api/client'
import type { ActionTag } from '../api/client'

export function TagSelect() {
  const navigate = useNavigate()
  const [tags, setTags] = useState<ActionTag[]>([])
  const [showNew, setShowNew] = useState(false)
  const [newName, setNewName] = useState('')

  const load = () => api.getTags().then(setTags)

  useEffect(() => {
    load()
  }, [])

  const handleSelect = (tag: ActionTag) => {
    navigate('/expense/new', { state: { tag } })
  }

  const handleCreate = async () => {
    if (!newName.trim()) return
    try {
      const tag = await api.createTag(newName.trim())
      setNewName('')
      setShowNew(false)
      await load()
      handleSelect(tag)
    } catch {
      alert('タグの作成に失敗しました')
    }
  }

  return (
    <div>
      <button className="btn-back" onClick={() => navigate('/')}>
        &larr; 戻る
      </button>
      <h1>支出理由を選択</h1>
      <div className="tag-grid">
        {tags.map((tag) => (
          <button key={tag.id} className="card tag-card" onClick={() => handleSelect(tag)}>
            {tag.name}
          </button>
        ))}
        {showNew ? (
          <div className="card">
            <div className="input-row">
              <input
                type="text"
                placeholder="新しい支出理由"
                value={newName}
                onChange={(e) => setNewName(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && handleCreate()}
                autoFocus
              />
              <button className="btn-primary" onClick={handleCreate}>
                追加
              </button>
              <button
                className="btn-secondary"
                onClick={() => {
                  setShowNew(false)
                  setNewName('')
                }}
              >
                取消
              </button>
            </div>
          </div>
        ) : (
          <button className="card tag-card" onClick={() => setShowNew(true)}>
            + 新しい支出理由
          </button>
        )}
      </div>
    </div>
  )
}
