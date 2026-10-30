import { createContext, useContext, useMemo, useState } from "react";
import { apiRequest } from "./api.js";

const AuthContext = createContext(null);

const TOKEN_KEY = "inventory.token";
const EMAIL_KEY = "inventory.email";

// Session tokens are kept in localStorage (approved trade-off with a noted XSS
// risk, docs/architecture.md §9). The backend issues stateless JWTs without
// server-side revocation, so logout means discarding the token here; the 24 h
// expiry bounds the exposure of a leaked one.
function saveSession(token, email) {
  localStorage.setItem(TOKEN_KEY, token);
  localStorage.setItem(EMAIL_KEY, email);
}

export function clearSession() {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(EMAIL_KEY);
}

export function AuthProvider({ children }) {
  const [token, setToken] = useState(() => localStorage.getItem(TOKEN_KEY));
  const [email, setEmail] = useState(() => localStorage.getItem(EMAIL_KEY) ?? "");

  const value = useMemo(
    () => ({
      token,
      email,
      isAuthenticated: Boolean(token),

      async register(email, password) {
        return apiRequest("/v1/auth/register", {
          method: "POST",
          body: { email, password },
        });
      },

      async login(email, password) {
        const data = await apiRequest("/v1/auth/login", {
          method: "POST",
          body: { email, password },
        });
        const normalized = email.trim().toLowerCase();
        saveSession(data.token, normalized);
        setToken(data.token);
        setEmail(normalized);
        return data;
      },

      async logout() {
        try {
          if (token) {
            // A 401 here means the token was already invalid/expired — still
            // fall through and clear the local session.
            await apiRequest("/v1/auth/logout", { method: "POST", token });
          }
        } finally {
          clearSession();
          setToken(null);
          setEmail("");
        }
      },
    }),
    [token, email]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within an AuthProvider");
  return ctx;
}
