import { useState } from "react";
import type { ActionTag } from "./api/client";
import { Home } from "./pages/Home";
import { TagSelect } from "./pages/TagSelect";
import { ExpenseInput } from "./pages/ExpenseInput";
import { TagManage } from "./pages/TagManage";

type Page = "home" | "tag-select" | "expense-input" | "tag-manage";

function App() {
  const [page, setPage] = useState<Page>("home");
  const [selectedTag, setSelectedTag] = useState<ActionTag | null>(null);

  switch (page) {
    case "home":
      return <Home onNavigate={(p) => setPage(p as Page)} />;
    case "tag-select":
      return (
        <TagSelect
          onSelect={(tag) => { setSelectedTag(tag); setPage("expense-input"); }}
          onBack={() => setPage("home")}
        />
      );
    case "expense-input":
      return (
        <ExpenseInput
          tag={selectedTag!}
          onBack={() => setPage("tag-select")}
          onSaved={() => setPage("home")}
        />
      );
    case "tag-manage":
      return <TagManage onBack={() => setPage("home")} />;
  }
}

export default App;
