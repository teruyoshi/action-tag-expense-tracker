interface Props {
  year: number;
  month: number;
  onChangeMonth: (delta: number) => void;
}

export function MonthNav({ year, month, onChangeMonth }: Props) {
  return (
    <div className="month-nav">
      <button onClick={() => onChangeMonth(-1)}>&lt;</button>
      <span>{year}年{month}月</span>
      <button onClick={() => onChangeMonth(1)}>&gt;</button>
    </div>
  );
}
