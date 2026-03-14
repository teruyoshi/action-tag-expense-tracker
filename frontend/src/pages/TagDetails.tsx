import { useEffect, useState } from "react";
import { api } from "../api/client";
import type { TagExpenseDetail } from "../api/client";

interface Props {
  tagId: number;
  tagName: string;
  year: number;
  month: number;
  onBack: () => void;
}

export function TagDetails({ tagId, tagName, year, month, onBack }: Props) {
  const [details, setDetails] = useState<TagExpenseDetail[]>([]);

  useEffect(() => {
    api.getTagExpenseDetails(year, month, tagId).then(setDetails);
  }, [year, month, tagId]);

  const total = details.reduce((sum, d) => sum + d.amount, 0);

  return (
    <div>
      <button className="btn-back" onClick={onBack}>&larr; 戻る</button>
      <h1>{tagName}</h1>
      <p className="detail-subtitle">{year}年{month}月</p>

      <div className="card">
        <h2>合計</h2>
        <p className="total">&yen;{total.toLocaleString()}</p>
      </div>

      <div className="card">
        <h2>明細</h2>
        {details.length === 0 ? (
          <p className="empty">データなし</p>
        ) : (
          <table className="detail-table">
            <thead>
              <tr>
                <th>日付</th>
                <th>項目</th>
                <th className="amount-col">金額</th>
              </tr>
            </thead>
            <tbody>
              {details.map((d, i) => (
                <tr key={i}>
                  <td>{d.date}</td>
                  <td>{d.item || "—"}</td>
                  <td className="amount-col">&yen;{d.amount.toLocaleString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
