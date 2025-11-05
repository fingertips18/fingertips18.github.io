import { Outlet } from 'react-router-dom';

import { Header } from '@/components/header';
import { SidebarProvider } from '@/components/shadcn/sidebar';
import { Sidebar } from '@/components/sidebar';

export function RootLayout() {
  return (
    <SidebarProvider>
      <Sidebar />
      <div className='w-full'>
        <Header />
        <main className='h-default'>
          <Outlet />
        </main>
      </div>
    </SidebarProvider>
  );
}
