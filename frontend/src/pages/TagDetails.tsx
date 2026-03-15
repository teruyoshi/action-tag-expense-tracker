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
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editItem, setEditItem] = useState("");
  const [editAmount, setEditAmount] = useState("");

  useEffect(() => {
    if (tagId && year && month) {
      api.getTagExpenseDetails(year, month, Number(tagId)).then(setDetails);
    }
  }, [year, month, tagId]);

  const total = details.reduce((sum, d) => sum + d.amount, 0);

  const startEdit = (d: TagExpenseDetail) => {
    setEditingId(d.id);
    setEditItem(d.item);
    setEditAmount(String(d.amount));
  };

  const cancelEdit = () => {
    setEditingId(null);
  };

  const saveEdit = async () => {
    if (editingId === null) return;
    const amount = Number(editAmount);
    if (!amount || amount <= 0) {
      alert("金額は1以上を入力してください");
      return;
    }
    try {
      await api.updateExpense(editingId, editItem, amount);
      setDetails((prev) =>
        prev.map((d) =>
          d.id === editingId ? { ...d, item: editItem, amount } : d
        )
      );
      setEditingId(null);
    } catch {
      alert("更新に失敗しました");
    }
  };

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
                <th></th>
              </tr>
            </thead>
            <tbody>
              {details.map((d) =>
                editingId === d.id ? (
                  <tr key={d.id}>
                    <td>{d.date}</td>
                    <td>
                      <input
                        type="text"
                        value={editItem}
                        onChange={(e) => setEditItem(e.target.value)}
                        className="edit-input"
                      />
                    </td>
                    <td>
                      <input
                        type="number"
                        value={editAmount}
                        onChange={(e) => setEditAmount(e.target.value)}
                        className="edit-input"
                      />
                    </td>
                    <td className="edit-actions">
                      <button className="btn-save" onClick={saveEdit}>保存</button>
                      <button className="btn-cancel" onClick={cancelEdit}>取消</button>
                    </td>
                  </tr>
                ) : (
                  <tr key={d.id} className="editable-row" onClick={() => startEdit(d)}>
                    <td>{d.date}</td>
                    <td>{d.item || "—"}</td>
                    <td className="amount-col">&yen;{d.amount.toLocaleString()}</td>
                    <td></td>
                  </tr>
                )
              )}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
