import { Link } from 'react-router-dom';

import { Socials } from './socials';

const Footer = () => {
  const now = new Date();

  return (
    <footer className='w-full bg-secondary/10 border-t flex flex-col items-center'>
      <div className='w-full flex-between p-4 lg:px-0 max-w-screen-lg mx-auto mt-4 gap-x-12'>
        <Socials />

        <div className='flex flex-wrap justify-end text-xs gap-x-1.5'>
          <p className='text-foreground/80'>Designed & Developed by</p>
          <Link
            to={'https://linkedin.com/in/ghiantan'}
            target='_blank'
            className='underline underline-offset-2 hover:drop-shadow-primary-glow transition-all'
          >
            Fingertips
          </Link>
        </div>
      </div>

      <div className='mt-4 py-1.5 bg-secondary/20 w-full flex-center'>
        <p className='text-xs text-muted-foreground'>
          © {now.getUTCFullYear().toString()} Ghian Carlos Tan. All rights
          reserved.
        </p>
      </div>
    </footer>
  );
};

export { Footer };
