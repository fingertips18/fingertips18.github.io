import { Outlet } from "react-router-dom";

import ToasterProvider from "@/lib/providers/toaster-provider";
import { Header } from "@/common/components/header";
import { Footer } from "@/common/components/footer";

const RootLayout = () => {
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
