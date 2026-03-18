import { useEffect, useState } from 'react'
import { api } from '../api/client'

export function useBalance() {
  const [balance, setBalance] = useState(0)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    api
      .getBalance()
      .then((b) => setBalance(b.amount))
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }, [])

  const updateBalance = async (amount: number): Promise<number> => {
    const b = await api.updateBalance(amount)
    setBalance(b.amount)
    return b.amount
  }

  return { balance, loading, error, updateBalance }
}
