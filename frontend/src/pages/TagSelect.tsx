import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { api } from "../api/client";
import type { ActionTag } from "../api/client";

export function TagSelect() {
  const navigate = useNavigate();
  const [tags, setTags] = useState<ActionTag[]>([]);

  useEffect(() => {
    api.getTags().then(setTags);
  }, []);

  const handleSelect = (tag: ActionTag) => {
    navigate("/expense/new", { state: { tag } });
  };

  return (
    <div>
      <button className="btn-back" onClick={() => navigate("/")}>&larr; 戻る</button>
      <h1>行動タグを選択</h1>
      {tags.length === 0 ? (
        <p className="empty">タグがありません。タグ管理から作成してください。</p>
      ) : (
        <div className="tag-grid">
          {tags.map((tag) => (
            <button key={tag.id} className="card tag-card" onClick={() => handleSelect(tag)}>
              {tag.name}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
