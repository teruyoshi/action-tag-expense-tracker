import { useEffect, useState } from "react";
import { useNavigate, useParams, useSearchParams } from "react-router-dom";
import { api } from "../api/client";
import type { TagExpenseDetail } from "../api/client";

export function TagDetails() {
  const navigate = useNavigate();
  const { tagId } = useParams<{ tagId: string }>();
  const [searchParams] = useSearchParams();
  const tagName = searchParams.get("name") ?? "";
  const year = Number(searchParams.get("year"));
  const month = Number(searchParams.get("month"));

  const [details, setDetails] = useState<TagExpenseDetail[]>([]);

  useEffect(() => {
    if (tagId && year && month) {
      api.getTagExpenseDetails(year, month, Number(tagId)).then(setDetails);
    }
  }, [year, month, tagId]);

  const total = details.reduce((sum, d) => sum + d.amount, 0);

  return (
    <div>
      <button className="btn-back" onClick={() => navigate("/")}>&larr; 戻る</button>
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
