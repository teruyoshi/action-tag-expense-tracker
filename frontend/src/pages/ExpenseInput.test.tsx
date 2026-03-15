import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { ExpenseInput } from './ExpenseInput'

const mockNavigate = vi.fn()
const mockLocation = { state: { tag: { id: 1, name: '通勤' } } }
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom')
  return {
    ...actual,
    useNavigate: () => mockNavigate,
    useLocation: () => mockLocation,
  }
})

vi.mock('../api/client', () => ({
  api: {
    createEvent: vi.fn(),
    createExpense: vi.fn(),
  },
}))

import { api } from '../api/client'
const mockApi = vi.mocked(api)

beforeEach(() => {
  vi.clearAllMocks()
  mockApi.createEvent.mockResolvedValue({ id: 1, date: '2026-03-14', action_tag_id: 1 })
  mockApi.createExpense.mockResolvedValue({ id: 1, event_id: 1, item: '', amount: 500 })
})

describe('ExpenseInput', () => {
  const renderExpenseInput = () =>
    render(
      <MemoryRouter>
        <ExpenseInput />
      </MemoryRouter>,
    )

  it('タグ名と日付が表示される', () => {
    renderExpenseInput()
    expect(screen.getByDisplayValue('通勤')).toBeInTheDocument()
    expect(screen.getByDisplayValue(new Date().toISOString().split('T')[0])).toBeInTheDocument()
  })

  it('項目追加ボタンで行が増える', async () => {
    renderExpenseInput()
    const user = userEvent.setup()

    expect(screen.getAllByPlaceholderText('金額')).toHaveLength(1)
    await user.click(screen.getByText('+ 項目追加'))
    expect(screen.getAllByPlaceholderText('金額')).toHaveLength(2)
  })

  it('金額を入力して保存できる', async () => {
    renderExpenseInput()
    const user = userEvent.setup()

    await user.type(screen.getByPlaceholderText('項目名'), '電車賃')
    await user.type(screen.getByPlaceholderText('金額'), '500')
    await user.click(screen.getByText('保存'))

    expect(mockApi.createEvent).toHaveBeenCalledWith(expect.any(String), 1)
    expect(mockApi.createExpense).toHaveBeenCalledWith(1, '電車賃', 500)
    expect(mockNavigate).toHaveBeenCalledWith('/')
  })

  it('金額が0以下の行は送信されない', async () => {
    renderExpenseInput()
    const user = userEvent.setup()

    await user.click(screen.getByText('保存'))
    expect(mockApi.createEvent).not.toHaveBeenCalled()
  })

  it('戻るボタンで/tags/selectに遷移する', async () => {
    renderExpenseInput()
    const user = userEvent.setup()
    await user.click(screen.getByText('← 戻る'))
    expect(mockNavigate).toHaveBeenCalledWith('/tags/select')
  })

  it('複数行の支出を保存できる', async () => {
    renderExpenseInput()
    const user = userEvent.setup()

    const itemInputs = screen.getAllByPlaceholderText('項目名')
    const amountInputs = screen.getAllByPlaceholderText('金額')
    await user.type(itemInputs[0], '電車賃')
    await user.type(amountInputs[0], '500')

    await user.click(screen.getByText('+ 項目追加'))
    const newAmountInputs = screen.getAllByPlaceholderText('金額')
    await user.type(screen.getAllByPlaceholderText('項目名')[1], 'バス代')
    await user.type(newAmountInputs[1], '300')

    await user.click(screen.getByText('保存'))

    expect(mockApi.createExpense).toHaveBeenCalledTimes(2)
    expect(mockApi.createExpense).toHaveBeenCalledWith(1, '電車賃', 500)
    expect(mockApi.createExpense).toHaveBeenCalledWith(1, 'バス代', 300)
  })
})
