import { Link } from 'react-router-dom';

import { Route } from '@/routes/route';

export function Logo() {
  return (
    <Link to={Route.root} className='flex-start gap-x-2'>
      <img
        src='/logo.svg'
        alt='Portfolio Console'
        loading='eager'
        className='size-6 object-contain'
      />
      <h1 className='text-lg font-bold'>Portfolio Console</h1>
    </Link>
  );
}
