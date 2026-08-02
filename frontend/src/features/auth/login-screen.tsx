import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import axios from "axios";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { login } from "@/features/auth/api";
import type { LoginResponse } from "@/types/auth";

const SESSION_KEY = "wishlist-session";

export function LoginScreen({ onSuccess }: { onSuccess: (session: LoginResponse) => void }) {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [errorText, setErrorText] = useState("");

  const loginMutation = useMutation({
    mutationFn: login,
    onSuccess: (session) => {
      localStorage.setItem(SESSION_KEY, JSON.stringify(session));
      onSuccess(session);
    },
    onError: (error) => {
      if (axios.isAxiosError(error) && error.response?.status === 401) {
        setErrorText("Неверный логин или пароль");
        return;
      }
      setErrorText("Ошибка входа. Проверь, что backend запущен.");
    },
  });

  return (
    <main className="grid min-h-screen place-items-center bg-background px-4">
      <Card className="w-full max-w-md">
        <CardContent>
          <h1 className="text-2xl font-semibold">Вход</h1>
          <p className="mt-2 text-sm text-muted">
            Используй тестовых пользователей из миграции: alice / secret или bob / secret
          </p>
          <form
            className="mt-6 space-y-4"
            onSubmit={(e) => {
              e.preventDefault();
              setErrorText("");
              loginMutation.mutate({ username, password });
            }}
          >
            <label className="block">
              <span className="mb-1 block text-sm">Username</span>
              <input
                className="w-full rounded-lg border border-border bg-white px-3 py-2 outline-none focus:border-gray-500"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="alice"
                required
              />
            </label>
            <label className="block">
              <span className="mb-1 block text-sm">Password</span>
              <input
                className="w-full rounded-lg border border-border bg-white px-3 py-2 outline-none focus:border-gray-500"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="secret"
                required
              />
            </label>
            {errorText && <p className="text-sm text-red-600">{errorText}</p>}
            <Button type="submit" className="w-full" disabled={loginMutation.isPending}>
              {loginMutation.isPending ? "Вхожу..." : "Войти"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </main>
  );
}
