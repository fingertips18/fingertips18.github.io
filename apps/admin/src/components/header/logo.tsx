import { Link, useLocation } from 'react-router-dom';

import { Hint } from '@/components/common/hint';
import { Button } from '@/components/shadcn/button';
import { Separator } from '@/components/shadcn/separator';
import { SidebarTrigger } from '@/components/shadcn/sidebar';
import { TABLE_TITLE } from '@/constants/tables';
import { Route } from '@/routes/route';

export function Logo() {
  const { pathname } = useLocation();

  // Look up the title for the current pathname, fallback to the root title
  const title =
    pathname in TABLE_TITLE
      ? TABLE_TITLE[pathname as keyof typeof TABLE_TITLE]
      : TABLE_TITLE[Route.root];

  return (
    <div className='flex-start gap-x-2 h-full'>
      <Hint label='Menu' asChild>
        <SidebarTrigger className='cursor-pointer' />
      </Hint>
      <Separator orientation='vertical' />
      <Button asChild variant='link' className='p-0 ml-2'>
        <Link to={pathname}>{title}</Link>
      </Button>
    </div>
  );
}
