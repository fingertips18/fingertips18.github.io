import { Outlet } from "react-router-dom";
import { useEffect } from "react";
import ReactGA from "react-ga4";

import ToasterProvider from "@/lib/providers/toaster-provider";
import { AppRoutes } from "@/routes/app-routes";
import { Header } from "@/components/header";
import { Footer } from "@/components/footer";

const RootLayout = () => {
  useEffect(() => {
    ReactGA.initialize(import.meta.env.VITE_GOOGLE_MEASUREMENT_ID);

    ReactGA.send({
      hitType: "pageview",
      page: AppRoutes.root,
      title: "Root View",
    });
  }, []);

  return (
    <ToasterProvider>
      <Header />
      <main className="h-full max-w-screen-lg mx-auto max-xl:overflow-x-hidden">
        <Outlet />
      </main>
      <Footer />
    </ToasterProvider>
  );
};

export default RootLayout;
