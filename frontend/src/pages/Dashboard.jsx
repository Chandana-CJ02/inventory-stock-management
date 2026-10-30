import { useEffect, useState } from "react";
import { apiRequest } from "../api.js";
import { useAuth } from "../auth.jsx";

// Dashboard page (T-207, rendering US-17 / FR-5.1…FR-5.3): inventory totals,
// the low-stock list, and the five most recent movements — all served by the
// backend's user-scoped summary endpoint.
export default function Dashboard() {
  const { token } = useAuth();
  const [summary, setSummary] = useState(null);
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [loadError, setLoadError] = useState("");

  async function load() {
    setLoading(true);
    setLoadError("");
    try {
      const [summaryData, prodData] = await Promise.all([
        apiRequest("/v1/dashboard/summary", { token }),
        apiRequest("/v1/products", { token }),
      ]);
      setSummary(summaryData);
      setProducts(prodData.products);
    } catch (err) {
      setLoadError(err.message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
  }, [token]);

  const productName = (id) => products.find((p) => p.id === id)?.name ?? `#${id}`;

  if (loading) return <p>Loading dashboard…</p>;
  if (loadError)
    return (
      <div>
        <h1>Dashboard</h1>
        <p className="error" role="alert">
          {loadError}
        </p>
      </div>
    );

  return (
    <div>
      <h1>Dashboard</h1>

      <div className="cards">
        <div className="card stat" data-testid="stat-active-products">
          <span className="stat-value">{summary.totals.active_products}</span>
          <span className="stat-label">Active products</span>
        </div>
        <div className="card stat" data-testid="stat-categories">
          <span className="stat-value">{summary.totals.categories}</span>
          <span className="stat-label">Categories</span>
        </div>
        <div className="card stat" data-testid="stat-low-stock">
          <span className="stat-value">{summary.totals.low_stock}</span>
          <span className="stat-label">Low-stock products</span>
        </div>
      </div>

      <h2>Low-stock products</h2>
      {summary.low_stock.length === 0 ? (
        <p className="empty-state">No low-stock products — everything is at or above its minimum.</p>
      ) : (
        <div className="table-wrap">
          <table className="table">
            <thead>
              <tr>
                <th>Name</th>
                <th>SKU</th>
                <th className="num">Current stock</th>
                <th className="num">Minimum stock</th>
              </tr>
            </thead>
            <tbody>
              {summary.low_stock.map((p) => (
                <tr key={p.id} className="row-low">
                  <td>{p.name}</td>
                  <td>{p.sku}</td>
                  <td className="num">{p.current_stock}</td>
                  <td className="num">{p.min_stock}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <h2>Recent transactions</h2>
      {summary.recent_transactions.length === 0 ? (
        <p className="empty-state">No stock movements recorded yet — record one on the Stock page.</p>
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
              {summary.recent_transactions.map((t) => (
                <tr key={t.id}>
                  <td>{productName(t.product_id)}</td>
                  <td>{t.type === "IN" ? "Stock In" : "Stock Out"}</td>
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
