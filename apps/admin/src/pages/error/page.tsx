import { RefreshCcw } from 'lucide-react';

import { Button } from '@/components/shadcn/button';

export default function ErrorPage() {
  return (
    <main className='h-dvh'>
      <section className='size-full flex-center flex-col gap-y-1'>
        <h1 className='lg:text-lg font-bold leading-none'>
          Something Went Wrong
        </h1>
        <p className='text-xs lg:text-sm text-muted-foreground'>
          Woops! You are not supposed to see this.
        </p>
        <Button
          onClick={() => window.location.reload()}
          className='rounded-full gap-x-2 mt-2 cursor-pointer'
        >
          <RefreshCcw aria-hidden='true' className='size-4' /> Refresh
        </Button>
      </section>
    </main>
  );
}
