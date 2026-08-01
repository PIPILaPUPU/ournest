import { Gift, Link2, ShieldCheck } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import type { ReactNode } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { getWishlist } from "@/lib/api";

export function LandingPage() {
  const { data, isLoading, isError } = useQuery({
    queryKey: ["wishlist"],
    queryFn: getWishlist,
  });

  return (
    <main className="min-h-screen bg-background text-foreground">
      <section className="mx-auto max-w-6xl px-6 pb-16 pt-14">
        <div className="grid gap-10 md:grid-cols-2 md:items-center">
          <div>
            <p className="mb-3 text-sm font-medium text-muted">Wish App</p>
            <h1 className="text-4xl font-semibold tracking-tight md:text-5xl">
              Список желаний для пары — спокойно и без лишнего
            </h1>
            <p className="mt-5 max-w-xl text-base leading-relaxed text-muted">
              Добавляйте идеи подарков, делитесь ссылкой и отмечайте, что уже
              выбрано. Минималистично, без сложного интерфейса и лишних шагов.
            </p>
            <div className="mt-8 flex gap-3">
              <Button size="lg">Создать желание</Button>
              <Button variant="outline" size="lg">
                Посмотреть список
              </Button>
            </div>
          </div>
          <Card>
            <CardContent className="space-y-4">
              <h2 className="text-lg font-medium">Как это работает</h2>
              <Feature
                icon={<Gift className="h-5 w-5" />}
                title="Добавь желание"
                text="Название, ссылка на товар и цена — достаточно для старта."
              />
              <Feature
                icon={<Link2 className="h-5 w-5" />}
                title="Поделись списком"
                text="Партнер видит актуальные идеи и выбирает подарок без переписок."
              />
              <Feature
                icon={<ShieldCheck className="h-5 w-5" />}
                title="Отметь статус"
                text="wanted, reserved, done — понятно, что еще актуально."
              />
            </CardContent>
          </Card>
        </div>
      </section>

      <section className="mx-auto max-w-6xl px-6 pb-20">
        <div className="mb-5 flex items-end justify-between">
          <h3 className="text-2xl font-semibold">Пример текущих желаний</h3>
          <span className="text-sm text-muted">данные с твоего backend API</span>
        </div>
        <div className="grid gap-4 md:grid-cols-2">
          {isLoading && <WishPlaceholder label="Загружаю wishlist..." />}
          {isError && (
            <WishPlaceholder label="Не удалось получить данные из backend" />
          )}
          {data?.map((wish) => (
            <Card key={wish.id}>
              <CardContent>
                <div className="flex items-start justify-between gap-3">
                  <h4 className="text-lg font-medium">{wish.title}</h4>
                  <span className="rounded-full border border-border px-3 py-1 text-xs">
                    {wish.status}
                  </span>
                </div>
                <p className="mt-2 text-sm text-muted">
                  {wish.description ?? "Описание не заполнено"}
                </p>
                <div className="mt-4 text-sm">
                  Цена: {wish.price ? `${wish.price} $` : "не указана"}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </section>
    </main>
  );
}

function Feature({
  icon,
  title,
  text,
}: {
  icon: ReactNode;
  title: string;
  text: string;
}) {
  return (
    <div className="flex items-start gap-3 rounded-lg border border-border p-4">
      <div className="mt-0.5">{icon}</div>
      <div>
        <p className="font-medium">{title}</p>
        <p className="mt-1 text-sm text-muted">{text}</p>
      </div>
    </div>
  );
}

function WishPlaceholder({ label }: { label: string }) {
  return (
    <Card className="md:col-span-2">
      <CardContent>
        <p className="text-sm text-muted">{label}</p>
      </CardContent>
    </Card>
  );
}
