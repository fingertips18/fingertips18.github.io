import { useEffect } from 'react';
import { Outlet } from 'react-router-dom';

import { Footer } from '@/components/footer';
import { Header } from '@/components/header';
import { AnalyticsService } from '@/lib/services/analytics';
import ToasterProvider from '@/providers/toaster-provider';
import { AppRoutes } from '@/routes/app-routes';

const RootLayout = () => {
  useEffect(() => {
    // Intentionally ignore the returned promise.
    void AnalyticsService.pageView({
      location: AppRoutes.root,
      title: 'Root View',
    });
  }, []);

  return (
    <ToasterProvider>
      <Header />
      <main className='h-full max-w-screen-lg mx-auto max-xl:overflow-x-hidden'>
        <Outlet />
      </main>
      <Footer />
    </ToasterProvider>
  );
};

export default RootLayout;
