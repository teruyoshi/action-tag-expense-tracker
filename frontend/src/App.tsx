import { Routes, Route } from "react-router-dom";
import { Home } from "./pages/Home";
import { TagSelect } from "./pages/TagSelect";
import { ExpenseInput } from "./pages/ExpenseInput";
import { TagManage } from "./pages/TagManage";
import { TagDetails } from "./pages/TagDetails";

function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/tags/select" element={<TagSelect />} />
      <Route path="/expense/new" element={<ExpenseInput />} />
      <Route path="/tags/manage" element={<TagManage />} />
      <Route path="/tags/:tagId/details" element={<TagDetails />} />
    </Routes>
  );
}

export default App;
