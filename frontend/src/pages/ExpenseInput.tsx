import { useState } from "react";
import { api } from "../api/client";
import type { ActionTag } from "../api/client";

interface ExpenseRow {
  item: string;
  amount: string;
}

interface Props {
  tag: ActionTag;
  onBack: () => void;
  onSaved: () => void;
}

export function ExpenseInput({ tag, onBack, onSaved }: Props) {
  const today = new Date().toISOString().split("T")[0];
  const [date, setDate] = useState(today);
  const [rows, setRows] = useState<ExpenseRow[]>([{ item: "", amount: "" }]);
  const [saving, setSaving] = useState(false);

  const updateRow = (index: number, field: keyof ExpenseRow, value: string) => {
    const next = [...rows];
    next[index] = { ...next[index], [field]: value };
    setRows(next);
  };

  const addRow = () => setRows([...rows, { item: "", amount: "" }]);

  const removeRow = (index: number) => {
    if (rows.length > 1) setRows(rows.filter((_, i) => i !== index));
  };

  const handleSave = async () => {
    const validRows = rows.filter((r) => Number(r.amount) > 0);
    if (validRows.length === 0) return;

    setSaving(true);
    try {
      const event = await api.createEvent(date, tag.id);
      for (const row of validRows) {
        await api.createExpense(event.id, row.item, Number(row.amount));
      }
      onSaved();
    } catch (e) {
      alert("保存に失敗しました: " + (e as Error).message);
    } finally {
      setSaving(false);
    }
  };

  return (
    <div>
      <button className="btn-back" onClick={onBack}>&larr; 戻る</button>
      <h1>支出入力</h1>

      <div className="card">
        <label>
          日付
          <input type="date" value={date} onChange={(e) => setDate(e.target.value)} />
        </label>
        <label>
          行動タグ
          <input type="text" value={tag.name} disabled />
        </label>
      </div>

      <div className="card">
        <h2>支出項目</h2>
        {rows.map((row, i) => (
          <div key={i} className="expense-row">
            <input
              type="text"
              placeholder="項目名"
              value={row.item}
              onChange={(e) => updateRow(i, "item", e.target.value)}
            />
            <input
              type="number"
              placeholder="金額"
              value={row.amount}
              onChange={(e) => updateRow(i, "amount", e.target.value)}
            />
            {rows.length > 1 && (
              <button className="btn-remove" onClick={() => removeRow(i)}>×</button>
            )}
          </div>
        ))}
        <button className="btn-secondary" onClick={addRow}>+ 項目追加</button>
      </div>

      <button className="btn-primary" onClick={handleSave} disabled={saving}>
        {saving ? "保存中..." : "保存"}
      </button>
    </div>
  );
}
