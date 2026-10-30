import { useEffect, useState } from "react";
import { apiRequest } from "../api.js";
import { useAuth } from "../auth.jsx";

// Product management page (T-206): create, edit, deactivate, list with live
// stock. Client-side validation mirrors the backend rules (AC-03 / TC-28
// frontend half); backend errors — duplicate SKU, foreign category, anything
// else — are surfaced readably (AC-15).
const blankForm = {
  name: "",
  sku: "",
  category_id: "",
  unit_price: "",
  opening_stock: "",
  min_stock: "",
};

export default function Products() {
  const { token } = useAuth();
  const [products, setProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [loadError, setLoadError] = useState("");
  const [form, setForm] = useState(blankForm);
  const [fieldErrors, setFieldErrors] = useState({});
  const [formError, setFormError] = useState("");
  const [formNotice, setFormNotice] = useState("");
  const [actionError, setActionError] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [editingId, setEditingId] = useState(null);
  const [busyId, setBusyId] = useState(null);

  async function load() {
    setLoading(true);
    setLoadError("");
    try {
      const [catData, prodData] = await Promise.all([
        apiRequest("/v1/categories", { token }),
        apiRequest("/v1/products", { token }),
      ]);
      setCategories(catData.categories);
      setProducts(prodData.products);
    } catch (err) {
      setLoadError(err.message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
  }, []);

  function setField(field, value) {
    setForm((f) => ({ ...f, [field]: value }));
  }

  function validate() {
    const p = form;
    const errs = {};
    if (p.name.trim() === "") errs.name = "Product name is required";
    if (p.sku.trim() === "") errs.sku = "SKU is required";
    if (p.category_id === "") errs.category_id = "Choose a category";
    const price = Number(p.unit_price);
    if (p.unit_price === "" || Number.isNaN(price) || price < 0) {
      errs.unit_price = "Unit price must be a non-negative number";
    }
    const opening = Number(p.opening_stock);
    if (p.opening_stock === "" || !Number.isInteger(opening) || opening < 0) {
      errs.opening_stock = "Opening stock must be a non-negative whole number";
    }
    const min = Number(p.min_stock);
    if (p.min_stock === "" || !Number.isInteger(min) || min < 0) {
      errs.min_stock = "Minimum stock must be a non-negative whole number";
    }
    setFieldErrors(errs);
    return Object.keys(errs).length === 0;
  }

  function startCreate() {
    setEditingId(null);
    setForm(blankForm);
    setFieldErrors({});
    setFormError("");
    setFormNotice("");
  }

  function startEdit(p) {
    setEditingId(p.id);
    setForm({
      name: p.name,
      sku: p.sku,
      category_id: String(p.category_id),
      unit_price: String(p.unit_price),
      opening_stock: String(p.opening_stock),
      min_stock: String(p.min_stock),
    });
    setFieldErrors({});
    setFormError("");
    setFormNotice("");
  }

  async function onSubmit(e) {
    e.preventDefault();
    setFormError("");
    setFormNotice("");
    if (!validate()) return;
    setSubmitting(true);
    const body = {
      name: form.name.trim(),
      sku: form.sku.trim(),
      category_id: Number(form.category_id),
      unit_price: Number(form.unit_price),
      opening_stock: Number(form.opening_stock),
      min_stock: Number(form.min_stock),
    };
    try {
      if (editingId) {
        await apiRequest(`/v1/products/${editingId}`, { method: "PUT", body, token });
        setFormNotice("Product updated.");
      } else {
        await apiRequest("/v1/products", { method: "POST", body, token });
        setFormNotice("Product created.");
      }
      setForm(blankForm);
      setEditingId(null);
      await load();
    } catch (err) {
      setFormError(err.message); // e.g. "sku already exists in your inventory"
    } finally {
      setSubmitting(false);
    }
  }

  async function setDeactivated(p, inactive) {
    setBusyId(p.id);
    setActionError("");
    try {
      if (inactive) {
        await apiRequest(`/v1/products/${p.id}/deactivate`, { method: "POST", token });
      } else {
        await apiRequest(`/v1/products/${p.id}`, {
          method: "PUT",
          body: { status: "active" },
          token,
        });
      }
      await load();
    } catch (err) {
      setActionError(err.message);
    } finally {
      setBusyId(null);
    }
  }

  const categoryName = (id) =>
    categories.find((c) => c.id === id)?.name ?? `category #${id}`;

  return (
    <div>
      <div className="page-head">
        <h1>Products</h1>
        <span className="chip">{products.length} product{products.length === 1 ? "" : "s"}</span>
      </div>

      <div className="panel">
      <form onSubmit={onSubmit} noValidate>
        <h2>{editingId ? "Edit product" : "Add product"}</h2>

        <label htmlFor="product-name">Name</label>
        <input
          id="product-name"
          type="text"
          value={form.name}
          onChange={(e) => setField("name", e.target.value)}
        />
        {fieldErrors.name && (
          <p className="error" role="alert">
            {fieldErrors.name}
          </p>
        )}

        <label htmlFor="product-sku">SKU</label>
        <input
          id="product-sku"
          type="text"
          value={form.sku}
          onChange={(e) => setField("sku", e.target.value)}
        />
        {fieldErrors.sku && (
          <p className="error" role="alert">
            {fieldErrors.sku}
          </p>
        )}

        <label htmlFor="product-category">Category</label>
        <select
          id="product-category"
          value={form.category_id}
          onChange={(e) => setField("category_id", e.target.value)}
        >
          <option value="">— choose a category —</option>
          {categories.map((c) => (
            <option key={c.id} value={c.id}>
              {c.name}
            </option>
          ))}
        </select>
        {fieldErrors.category_id && (
          <p className="error" role="alert">
            {fieldErrors.category_id}
          </p>
        )}

        <label htmlFor="product-price">Unit price</label>
        <input
          id="product-price"
          type="number"
          step="0.01"
          min="0"
          value={form.unit_price}
          onChange={(e) => setField("unit_price", e.target.value)}
        />
        {fieldErrors.unit_price && (
          <p className="error" role="alert">
            {fieldErrors.unit_price}
          </p>
        )}

        <label htmlFor="product-opening">Opening stock</label>
        <input
          id="product-opening"
          type="number"
          step="1"
          min="0"
          value={form.opening_stock}
          onChange={(e) => setField("opening_stock", e.target.value)}
        />
        <p className="hint">Opening stock is recorded once at creation and is not changed by edits — stock movements update the current stock.</p>
        {fieldErrors.opening_stock && (
          <p className="error" role="alert">
            {fieldErrors.opening_stock}
          </p>
        )}

        <label htmlFor="product-min">Minimum stock</label>
        <input
          id="product-min"
          type="number"
          step="1"
          min="0"
          value={form.min_stock}
          onChange={(e) => setField("min_stock", e.target.value)}
        />
        {fieldErrors.min_stock && (
          <p className="error" role="alert">
            {fieldErrors.min_stock}
          </p>
        )}

        {formError && (
          <p className="error" role="alert">
            {formError}
          </p>
        )}
        {formNotice && <p className="notice">{formNotice}</p>}

        <div className="actions">
          <button type="submit" disabled={submitting}>
            {submitting
              ? "Saving…"
              : editingId
                ? "Save changes"
                : "Add product"}
          </button>
          {editingId && (
            <button
              type="button"
              className="secondary"
              onClick={startCreate}
              disabled={submitting}
            >
              Cancel edit
            </button>
          )}
        </div>
      </form>
      </div>

      {actionError && (
        <p className="error" role="alert">
          {actionError}
        </p>
      )}

      {loading ? (
        <p>Loading products…</p>
      ) : loadError ? (
        <p className="error" role="alert">
          {loadError}
        </p>
      ) : products.length === 0 ? (
        <p className="empty-state">No products yet — add your first one above.</p>
      ) : (
        <div className="table-wrap">
          <table className="table">
            <thead>
              <tr>
                <th>Name</th>
                <th>SKU</th>
                <th>Category</th>
                <th className="num">Unit price</th>
                <th className="num">Current stock</th>
                <th className="num">Min stock</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {products.map((p) => {
                const low = p.status === "active" && p.current_stock <= p.min_stock;
                return (
                  <tr key={p.id} className={low ? "row-low" : undefined}>
                    <td>{p.name}</td>
                    <td>{p.sku}</td>
                    <td>{categoryName(p.category_id)}</td>
                    <td className="num">{p.unit_price.toFixed(2)}</td>
                    <td className="num">
                      {p.current_stock}
                      {low && <span className="badge badge-low">low</span>}
                    </td>
                    <td className="num">{p.min_stock}</td>
                    <td>
                      <span className={`badge ${p.status === "active" ? "badge-active" : "badge-inactive"}`}>
                        {p.status}
                      </span>
                    </td>
                    <td>
                      <button
                        type="button"
                        className="secondary"
                        onClick={() => startEdit(p)}
                        disabled={busyId === p.id}
                      >
                        Edit
                      </button>{" "}
                      {p.status === "active" ? (
                        <button
                          type="button"
                          onClick={() => setDeactivated(p, true)}
                          disabled={busyId === p.id}
                        >
                          Deactivate
                        </button>
                      ) : (
                        <button
                          type="button"
                          onClick={() => setDeactivated(p, false)}
                          disabled={busyId === p.id}
                        >
                          Reactivate
                        </button>
                      )}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
