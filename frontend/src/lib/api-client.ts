import axios from "axios";

const baseURL = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

export const apiClient = axios.create({
  baseURL,
  timeout: 10000,
});

apiClient.interceptors.request.use((config) => {
  const raw = localStorage.getItem("wishlist-session");
  if (!raw) return config;

  try {
    const session = JSON.parse(raw) as { token?: string };
    if (session.token) {
      config.headers.Authorization = `Bearer ${session.token}`;
    }
  } catch {
    localStorage.removeItem("wishlist-session");
  }

  return config;
});

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401 && localStorage.getItem("wishlist-session")) {
      localStorage.removeItem("wishlist-session");
      window.location.assign("/");
    }
    return Promise.reject(error);
  },
);
