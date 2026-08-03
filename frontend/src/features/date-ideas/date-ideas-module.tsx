import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  CalendarHeart,
  Check,
  PackageOpen,
  RotateCcw,
  Sparkles,
  Trash2,
  X,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  SuccessToast,
  useSuccessToast,
} from "@/components/ui/success-toast";
import {
  createDateIdea,
  deleteDateIdea,
  getDateIdeas,
  getRandomSecretDateIdea,
  updateDateIdea,
} from "@/features/date-ideas/api";
import type { DateIdea } from "@/features/date-ideas/types";

const queryKey = ["date-ideas"];

export function DateIdeasModule() {
  const queryClient = useQueryClient();
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [isSecret, setIsSecret] = useState(false);
  const [errorText, setErrorText] = useState("");
  const [revealedIdea, setRevealedIdea] = useState<DateIdea | null>(null);
  const [isEmptyBoxOpen, setIsEmptyBoxOpen] = useState(false);
  const { toast, showSuccessToast, dismissSuccessToast } = useSuccessToast();

  const { data = [], isLoading, isError } = useQuery({
    queryKey,
    queryFn: getDateIdeas,
  });

  const refresh = () => queryClient.invalidateQueries({ queryKey });

  const createMutation = useMutation({
    mutationFn: createDateIdea,
    onSuccess: () => {
      setTitle("");
      setDescription("");
      setIsSecret(false);
      setErrorText("");
      showSuccessToast("Идея для свидания успешно добавлена");
      refresh();
    },
    onError: () => setErrorText("Не удалось добавить идею"),
  });

  const updateMutation = useMutation({
    mutationFn: ({ idea, status }: { idea: DateIdea; status: DateIdea["status"] }) =>
      updateDateIdea(idea.id, { status }),
    onSuccess: refresh,
  });

  const deleteMutation = useMutation({
    mutationFn: deleteDateIdea,
    onSuccess: () => {
      showSuccessToast("Идея для свидания успешно удалена");
      refresh();
    },
  });

  const revealMutation = useMutation({
    mutationFn: getRandomSecretDateIdea,
    onSuccess: (idea) => {
      setRevealedIdea(idea);
      setIsEmptyBoxOpen(idea === null);
    },
  });

  const isBusy =
    createMutation.isPending ||
    updateMutation.isPending ||
    deleteMutation.isPending ||
    revealMutation.isPending;

  return (
    <div className="mx-auto max-w-5xl space-y-6">
      <Card>
        <CardContent>
          <div className="flex items-center gap-3">
            <CalendarHeart className="h-5 w-5" />
            <div>
              <h2 className="text-xl font-semibold">Идеи для свиданий</h2>
              <p className="text-sm text-muted">Сохраняйте идеи и отмечайте реализованные.</p>
            </div>
          </div>

          <form
            className="mt-5 grid gap-3"
            onSubmit={(event) => {
              event.preventDefault();
              if (!title.trim()) {
                setErrorText("Введите название идеи");
                return;
              }
              createMutation.mutate({
                title: title.trim(),
                description: description.trim() || undefined,
                is_secret: isSecret,
              });
            }}
          >
            <input
              className="w-full rounded-lg border border-border bg-white px-3 py-2 outline-none focus:border-gray-500"
              placeholder="Например: вечерняя прогулка у воды"
              value={title}
              onChange={(event) => setTitle(event.target.value)}
            />
            <textarea
              className="min-h-24 w-full resize-y rounded-lg border border-border bg-white px-3 py-2 outline-none focus:border-gray-500"
              placeholder="Небольшое описание"
              value={description}
              onChange={(event) => setDescription(event.target.value)}
            />
            <label className="flex cursor-pointer items-center gap-3 rounded-lg border border-border px-3 py-2 text-sm">
              <input
                type="checkbox"
                checked={isSecret}
                onChange={(event) => setIsSecret(event.target.checked)}
                className="h-4 w-4 accent-gray-900"
              />
              <span>
                Скрыть в шкатулке
                <span className="ml-2 text-muted">Идея не появится в общем списке</span>
              </span>
            </label>
            {errorText && <p className="text-sm text-red-600">{errorText}</p>}
            <div className="flex flex-wrap gap-2">
              <Button type="submit" disabled={isBusy}>
                Добавить идею
              </Button>
              <Button
                type="button"
                variant="outline"
                disabled={isBusy}
                onClick={() => revealMutation.mutate()}
              >
                <PackageOpen className="mr-2 h-4 w-4" />
                {revealMutation.isPending ? "Открываю..." : "Достать из шкатулки"}
              </Button>
            </div>
            {revealMutation.isError && (
              <p className="text-sm text-red-600">Не удалось открыть шкатулку</p>
            )}
          </form>
        </CardContent>
      </Card>

      <div className="grid gap-4 sm:grid-cols-2">
        {isLoading && <p className="text-sm text-muted">Загружаю идеи...</p>}
        {isError && <p className="text-sm text-red-600">Не удалось загрузить идеи</p>}
        {!isLoading && data.length === 0 && (
          <p className="text-sm text-muted">Пока нет идей. Добавьте первую выше.</p>
        )}
        {data.map((idea) => (
          <Card key={idea.id} className={idea.status === "done" ? "opacity-70" : ""}>
            <CardContent className="flex h-full flex-col">
              <div className="flex items-start justify-between gap-3">
                <h3 className="font-semibold">{idea.title}</h3>
                <span className="rounded-full border border-border px-3 py-1 text-xs">
                  {idea.status === "done" ? "реализовано" : "запланировано"}
                </span>
              </div>
              <p className="mt-3 flex-1 text-sm leading-relaxed text-muted">
                {idea.description ?? "Без описания"}
              </p>
              <p className="mt-4 text-xs text-muted">Добавил: {idea.author_username}</p>
              <div className="mt-4 flex flex-wrap gap-2">
                <Button
                  variant="outline"
                  disabled={isBusy}
                  onClick={() =>
                    updateMutation.mutate({
                      idea,
                      status: idea.status === "done" ? "planned" : "done",
                    })
                  }
                >
                  {idea.status === "done" ? (
                    <RotateCcw className="mr-2 h-4 w-4" />
                  ) : (
                    <Check className="mr-2 h-4 w-4" />
                  )}
                  {idea.status === "done" ? "Вернуть" : "Готово"}
                </Button>
                <Button
                  variant="outline"
                  disabled={isBusy}
                  onClick={() => deleteMutation.mutate(idea.id)}
                >
                  <Trash2 className="mr-2 h-4 w-4" />
                  Удалить
                </Button>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {(revealedIdea || isEmptyBoxOpen) && (
        <SecretIdeaModal
          idea={revealedIdea}
          onClose={() => {
            setRevealedIdea(null);
            setIsEmptyBoxOpen(false);
          }}
          onTryAgain={() => revealMutation.mutate()}
          isLoading={revealMutation.isPending}
        />
      )}
      <SuccessToast toast={toast} onClose={dismissSuccessToast} />
    </div>
  );
}

function SecretIdeaModal({
  idea,
  onClose,
  onTryAgain,
  isLoading,
}: {
  idea: DateIdea | null;
  onClose: () => void;
  onTryAgain: () => void;
  isLoading: boolean;
}) {
  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/30 p-4 backdrop-blur-sm"
      role="dialog"
      aria-modal="true"
      aria-labelledby="secret-idea-title"
      onClick={onClose}
    >
      <div
        className="secret-idea-reveal w-full max-w-md rounded-2xl border border-border bg-white p-6 shadow-soft"
        onClick={(event) => event.stopPropagation()}
      >
        <div className="flex items-start justify-between gap-4">
          <div className="flex h-10 w-10 items-center justify-center rounded-full bg-gray-100">
            <Sparkles className="h-5 w-5" />
          </div>
          <button
            type="button"
            onClick={onClose}
            className="rounded-lg p-2 text-muted transition hover:bg-gray-100"
            aria-label="Закрыть"
          >
            <X className="h-4 w-4" />
          </button>
        </div>

        {idea ? (
          <>
            <p className="mt-5 text-xs uppercase tracking-widest text-muted">
              Сегодня из шкатулки
            </p>
            <h3 id="secret-idea-title" className="mt-2 text-2xl font-semibold">
              {idea.title}
            </h3>
            <p className="mt-3 leading-relaxed text-muted">
              {idea.description ?? "Пусть детали останутся маленьким сюрпризом."}
            </p>
            <p className="mt-5 text-xs text-muted">Добавил: {idea.author_username}</p>
          </>
        ) : (
          <>
            <h3 id="secret-idea-title" className="mt-5 text-xl font-semibold">
              Шкатулка пока пуста
            </h3>
            <p className="mt-2 text-sm text-muted">
              Добавьте хотя бы одну скрытую идею, и она сможет появиться здесь.
            </p>
          </>
        )}

        <div className="mt-6 flex justify-end gap-2">
          {idea && (
            <Button variant="outline" onClick={onTryAgain} disabled={isLoading}>
              Ещё идея
            </Button>
          )}
          <Button onClick={onClose}>Закрыть</Button>
        </div>
      </div>
    </div>
  );
}
