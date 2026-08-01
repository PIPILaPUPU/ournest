import { useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import axios from "axios";
import { Menu, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { createWish, deleteWish, getWishlist, login, updateWish } from "@/lib/api";
import type { LoginResponse } from "@/types/auth";
import type { GroupColor, Wish } from "@/types/wish";

const SESSION_KEY = "wishlist-session";
type GroupName = "Техника" | "Еда" | "Учеба" | "Работа" | "Общее";

const FIXED_GROUPS: Array<{ name: GroupName; color: GroupColor }> = [
  { name: "Техника", color: "blue" },
  { name: "Еда", color: "green" },
  { name: "Учеба", color: "violet" },
  { name: "Работа", color: "rose" },
  { name: "Общее", color: "slate" },
];

const GROUP_STYLE: Record<GroupColor, { border: string; background: string }> = {
  slate: { border: "border-slate-300", background: "bg-slate-50" },
  blue: { border: "border-blue-300", background: "bg-blue-50" },
  green: { border: "border-emerald-300", background: "bg-emerald-50" },
  violet: { border: "border-violet-300", background: "bg-violet-50" },
  rose: { border: "border-rose-300", background: "bg-rose-50" },
};

type ModuleKey = "wishlist";

export function AppPage() {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [activeModule, setActiveModule] = useState<ModuleKey>("wishlist");
  const [session, setSession] = useState<LoginResponse | null>(() => {
    const raw = localStorage.getItem(SESSION_KEY);
    if (!raw) return null;
    try {
      return JSON.parse(raw) as LoginResponse;
    } catch {
      return null;
    }
  });

  const logout = () => {
    localStorage.removeItem(SESSION_KEY);
    setSession(null);
  };

  if (!session) {
    return <LoginScreen onSuccess={setSession} />;
  }

  return (
    <div className="min-h-screen bg-background text-foreground">
      <header className="sticky top-0 z-20 border-b border-border bg-white/90 backdrop-blur">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-4 py-3">
          <div className="flex items-center gap-3">
            <Button
              variant="outline"
              size="default"
              onClick={() => setIsSidebarOpen((prev) => !prev)}
              aria-label="Toggle menu"
            >
              {isSidebarOpen ? <X className="h-4 w-4" /> : <Menu className="h-4 w-4" />}
            </Button>
            <h1 className="text-lg font-semibold">Wish App</h1>
          </div>
          <div className="flex items-center gap-3 text-sm text-muted">
            <span>{session.username}</span>
            <Button variant="outline" onClick={logout}>
              Выйти
            </Button>
          </div>
        </div>
      </header>

      <div className="mx-auto flex max-w-7xl">
        <aside
          className={[
            "fixed left-0 top-[61px] z-10 h-[calc(100vh-61px)] w-72 border-r border-border bg-white p-4 transition-transform",
            isSidebarOpen ? "translate-x-0" : "-translate-x-full",
          ].join(" ")}
        >
          <p className="mb-2 px-2 text-xs uppercase tracking-wide text-muted">Микросервисы</p>
          <button
            type="button"
            onClick={() => {
              setActiveModule("wishlist");
              setIsSidebarOpen(false);
            }}
            className={[
              "w-full rounded-lg px-3 py-2 text-left text-sm",
              activeModule === "wishlist" ? "bg-gray-900 text-white" : "hover:bg-gray-100",
            ].join(" ")}
          >
            Список желаний
          </button>
        </aside>

        <main className="w-full px-4 py-6 md:px-6">
          {activeModule === "wishlist" && <WishlistModule currentUserId={session.user_id} />}
        </main>
      </div>
    </div>
  );
}

function LoginScreen({ onSuccess }: { onSuccess: (session: LoginResponse) => void }) {
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

function WishlistModule({ currentUserId }: { currentUserId: number }) {
  const queryClient = useQueryClient();
  const [activeGroup, setActiveGroup] = useState<GroupName | null>(null);

  const { data: wishes, isLoading, isError } = useQuery({
    queryKey: ["wishlist"],
    queryFn: getWishlist,
  });

  const createMutation = useMutation({
    mutationFn: createWish,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["wishlist"] }),
  });

  const deleteMutation = useMutation({
    mutationFn: deleteWish,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["wishlist"] }),
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: { description?: string; url?: string; price?: number } }) =>
      updateWish(id, payload),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["wishlist"] }),
  });

  const wishesByGroup = useMemo(() => {
    const byGroup: Record<GroupName, Wish[]> = {
      Техника: [],
      Еда: [],
      Учеба: [],
      Работа: [],
      Общее: [],
    };
    for (const wish of wishes ?? []) {
      if (wish.group_name in byGroup) {
        byGroup[wish.group_name as GroupName].push(wish);
      } else {
        byGroup.Общее.push(wish);
      }
    }
    for (const group of Object.keys(byGroup) as GroupName[]) {
      byGroup[group].sort((a, b) => b.id - a.id);
    }
    return byGroup;
  }, [wishes]);

  return (
    <div className="mx-auto max-w-6xl">
      <Card>
        <CardContent>
          <h2 className="text-xl font-semibold">Список желаний по группам</h2>
          <p className="mt-1 text-sm text-muted">
            Нажми на группу, чтобы открыть детали и действия по желаниям.
          </p>

          <div className="mt-5 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {FIXED_GROUPS.map((group) => (
              <GroupTile
                key={group.name}
                groupName={group.name}
                groupColor={group.color}
                wishes={wishesByGroup[group.name]}
                onClick={() => setActiveGroup(group.name)}
              />
            ))}
          </div>

          {isLoading && <p className="mt-4 text-sm text-muted">Загружаю...</p>}
          {isError && <p className="mt-4 text-sm text-red-600">Ошибка загрузки списка</p>}
        </CardContent>
      </Card>

      {activeGroup && (
        <GroupModal
          groupName={activeGroup}
          groupColor={FIXED_GROUPS.find((g) => g.name === activeGroup)?.color ?? "slate"}
          wishes={wishesByGroup[activeGroup]}
          currentUserId={currentUserId}
          onClose={() => setActiveGroup(null)}
          onCreate={(payload) => createMutation.mutate(payload)}
          onDelete={(id) => deleteMutation.mutate(id)}
          onUpdate={(id, payload) => updateMutation.mutate({ id, payload })}
          isBusy={createMutation.isPending || deleteMutation.isPending || updateMutation.isPending}
        />
      )}
    </div>
  );
}

function GroupTile({
  groupName,
  groupColor,
  wishes,
  onClick,
}: {
  groupName: GroupName;
  groupColor: GroupColor;
  wishes: Wish[];
  onClick: () => void;
}) {
  const previewTitles = wishes.slice(0, 4).map((wish) => wish.title);
  const opacities = ["opacity-80", "opacity-60", "opacity-45", "opacity-30"];

  return (
    <button
      type="button"
      onClick={onClick}
      className={`aspect-square w-full rounded-xl border p-4 text-left transition hover:shadow-soft ${GROUP_STYLE[groupColor].border} ${GROUP_STYLE[groupColor].background}`}
    >
      <div className="flex h-full flex-col">
        <div className="mb-2 flex items-center justify-between">
          <h3 className="font-semibold">{groupName}</h3>
          <span className="text-xs text-muted">{wishes.length}</span>
        </div>
        <div className="mt-2 space-y-2 text-sm text-gray-400">
          {previewTitles.length === 0 && <p className="text-gray-400">Пока пусто</p>}
          {previewTitles.map((title, idx) => (
            <p key={`${title}-${idx}`} className={`truncate ${opacities[idx] ?? "opacity-20"}`}>
              {title}
            </p>
          ))}
        </div>
      </div>
    </button>
  );
}

function GroupModal({
  groupName,
  groupColor,
  wishes,
  currentUserId,
  onClose,
  onCreate,
  onDelete,
  onUpdate,
  isBusy,
}: {
  groupName: GroupName;
  groupColor: GroupColor;
  wishes: Wish[];
  currentUserId: number;
  onClose: () => void;
  onCreate: (payload: {
    owner_id: number;
    title: string;
    description?: string;
    url?: string;
    price?: number;
    status?: "wanted" | "reserved" | "done";
    group_name?: string;
  }) => void;
  onDelete: (id: number) => void;
  onUpdate: (id: number, payload: { description?: string; url?: string; price?: number }) => void;
  isBusy: boolean;
}) {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [url, setURL] = useState("");
  const [price, setPrice] = useState("");
  const [errorText, setErrorText] = useState("");

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/45 p-4">
      <div className="max-h-[90vh] w-full max-w-3xl overflow-auto rounded-xl bg-white shadow-2xl">
        <div
          className={`sticky top-0 z-10 flex items-center justify-between border-b p-4 ${GROUP_STYLE[groupColor].border} ${GROUP_STYLE[groupColor].background}`}
        >
          <h3 className="text-lg font-semibold">{groupName}</h3>
          <Button variant="outline" onClick={onClose}>
            Закрыть
          </Button>
        </div>

        <div className="space-y-5 p-4">
          <Card>
            <CardContent>
              <h4 className="text-base font-semibold">Добавить желание в группу</h4>
              <form
                className="mt-3 grid gap-3"
                onSubmit={(e) => {
                  e.preventDefault();
                  if (!title.trim()) {
                    setErrorText("Введите название");
                    return;
                  }
                  onCreate({
                    owner_id: currentUserId,
                    title: title.trim(),
                    description: description.trim() || undefined,
                    url: url.trim() || undefined,
                    price: price ? Number(price) : undefined,
                    status: "wanted",
                    group_name: groupName,
                  });
                  setTitle("");
                  setDescription("");
                  setURL("");
                  setPrice("");
                  setErrorText("");
                }}
              >
                <input
                  className="w-full rounded-lg border border-border bg-white px-3 py-2 outline-none focus:border-gray-500"
                  placeholder="Название"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  required
                />
                <input
                  className="w-full rounded-lg border border-border bg-white px-3 py-2 outline-none focus:border-gray-500"
                  placeholder="Описание"
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                />
                <input
                  className="w-full rounded-lg border border-border bg-white px-3 py-2 outline-none focus:border-gray-500"
                  placeholder="Ссылка"
                  value={url}
                  onChange={(e) => setURL(e.target.value)}
                />
                <input
                  className="w-full rounded-lg border border-border bg-white px-3 py-2 outline-none focus:border-gray-500"
                  placeholder="Цена"
                  type="number"
                  step="0.01"
                  min="0"
                  value={price}
                  onChange={(e) => setPrice(e.target.value)}
                />
                {errorText && <p className="text-sm text-red-600">{errorText}</p>}
                <Button type="submit" disabled={isBusy}>
                  {isBusy ? "Сохраняю..." : "Добавить"}
                </Button>
              </form>
            </CardContent>
          </Card>

          <div className="space-y-3">
            {wishes.length === 0 && (
              <p className="text-sm text-muted">В этой группе пока нет желаний.</p>
            )}
            {wishes.map((wish) => (
              <WishRow
                key={wish.id}
                wish={wish}
                onDelete={() => onDelete(wish.id)}
                onUpdate={(payload) => onUpdate(wish.id, payload)}
                isBusy={isBusy}
              />
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}

function WishRow({
  wish,
  onDelete,
  onUpdate,
  isBusy,
}: {
  wish: Wish;
  onDelete: () => void;
  onUpdate: (payload: { description?: string; url?: string; price?: number }) => void;
  isBusy: boolean;
}) {
  const [isEditing, setIsEditing] = useState(false);
  const [description, setDescription] = useState(wish.description ?? "");
  const [url, setURL] = useState(wish.url ?? "");
  const [price, setPrice] = useState(wish.price != null ? String(wish.price) : "");

  return (
    <div className="rounded-lg border border-border bg-white p-4">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <p className="font-medium">{wish.title}</p>
          <p className="mt-1 text-sm text-muted">{wish.description ?? "Без описания"}</p>
          <p className="mt-1 text-xs text-muted">
            {wish.url ?? "Ссылка не указана"} {wish.price ? `• ${wish.price}$` : ""}
          </p>
          <p className="mt-1 text-xs text-muted">Добавил: {wish.owner_username}</p>
        </div>
        <div className="flex items-center gap-2">
          <span className="rounded-full border border-border px-3 py-1 text-xs">{wish.status}</span>
          <Button variant="outline" onClick={() => setIsEditing((prev) => !prev)} disabled={isBusy}>
            {isEditing ? "Скрыть" : "Изменить"}
          </Button>
          <Button variant="outline" onClick={onDelete} disabled={isBusy}>
            Удалить
          </Button>
        </div>
      </div>
      {isEditing && (
        <form
          className="mt-4 grid gap-2 border-t border-border pt-3"
          onSubmit={(e) => {
            e.preventDefault();
            onUpdate({
              description: description.trim() || undefined,
              url: url.trim() || undefined,
              price: price.trim() ? Number(price) : undefined,
            });
            setIsEditing(false);
          }}
        >
          <input
            className="w-full rounded-lg border border-border bg-white px-3 py-2 text-sm outline-none focus:border-gray-500"
            placeholder="Описание"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
          />
          <input
            className="w-full rounded-lg border border-border bg-white px-3 py-2 text-sm outline-none focus:border-gray-500"
            placeholder="Ссылка"
            value={url}
            onChange={(e) => setURL(e.target.value)}
          />
          <input
            className="w-full rounded-lg border border-border bg-white px-3 py-2 text-sm outline-none focus:border-gray-500"
            placeholder="Цена"
            type="number"
            step="0.01"
            min="0"
            value={price}
            onChange={(e) => setPrice(e.target.value)}
          />
          <Button type="submit" disabled={isBusy}>
            Сохранить изменения
          </Button>
        </form>
      )}
    </div>
  );
}
