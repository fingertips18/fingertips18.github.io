import { RefreshCcw } from 'lucide-react';
import { useEffect } from 'react';

import { Button } from '@/components/shadcn/button';
import { AnalyticsService } from '@/lib/services/analytics';
import { AppRoutes } from '@/routes/app-routes';

const ErrorPage = () => {
  useEffect(() => {
    // Intentionally ignore the returned promise.
    void AnalyticsService.pageView({
      location: AppRoutes.github404,
      title: 'Error View',
    });
  }, []);

  return (
    <section className='h-[calc(100dvh_-_56px)] flex-center flex-col gap-y-2 lg:gap-y-1.5 leading-tight'>
      <h6 className='lg:text-lg font-bold'>Something Went Wrong</h6>
      <p className='text-xs lg:text-sm text-muted-foreground'>
        Woops! You are not supposed to see this.
      </p>
      <Button
        onClick={() => window.location.reload()}
        className='rounded-full gap-x-2'
      >
        <RefreshCcw className='w-4 h-4' /> Refresh
      </Button>
    </section>
  );
};

export default ErrorPage;
