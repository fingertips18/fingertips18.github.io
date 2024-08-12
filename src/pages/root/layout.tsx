import { Outlet } from "react-router-dom";

import { ModeToggle } from "@/common/components/mode-toggle";
import { Header } from "@/common/components/header";

const RootLayout = () => {
  return (
    <>
      <Header />
      <main className="h-full max-w-screen-lg mx-auto max-xl:overflow-x-hidden">
        <Outlet />
      </main>
      <div className="fixed bottom-6 right-6 lg:hidden">
        <ModeToggle />
      </div>
    </>
  );
};

export default RootLayout;
