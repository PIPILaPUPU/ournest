import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { CalendarHeart, Check, RotateCcw, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  createDateIdea,
  deleteDateIdea,
  getDateIdeas,
  updateDateIdea,
} from "@/features/date-ideas/api";
import type { DateIdea } from "@/features/date-ideas/types";

const queryKey = ["date-ideas"];

export function DateIdeasModule({ currentUserId }: { currentUserId: number }) {
  const queryClient = useQueryClient();
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [errorText, setErrorText] = useState("");

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
      setErrorText("");
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
    onSuccess: refresh,
  });

  const isBusy =
    createMutation.isPending || updateMutation.isPending || deleteMutation.isPending;

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
                author_id: currentUserId,
                title: title.trim(),
                description: description.trim() || undefined,
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
            {errorText && <p className="text-sm text-red-600">{errorText}</p>}
            <Button type="submit" disabled={isBusy}>
              Добавить идею
            </Button>
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
    </div>
  );
}
