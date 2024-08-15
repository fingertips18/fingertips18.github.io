import { Terminal } from "lucide-react";

import { QUERYELEMENTS, ROOTSECTIONS } from "@/constants/enums";
import { PROJECTS } from "@/constants/projects";
import { cn } from "@/lib/utils";

import { ProjectItem } from "./project-item";

const Projects = () => {
  return (
    <section
      className={cn(
        "min-h-dvh flex flex-col gap-y-2 lg:gap-y-6 border-b pt-14 pb-6 px-2 lg:px-0",
        QUERYELEMENTS.rootSection
      )}
      id={ROOTSECTIONS.projects}
    >
      <div className="flex items-center justify-end gap-x-2 w-full pt-6 lg:relative">
        <h2 className="text-lg lg:text-4xl font-bold">PROJECTS</h2>
        <span className="w-[32px] lg:w-[128px] h-1 rounded-full bg-muted-foreground tracking-widest" />
        <Terminal className="w-5 lg:w-8 h-5 lg:h-8 sm:absolute xs:left-6 lg:left-0 opacity-50" />
      </div>

      <p className="text-xs lg:text-sm text-muted-foreground text-center lg:mt-2 w-3/4 lg:w-full">
        I’ve developed various projects, ranging from web applications to
        Android apps. Here are a few highlights.
      </p>

      <div className="w-full grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 mt-8 gap-4">
        {PROJECTS.map((p) => (
          <ProjectItem key={p.name} {...p} />
        ))}
      </div>
    </section>
  );
};

export { Projects };
