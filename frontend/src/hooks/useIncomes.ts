import { useCallback, useEffect, useState } from 'react'
import { api } from '../api/client'
import type { Income } from '../api/client'

export function useIncomes(year: number, month: number) {
  const [incomes, setIncomes] = useState<Income[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchIncomes = useCallback(async () => {
    try {
      setLoading(true)
      const data = await api.getIncomes(year, month)
      setIncomes(data)
    } catch (e) {
      setError((e as Error).message)
    } finally {
      setLoading(false)
    }
  }, [year, month])

  useEffect(() => {
    fetchIncomes()
  }, [fetchIncomes])

  const createIncome = async (
    income_category_id: number,
    date: string,
    description: string,
    amount: number,
  ) => {
    await api.createIncome(income_category_id, date, description, amount)
    await fetchIncomes()
  }

  const updateIncome = async (
    id: number,
    income_category_id: number,
    date: string,
    description: string,
    amount: number,
  ) => {
    await api.updateIncome(id, income_category_id, date, description, amount)
    await fetchIncomes()
  }

  const deleteIncome = async (id: number) => {
    await api.deleteIncome(id)
    await fetchIncomes()
  }

  return { incomes, loading, error, fetchIncomes, createIncome, updateIncome, deleteIncome }
}
