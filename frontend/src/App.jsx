import { useState } from "react";
import { useAuth } from "./auth.jsx";
import Login from "./pages/Login.jsx";
import Register from "./pages/Register.jsx";
import Categories from "./pages/Categories.jsx";
import Products from "./pages/Products.jsx";
import Stock from "./pages/Stock.jsx";
import Dashboard from "./pages/Dashboard.jsx";

export default function App() {
  const { isAuthenticated, email, logout } = useAuth();
  const [view, setView] = useState("login");
  const [page, setPage] = useState("dashboard");
  const [loggingOut, setLoggingOut] = useState(false);
  const [logoutError, setLogoutError] = useState("");

  if (isAuthenticated) {
    return (
      <main className="page">
        <header className="topbar">
          <span className="brand">Inventory &amp; Stock Management</span>
          <div className="topbar-right">
            <span className="muted">Signed in as {email}</span>
            <button
              type="button"
              disabled={loggingOut}
              onClick={async () => {
                setLogoutError("");
                setLoggingOut(true);
                try {
                  await logout();
                } catch (err) {
                  setLogoutError(err.message);
                } finally {
                  setLoggingOut(false);
                }
              }}
            >
              {loggingOut ? "Logging out…" : "Log out"}
            </button>
          </div>
        </header>
        <nav className="nav">
          {logoutError && (
            <span className="error" role="alert">
              {logoutError}
            </span>
          )}
          {[
            ["dashboard", "Dashboard"],
            ["products", "Products"],
            ["categories", "Categories"],
            ["stock", "Stock"],
          ].map(([key, label]) => (
            <button
              key={key}
              type="button"
              className={page === key ? "active" : ""}
              onClick={() => setPage(key)}
            >
              {label}
            </button>
          ))}
        </nav>
        {page === "dashboard" && <Dashboard />}
        {page === "products" && <Products />}
        {page === "categories" && <Categories />}
        {page === "stock" && <Stock />}
      </main>
    );
  }

  return (
    <main className="page">
      <header className="topbar">
        <span className="brand">Inventory &amp; Stock Management</span>
      </header>
      <section className="card">
        {view === "login" ? (
          <Login onSwitch={() => setView("register")} />
        ) : (
          <Register onSwitch={() => setView("login")} />
        )}
      </section>
    </main>
  );
}
