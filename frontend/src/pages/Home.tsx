import { useEffect, useState } from "react";
import { api } from "../api/client";
import type { TagSummary } from "../api/client";

interface Props {
  onNavigate: (page: string) => void;
}

export function Home({ onNavigate }: Props) {
  const now = new Date();
  const [year, setYear] = useState(now.getFullYear());
  const [month, setMonth] = useState(now.getMonth() + 1);
  const [total, setTotal] = useState(0);
  const [tagTotals, setTagTotals] = useState<TagSummary[]>([]);

  useEffect(() => {
    api.getMonthTotal(year, month).then((r) => setTotal(r.total));
    api.getTagTotals(year, month).then(setTagTotals);
  }, [year, month]);

  const changeMonth = (delta: number) => {
    let m = month + delta;
    let y = year;
    if (m < 1) { m = 12; y--; }
    if (m > 12) { m = 1; y++; }
    setYear(y);
    setMonth(m);
  };

  return (
    <div>
      <h1>家計簿</h1>

      <div className="month-nav">
        <button onClick={() => changeMonth(-1)}>&lt;</button>
        <span>{year}年{month}月</span>
        <button onClick={() => changeMonth(1)}>&gt;</button>
      </div>

      <div className="card">
        <h2>月支出合計</h2>
        <p className="total">&yen;{total.toLocaleString()}</p>
      </div>

      <div className="card">
        <h2>行動タグ別</h2>
        {tagTotals.length === 0 ? (
          <p className="empty">データなし</p>
        ) : (
          <ul className="tag-list">
            {tagTotals.map((t) => (
              <li key={t.tag}>
                <span>{t.tag}</span>
                <span>&yen;{t.total.toLocaleString()}</span>
              </li>
            ))}
          </ul>
        )}
      </div>

      <div className="actions">
        <button className="btn-primary" onClick={() => onNavigate("tag-select")}>
          支出入力
        </button>
        <button className="btn-secondary" onClick={() => onNavigate("tag-manage")}>
          タグ管理
        </button>
      </div>
    </div>
  );
}
