import {
  createBrowserRouter,
  createRoutesFromElements,
  Route,
} from "react-router-dom";

import NotFoundPage from "@/pages/not-found/page";
import RootLayout from "@/pages/root/layout";
import RootPage from "@/pages/root/page";

import { AppRoutes } from "./app-routes";

export const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path={AppRoutes.root} element={<RootLayout />}>
      {/* Root */}
      <Route index element={<RootPage />} />

      {/* Not Found */}
      <Route path={AppRoutes.notFound} element={<NotFoundPage />} />
    </Route>
  )
);
