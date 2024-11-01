import { QUERYELEMENT, ROOTSECTION } from "@/constants/enums";
import { useMounted } from "@/lib/hooks/useMounted";
import { BUILDS } from "@/constants/collections";
import { WAVE } from "@/constants/assets";
import { cn } from "@/lib/utils";

import { ProfilePicture } from "./profile-picture";
import { ResumeButton } from "./resume-button";
import { Introduction } from "./introduction";
import { TypingTexts } from "./typing-texts";
import SocialButtons from "./social-buttons";

const Hero = () => {
  const isMounted = useMounted();

  return (
    <section
      className={cn(
        "min-h-dvh flex-center flex-col gap-y-12 lg:gap-y-24 p-6 lg:py-6 relative border-b lg:px-4 xl:px-0",
        QUERYELEMENT.rootSection
      )}
      id={ROOTSECTION.about}
    >
      <div className="mt-14 flex-center lg:flex-between flex-col-reverse lg:flex-row gap-y-4 lg:gap-y-8 gap-x-24 w-full">
        <div
          className={cn(
            "flex items-center lg:items-start flex-col lg:gap-2 transition-opacity duration-500 ease-in-out",
            isMounted ? "opacity-100" : "opacity-0"
          )}
        >
          <div className="flex items-start justify-center gap-x-2 relative">
            <p className="lg:text-xl font-semibold">Hi there!</p>
            <img
              src={WAVE}
              alt="Wave"
              width={181}
              height={193}
              className="w-6 lg:w-10 h-5 lg:h-8 relative -top-0.5 lg:-top-1.5"
            />
          </div>
          <h1 className="text-2xl lg:text-4xl font-bold flex items-center flex-col lg:flex-row">
            I'm Ghian Carlos Tan{" "}
            <span className="text-sm lg:text-lg font-semibold text-muted-foreground lg:ml-2">
              (Fingertips)
            </span>
          </h1>
          <TypingTexts texts={BUILDS} />
          <Introduction />
          <ResumeButton />
        </div>

        <ProfilePicture />
      </div>

      <SocialButtons isMounted={isMounted} />
    </section>
  );
};

export { Hero };
