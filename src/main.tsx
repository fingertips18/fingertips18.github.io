import { RouterProvider } from "react-router-dom";
import ReactDOM from "react-dom/client";
import React from "react";

import { ThemeProvider } from "./lib/providers/theme-provider";
import { router } from "./routes/router";
import "./index.css";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ThemeProvider>
      <RouterProvider router={router} />
    </ThemeProvider>
  </React.StrictMode>
);
