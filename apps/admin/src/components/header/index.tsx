import { Hint } from '@/components/common/hint';
import { Separator } from '@/components/shadcn/separator';
import { SidebarTrigger } from '@/components/shadcn/sidebar';

import { Breadcrumbs } from './breadcrumbs';

export function Header() {
  return (
    <header className='h-14 w-full max-w-7xl mx-auto flex-between px-4 py-3 border-b'>
      <div className='flex-start gap-x-2 h-full'>
        <Hint label='Menu' asChild>
          <SidebarTrigger className='cursor-pointer' />
        </Hint>
        <Separator orientation='vertical' />
        <Breadcrumbs />
      </div>
    </header>
  );
}
