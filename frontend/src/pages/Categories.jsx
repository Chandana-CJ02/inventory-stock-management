import { useEffect, useState } from "react";
import { apiRequest } from "../api.js";
import { useAuth } from "../auth.jsx";

// Category management page (T-206): create, rename, delete, list.
// Blank names are blocked client-side (AC-03 / TC-19); backend messages are
// surfaced readably (AC-15) — e.g. duplicate names and the delete guard for
// categories that still have products (AC-08).
export default function Categories() {
  const { token } = useAuth();
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [loadError, setLoadError] = useState("");
  const [name, setName] = useState("");
  const [nameError, setNameError] = useState("");
  const [actionError, setActionError] = useState("");
  const [busy, setBusy] = useState(false);
  const [editingId, setEditingId] = useState(null);
  const [editName, setEditName] = useState("");
  const [editError, setEditError] = useState("");

  async function load() {
    setLoading(true);
    setLoadError("");
    try {
      const data = await apiRequest("/v1/categories", { token });
      setCategories(data.categories);
    } catch (err) {
      setLoadError(err.message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
  }, []);

  async function onCreate(e) {
    e.preventDefault();
    setNameError("");
    setActionError("");
    if (name.trim() === "") {
      setNameError("Category name is required");
      return;
    }
    setBusy(true);
    try {
      await apiRequest("/v1/categories", { method: "POST", body: { name }, token });
      setName("");
      await load();
    } catch (err) {
      setActionError(err.message);
    } finally {
      setBusy(false);
    }
  }

  function startEdit(c) {
    setEditingId(c.id);
    setEditName(c.name);
    setEditError("");
  }

  async function onSaveEdit(c) {
    setEditError("");
    if (editName.trim() === "") {
      setEditError("Category name is required");
      return;
    }
    try {
      await apiRequest(`/v1/categories/${c.id}`, {
        method: "PUT",
        body: { name: editName },
        token,
      });
      setEditingId(null);
      await load();
    } catch (err) {
      setEditError(err.message);
    }
  }

  async function onDelete(c) {
    if (!window.confirm(`Delete category "${c.name}"?`)) return;
    setActionError("");
    try {
      await apiRequest(`/v1/categories/${c.id}`, { method: "DELETE", token });
      await load();
    } catch (err) {
      setActionError(err.message); // e.g. "category has products and cannot be deleted"
    }
  }

  return (
    <div>
      <div className="page-head">
        <h1>Categories</h1>
        <span className="chip">{categories.length} categor{categories.length === 1 ? "y" : "ies"}</span>
      </div>

      <div className="panel">
      <form onSubmit={onCreate} noValidate>
        <label htmlFor="category-name">New category name</label>
        <input
          id="category-name"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        {nameError && (
          <p className="error" role="alert">
            {nameError}
          </p>
        )}
        {actionError && (
          <p className="error" role="alert">
            {actionError}
          </p>
        )}
        <button type="submit" disabled={busy}>
          {busy ? "Adding…" : "Add category"}
        </button>
      </form>
      </div>

      {loading ? (
        <p>Loading categories…</p>
      ) : loadError ? (
        <p className="error" role="alert">
          {loadError}
        </p>
      ) : categories.length === 0 ? (
        <p className="empty-state">No categories yet — add your first one above.</p>
      ) : (
        <div className="table-wrap">
          <table className="table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Created</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {categories.map((c) => (
              <tr key={c.id}>
                <td>
                  {editingId === c.id ? (
                    <>
                      <input
                        aria-label={`Rename category ${c.name}`}
                        value={editName}
                        onChange={(e) => setEditName(e.target.value)}
                      />
                      {editError && (
                        <p className="error" role="alert">
                          {editError}
                        </p>
                      )}
                    </>
                  ) : (
                    c.name
                  )}
                </td>
                <td>{new Date(c.created_at).toLocaleString()}</td>
                <td>
                  {editingId === c.id ? (
                    <>
                      <button type="button" onClick={() => onSaveEdit(c)}>
                        Save
                      </button>{" "}
                      <button
                        type="button"
                        className="secondary"
                        onClick={() => setEditingId(null)}
                      >
                        Cancel
                      </button>
                    </>
                  ) : (
                    <>
                      <button
                        type="button"
                        className="secondary"
                        onClick={() => startEdit(c)}
                      >
                        Rename
                      </button>{" "}
                      <button type="button" onClick={() => onDelete(c)}>
                        Delete
                      </button>
                    </>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
