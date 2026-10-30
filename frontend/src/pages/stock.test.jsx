import { beforeEach, describe, expect, it } from "vitest";
import { render, screen, within, fireEvent, waitFor } from "@testing-library/react";
import Stock from "./Stock.jsx";
import { AuthProvider } from "../auth.jsx";
import { fetchRouter, callsTo } from "../test-api-mock.js";

const products = [
  { id: 11, category_id: 5, name: "Widget", sku: "WID-1", unit_price: 1, opening_stock: 10, current_stock: 15, min_stock: 3, status: "active" },
  { id: 12, category_id: 5, name: "Gadget", sku: "GDG-1", unit_price: 1, opening_stock: 4, current_stock: 2, min_stock: 1, status: "active" },
];
let movements = [
  { id: 21, product_id: 11, type: "IN", quantity: 5, remarks: null, created_at: "2026-10-28T09:00:00Z" },
  { id: 22, product_id: 11, type: "OUT", quantity: 2, remarks: "sold", created_at: "2026-10-28T10:00:00Z" },
];

function renderStock(extraRoutes = []) {
  const spy = fetchRouter([
    { method: "GET", prefix: "/api/v1/products", body: { products } },
    {
      method: "GET",
      prefix: "/api/v1/transactions",
      body: () => ({ transactions: movements }),
    },
    {
      method: "POST",
      prefix: "/api/v1/transactions",
      status: 201,
      body: (req) => {
        const b = JSON.parse(req.opts.body);
        const created = {
          id: 99 + movements.length,
          product_id: b.product_id,
          type: b.type,
          quantity: b.quantity,
          remarks: b.remarks ?? null,
          created_at: "2026-10-28T11:00:00Z",
          current_stock: 99,
        };
        movements = [...movements, created];
        return created;
      },
    },
    ...extraRoutes,
  ]);
  render(
    <AuthProvider>
      <Stock />
    </AuthProvider>
  );
  return spy;
}

function fill(label, value) {
  fireEvent.change(screen.getByLabelText(label), { target: { value } });
}

function chooseOption(label, value) {
  fireEvent.change(screen.getByLabelText(label), { target: { value } });
}

beforeEach(() => {
  localStorage.setItem("inventory.token", "test-token");
  localStorage.setItem("inventory.email", "user@example.test");
  movements = [
    { id: 21, product_id: 11, type: "IN", quantity: 5, remarks: null, created_at: "2026-10-28T09:00:00Z" },
    { id: 22, product_id: 11, type: "OUT", quantity: 2, remarks: "sold", created_at: "2026-10-28T10:00:00Z" },
  ];
});

describe("Stock page (T-207)", () => {
  it("renders the movement form, product options, and history rows", async () => {
    renderStock();
    await waitFor(() =>
      expect(within(screen.getByLabelText("Product")).getAllByRole("option", { name: /Widget/ }).length).toBeGreaterThan(0)
    );
    const historyTable = screen.getByRole("table");
    expect(within(historyTable).getAllByText("Stock In").length).toBeGreaterThan(0);
    expect(within(historyTable).getAllByText("Stock Out").length).toBeGreaterThan(0);
    expect(within(historyTable).getByText(/sold/)).toBeInTheDocument();
  });

  it("blocks submission with an empty form and never calls the API", async () => {
    const spy = renderStock();
    const before = callsTo(spy, "POST", "/api/v1/transactions").length;
    fireEvent.click(screen.getByRole("button", { name: /record movement/i }));
    expect(await screen.findByText("Choose a product")).toBeInTheDocument();
    expect(screen.getByText("Quantity must be a positive whole number")).toBeInTheDocument();
    expect(callsTo(spy, "POST", "/api/v1/transactions").length).toBe(before);
  });

  it("blocks zero, negative, and fractional quantities", async () => {
    renderStock();
    await screen.findByText("Record movement");
    for (const qty of ["0", "-2", "2.5"]) {
      chooseOption("Product", "11");
      fireEvent.change(screen.getByLabelText("Type", { selector: "#movement-type" }), { target: { value: "IN" } });
      fill("Quantity", qty);
      fireEvent.click(screen.getByRole("button", { name: /record movement/i }));
      await waitFor(() =>
        expect(screen.getByText("Quantity must be a positive whole number")).toBeInTheDocument()
      );
    }
  });

  it("records a Stock In with remarks and refreshes the history", async () => {
    const spy = renderStock();
    await screen.findByText("Record movement");
    chooseOption("Product", "11");
    fill("Quantity", "3");
    fill("Remarks (optional)", "delivery");
    fireEvent.click(screen.getByRole("button", { name: /record movement/i }));
    await waitFor(() => expect(screen.getByText(/delivery/)).toBeInTheDocument());
    const posts = callsTo(spy, "POST", "/api/v1/transactions");
    expect(posts.length).toBe(1);
    expect(JSON.parse(posts[0][1].body)).toEqual({
      product_id: 11,
      type: "IN",
      quantity: 3,
      remarks: "delivery",
    });
  });

  it("surfaces the backend's insufficient-stock message (AC-15)", async () => {
    renderStock([
      {
        method: "POST",
        prefix: "/api/v1/transactions",
        status: 409,
        body: { error: "insufficient stock" },
      },
    ]);
    await screen.findByText("Record movement");
    chooseOption("Product", "11");
    fill("Quantity", "999");
    fireEvent.click(screen.getByRole("button", { name: /record movement/i }));
    await waitFor(() => expect(screen.getByText(/insufficient stock/i)).toBeInTheDocument());
  });

  it("filters the history by type (AC-15/FR-4.4)", async () => {
    const spy = renderStock();
    await screen.findByText("Record movement");
    fireEvent.change(screen.getByLabelText("Type", { selector: "#filter-type" }), {
      target: { value: "OUT" },
    });
    fireEvent.click(screen.getByRole("button", { name: /apply filters/i }));
    await waitFor(() =>
      expect(callsTo(spy, "GET", "/api/v1/transactions?type=OUT").length).toBe(1)
    );
  });
});
