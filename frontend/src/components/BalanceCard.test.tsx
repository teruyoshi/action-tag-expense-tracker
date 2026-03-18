import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { BalanceCard } from './BalanceCard'

beforeEach(() => {
  vi.clearAllMocks()
})

describe('BalanceCard', () => {
  it('所持金を表示する', () => {
    render(<BalanceCard balance={100000} onBalanceUpdate={vi.fn()} />)
    expect(screen.getByText('¥100,000')).toBeInTheDocument()
  })

  it('設定ボタンでモーダルが開く', async () => {
    render(<BalanceCard balance={100000} onBalanceUpdate={vi.fn()} />)
    const user = userEvent.setup()
    await user.click(screen.getByText('設定'))
    expect(screen.getByText('所持金を設定')).toBeInTheDocument()
  })

  it('モーダルで金額を更新できる', async () => {
    const onBalanceUpdate = vi.fn().mockResolvedValue(undefined)
    render(<BalanceCard balance={100000} onBalanceUpdate={onBalanceUpdate} />)
    const user = userEvent.setup()

    await user.click(screen.getByText('設定'))
    const input = screen.getByPlaceholderText('金額を入力')
    await user.clear(input)
    await user.type(input, '50000')
    await user.click(screen.getByText('保存'))

    expect(onBalanceUpdate).toHaveBeenCalledWith(50000)
    await waitFor(() => {
      expect(screen.queryByText('所持金を設定')).not.toBeInTheDocument()
    })
  })

  it('キャンセルでモーダルが閉じる', async () => {
    render(<BalanceCard balance={100000} onBalanceUpdate={vi.fn()} />)
    const user = userEvent.setup()
    await user.click(screen.getByText('設定'))
    await user.click(screen.getByText('キャンセル'))
    expect(screen.queryByText('所持金を設定')).not.toBeInTheDocument()
  })
})
