import { useLenis } from 'lenis/react';
import { LucideMenu, MoveLeft } from 'lucide-react';
import { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';

import { Hint } from '@/components/common/hint';
import { Button } from '@/components/shadcn/button';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/shadcn/sheet';
import { ROOTMENU } from '@/constants/collections';
import { ROOTSECTION } from '@/constants/enums';
import { useResize } from '@/lib/hooks/useResize';
import { cn } from '@/lib/utils';
import { AppRoutes } from '@/routes/app-routes';

import { ModeToggle } from './mode-toggle';

interface SheetMenuProps {
  active: ROOTSECTION;
}

const SheetMenu = ({ active }: SheetMenuProps) => {
  const lenis = useLenis();
  const location = useLocation();
  // Track user's intent to open the sheet
  const [userOpen, setUserOpen] = useState(false);
  const { width } = useResize();

  // Derive the actual open state: only open if user wants it AND on mobile
  const open = width <= 1024 && userOpen;

  const onOpenChange = (open: boolean) => {
    if (!lenis) return;

    if (open) {
      lenis.stop();
    } else {
      lenis.start();
    }
  };

  const onClick = (id: string) => {
    const section = document.getElementById(id);
    if (section) {
      section.scrollIntoView({
        behavior: 'smooth',
        block: 'start',
        inline: 'nearest',
      });
    }
  };

  return (
    <Sheet
      open={open}
      onOpenChange={(open) => {
        onOpenChange(open);
        setUserOpen(open);
      }}
    >
      <Hint asChild label='Menu'>
        <SheetTrigger asChild>
          <Button
            variant={'ghost'}
            size={'icon'}
            aria-label='menu-toggle'
            className='lg:hidden hover:drop-shadow-primary-glow'
          >
            <LucideMenu className='w-5 h-5' />
          </Button>
        </SheetTrigger>
      </Hint>
      <SheetContent data-lenis-prevent className='overflow-y-auto no-scrollbar'>
        <SheetHeader className='mt-4 items-start!'>
          <SheetTitle className='text-sm'>Menu</SheetTitle>
          <SheetDescription className='text-xs text-start'>
            Discover my portfolio, skills, projects, and how to connect.
          </SheetDescription>
        </SheetHeader>

        <nav className='w-full flex justify-end mt-10 flex-1'>
          {location.pathname === '/' ? (
            <ul className='space-y-6 text-end'>
              {ROOTMENU.map((m, i) => (
                <li
                  key={`${m.label}-${i}`}
                  className={cn(
                    'capitalize font-semibold leading-none hover:scale-95 transition-all cursor-pointer hover:drop-shadow-primary-glow lg:hover:text-accent',
                    active === m.label && 'text-accent',
                  )}
                  onClick={() => onClick(m.label)}
                >
                  {m.label}
                </li>
              ))}
            </ul>
          ) : (
            <Link
              to={AppRoutes.root}
              className='flex items-center gap-x-2 hover:text-accent'
            >
              <MoveLeft className='size-4' /> Go home
            </Link>
          )}
        </nav>

        <SheetFooter className='fixed bottom-4 right-4'>
          <ModeToggle />
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
};

export { SheetMenu };
