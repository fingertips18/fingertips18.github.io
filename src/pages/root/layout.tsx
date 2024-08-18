import { Outlet } from "react-router-dom";
import { useEffect } from "react";
import ReactGA from "react-ga4";

import ToasterProvider from "@/lib/providers/toaster-provider";
import { Header } from "@/common/components/header";
import { Footer } from "@/common/components/footer";

const RootLayout = () => {
  useEffect(() => {
    ReactGA.initialize(import.meta.env.VITE_GOOGLE_MEASUREMENT_ID);
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
