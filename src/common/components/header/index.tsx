import { ResumeButton } from "./resume-button";
import { Logo } from "./logo";
import { Menu } from "./menu";

const Header = () => {
  return (
    <header className="h-14 flex-center bg-light-background border-b px-4 md:px-8 lg:px-0">
      <div className="flex-between h-full w-full max-w-screen-lg">
        <Logo />
        <Menu />
        <ResumeButton />
      </div>
    </header>
  );
};

export { Header };
