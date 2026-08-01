import { createBrowserRouter } from "react-router-dom";
import { AppPage } from "@/pages/app-page";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <AppPage />,
  },
]);
