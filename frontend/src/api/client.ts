const BASE_URL = ''

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  if (!res.ok) {
    const text = await res.text()
    let message = text || res.statusText
    try {
      const json = JSON.parse(text)
      if (json.error) message = json.error
    } catch {
      // plain text fallback
    }
    throw new Error(message)
  }
  if (res.status === 204) return undefined as T
  return res.json()
}

export interface ActionTag {
  id: number
  name: string
}

export interface Event {
  id: number
  date: string
  action_tag_id: number
}

export interface Expense {
  id: number
  event_id: number
  item: string
  amount: number
}

export interface TagSummary {
  tag_id: number
  tag: string
  total: number
}

export interface TagSummaryWithDiff {
  tag_id: number
  tag: string
  total: number
  prev_total: number
  diff: number
}

export interface TagExpenseDetail {
  id: number
  date: string
  item: string
  amount: number
}

export interface Balance {
  id: number
  amount: number
  updated_at: string
}

export const api = {
  getTags: () => request<ActionTag[]>('/tags'),
  createTag: (name: string) =>
    request<ActionTag>('/tags', { method: 'POST', body: JSON.stringify({ name }) }),
  updateTag: (id: number, name: string) =>
    request<ActionTag>(`/tags/${id}`, { method: 'PUT', body: JSON.stringify({ name }) }),
  deleteTag: (id: number) => request<void>(`/tags/${id}`, { method: 'DELETE' }),

  createEvent: (date: string, action_tag_id: number) =>
    request<Event>('/events', { method: 'POST', body: JSON.stringify({ date, action_tag_id }) }),

  createExpense: (event_id: number, item: string, amount: number) =>
    request<Expense>('/expenses', {
      method: 'POST',
      body: JSON.stringify({ event_id, item: item || undefined, amount }),
    }),

  updateExpense: (id: number, item: string, amount: number) =>
    request<Expense>(`/expenses/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ item: item || undefined, amount }),
    }),

  getMonthTotal: (year: number, month: number) =>
    request<{ total: number }>(`/summary/month?year=${year}&month=${month}`),

  getTagTotals: (year: number, month: number) =>
    request<TagSummary[]>(`/summary/tag?year=${year}&month=${month}`),

  getTagTotalsWithDiff: (year: number, month: number) =>
    request<TagSummaryWithDiff[]>(`/summary/tag/diff?year=${year}&month=${month}`),

  getTagExpenseDetails: (year: number, month: number, tagId: number) =>
    request<TagExpenseDetail[]>(`/summary/tag/details?year=${year}&month=${month}&tag_id=${tagId}`),

  getBalance: () => request<Balance>('/balance'),
  updateBalance: (amount: number) =>
    request<Balance>('/balance', { method: 'PUT', body: JSON.stringify({ amount }) }),
}
