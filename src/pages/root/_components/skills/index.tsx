import { useLenis } from "lenis/react";
import { useRef } from "react";
import { Link } from "react-router-dom";

import { QUERYELEMENT, ROOTSECTION } from "@/constants/enums";
import { useObserver } from "@/lib/hooks/useObserver";
import { cn } from "@/lib/utils";
import { AppRoutes } from "@/routes/app-routes";

import { Backend } from "./backend";
import { Frontend } from "./frontend";
import { GradientOverlay } from "./gradient-overlay";
import { Others } from "./others";
import { Tools } from "./tools";

const Skills = () => {
  const sectionRef = useRef<HTMLElement | null>(null);
  const { isVisible } = useObserver({ elementRef: sectionRef });
  const lenis = useLenis();

  const handleScroll = () => {
    if (!lenis) return;

    lenis.scrollTo(0);
  };

  return (
    <section
      id={ROOTSECTION.skills}
      ref={sectionRef}
      className={cn(
        "min-h-dvh h-dvh pt-14 flex-between flex-col gap-y-6 border-b",
        QUERYELEMENT.rootSection
      )}
    >
      <div className="leading-none flex-center flex-col">
        <h4 className="text-xs lg:text-sm font-bold text-center tracking-widest pt-6 lg:pb-2">
          SKILLS
        </h4>
        <p className="text-xl lg:text-5xl text-center">
          Innovate, Implement, <span className="text-primary">Repeat.</span>
        </p>
        <p className="text-xs lg:text-sm text-muted-foreground text-center lg:mt-2 w-3/4 lg:w-full">
          Showcasing the skills I've developed and refined over the past 3
          years.
        </p>
        <Link
          to={AppRoutes.skills}
          onClick={handleScroll}
          className="mt-2 text-sm hover:text-accent hover:drop-shadow-purple-glow underline-offset-4 hover:underline"
        >
          View All
        </Link>
      </div>
      <div
        className={cn(
          "w-full h-full flex-center flex-col gap-y-4 relative transition-opacity duration-1000 ease-in-out",
          isVisible ? "opacity-100" : "opacity-0"
        )}
      >
        {isVisible && (
          <>
            <Frontend />
            <Backend />
            <Others />
            <Tools />
            <GradientOverlay />
          </>
        )}
      </div>
      <p className="text-xs text-muted-foreground text-center max-w-screen-sm mx-auto w-4/5 lg:w-full lg:mt-6 pb-6">
        Currently expanding my skill set by delving into{" "}
        <span className="text-foreground/80">DevOps</span> practices, focusing
        on automation, CI/CD, and infrastructure management to enhance
        development and operational efficiency.
      </p>
    </section>
  );
};

export { Skills };
