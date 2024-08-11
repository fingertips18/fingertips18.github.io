import { Link } from "react-router-dom";

import { AppRoutes } from "@/routes/app-routes";
import { LIGHTLOGO } from "@/constants/assets";

const Logo = () => {
  return (
    <Link
      to={AppRoutes.root}
      className="hover:scale-95 transition-all hover:drop-shadow-glow"
    >
      <img src={LIGHTLOGO} alt="Logo" className="h-6" />
    </Link>
  );
};

export default Logo;
