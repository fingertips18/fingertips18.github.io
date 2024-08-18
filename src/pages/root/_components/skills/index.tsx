import { QUERYELEMENTS, ROOTSECTIONS } from "@/constants/enums";
import { cn } from "@/lib/utils";

import { GradientOverlay } from "./gradient-overlay";
import { Frontend } from "./frontend";
import { Backend } from "./backend";
import { Others } from "./others";
import { Tools } from "./tools";

const Skills = () => {
  return (
    <section
      className={cn(
        "min-h-dvh h-dvh pt-14 flex-between flex-col gap-y-6 border-b",
        QUERYELEMENTS.rootSection
      )}
      id={ROOTSECTIONS.skills}
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
      </div>
      <div className="w-full flex-center flex-col gap-y-4 relative">
        <Frontend />
        <Backend />
        <Others />
        <Tools />
        <GradientOverlay />
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
