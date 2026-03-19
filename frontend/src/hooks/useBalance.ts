import { useEffect, useState } from 'react'
import { api } from '../api/client'

export function useBalance() {
  const [balance, setBalance] = useState(0)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let ignore = false
    api
      .getBalance()
      .then((b) => {
        if (!ignore) setBalance(b.amount)
      })
      .catch((e) => {
        if (!ignore) setError(e.message)
      })
      .finally(() => {
        if (!ignore) setLoading(false)
      })
    return () => {
      ignore = true
    }
  }, [])

  const updateBalance = async (amount: number): Promise<number> => {
    const b = await api.updateBalance(amount)
    setBalance(b.amount)
    return b.amount
  }

  return { balance, loading, error, updateBalance }
}
