import { RouterProvider } from "react-router-dom";
import ReactDOM from "react-dom/client";
import React from "react";

import { ThemeProvider } from "./lib/providers/theme-provider";
import LenisProvider from "./lib/providers/lenis-provider";
import { router } from "./routes/router";
import "./index.css";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ThemeProvider>
      <LenisProvider>
        <RouterProvider router={router} />
      </LenisProvider>
    </ThemeProvider>
  </React.StrictMode>
);
