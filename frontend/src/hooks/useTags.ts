import { useEffect, useState } from 'react'
import { api } from '../api/client'
import type { ActionTag } from '../api/client'

export function useTags() {
  const [tags, setTags] = useState<ActionTag[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchTags = () => {
    api
      .getTags()
      .then(setTags)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    fetchTags()
  }, [])

  const createTag = async (name: string): Promise<ActionTag> => {
    const tag = await api.createTag(name)
    fetchTags()
    return tag
  }

  const updateTag = async (id: number, name: string) => {
    await api.updateTag(id, name)
    fetchTags()
  }

  const deleteTag = async (id: number) => {
    await api.deleteTag(id)
    fetchTags()
  }

  return { tags, loading, error, createTag, updateTag, deleteTag }
}
