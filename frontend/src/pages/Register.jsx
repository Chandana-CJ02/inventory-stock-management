import { useState } from "react";
import { useAuth } from "../auth.jsx";

const EMAIL_RE = /^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$/;

// Registration form. Client-side validation (AC-03) mirrors the backend rules:
// valid email, password between 8 and 72 characters (bcrypt's input limit).
// Registration does not auto-login; the user proceeds to the login form.
export default function Register({ onSwitch }) {
  const { register } = useAuth();
  const [form, setForm] = useState({ email: "", password: "" });
  const [fieldErrors, setFieldErrors] = useState({});
  const [apiError, setApiError] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [registeredEmail, setRegisteredEmail] = useState("");

  function validate() {
    const errs = {};
    if (!EMAIL_RE.test(form.email.trim())) {
      errs.email = "Enter a valid email address";
    }
    if (form.password.length < 8) {
      errs.password = "Password must be at least 8 characters";
    } else if (form.password.length > 72) {
      errs.password = "Password must be at most 72 characters";
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
      await register(form.email.trim().toLowerCase(), form.password);
      setRegisteredEmail(form.email.trim().toLowerCase());
    } catch (err) {
      setApiError(err.message);
    } finally {
      setSubmitting(false);
    }
  }

  if (registeredEmail) {
    return (
      <div>
        <h1>Account created</h1>
        <p>
          <strong>{registeredEmail}</strong> is registered. Log in to continue.
        </p>
        <button type="button" onClick={onSwitch}>
          Go to login
        </button>
      </div>
    );
  }

  return (
    <form onSubmit={onSubmit} noValidate>
      <h1>Register</h1>

      <label htmlFor="register-email">Email</label>
      <input
        id="register-email"
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

      <label htmlFor="register-password">Password</label>
      <input
        id="register-password"
        type="password"
        autoComplete="new-password"
        value={form.password}
        onChange={(e) => setForm({ ...form, password: e.target.value })}
      />
      <p className="hint">Between 8 and 72 characters.</p>
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
        {submitting ? "Creating account…" : "Register"}
      </button>
      <p className="muted">
        Already have an account?{" "}
        <button type="button" className="link" onClick={onSwitch}>
          Log in
        </button>
      </p>
    </form>
  );
}
