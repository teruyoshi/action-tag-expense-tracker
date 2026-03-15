import { describe, it, expect, vi, beforeEach } from 'vitest'

const mockFetch = vi.fn()
globalThis.fetch = mockFetch

// client.ts を動的にインポートするため、fetch モックの後にインポート
const { api } = await import('./client')

beforeEach(() => {
  mockFetch.mockReset()
})

function jsonResponse(data: unknown, status = 200) {
  return Promise.resolve({
    ok: true,
    status,
    json: () => Promise.resolve(data),
    text: () => Promise.resolve(JSON.stringify(data)),
  })
}

function errorResponse(status: number, body: string) {
  return Promise.resolve({
    ok: false,
    status,
    statusText: 'Bad Request',
    text: () => Promise.resolve(body),
  })
}

function noContentResponse() {
  return Promise.resolve({
    ok: true,
    status: 204,
    json: () => Promise.reject(new Error('no body')),
    text: () => Promise.resolve(''),
  })
}

describe('api.getTags', () => {
  it('returns tags on success', async () => {
    const tags = [{ id: 1, name: '通勤' }]
    mockFetch.mockReturnValueOnce(jsonResponse(tags))

    const result = await api.getTags()
    expect(result).toEqual(tags)
    expect(mockFetch).toHaveBeenCalledWith(
      'http://localhost:8080/tags',
      expect.objectContaining({ headers: { 'Content-Type': 'application/json' } }),
    )
  })
})

describe('api.createTag', () => {
  it('sends POST with name', async () => {
    const tag = { id: 1, name: '外食' }
    mockFetch.mockReturnValueOnce(jsonResponse(tag))

    const result = await api.createTag('外食')
    expect(result).toEqual(tag)
    expect(mockFetch).toHaveBeenCalledWith(
      'http://localhost:8080/tags',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ name: '外食' }),
      }),
    )
  })
})

describe('api.updateTag', () => {
  it('sends PUT with id and name', async () => {
    const tag = { id: 1, name: '買い物' }
    mockFetch.mockReturnValueOnce(jsonResponse(tag))

    const result = await api.updateTag(1, '買い物')
    expect(result).toEqual(tag)
    expect(mockFetch).toHaveBeenCalledWith(
      'http://localhost:8080/tags/1',
      expect.objectContaining({
        method: 'PUT',
        body: JSON.stringify({ name: '買い物' }),
      }),
    )
  })
})

describe('api.deleteTag', () => {
  it('sends DELETE and handles 204', async () => {
    mockFetch.mockReturnValueOnce(noContentResponse())

    const result = await api.deleteTag(1)
    expect(result).toBeUndefined()
    expect(mockFetch).toHaveBeenCalledWith(
      'http://localhost:8080/tags/1',
      expect.objectContaining({ method: 'DELETE' }),
    )
  })
})

describe('api.createEvent', () => {
  it('sends POST with date and action_tag_id', async () => {
    const event = { id: 1, date: '2026-03-14', action_tag_id: 2 }
    mockFetch.mockReturnValueOnce(jsonResponse(event))

    const result = await api.createEvent('2026-03-14', 2)
    expect(result).toEqual(event)
    expect(mockFetch).toHaveBeenCalledWith(
      'http://localhost:8080/events',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ date: '2026-03-14', action_tag_id: 2 }),
      }),
    )
  })
})

describe('api.createExpense', () => {
  it('sends POST with item and amount', async () => {
    const expense = { id: 1, event_id: 1, item: '電車賃', amount: 500 }
    mockFetch.mockReturnValueOnce(jsonResponse(expense))

    const result = await api.createExpense(1, '電車賃', 500)
    expect(result).toEqual(expense)
  })

  it('sends undefined item when empty string', async () => {
    const expense = { id: 2, event_id: 1, item: '', amount: 300 }
    mockFetch.mockReturnValueOnce(jsonResponse(expense))

    await api.createExpense(1, '', 300)
    const body = JSON.parse(mockFetch.mock.calls[0][1].body)
    expect(body.item).toBeUndefined()
  })
})

describe('api.getMonthTotal', () => {
  it('fetches with year and month query params', async () => {
    mockFetch.mockReturnValueOnce(jsonResponse({ total: 15000 }))

    const result = await api.getMonthTotal(2026, 3)
    expect(result).toEqual({ total: 15000 })
    expect(mockFetch).toHaveBeenCalledWith(
      'http://localhost:8080/summary/month?year=2026&month=3',
      expect.any(Object),
    )
  })
})

describe('api.getTagTotals', () => {
  it('fetches tag summaries', async () => {
    const data = [{ tag: '通勤', total: 5000 }]
    mockFetch.mockReturnValueOnce(jsonResponse(data))

    const result = await api.getTagTotals(2026, 3)
    expect(result).toEqual(data)
  })
})

describe('error handling', () => {
  it('throws error with response body on failure', async () => {
    mockFetch.mockReturnValueOnce(errorResponse(400, 'name is required'))

    await expect(api.createTag('')).rejects.toThrow('name is required')
  })

  it('throws statusText when body is empty', async () => {
    mockFetch.mockReturnValueOnce(
      Promise.resolve({
        ok: false,
        status: 400,
        statusText: 'Bad Request',
        text: () => Promise.resolve(''),
      }),
    )

    await expect(api.getTags()).rejects.toThrow('Bad Request')
  })
})
