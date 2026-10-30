import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// The dev server proxies /api to the Go backend so the browser only talks to
// its own origin — no CORS handling needed in development. In production both
// would sit behind the same reverse proxy.
export default defineConfig({
  plugins: [react()],
  // Use the automatic JSX runtime in every transform path (including test
  // files), so components don't need a React import.
  esbuild: { jsx: "automatic" },
  server: {
    port: 5173,
    proxy: {
      "/api": {
        target: "http://localhost:8081",
        changeOrigin: true,
      },
    },
  },
  test: {
    environment: "jsdom",
    globals: true,
    setupFiles: ["./src/test-setup.js"],
  },
});
