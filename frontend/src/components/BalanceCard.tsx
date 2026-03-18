import { useState } from 'react'

interface Props {
  balance: number
  onBalanceUpdate: (amount: number) => Promise<void>
}

export function BalanceCard({ balance, onBalanceUpdate }: Props) {
  const [showModal, setShowModal] = useState(false)
  const [balanceInput, setBalanceInput] = useState('')

  const handleOpen = () => {
    setBalanceInput(String(balance))
    setShowModal(true)
  }

  const handleSubmit = () => {
    const amount = parseInt(balanceInput, 10)
    if (isNaN(amount)) return
    onBalanceUpdate(amount).then(() => {
      setShowModal(false)
      setBalanceInput('')
    })
  }

  return (
    <>
      <div className="card">
        <h2>所持金</h2>
        <p className="total">&yen;{balance.toLocaleString()}</p>
        <button className="btn-secondary" onClick={handleOpen}>
          設定
        </button>
      </div>

      {showModal && (
        <div className="modal-overlay" onClick={() => setShowModal(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h2>所持金を設定</h2>
            <input
              type="number"
              value={balanceInput}
              onChange={(e) => setBalanceInput(e.target.value)}
              placeholder="金額を入力"
              autoFocus
            />
            <div className="modal-actions">
              <button className="btn-primary" onClick={handleSubmit}>
                保存
              </button>
              <button className="btn-secondary" onClick={() => setShowModal(false)}>
                キャンセル
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  )
}
