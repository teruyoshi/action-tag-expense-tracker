import { useState } from "react";
import type { ActionTag } from "./api/client";
import { Home } from "./pages/Home";
import { TagSelect } from "./pages/TagSelect";
import { ExpenseInput } from "./pages/ExpenseInput";
import { TagManage } from "./pages/TagManage";
import { TagDetails } from "./pages/TagDetails";

type Page = "home" | "tag-select" | "expense-input" | "tag-manage" | "tag-details";

interface TagDetailParams {
  tagId: number;
  tagName: string;
  year: number;
  month: number;
}

function App() {
  const [page, setPage] = useState<Page>("home");
  const [selectedTag, setSelectedTag] = useState<ActionTag | null>(null);
  const [tagDetailParams, setTagDetailParams] = useState<TagDetailParams | null>(null);

  switch (page) {
    case "home":
      return (
        <Home
          onNavigate={(p) => setPage(p as Page)}
          onTagDetail={(params) => { setTagDetailParams(params); setPage("tag-details"); }}
        />
      );
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
    case "tag-details":
      return (
        <TagDetails
          tagId={tagDetailParams!.tagId}
          tagName={tagDetailParams!.tagName}
          year={tagDetailParams!.year}
          month={tagDetailParams!.month}
          onBack={() => setPage("home")}
        />
      );
  }
}

export default App;
