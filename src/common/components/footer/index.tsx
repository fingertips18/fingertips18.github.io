import { Link } from "react-router-dom";

import { SOCIALS } from "@/constants/collections";
import { Hint } from "@/common/components/hint";

const Footer = () => {
  return (
    <footer className="w-full bg-secondary/10 border-t">
      <div className="flex-between p-4 lg:px-0 max-w-screen-lg mx-auto mt-4">
        <ul className="flex items-start gap-x-2">
          {SOCIALS.filter((s) => s.label !== "LinkedIn").map((s) => (
            <Hint key={`footer-${s.label}`} asChild label={s.label} side="top">
              <Link
                to={s.href}
                className="hover:drop-shadow-primary-glow transition-all"
              >
                <li className="rounded-full border border-primary/50 hover:border-primary bg-primary/20 hover:bg-primary/50 p-2.5">
                  <s.icon className="w-4 h-4 pointer-events-none" />
                </li>
              </Link>
            </Hint>
          ))}
        </ul>

        <div className="flex items-end text-xs gap-x-1.5">
          <p className="text-foreground/80">Designed & Developed by</p>
          <Link
            to={"https://linkedin.com/in/ghiantan"}
            className="font-semibold underline underline-offset-2 hover:drop-shadow-primary-glow transition-all"
          >
            Ghian Carlos Tan
          </Link>
        </div>
      </div>
    </footer>
  );
};

export { Footer };
