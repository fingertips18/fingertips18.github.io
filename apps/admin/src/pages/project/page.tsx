import { FolderCode, Plus } from 'lucide-react';
import { Link } from 'react-router-dom';

import { Button } from '@/components/shadcn/button';
import { Route } from '@/routes/route';

import { List } from './_components/list';

export default function ProjectPage() {
  return (
    <section className='content padding flex flex-col'>
      <div className='flex-between gap-x-4'>
        <h1 className='flex-start gap-x-2 lg:text-2xl'>
          <FolderCode aria-hidden='true' className='size-6' />
          Projects
        </h1>
        <Button asChild>
          <Link to={`${Route.project}/add`} className='flex-center gap-x-2'>
            <Plus aria-hidden='true' className='size-4' />
            Add Project
          </Link>
        </Button>
      </div>

      <p className='text-muted-foreground text-sm'>
        A central hub for all the projects you’ve created or submitted.
      </p>

      <List />
    </section>
  );
}
