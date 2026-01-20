import { useState } from 'react';
import { Blurhash } from 'react-blurhash';
import { Link } from 'react-router-dom';

import { Button } from '@/components/shadcn/button';
import { cn } from '@/lib/utils';
import type { Project } from '@/types/project';

interface CardProps {
  project: Project;
}

export function Card({ project }: CardProps) {
  const [imageLoaded, setImageLoaded] = useState<boolean>(false);

  return (
    <div className='flex flex-col aspect-square rounded-md border overflow-hidden bg-gray-100 transition-all duration-400 ease-in-out hover:scale-95 hover:shadow-2xl'>
      <div className='relative aspect-video'>
        {project.previews[0]?.url && (
          <img
            src={project.previews[0].url}
            alt={`${project.title} preview`}
            sizes='(min-width: 1024px) 25vw, (min-width: 640px) 50vw, 100vw'
            onLoad={() => setImageLoaded(true)}
            className='absolute object-center object-cover size-full'
          />
        )}
        <Blurhash
          hash={project.blurhash}
          width='100%'
          height='100%'
          className={cn(
            'object-cover transition-opacity duration-400 ease-in-out',
            imageLoaded && 'opacity-0',
          )}
        />
      </div>
      <div className='flex-1 flex justify-between flex-col text-ellipsis overflow-hidden p-2'>
        <div className='space-y-1'>
          <h2 className='text-2xl sm:text-lg font-bold line-clamp-1'>
            {project.title}
          </h2>
          <p className='line-clamp-2 text-foreground/60 text-base sm:text-xs'>
            {project.description}
          </p>
        </div>
        <div className='flex-start gap-x-1 text-sm w-full'>
          <span className='font-semibold text-foreground/60'>Live:</span>
          <Button asChild variant='link' className='p-0 h-auto text-sm'>
            <Link to={project.link} target='_blank' rel='noopener noreferrer'>
              {project.link}
            </Link>
          </Button>
        </div>
      </div>
    </div>
  );
}
