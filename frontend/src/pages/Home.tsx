import { useEffect, useState } from "react";
import { api } from "../api/client";
import type { TagSummary } from "../api/client";

interface Props {
  onNavigate: (page: string) => void;
  onTagDetail: (params: { tagId: number; tagName: string; year: number; month: number }) => void;
}

export function Home({ onNavigate, onTagDetail }: Props) {
  const now = new Date();
  const [year, setYear] = useState(now.getFullYear());
  const [month, setMonth] = useState(now.getMonth() + 1);
  const [total, setTotal] = useState(0);
  const [tagTotals, setTagTotals] = useState<TagSummary[]>([]);
  const [balance, setBalance] = useState(0);
  const [showBalanceModal, setShowBalanceModal] = useState(false);
  const [balanceInput, setBalanceInput] = useState("");

  useEffect(() => {
    api.getMonthTotal(year, month).then((r) => setTotal(r.total));
    api.getTagTotals(year, month).then(setTagTotals);
  }, [year, month]);

  useEffect(() => {
    api.getBalance().then((b) => setBalance(b.amount));
  }, []);

  const changeMonth = (delta: number) => {
    let m = month + delta;
    let y = year;
    if (m < 1) { m = 12; y--; }
    if (m > 12) { m = 1; y++; }
    setYear(y);
    setMonth(m);
  };

  const handleBalanceSubmit = () => {
    const amount = parseInt(balanceInput, 10);
    if (isNaN(amount)) return;
    api.updateBalance(amount).then((b) => {
      setBalance(b.amount);
      setShowBalanceModal(false);
      setBalanceInput("");
    });
  };

  return (
    <div>
      <h1>家計簿</h1>

      <div className="card">
        <h2>所持金</h2>
        <p className="total">&yen;{balance.toLocaleString()}</p>
        <button className="btn-secondary" onClick={() => {
          setBalanceInput(String(balance));
          setShowBalanceModal(true);
        }}>
          設定
        </button>
      </div>

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
              <li key={t.tag} className="tag-list-clickable" onClick={() => onTagDetail({ tagId: t.tag_id, tagName: t.tag, year, month })}>
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

      {showBalanceModal && (
        <div className="modal-overlay" onClick={() => setShowBalanceModal(false)}>
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
              <button className="btn-primary" onClick={handleBalanceSubmit}>
                保存
              </button>
              <button className="btn-secondary" onClick={() => setShowBalanceModal(false)}>
                キャンセル
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
