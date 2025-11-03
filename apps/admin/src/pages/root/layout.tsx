import { Outlet } from 'react-router-dom';

import { Header } from '@/components/header';
import { SidebarProvider, SidebarTrigger } from '@/components/shadcn/sidebar';
import { Sidebar } from '@/components/sidebar';

export function RootLayout() {
  return (
    <SidebarProvider className='flex-col'>
      <Header />
      <Sidebar />
      <main className='h-default max-w-7xl mx-auto'>
        <SidebarTrigger />
        <Outlet />
      </main>
    </SidebarProvider>
  );
}
