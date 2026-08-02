import { useState } from "react";
import { CalendarHeart, Gift, Menu, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { LoginScreen } from "@/features/auth/login-screen";
import { DateIdeasModule } from "@/features/date-ideas/date-ideas-module";
import { WishlistModule } from "@/features/wishlist/wishlist-module";
import type { LoginResponse } from "@/types/auth";

const SESSION_KEY = "wishlist-session";
type ModuleKey = "wishlist" | "date-ideas";

const modules: Array<{
  key: ModuleKey;
  label: string;
  icon: typeof Gift;
}> = [
  { key: "wishlist", label: "Список желаний", icon: Gift },
  { key: "date-ideas", label: "Идеи для свиданий", icon: CalendarHeart },
];

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
              onClick={() => setIsSidebarOpen((open) => !open)}
              aria-label="Открыть меню"
            >
              {isSidebarOpen ? <X className="h-4 w-4" /> : <Menu className="h-4 w-4" />}
            </Button>
            <h1 className="text-lg font-semibold">Our Nest</h1>
          </div>
          <div className="flex items-center gap-3 text-sm text-muted">
            <span>{session.username}</span>
            <Button variant="outline" onClick={logout}>
              Выйти
            </Button>
          </div>
        </div>
      </header>

      <aside
        className={[
          "fixed left-0 top-[61px] z-10 h-[calc(100vh-61px)] w-72 border-r border-border bg-white p-4 transition-transform",
          isSidebarOpen ? "translate-x-0" : "-translate-x-full",
        ].join(" ")}
      >
        <p className="mb-2 px-2 text-xs uppercase tracking-wide text-muted">Разделы</p>
        <nav className="space-y-1">
          {modules.map((module) => {
            const Icon = module.icon;
            return (
              <button
                key={module.key}
                type="button"
                onClick={() => {
                  setActiveModule(module.key);
                  setIsSidebarOpen(false);
                }}
                className={[
                  "flex w-full items-center gap-3 rounded-lg px-3 py-2 text-left text-sm",
                  activeModule === module.key
                    ? "bg-gray-900 text-white"
                    : "hover:bg-gray-100",
                ].join(" ")}
              >
                <Icon className="h-4 w-4" />
                {module.label}
              </button>
            );
          })}
        </nav>
      </aside>

      <main className="mx-auto w-full max-w-7xl px-4 py-6 md:px-6">
        {activeModule === "wishlist" && (
          <WishlistModule currentUserId={session.user_id} />
        )}
        {activeModule === "date-ideas" && (
          <DateIdeasModule currentUserId={session.user_id} />
        )}
      </main>
    </div>
  );
}
