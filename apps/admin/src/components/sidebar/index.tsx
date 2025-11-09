import { Link, useLocation } from 'react-router-dom';

import { Button } from '@/components/shadcn/button';
import {
  Sidebar as ShadcnSidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/shadcn/sidebar';
import { TABLES } from '@/constants/tables';
import { Route } from '@/routes/route';

export function Sidebar() {
  const { pathname } = useLocation();

  // Split and filter to remove empty parts
  const parts = pathname.split('/').filter(Boolean); // e.g. ['project', 'education', 'skill']

  return (
    <ShadcnSidebar className='overflow-x-hidden'>
      <SidebarHeader className='h-14 items-start justify-center'>
        <Button
          asChild
          variant='ghost'
          className='py-0 px-2 w-full hover:bg-sidebar-accent/10! text-sidebar-foreground!'
        >
          <Link to={Route.root} className='flex-start gap-x-2'>
            <img
              src='/logo.svg'
              alt='Portfolio Console'
              className='size-6 object-contain'
            />
            <h6 className='flex flex-col tracking-widest'>
              <span className='text-base font-bold leading-none'>
                Portfolio
              </span>
              <span className='text-xs font-medium leading-none'>Console</span>
            </h6>
          </Link>
        </Button>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Tables</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu className='space-y-1'>
              {TABLES.map((table) => {
                const Icon = table.icon;

                const isActive = parts.includes(table.label.toLowerCase());

                return (
                  <SidebarMenuItem key={table.label}>
                    <SidebarMenuButton
                      asChild
                      isActive={isActive}
                      className='hover:bg-sidebar-accent/10 hover:text-sidebar-foreground dark:text-sidebar-accent-foreground'
                    >
                      <Link to={table.url}>
                        <Icon />
                        <span>{table.label}</span>
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                );
              })}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        <SidebarGroup />
      </SidebarContent>
      <SidebarFooter />
    </ShadcnSidebar>
  );
}
