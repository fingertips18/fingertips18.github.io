import { Link } from "react-router-dom";

import { SOCIALS } from "@/constants/collections";
import { Hint } from "@/common/components/hint";

const Socials = () => {
  return (
    <ul className="flex items-start gap-x-2">
      {SOCIALS.filter((s) => s.label !== "LinkedIn").map((s) => (
        <Hint key={`footer-${s.label}`} asChild label={s.label} side="top">
          <Link
            to={s.href}
            target="_blank"
            className="hover:drop-shadow-primary-glow transition-all"
          >
            <li className="rounded-full border border-primary/50 hover:border-primary bg-primary/20 hover:bg-primary/50 p-1.5 lg:p-2.5">
              <s.icon className="w-4 h-4 pointer-events-none" />
            </li>
          </Link>
        </Hint>
      ))}
    </ul>
  );
};

export { Socials };
