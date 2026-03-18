import { api } from '../api/client'
import type { Expense } from '../api/client'

export function useExpenses() {
  const createEventWithExpenses = async (
    date: string,
    tagId: number,
    rows: { item: string; amount: number }[],
  ) => {
    const event = await api.createEvent(date, tagId)
    for (const row of rows) {
      await api.createExpense(event.id, row.item, row.amount)
    }
  }

  const updateExpense = async (id: number, item: string, amount: number): Promise<Expense> => {
    return api.updateExpense(id, item, amount)
  }

  return { createEventWithExpenses, updateExpense }
}
