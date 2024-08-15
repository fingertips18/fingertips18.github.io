import { Bolt, GraduationCap } from "lucide-react";

import { QUERYELEMENTS, ROOTSECTIONS } from "@/constants/enums";
import { cn } from "@/lib/utils";

const Education = () => {
  return (
    <section
      className={cn(
        "min-h-dvh h-dvh flex flex-col gap-y-2 lg:gap-y-6 border-b pt-14 pb-6 px-2 lg:px-0",
        QUERYELEMENTS.rootSection
      )}
      id={ROOTSECTIONS.education}
    >
      <div className="flex items-center gap-x-2 w-full pt-6 lg:relative">
        <span className="w-[32px] lg:w-[128px] h-1 rounded-full bg-muted-foreground tracking-widest" />
        <h2 className="text-lg lg:text-4xl font-bold">EDUCATION</h2>
        <GraduationCap className="w-5 lg:w-8 h-5 lg:h-8 sm:absolute xs:right-6 lg:right-0 opacity-50" />
      </div>
      <p className="text-xs lg:text-sm text-muted-foreground text-center lg:mt-2 w-3/4 mx-auto">
        Throughout my academic journey, each experience has played a distinct
        role in my development. Here’s an overview of the key milestones in my
        academic journey.
      </p>

      <div className="flex-center gap-x-4 w-full h-full">
        <Bolt className="w-8 h-8 animate-spin" />
        <p className="text-lg font-bold">Under Construction</p>
      </div>
      {/* <VerticalTimeline
        lineColor="hsl(var(--foreground) / 0.6)"
        className="mt-4 lg:mt-20"
      >
        <VerticalTimelineElement></VerticalTimelineElement>
      </VerticalTimeline> */}
    </section>
  );
};

export { Education };
