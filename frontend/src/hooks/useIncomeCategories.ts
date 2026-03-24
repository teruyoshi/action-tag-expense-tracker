import { useEffect, useState } from 'react'
import { api } from '../api/client'
import type { IncomeCategory } from '../api/client'

export function useIncomeCategories() {
  const [categories, setCategories] = useState<IncomeCategory[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchCategories = async () => {
    try {
      const data = await api.getIncomeCategories()
      setCategories(data)
    } catch (e) {
      setError((e as Error).message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchCategories()
  }, [])

  const createCategory = async (name: string): Promise<IncomeCategory> => {
    const category = await api.createIncomeCategory(name)
    await fetchCategories()
    return category
  }

  const updateCategory = async (id: number, name: string) => {
    await api.updateIncomeCategory(id, name)
    await fetchCategories()
  }

  const deleteCategory = async (id: number) => {
    await api.deleteIncomeCategory(id)
    await fetchCategories()
  }

  return { categories, loading, error, createCategory, updateCategory, deleteCategory }
}
