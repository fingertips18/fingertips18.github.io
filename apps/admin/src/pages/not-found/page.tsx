import { Link } from 'react-router-dom';

import { Button } from '@/components/shadcn/button';
import { Route } from '@/routes/route';

export default function NotFoundPage() {
  return (
    <section className='content flex-center flex-col gap-y-1 leading-tight'>
      <h6 className='lg:text-lg font-bold leading-none'>404 Page Not Found</h6>
      <p className='text-xs lg:text-sm text-muted-foreground'>
        Woops! Looks like this page doesn't exist.
      </p>
      <Button asChild className='rounded-full mt-2 cursor-pointer'>
        <Link to={Route.root}>Go back</Link>
      </Button>
    </section>
  );
}
