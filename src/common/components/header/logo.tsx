import { Link } from "react-router-dom";

import { DARKLOGO, LIGHTLOGO } from "@/constants/assets";
import { Theme, useTheme } from "@/lib/hooks/use-theme";
import { AppRoutes } from "@/routes/app-routes";

const Logo = () => {
  const { theme } = useTheme();

  const logo =
    theme === Theme.dark || theme === Theme.system ? DARKLOGO : LIGHTLOGO;

  return (
    <Link
      to={AppRoutes.root}
      className="hover:scale-95 transition-all hover:drop-shadow-primary-glow"
    >
      <img src={logo} alt="Logo" className="h-4 lg:h-6" />
    </Link>
  );
};

export { Logo };
