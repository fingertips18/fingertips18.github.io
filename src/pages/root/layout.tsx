import { Outlet } from "react-router-dom";

import { Header } from "@/common/components/header";

import { ModeToggle } from "./_components/mode-toggle";

const RootLayout = () => {
  return (
    <>
      <Header />
      <main className="h-[calc(100dvh_-_56px)] max-w-screen-lg mx-auto py-8">
        <Outlet />
      </main>
      <ModeToggle />
    </>
  );
};

export default RootLayout;
