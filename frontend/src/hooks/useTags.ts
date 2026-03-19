import { useEffect, useState } from 'react'
import { api } from '../api/client'
import type { ActionTag } from '../api/client'

export function useTags() {
  const [tags, setTags] = useState<ActionTag[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchTags = async () => {
    try {
      const data = await api.getTags()
      setTags(data)
    } catch (e) {
      setError((e as Error).message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchTags()
  }, [])

  const createTag = async (name: string): Promise<ActionTag> => {
    const tag = await api.createTag(name)
    await fetchTags()
    return tag
  }

  const updateTag = async (id: number, name: string) => {
    await api.updateTag(id, name)
    await fetchTags()
  }

  const deleteTag = async (id: number) => {
    await api.deleteTag(id)
    await fetchTags()
  }

  return { tags, loading, error, createTag, updateTag, deleteTag }
}
