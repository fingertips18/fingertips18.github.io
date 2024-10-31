import { Link } from "react-router-dom";
import { useLenis } from "lenis/react";

import { DARKLOGO, LIGHTLOGO } from "@/constants/assets";
import { Theme, useTheme } from "@/lib/hooks/useTheme";
import { Skeleton } from "@/components/shadcn/skeleton";
import { useMounted } from "@/lib/hooks/useMounted";
import { useResize } from "@/lib/hooks/useResize";
import { AppRoutes } from "@/routes/app-routes";

const Logo = () => {
  const { theme } = useTheme();
  const isMounted = useMounted();
  const lenis = useLenis();
  const { width } = useResize();

  if (!isMounted) {
    return <Skeleton className="w-20 lg:w-28 h-6 lg:h-8" />;
  }

  const logo =
    theme === Theme.dark || theme === Theme.system ? DARKLOGO : LIGHTLOGO;

  const onClick = () => lenis?.scrollTo(0);

  const lg = width > 1024;

  return (
    <Link
      to={AppRoutes.root}
      className="hover:scale-95 transition-all hover:drop-shadow-primary-glow"
      onClick={onClick}
    >
      <img
        src={logo}
        alt="Logo"
        width={lg ? 102 : 68}
        height={lg ? 24 : 16}
        className="h-4 lg:h-6 object-contain"
        loading="eager"
      />
    </Link>
  );
};

export { Logo };
