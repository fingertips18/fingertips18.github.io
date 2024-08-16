import { Bolt } from "lucide-react";

import { QUERYELEMENTS, ROOTSECTIONS } from "@/constants/enums";
import { cn } from "@/lib/utils";

const Contact = () => {
  return (
    <section
      className={cn("min-h-dvh h-dvh", QUERYELEMENTS.rootSection)}
      id={ROOTSECTIONS.contact}
    >
      <div className="flex-center gap-x-4 w-full h-full">
        <Bolt className="w-8 h-8 animate-spin" />
        <p className="text-lg font-bold">Contact Under Construction</p>
      </div>
    </section>
  );
};

export { Contact };
