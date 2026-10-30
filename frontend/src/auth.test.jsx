import { beforeEach, describe, expect, it, vi } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import App from "./App.jsx";
import { AuthProvider } from "./auth.jsx";

const TOKEN_KEY = "inventory.token";
const EMAIL_KEY = "inventory.email";

function renderApp() {
  return render(
    <AuthProvider>
      <App />
    </AuthProvider>
  );
}

function fill(fieldLabel, value) {
  fireEvent.change(screen.getByLabelText(fieldLabel), {
    target: { value },
  });
}

function submit(label) {
  fireEvent.click(screen.getByRole("button", { name: label }));
}

// Router mock: the signed-in view mounts the Products page, which fetches
// categories and products; tests override specific calls via `overrides`.
function apiMock({ login, logout, overrides = [] } = {}) {
  return vi.fn(async (url, opts = {}) => {
    const method = String(opts.method ?? "GET").toUpperCase();
    const respond = (status, body) => ({
      ok: status < 400,
      status,
      json: async () => body,
    });
    if (method === "GET" && url === "/api/v1/categories") {
      return respond(200, { categories: [] });
    }
    if (method === "GET" && url === "/api/v1/products") {
      return respond(200, { products: [] });
    }
    if (method === "POST" && url === "/api/v1/auth/login") {
      return login ? login() : respond(200, { token: "jwt-token-from-login", expires_at: "x" });
    }
    if (method === "POST" && url === "/api/v1/auth/logout") {
      return logout ? logout() : respond(200, { message: "logged out" });
    }
    for (const o of overrides) {
      if (method === o.method && url === o.url) return respond(o.status, o.body);
    }
    return respond(404, { error: "not found" });
  });
}

beforeEach(() => {
  localStorage.clear();
  vi.restoreAllMocks();
});

describe("registration form (client-side validation, AC-03)", () => {
  function openRegister() {
    renderApp();
    submit(/register/i); // switch to the register view
  }

  it("blocks submission with an invalid email and never calls the API", async () => {
    const spy = apiMock();
    global.fetch = spy;
    openRegister();
    fill("Email", "not-an-email");
    fill("Password", "GoodPassword123");
    submit("Register");
    await waitFor(() =>
      expect(screen.getAllByText(/enter a valid email address/i).length).toBeGreaterThan(0)
    );
    expect(spy.mock.calls.filter(([u, o]) => o?.method === "POST")).toHaveLength(0);
  });

  it("blocks submission with a short password and never calls the API", async () => {
    const spy = apiMock();
    global.fetch = spy;
    openRegister();
    fill("Email", "newuser@example.test");
    fill("Password", "short");
    submit("Register");
    await waitFor(() =>
      expect(screen.getAllByText(/password must be at least 8 characters/i).length).toBeGreaterThan(0)
    );
    expect(spy.mock.calls.filter(([u, o]) => o?.method === "POST")).toHaveLength(0);
  });

  it("blocks registration with both fields blank and never calls the API", async () => {
    const spy = apiMock();
    global.fetch = spy;
    openRegister();
    submit("Register"); // both fields left empty
    await waitFor(() =>
      expect(screen.getAllByText(/enter a valid email address/i).length).toBeGreaterThan(0)
    );
    expect(screen.getByText(/password must be at least 8 characters/i)).toBeInTheDocument();
    expect(spy.mock.calls.filter(([u, o]) => o?.method === "POST")).toHaveLength(0);
  });

  it("shows the backend's duplicate-email message (AC-15)", async () => {
    global.fetch = apiMock({
      overrides: [
        {
          method: "POST",
          url: "/api/v1/auth/register",
          status: 409,
          body: { error: "email is already registered" },
        },
      ],
    });
    openRegister();
    fill("Email", "taken@example.test");
    fill("Password", "GoodPassword123");
    submit("Register");
    await waitFor(() =>
      expect(screen.getByText(/email is already registered/i)).toBeInTheDocument()
    );
  });

  it("shows a success panel after a successful registration", async () => {
    global.fetch = apiMock({
      overrides: [
        {
          method: "POST",
          url: "/api/v1/auth/register",
          status: 201,
          body: { id: 1, email: "newuser@example.test", created_at: "2026-10-28T00:00:00Z" },
        },
      ],
    });
    openRegister();
    fill("Email", "newuser@example.test");
    fill("Password", "GoodPassword123");
    submit("Register");
    await waitFor(() => expect(screen.getByText(/account created/i)).toBeInTheDocument());
    expect(screen.getByText(/newuser@example.test/)).toBeInTheDocument();
  });
});

describe("login and logout flow", () => {
  it("logs in, stores the token, shows the signed-in view, and logs out", async () => {
    const spy = apiMock();
    global.fetch = spy;

    renderApp();
    fill("Email", "User@Example.Test");
    fill("Password", "TestPassword123");
    submit(/log in$/i);

    await waitFor(() =>
      expect(screen.getByText(/signed in as user@example.test/i)).toBeInTheDocument()
    );
    expect(localStorage.getItem(TOKEN_KEY)).toBe("jwt-token-from-login");
    expect(localStorage.getItem(EMAIL_KEY)).toBe("user@example.test");
    const loginCalls = spy.mock.calls.filter(([u]) => u === "/api/v1/auth/login");
    expect(loginCalls.length).toBe(1);

    submit(/log out/i);
    await waitFor(() => expect(screen.getByText(/need an account/i)).toBeInTheDocument());
    expect(localStorage.getItem(TOKEN_KEY)).toBeNull();
    const logoutCalls = spy.mock.calls.filter(([u, o]) => u === "/api/v1/auth/logout" && (o?.method ?? "GET") === "POST");
    expect(logoutCalls.length).toBe(1);
    expect(logoutCalls[0][1].headers.Authorization).toBe("Bearer jwt-token-from-login");
  });

  it("shows the server's message for a failed login and stays on the form", async () => {
    global.fetch = apiMock({
      login: () => ({ ok: false, status: 401, json: async () => ({ error: "invalid email or password" }) }),
    });
    renderApp();
    fill("Email", "user@example.test");
    fill("Password", "WrongPassword999");
    submit(/log in$/i);
    await waitFor(() => expect(screen.getByText(/invalid email or password/i)).toBeInTheDocument());
    expect(screen.getByLabelText("Email")).toBeInTheDocument(); // still on the login form
    expect(localStorage.getItem(TOKEN_KEY)).toBeNull();
  });

  it("blocks login with both fields blank and never calls the API", async () => {
    const spy = apiMock();
    global.fetch = spy;
    renderApp();
    submit(/log in$/i); // both fields left empty
    await waitFor(() =>
      expect(screen.getAllByText(/enter a valid email address/i).length).toBeGreaterThan(0)
    );
    expect(screen.getByText(/password is required/i)).toBeInTheDocument();
    expect(spy.mock.calls.filter(([u, o]) => o?.method === "POST")).toHaveLength(0);
  });

  it("shows a connection error when the server is unreachable", async () => {
    global.fetch = vi.fn().mockRejectedValue(new Error("network down"));
    renderApp();
    fill("Email", "user@example.test");
    fill("Password", "TestPassword123");
    submit(/log in$/i);
    await waitFor(() =>
      expect(screen.getByText(/cannot reach the server/i)).toBeInTheDocument()
    );
  });

  it("discards a stale local session when a protected call returns 401 (expired token)", async () => {
    localStorage.setItem(TOKEN_KEY, "expired-token");
    localStorage.setItem(EMAIL_KEY, "user@example.test");
    global.fetch = apiMock({
      logout: () => ({ ok: false, status: 401, json: async () => ({ error: "invalid or expired token" }) }),
    });

    renderApp();
    expect(screen.getByText(/signed in as user@example.test/i)).toBeInTheDocument();

    submit(/log out/i);
    await waitFor(() => expect(screen.getByText(/need an account/i)).toBeInTheDocument());
    expect(localStorage.getItem(TOKEN_KEY)).toBeNull();
  });
});
