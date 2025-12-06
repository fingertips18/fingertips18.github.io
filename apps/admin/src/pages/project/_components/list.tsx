import { TriangleAlert } from 'lucide-react';
import { useEffect, useMemo } from 'react';
import { useSearchParams } from 'react-router-dom';

import { Skeleton } from '@/components/shadcn/skeleton';
import { APIRoute } from '@/constants/api';
import { useFetch } from '@/hooks/useFetch';
import { mapProject, type Project } from '@/types/project';

import { Card } from './card';

export function List() {
  const [searchParams, setSearchParams] = useSearchParams();
  const page = searchParams.get('page');
  const pageSize = searchParams.get('page_size');
  const sortOrder = searchParams.get('sort_order');
  const type = searchParams.get('type');

  const newSearchParams = useMemo(() => {
    const params = new URLSearchParams();
    params.set('page', page || '1');
    params.set('page_size', pageSize || '10');
    params.set('sort_by', 'created_at'); // Force created_at only for now

    if (sortOrder) {
      params.set('sort_ascending', sortOrder === 'asc' ? 'true' : 'false');
    }

    if (type) {
      params.set('type', type);
    }

    return params;
  }, [page, pageSize, sortOrder, type]);

  useEffect(() => {
    if (searchParams.keys.length !== 0) return;

    const readableSearchParams = new URLSearchParams(newSearchParams);
    const sortAscending = readableSearchParams.get('sort_ascending');
    readableSearchParams.delete('sort_ascending');
    readableSearchParams.set(
      'sort_order',
      sortAscending === 'true' ? 'asc' : 'desc',
    );
    readableSearchParams.set('sort_by', sortAscending ? 'oldest' : 'latest');

    setSearchParams(readableSearchParams);
  }, [newSearchParams, searchParams.keys.length, setSearchParams]);

  const { data, loading, error } = useFetch<Project[]>({
    url: `${APIRoute.project}s?${newSearchParams.toString()}`,
    method: 'GET',
    toastOptions: {
      errorTitle: 'Failed to list projects',
      errorMessage: 'Projects could not be loaded. Try again later.',
    },
  });

  if (loading) {
    return <ListSkeleton />;
  }

  if (error) {
    return (
      <div className='m-auto flex-center flex-col gap-y-1'>
        <TriangleAlert
          aria-hidden='true'
          className='size-6 lg:size-8 text-foreground/50'
        />
        <p className='placeholder'>
          We couldn’t load your projects. Please try again later.
        </p>
      </div>
    );
  }

  if (!data) {
    return <ListSkeleton />;
  }

  if (data.length === 0) {
    return (
      <p className='placeholder m-auto'>You haven't added any projects yet.</p>
    );
  }

  const projects = data.map((p) => mapProject(p));

  return (
    <div className='mt-4 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4'>
      {projects.map((project) => (
        <Card key={project.id} project={project} />
      ))}
    </div>
  );
}

function ListSkeleton() {
  return (
    <div className='mt-4 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4'>
      {[...Array(12)].map((_, i) => (
        <Skeleton key={i} className='aspect-square bg-primary/15' />
      ))}
    </div>
  );
}
