const BASE_URL = "http://localhost:8080";

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || res.statusText);
  }
  if (res.status === 204) return undefined as T;
  return res.json();
}

export interface ActionTag {
  id: number;
  name: string;
}

export interface Event {
  id: number;
  date: string;
  action_tag_id: number;
}

export interface Expense {
  id: number;
  event_id: number;
  item: string;
  amount: number;
}

export interface TagSummary {
  tag: string;
  total: number;
}

export const api = {
  getTags: () => request<ActionTag[]>("/tags"),
  createTag: (name: string) => request<ActionTag>("/tags", { method: "POST", body: JSON.stringify({ name }) }),
  updateTag: (id: number, name: string) => request<ActionTag>(`/tags/${id}`, { method: "PUT", body: JSON.stringify({ name }) }),
  deleteTag: (id: number) => request<void>(`/tags/${id}`, { method: "DELETE" }),

  createEvent: (date: string, action_tag_id: number) =>
    request<Event>("/events", { method: "POST", body: JSON.stringify({ date, action_tag_id }) }),

  createExpense: (event_id: number, item: string, amount: number) =>
    request<Expense>("/expenses", { method: "POST", body: JSON.stringify({ event_id, item: item || undefined, amount }) }),

  getMonthTotal: (year: number, month: number) =>
    request<{ total: number }>(`/summary/month?year=${year}&month=${month}`),

  getTagTotals: (year: number, month: number) =>
    request<TagSummary[]>(`/summary/tag?year=${year}&month=${month}`),
};
