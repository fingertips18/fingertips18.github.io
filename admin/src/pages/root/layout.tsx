import { Outlet } from 'react-router-dom';

import { Header } from '@/components/header';

export function RootLayout() {
  return (
    <>
      <Header />
      <main className='h-default max-w-7xl mx-auto'>
        <Outlet />
      </main>
    </>
  );
}
