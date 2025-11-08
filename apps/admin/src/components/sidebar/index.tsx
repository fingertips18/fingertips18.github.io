import { Link } from 'react-router-dom';

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
  return (
    <ShadcnSidebar className='overflow-x-hidden'>
      <SidebarHeader className='h-14 items-start justify-center'>
        <Button
          asChild
          variant='ghost'
          className='py-0 px-2 w-full hover:bg-sidebar-accent! text-sidebar-foreground!'
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
            <SidebarMenu>
              {TABLES.map((table) => {
                const Icon = table.icon;

                return (
                  <SidebarMenuItem key={table.label}>
                    <SidebarMenuButton asChild>
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
