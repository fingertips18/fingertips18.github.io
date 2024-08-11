import { Outlet } from "react-router-dom";

import Header from "../../common/components/header";

const RootLayout = () => {
  return (
    <>
      <Header />
      <main className="h-[calc(100dvg_-_56px)] max-w-screen-lg mx-auto py-8">
        <Outlet />
      </main>
    </>
  );
};

export default RootLayout;
