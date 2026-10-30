// Thin fetch wrapper around the backend REST API. The API base defaults to
// /api (same-origin via the Vite dev proxy); it can be overridden with
// VITE_API_BASE for other environments. No secrets ever live in this bundle —
// the JWT_SECRET stays on the server, and the session JWT lives only in
// localStorage (approved trade-off, docs/architecture.md §9/Q11).

const API_BASE = import.meta.env.VITE_API_BASE ?? "/api";

export class ApiError extends Error {
  constructor(status, message) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

export async function apiRequest(path, { method = "GET", body, token } = {}) {
  const headers = {};
  if (body !== undefined) headers["Content-Type"] = "application/json";
  if (token) headers["Authorization"] = `Bearer ${token}`;

  let res;
  try {
    res = await fetch(`${API_BASE}${path}`, {
      method,
      headers,
      body: body === undefined ? undefined : JSON.stringify(body),
    });
  } catch {
    throw new ApiError(0, "cannot reach the server — check your connection");
  }

  let data = null;
  try {
    data = await res.json();
  } catch {
    // non-JSON body: fall through with data = null
  }

  if (!res.ok) {
    throw new ApiError(res.status, data?.error ?? `request failed (${res.status})`);
  }
  return data;
}
