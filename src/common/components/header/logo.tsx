import { Link } from "react-router-dom";
import { useLenis } from "lenis/react";

import { DARKLOGO, LIGHTLOGO } from "@/constants/assets";
import { Theme, useTheme } from "@/lib/hooks/use-theme";
import { useClient } from "@/lib/hooks/use-client";
import { AppRoutes } from "@/routes/app-routes";
import { Skeleton } from "../shadcn/skeleton";

const Logo = () => {
  const { theme } = useTheme();
  const isMounted = useClient();
  const lenis = useLenis();

  if (!isMounted) {
    return <Skeleton className="w-20 lg:w-28 h-6 lg:h-8" />;
  }

  const logo =
    theme === Theme.dark || theme === Theme.system ? DARKLOGO : LIGHTLOGO;

  const onClick = () => lenis?.scrollTo(0);

  return (
    <Link
      to={AppRoutes.root}
      className="hover:scale-95 transition-all hover:drop-shadow-primary-glow"
      onClick={onClick}
    >
      <img src={logo} alt="Logo" className="h-4 lg:h-6" />
    </Link>
  );
};

export { Logo };
