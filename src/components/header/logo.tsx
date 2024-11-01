import { Link } from "react-router-dom";
import { useLenis } from "lenis/react";

import { DARKLOGO, LIGHTLOGO } from "@/constants/assets";
import { Skeleton } from "@/components/shadcn/skeleton";
import { Theme, useTheme } from "@/lib/hooks/useTheme";
import { useMounted } from "@/lib/hooks/useMounted";
import { AppRoutes } from "@/routes/app-routes";

const Logo = () => {
  const { theme } = useTheme();
  const isMounted = useMounted();
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
      className="hover:scale-105 transition-all duration-500 ease-in-out hover:drop-shadow-primary-glow"
      onClick={onClick}
    >
      <img
        src={logo}
        alt="Logo"
        width={899}
        height={212}
        className="w-[68px] lg:[102px] h-4 lg:h-6 object-contain"
        loading="eager"
      />
    </Link>
  );
};

export { Logo };
