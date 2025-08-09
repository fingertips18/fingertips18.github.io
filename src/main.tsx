import "./index.css";

import React from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider } from "react-router-dom";

import LenisProvider from "./providers/lenis-provider";
import { ThemeProvider } from "./providers/theme-provider";
import { router } from "./routes/router";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ThemeProvider>
      <LenisProvider>
        <RouterProvider router={router} />
      </LenisProvider>
    </ThemeProvider>
  </React.StrictMode>
);
