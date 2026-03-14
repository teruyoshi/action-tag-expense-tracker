import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { TagSummaryCard } from "./TagSummaryCard";

describe("TagSummaryCard", () => {
  const tagTotals = [
    { tag_id: 1, tag: "通勤", total: 5000 },
    { tag_id: 2, tag: "外食", total: 10000 },
  ];

  it("タグ別合計を表示する", () => {
    render(<TagSummaryCard tagTotals={tagTotals} onTagDetail={vi.fn()} year={2026} month={3} />);
    expect(screen.getByText("通勤")).toBeInTheDocument();
    expect(screen.getByText("¥5,000")).toBeInTheDocument();
    expect(screen.getByText("外食")).toBeInTheDocument();
    expect(screen.getByText("¥10,000")).toBeInTheDocument();
  });

  it("データがない場合は「データなし」を表示する", () => {
    render(<TagSummaryCard tagTotals={[]} onTagDetail={vi.fn()} year={2026} month={3} />);
    expect(screen.getByText("データなし")).toBeInTheDocument();
  });

  it("タグクリックでonTagDetailが呼ばれる", async () => {
    const onTagDetail = vi.fn();
    render(<TagSummaryCard tagTotals={tagTotals} onTagDetail={onTagDetail} year={2026} month={3} />);
    const user = userEvent.setup();
    await user.click(screen.getByText("通勤"));
    expect(onTagDetail).toHaveBeenCalledWith({
      tagId: 1,
      tagName: "通勤",
      year: 2026,
      month: 3,
    });
  });
});
