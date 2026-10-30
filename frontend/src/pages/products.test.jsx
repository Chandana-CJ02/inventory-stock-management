import { beforeEach, describe, expect, it } from "vitest";
import { render, screen, within, fireEvent, waitFor } from "@testing-library/react";
import Products from "./Products.jsx";
import { AuthProvider } from "../auth.jsx";
import { fetchRouter, callsTo } from "../test-api-mock.js";

const category = { id: 5, name: "Electronics", created_at: "x", updated_at: "x" };
const product = {
  id: 7,
  category_id: 5,
  name: "Widget",
  sku: "WID-1",
  unit_price: 10.5,
  opening_stock: 10,
  current_stock: 7,
  min_stock: 3,
  status: "active",
  created_at: "x",
  updated_at: "x",
};

function baseRoutes(extra = []) {
  return [
    { method: "GET", prefix: "/api/v1/categories", body: { categories: [category] } },
    { method: "GET", prefix: "/api/v1/products", body: { products: [product] } },
    ...extra,
  ];
}

function renderProducts(extraRoutes = []) {
  const spy = fetchRouter([...baseRoutes(), ...extraRoutes]);
  render(
    <AuthProvider>
      <Products />
    </AuthProvider>
  );
  return spy;
}

function fill(label, value) {
  fireEvent.change(screen.getByLabelText(label), { target: { value } });
}

async function submitProduct(body, extraRoutes = []) {
  const spy = renderProducts([
    {
      method: "POST",
      prefix: "/api/v1/products",
      status: 201,
      body: { ...body, id: 99, current_stock: body.opening_stock, status: "active" },
    },
    ...extraRoutes,
  ]);
  // Wait for the categories GET to resolve so the select has its options.
  await waitFor(() =>
    expect(
      within(screen.getByLabelText("Category")).getAllByRole("option", { name: "Electronics" }).length
    ).toBeGreaterThan(0)
  );
  fill("Name", body.name);
  fill("SKU", body.sku);
  fireEvent.change(screen.getByLabelText("Category"), { target: { value: String(body.category_id) } });
  fill("Unit price", String(body.unit_price));
  fill("Opening stock", String(body.opening_stock));
  fill("Minimum stock", String(body.min_stock));
  fireEvent.click(screen.getByRole("button", { name: /add product/i }));
  return spy;
}

beforeEach(() => {
  localStorage.setItem("inventory.token", "test-token");
  localStorage.setItem("inventory.email", "user@example.test");
});

describe("Products page", () => {
  it("lists products with current stock, category name, and low-stock badge", async () => {
    renderProducts();
    await waitFor(() => expect(screen.getByText("Widget")).toBeInTheDocument());
    expect(screen.getByText("WID-1")).toBeInTheDocument();
    // Category name appears in the select option and the table cell.
    expect(screen.getAllByText("Electronics").length).toBeGreaterThanOrEqual(1);
    expect(screen.getByText("7")).toBeInTheDocument();
    expect(screen.getByText("10.50")).toBeInTheDocument();
    // current (7) > min (3): no low badge
    expect(screen.queryByText("low")).not.toBeInTheDocument();
  });

  it("TC-28 (frontend): blocks creation when required fields are missing and never calls the API", async () => {
    const spy = renderProducts();
    const before = spy.mock.calls.length;
    fireEvent.click(screen.getByRole("button", { name: /add product/i }));
    await waitFor(() =>
      expect(screen.getAllByText(/is required|must be a non-negative/i).length).toBeGreaterThan(0)
    );
    expect(callsTo(spy, "POST", "/api/v1/products").length).toBe(0);
    expect(spy.mock.calls.length).toBe(before);
  });

  it("TC-28 (frontend): blocks negative price and fractional opening stock", async () => {
    const spy = renderProducts();
    await waitFor(() =>
      expect(within(screen.getByLabelText("Category")).getAllByRole("option", { name: "Electronics" }).length).toBeGreaterThan(0)
    ); // category options loaded
    fill("Name", "P");
    fill("SKU", "P-1");
    fireEvent.change(screen.getByLabelText("Category"), { target: { value: "5" } });
    fill("Unit price", "-1");
    fill("Opening stock", "3");
    fill("Minimum stock", "0");
    fireEvent.click(screen.getByRole("button", { name: /add product/i }));
    await waitFor(() =>
      expect(screen.getByText(/unit price must be a non-negative/i)).toBeInTheDocument()
    );
    fill("Unit price", "1");
    fill("Opening stock", "2.5");
    fireEvent.click(screen.getByRole("button", { name: /add product/i }));
    await waitFor(() =>
      expect(screen.getByText(/opening stock must be a non-negative whole/i)).toBeInTheDocument()
    );
    expect(callsTo(spy, "POST", "/api/v1/products").length).toBe(0);
  });

  it("creates a product with the exact payload and refreshes the list", async () => {
    const spy = await submitProduct({
      name: "Gadget",
      sku: "GDG-1",
      category_id: 5,
      unit_price: 4,
      opening_stock: 6,
      min_stock: 2,
    });
    await waitFor(() => expect(screen.getByText(/product created/i)).toBeInTheDocument());
    const posts = callsTo(spy, "POST", "/api/v1/products");
    expect(posts.length).toBe(1);
    expect(JSON.parse(posts[0][1].body)).toEqual({
      name: "Gadget",
      sku: "GDG-1",
      category_id: 5,
      unit_price: 4,
      opening_stock: 6,
      min_stock: 2,
    });
  });

  it("surfaces the backend's duplicate-SKU message (AC-15)", async () => {
    renderProducts([
      {
        method: "POST",
        prefix: "/api/v1/products",
        status: 409,
        body: { error: "sku already exists in your inventory" },
      },
    ]);
    await waitFor(() =>
      expect(within(screen.getByLabelText("Category")).getAllByRole("option", { name: "Electronics" }).length).toBeGreaterThan(0)
    ); // category options loaded
    fill("Name", "Dup");
    fill("SKU", "WID-1");
    fireEvent.change(screen.getByLabelText("Category"), { target: { value: "5" } });
    fill("Unit price", "1");
    fill("Opening stock", "1");
    fill("Minimum stock", "1");
    fireEvent.click(screen.getByRole("button", { name: /add product/i }));
    await waitFor(() =>
      expect(screen.getByText(/sku already exists in your inventory/i)).toBeInTheDocument()
    );
  });

  it("deactivates a product via the dedicated endpoint", async () => {
    const spy = renderProducts([
      {
        method: "POST",
        prefix: "/api/v1/products/7/deactivate",
        body: { ...product, status: "inactive" },
      },
    ]);
    await waitFor(() => expect(screen.getByText("Widget")).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: /^deactivate$/i }));
    await waitFor(() =>
      expect(callsTo(spy, "POST", "/api/v1/products/7/deactivate").length).toBe(1)
    );
  });
});
