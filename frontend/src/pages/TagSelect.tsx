import { useEffect, useState } from "react";
import { api } from "../api/client";
import type { ActionTag } from "../api/client";

interface Props {
  onSelect: (tag: ActionTag) => void;
  onBack: () => void;
}

export function TagSelect({ onSelect, onBack }: Props) {
  const [tags, setTags] = useState<ActionTag[]>([]);

  useEffect(() => {
    api.getTags().then(setTags);
  }, []);

  return (
    <div>
      <button className="btn-back" onClick={onBack}>&larr; 戻る</button>
      <h1>行動タグを選択</h1>
      {tags.length === 0 ? (
        <p className="empty">タグがありません。タグ管理から作成してください。</p>
      ) : (
        <div className="tag-grid">
          {tags.map((tag) => (
            <button key={tag.id} className="card tag-card" onClick={() => onSelect(tag)}>
              {tag.name}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
