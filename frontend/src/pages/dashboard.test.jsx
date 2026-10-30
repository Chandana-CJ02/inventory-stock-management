import { beforeEach, describe, expect, it } from "vitest";
import { render, screen, within, waitFor } from "@testing-library/react";
import Dashboard from "./Dashboard.jsx";
import { AuthProvider } from "../auth.jsx";
import { fetchRouter } from "../test-api-mock.js";

const summaryFull = {
  totals: { active_products: 3, categories: 2, low_stock: 2 },
  low_stock: [
    { id: 12, category_id: 6, name: "Low Widget", sku: "LW-1", current_stock: 2, min_stock: 5 },
    { id: 14, category_id: 6, name: "Boundary Item", sku: "BND-1", current_stock: 5, min_stock: 5 },
  ],
  recent_transactions: [
    { id: 31, product_id: 11, type: "IN", quantity: 3, remarks: null, created_at: "2026-10-28T11:00:00Z" },
    { id: 30, product_id: 11, type: "OUT", quantity: 1, remarks: "sold", created_at: "2026-10-28T10:30:00Z" },
  ],
};

const summaryEmpty = {
  totals: { active_products: 0, categories: 0, low_stock: 0 },
  low_stock: [],
  recent_transactions: [],
};

function renderDashboard(summary) {
  const spy = fetchRouter([
    { method: "GET", prefix: "/api/v1/dashboard/summary", body: summary },
    {
      method: "GET",
      prefix: "/api/v1/products",
      body: {
        products: [
          { id: 11, name: "Widget", category_id: 1, sku: "W", unit_price: 1, opening_stock: 0, current_stock: 0, min_stock: 0, status: "active" },
          { id: 12, name: "Low Widget", category_id: 1, sku: "LW", unit_price: 1, opening_stock: 0, current_stock: 0, min_stock: 0, status: "active" },
        ],
      },
    },
  ]);
  render(
    <AuthProvider>
      <Dashboard />
    </AuthProvider>
  );
  return spy;
}

beforeEach(() => {
  localStorage.setItem("inventory.token", "test-token");
  localStorage.setItem("inventory.email", "user@example.test");
});

describe("Dashboard page (T-207, US-17 UI)", () => {
  it("renders the three totals (FR-5.1)", async () => {
    renderDashboard(summaryFull);
    await waitFor(() =>
      expect(screen.getByTestId("stat-active-products")).toBeInTheDocument()
    );
    expect(
      within(screen.getByTestId("stat-active-products")).getByText("3")
    ).toBeInTheDocument();
    expect(
      within(screen.getByTestId("stat-categories")).getByText("2")
    ).toBeInTheDocument();
    expect(
      within(screen.getByTestId("stat-low-stock")).getByText("2")
    ).toBeInTheDocument();
    expect(screen.getByRole("heading", { name: "Low-stock products" })).toBeInTheDocument();
  });

  it("renders the low-stock list with boundary item included (FR-5.2)", async () => {
    renderDashboard(summaryFull);
    await waitFor(() =>
      expect(screen.getByRole("heading", { name: "Low-stock products" })).toBeInTheDocument()
    );
    expect(screen.getByText("Low Widget")).toBeInTheDocument();
    expect(screen.getByText("Boundary Item")).toBeInTheDocument();
    expect(screen.getByText("LW-1")).toBeInTheDocument();
    expect(screen.getByText("BND-1")).toBeInTheDocument();
  });

  it("renders the recent movements with product names mapped (FR-5.3)", async () => {
    renderDashboard(summaryFull);
    await waitFor(() =>
      expect(screen.getByRole("heading", { name: "Recent transactions" })).toBeInTheDocument()
    );
    // Both recent movements belong to product 11 → both rows read "Widget".
    expect(screen.getAllByText("Widget").length).toBe(2);
    expect(screen.getByText("sold")).toBeInTheDocument();
    const ins = screen.getAllByText("Stock In");
    expect(ins.length).toBe(1);
    const outs = screen.getAllByText("Stock Out");
    expect(outs.length).toBe(1);
  });

  it("renders an empty dashboard without errors", async () => {
    renderDashboard(summaryEmpty);
    await waitFor(() =>
      expect(screen.getByText("Active products")).toBeInTheDocument()
    );
    expect(screen.getByText("No low-stock products — everything is at or above its minimum.")).toBeInTheDocument();
    expect(screen.getByText(/no stock movements recorded yet/i)).toBeInTheDocument();
  });

  it("shows a readable error when the summary cannot be loaded", async () => {
    const spy = fetchRouter([
      {
        method: "GET",
        prefix: "/api/v1/dashboard/summary",
        status: 500,
        body: { error: "internal error" },
      },
    ]);
    render(
      <AuthProvider>
        <Dashboard />
      </AuthProvider>
    );
    await waitFor(() => expect(screen.getByText(/internal error/i)).toBeInTheDocument());
    expect(spy).toHaveBeenCalled();
  });
});
