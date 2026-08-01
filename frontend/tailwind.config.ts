import type { Config } from "tailwindcss";

export default {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        background: "#f7f7f5",
        foreground: "#1f2937",
        muted: "#6b7280",
        card: "#ffffff",
        border: "#e5e7eb",
        primary: "#1f2937",
        "primary-foreground": "#ffffff",
      },
      boxShadow: {
        soft: "0 8px 24px rgba(15, 23, 42, 0.06)",
      },
    },
  },
  plugins: [],
} satisfies Config;
