import { useEffect, useState } from "react";
import { api } from "../api/client";
import type { TagSummary } from "../api/client";
import { MonthNav } from "../components/MonthNav";
import { BalanceCard } from "../components/BalanceCard";
import { TagSummaryCard } from "../components/TagSummaryCard";

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

  return (
    <div>
      <h1>家計簿</h1>

      <BalanceCard balance={balance} onBalanceUpdate={setBalance} />

      <MonthNav year={year} month={month} onChangeMonth={changeMonth} />

      <div className="card">
        <h2>月支出合計</h2>
        <p className="total">&yen;{total.toLocaleString()}</p>
      </div>

      <TagSummaryCard tagTotals={tagTotals} onTagDetail={onTagDetail} year={year} month={month} />

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
