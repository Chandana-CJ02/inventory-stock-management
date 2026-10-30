import { useState } from "react";
import { useAuth } from "../auth.jsx";

const EMAIL_RE = /^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$/;

// Login form. Client-side validation (AC-03) blocks submission before any
// API call; server errors are surfaced readably (AC-15).
export default function Login({ onSwitch }) {
  const { login } = useAuth();
  const [form, setForm] = useState({ email: "", password: "" });
  const [fieldErrors, setFieldErrors] = useState({});
  const [apiError, setApiError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  function validate() {
    const errs = {};
    if (!EMAIL_RE.test(form.email.trim())) {
      errs.email = "Enter a valid email address";
    }
    if (form.password.length === 0) {
      errs.password = "Password is required";
    } else if (form.password.length < 8) {
      errs.password = "Password must be at least 8 characters";
    }
    setFieldErrors(errs);
    return Object.keys(errs).length === 0;
  }

  async function onSubmit(e) {
    e.preventDefault();
    setApiError("");
    if (!validate()) return;
    setSubmitting(true);
    try {
      await login(form.email, form.password);
    } catch (err) {
      setApiError(err.message);
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <form onSubmit={onSubmit} noValidate>
      <h1>Log in</h1>

      <label htmlFor="login-email">Email</label>
      <input
        id="login-email"
        type="email"
        autoComplete="email"
        value={form.email}
        onChange={(e) => setForm({ ...form, email: e.target.value })}
      />
      {fieldErrors.email && (
        <p className="error" role="alert">
          {fieldErrors.email}
        </p>
      )}

      <label htmlFor="login-password">Password</label>
      <input
        id="login-password"
        type="password"
        autoComplete="current-password"
        value={form.password}
        onChange={(e) => setForm({ ...form, password: e.target.value })}
      />
      {fieldErrors.password && (
        <p className="error" role="alert">
          {fieldErrors.password}
        </p>
      )}

      {apiError && (
        <p className="error" role="alert">
          {apiError}
        </p>
      )}

      <button type="submit" disabled={submitting}>
        {submitting ? "Logging in…" : "Log in"}
      </button>
      <p className="muted">
        Need an account?{" "}
        <button type="button" className="link" onClick={onSwitch}>
          Register
        </button>
      </p>
    </form>
  );
}
