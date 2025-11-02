import { Logo } from './logo';

export function Header() {
  return (
    <header className='h-14 w-full max-w-7xl mx-auto flex-between px-4 border-b'>
      <Logo />
    </header>
  );
}
