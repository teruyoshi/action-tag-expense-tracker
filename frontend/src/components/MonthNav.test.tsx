import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MonthNav } from "./MonthNav";

describe("MonthNav", () => {
  it("年月を表示する", () => {
    render(<MonthNav year={2026} month={3} onChangeMonth={vi.fn()} />);
    expect(screen.getByText("2026年3月")).toBeInTheDocument();
  });

  it("前月ボタンでonChangeMonth(-1)が呼ばれる", async () => {
    const onChangeMonth = vi.fn();
    render(<MonthNav year={2026} month={3} onChangeMonth={onChangeMonth} />);
    const user = userEvent.setup();
    await user.click(screen.getByText("<"));
    expect(onChangeMonth).toHaveBeenCalledWith(-1);
  });

  it("次月ボタンでonChangeMonth(1)が呼ばれる", async () => {
    const onChangeMonth = vi.fn();
    render(<MonthNav year={2026} month={3} onChangeMonth={onChangeMonth} />);
    const user = userEvent.setup();
    await user.click(screen.getByText(">"));
    expect(onChangeMonth).toHaveBeenCalledWith(1);
  });
});
