import { beforeEach, describe, expect, it, vi } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import Categories from "./Categories.jsx";
import { AuthProvider } from "../auth.jsx";
import { fetchRouter, callsTo } from "../test-api-mock.js";

const existing = {
  id: 1,
  name: "Electronics",
  created_at: "2026-10-28T00:00:00Z",
  updated_at: "2026-10-28T00:00:00Z",
};

// `list` is the mutable backing store served by the GET route; tests mutate
// it inside route body functions to simulate server-side state changes.
let list;

function renderCategories(routes = []) {
  const spy = fetchRouter([
    {
      method: "GET",
      prefix: "/api/v1/categories",
      body: () => ({ categories: list }),
    },
    ...routes,
  ]);
  list = [existing];
  render(
    <AuthProvider>
      <Categories />
    </AuthProvider>
  );
  return spy;
}

function fill(label, value) {
  fireEvent.change(screen.getByLabelText(label), { target: { value } });
}

beforeEach(() => {
  localStorage.setItem("inventory.token", "test-token");
  localStorage.setItem("inventory.email", "user@example.test");
  window.confirm = vi.fn(() => true);
});

describe("Categories page", () => {
  it("loads and lists the user's categories", async () => {
    renderCategories();
    await waitFor(() => expect(screen.getByText("Electronics")).toBeInTheDocument());
  });

  it("TC-19: blocks blank category names before any API call", async () => {
    const spy = renderCategories();
    const before = spy.mock.calls.length;
    fireEvent.click(screen.getByRole("button", { name: /add category/i }));
    expect(await screen.findByText(/category name is required/i)).toBeInTheDocument();
    expect(callsTo(spy, "POST", "/api/v1/categories").length).toBe(0);
    expect(spy.mock.calls.length).toBe(before);
  });

  it("creates a category and refreshes the list", async () => {
    const created = { id: 2, name: "Stationery", created_at: "x", updated_at: "x" };
    const spy = renderCategories([
      {
        method: "POST",
        prefix: "/api/v1/categories",
        status: 201,
        body: () => {
          list = [...list, created];
          return created;
        },
      },
    ]);
    fill("New category name", "Stationery");
    fireEvent.click(screen.getByRole("button", { name: /add category/i }));
    await waitFor(() => expect(screen.getByText("Stationery")).toBeInTheDocument());
    const posts = callsTo(spy, "POST", "/api/v1/categories");
    expect(posts.length).toBe(1);
    expect(JSON.parse(posts[0][1].body)).toEqual({ name: "Stationery" });
    expect(list).toHaveLength(2);
  });

  it("surfaces the backend's duplicate-name message (AC-15)", async () => {
    renderCategories([
      {
        method: "POST",
        prefix: "/api/v1/categories",
        status: 409,
        body: { error: "a category with this name already exists" },
      },
    ]);
    await waitFor(() => expect(screen.getByText("Electronics")).toBeInTheDocument());
    fill("New category name", "Electronics");
    fireEvent.click(screen.getByRole("button", { name: /add category/i }));
    await waitFor(() =>
      expect(screen.getByText(/a category with this name already exists/i)).toBeInTheDocument()
    );
    // the failed create did not add a row
    expect(screen.getAllByText("Electronics").length).toBe(1);
  });

  it("TC-18: surfaces the backend's delete-guard message when products exist", async () => {
    renderCategories([
      {
        method: "DELETE",
        prefix: "/api/v1/categories/1",
        status: 409,
        body: { error: "category has products and cannot be deleted" },
      },
    ]);
    await waitFor(() => expect(screen.getByText("Electronics")).toBeInTheDocument());
    fireEvent.click(screen.getAllByRole("button", { name: /delete/i })[0]);
    await waitFor(() =>
      expect(screen.getByText(/category has products and cannot be deleted/i)).toBeInTheDocument()
    );
    // The category is still listed after the failed delete.
    expect(screen.getByText("Electronics")).toBeInTheDocument();
  });

  it("renames a category via PUT", async () => {
    const renamed = { ...existing, name: "Electronics & Gadgets" };
    const spy = renderCategories([
      {
        method: "PUT",
        prefix: "/api/v1/categories/1",
        body: () => {
          list = [renamed];
          return renamed;
        },
      },
    ]);
    await waitFor(() => expect(screen.getByText("Electronics")).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: /rename/i }));
    fireEvent.change(screen.getByLabelText(/rename category electronics/i), {
      target: { value: "Electronics & Gadgets" },
    });
    fireEvent.click(screen.getByRole("button", { name: /save/i }));
    await waitFor(() =>
      expect(screen.getByText("Electronics & Gadgets")).toBeInTheDocument()
    );
    const puts = callsTo(spy, "PUT", "/api/v1/categories/1");
    expect(puts.length).toBe(1);
    expect(JSON.parse(puts[0][1].body)).toEqual({ name: "Electronics & Gadgets" });
  });
});
