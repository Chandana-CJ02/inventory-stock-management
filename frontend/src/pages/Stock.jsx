import { useEffect, useState } from "react";
import { apiRequest } from "../api.js";
import { useAuth } from "../auth.jsx";

// Stock operations + transaction history page (T-207): record Stock In/Out
// movements (AC-09 client-side validation) and browse the immutable history
// with product/type filters (FR-4.4). Movements are created through the
// backend's atomic, concurrency-safe endpoint; history rows are read-only
// (AC-13 — no edit/delete controls exist anywhere in this UI).
const blankForm = { product_id: "", type: "IN", quantity: "", remarks: "" };

export default function Stock() {
  const { token } = useAuth();
  const [products, setProducts] = useState([]);
  const [transactions, setTransactions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [loadError, setLoadError] = useState("");
  const [filterProduct, setFilterProduct] = useState("");
  const [filterType, setFilterType] = useState("");
  const [form, setForm] = useState(blankForm);
  const [fieldErrors, setFieldErrors] = useState({});
  const [formError, setFormError] = useState("");
  const [formNotice, setFormNotice] = useState("");
  const [submitting, setSubmitting] = useState(false);

  async function loadAll() {
    setLoading(true);
    setLoadError("");
    try {
      const [prodData, txData] = await Promise.all([
        apiRequest("/v1/products", { token }),
        apiRequest("/v1/transactions", { token }),
      ]);
      setProducts(prodData.products);
      setTransactions(txData.transactions);
    } catch (err) {
      setLoadError(err.message);
    } finally {
      setLoading(false);
    }
  }

  async function loadHistory(extraQuery = "") {
    setLoadError("");
    try {
      const txData = await apiRequest(`/v1/transactions${extraQuery}`, { token });
      setTransactions(txData.transactions);
    } catch (err) {
      setLoadError(err.message);
    }
  }

  useEffect(() => {
    loadAll();
  }, [token]);

  function applyFilters() {
    const qs = new URLSearchParams();
    if (filterProduct) qs.set("product_id", filterProduct);
    if (filterType) qs.set("type", filterType);
    const query = qs.toString();
    loadHistory(query ? `?${query}` : "");
  }

  function setField(field, value) {
    setForm((f) => ({ ...f, [field]: value }));
  }

  function validate() {
    const p = form;
    const errs = {};
    if (p.product_id === "") errs.product_id = "Choose a product";
    if (p.type !== "IN" && p.type !== "OUT") errs.type = "Type must be IN or OUT";
    const qty = Number(p.quantity);
    if (p.quantity === "" || !Number.isInteger(qty) || qty <= 0) {
      errs.quantity = "Quantity must be a positive whole number";
    }
    if (p.remarks && p.remarks.length > 500) {
      errs.remarks = "Remarks must be at most 500 characters";
    }
    setFieldErrors(errs);
    return Object.keys(errs).length === 0;
  }

  async function onSubmit(e) {
    e.preventDefault();
    setFormError("");
    setFormNotice("");
    if (!validate()) return;
    setSubmitting(true);
    const body = {
      product_id: Number(form.product_id),
      type: form.type,
      quantity: Number(form.quantity),
    };
    if (form.remarks.trim() !== "") body.remarks = form.remarks.trim();
    try {
      const data = await apiRequest("/v1/transactions", { method: "POST", body, token });
      setFormNotice(`Movement recorded — ${form.type} ${form.quantity}, new stock ${data.current_stock}.`);
      setForm((f) => ({ ...f, quantity: "", remarks: "" }));
      await loadHistory(buildFilterQuery());
    } catch (err) {
      setFormError(err.message); // e.g. "insufficient stock"
    } finally {
      setSubmitting(false);
    }
  }

  function buildFilterQuery() {
    const qs = new URLSearchParams();
    if (filterProduct) qs.set("product_id", filterProduct);
    if (filterType) qs.set("type", filterType);
    const query = qs.toString();
    return query ? `?${query}` : "";
  }

  const productName = (id) => products.find((p) => p.id === id)?.name ?? `#${id}`;

  return (
    <div>
      <div className="page-head">
        <h1>Stock movements</h1>
      </div>

      <div className="panel">
      <form onSubmit={onSubmit} noValidate>
        <h2>Record a movement</h2>

        <label htmlFor="movement-product">Product</label>
        <select
          id="movement-product"
          value={form.product_id}
          onChange={(e) => setField("product_id", e.target.value)}
        >
          <option value="">— choose a product —</option>
          {products.map((p) => (
            <option key={p.id} value={p.id}>
              {p.name} (stock: {p.current_stock})
            </option>
          ))}
        </select>
        {fieldErrors.product_id && (
          <p className="error" role="alert">
            {fieldErrors.product_id}
          </p>
        )}

        <label htmlFor="movement-type">Type</label>
        <select
          id="movement-type"
          value={form.type}
          onChange={(e) => setField("type", e.target.value)}
        >
          <option value="IN">Stock In</option>
          <option value="OUT">Stock Out</option>
        </select>
        {fieldErrors.type && (
          <p className="error" role="alert">
            {fieldErrors.type}
          </p>
        )}

        <label htmlFor="movement-quantity">Quantity</label>
        <input
          id="movement-quantity"
          type="number"
          step="1"
          min="1"
          value={form.quantity}
          onChange={(e) => setField("quantity", e.target.value)}
        />
        {fieldErrors.quantity && (
          <p className="error" role="alert">
            {fieldErrors.quantity}
          </p>
        )}

        <label htmlFor="movement-remarks">Remarks (optional)</label>
        <input
          id="movement-remarks"
          type="text"
          value={form.remarks}
          onChange={(e) => setField("remarks", e.target.value)}
        />
        {fieldErrors.remarks && (
          <p className="error" role="alert">
            {fieldErrors.remarks}
          </p>
        )}

        {formError && (
          <p className="error" role="alert">
            {formError}
          </p>
        )}
        {formNotice && <p className="notice">{formNotice}</p>}

        <button type="submit" disabled={submitting}>
          {submitting ? "Recording…" : "Record movement"}
        </button>
      </form>
      </div>

      <h2>Transaction history</h2>
      <div className="filters">
        <label htmlFor="filter-product">Filter by product</label>
        <select
          id="filter-product"
          value={filterProduct}
          onChange={(e) => setFilterProduct(e.target.value)}
        >
          <option value="">All products</option>
          {products.map((p) => (
            <option key={p.id} value={p.id}>
              {p.name}
            </option>
          ))}
        </select>
        <label htmlFor="filter-type">Type</label>
        <select
          id="filter-type"
          value={filterType}
          onChange={(e) => setFilterType(e.target.value)}
        >
          <option value="">All types</option>
          <option value="IN">Stock In</option>
          <option value="OUT">Stock Out</option>
        </select>
        <button type="button" className="secondary" onClick={applyFilters}>
          Apply filters
        </button>
      </div>

      {loading ? (
        <p>Loading history…</p>
      ) : loadError ? (
        <p className="error" role="alert">
          {loadError}
        </p>
      ) : transactions.length === 0 ? (
        <p className="empty-state">No transactions match the current filters.</p>
      ) : (
        <div className="table-wrap">
          <table className="table">
            <thead>
              <tr>
                <th>Product</th>
                <th>Type</th>
                <th className="num">Quantity</th>
                <th>Remarks</th>
                <th>Recorded</th>
              </tr>
            </thead>
            <tbody>
              {transactions.map((t) => (
                <tr key={t.id}>
                  <td>{productName(t.product_id)}</td>
                  <td>
                    <span className={`badge ${t.type === "IN" ? "badge-in" : "badge-out"}`}>
                      {t.type === "IN" ? "Stock In" : "Stock Out"}
                    </span>
                  </td>
                  <td className="num">{t.quantity}</td>
                  <td>{t.remarks ?? "—"}</td>
                  <td>{new Date(t.created_at).toLocaleString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
