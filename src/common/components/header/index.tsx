import { ModeToggle } from "./mode-toggle";
import { Navbar } from "./navbar";
import { Logo } from "./logo";

const Header = () => {
  return (
    <header className="h-14 w-full fixed z-50 top-0 flex-center bg-background/50 backdrop-blur-lg border-b px-4 md:px-8 lg:px-0 blur-performance">
      <div className="flex-between h-full w-full max-w-screen-lg">
        <Logo />
        <Navbar />
        <div className="hidden lg:flex lg:items-end">
          <ModeToggle />
        </div>
      </div>
    </header>
  );
};

export { Header };
