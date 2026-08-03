self.addEventListener("push", (event) => {
  let payload = {
    title: "Our Nest",
    body: "Появилось новое событие",
    url: "/",
    type: "notification",
  };

  if (event.data) {
    try {
      payload = { ...payload, ...event.data.json() };
    } catch {
      payload.body = event.data.text();
    }
  }

  event.waitUntil(
    self.registration.showNotification(payload.title, {
      body: payload.body,
      icon: "/icons/app-icon-192.png",
      badge: "/icons/app-icon-192.png",
      data: {
        url: payload.url,
      },
    }),
  );
});

self.addEventListener("notificationclick", (event) => {
  event.notification.close();
  const targetURL = new URL(event.notification.data?.url || "/", self.location.origin);

  event.waitUntil(
    self.clients
      .matchAll({ type: "window", includeUncontrolled: true })
      .then(async (clients) => {
        for (const client of clients) {
          if (new URL(client.url).origin === targetURL.origin) {
            await client.navigate(targetURL.href);
            return client.focus();
          }
        }
        return self.clients.openWindow(targetURL.href);
      }),
  );
});
