import { Outlet } from 'react-router-dom';

import { Header } from '@/components/header';
import { Sidebar } from '@/components/sidebar';
import Providers from '@/providers';

export default function RootLayout() {
  return (
    <Providers>
      <Sidebar />
      <div className='w-full'>
        <Header />
        <main className='h-default'>
          <Outlet />
        </main>
      </div>
    </Providers>
  );
}
