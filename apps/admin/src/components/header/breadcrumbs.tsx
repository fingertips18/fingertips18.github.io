import { useLocation } from 'react-router-dom';

import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/shadcn/breadcrumb';
import { TABLE_BREADCRUMB_MAP } from '@/constants/tables';
import { Route } from '@/routes/route';

export function Breadcrumbs() {
  const { pathname } = useLocation();

  // Split and filter to remove empty parts
  const parts = pathname.split('/').filter(Boolean); // e.g. ['project', 'add']

  // Get the main route (e.g. 'project', 'education', 'skill')
  const mainPath = `/${parts[0] ?? ''}`; // e.g. '/project'
  const subPath = parts[1]; // e.g. 'add' or 'edit'

  const mainBreadcrumb =
    TABLE_BREADCRUMB_MAP[mainPath as keyof typeof TABLE_BREADCRUMB_MAP];
  const subBreadcrumb = mainBreadcrumb?.subPaths.find((sp) =>
    subPath ? sp.url.endsWith(`/${subPath}`) : false,
  );

  // Use link only if there's a sub path
  const MainComp = subBreadcrumb ? BreadcrumbLink : BreadcrumbPage;

  return (
    <Breadcrumb className='ml-2'>
      <BreadcrumbList>
        <BreadcrumbItem>
          <BreadcrumbLink
            href={Route.root}
            className='outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 ring-offset-background rounded-xs'
          >
            Dashboard
          </BreadcrumbLink>
        </BreadcrumbItem>

        {mainBreadcrumb && (
          <>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <MainComp
                href={mainBreadcrumb.url}
                className='outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 ring-offset-background rounded-xs'
              >
                {mainBreadcrumb.label}
              </MainComp>
            </BreadcrumbItem>
          </>
        )}

        {subBreadcrumb && (
          <>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <BreadcrumbPage>{subBreadcrumb.label}</BreadcrumbPage>
            </BreadcrumbItem>
          </>
        )}
      </BreadcrumbList>
    </Breadcrumb>
  );
}
