import { vi } from "vitest";

// Small fetch mock helper shared by frontend tests: routes requests to
// canned responses by method + URL prefix. Later routes override earlier
// ones with the same method + prefix, so tests can append overrides.
export function fetchRouter(routes) {
  const spy = vi.fn(async (url, opts = {}) => {
    const method = String(opts.method ?? "GET").toUpperCase();
    const u = String(url);
    let match;
    for (const r of routes) {
      if (r.method.toUpperCase() === method && u.startsWith(r.prefix)) {
        match = r; // keep the last match
      }
    }
    if (!match) {
      return { ok: false, status: 404, json: async () => ({ error: "not found" }) };
    }
    const status = match.status ?? 200;
    const body = typeof match.body === "function" ? match.body({ url: u, opts }) : match.body;
    return { ok: status < 400, status, json: async () => body };
  });
  global.fetch = spy;
  return spy;
}

// All recorded calls to the spy for one method + URL prefix.
export function callsTo(spy, method, prefix) {
  return spy.mock.calls.filter(
    ([u, o]) =>
      u === prefix &&
      String(o?.method ?? "GET").toUpperCase() === method.toUpperCase()
  );
}
