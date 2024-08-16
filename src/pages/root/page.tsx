import { Experience } from "./_components/experience";
import { Education } from "./_components/education";
import { Projects } from "./_components/projects";
import { Contact } from "./_components/contact";
import { Skills } from "./_components/skills";
import { Hero } from "./_components/hero";

const RootPage = () => {
  return (
    <>
      <Hero />
      <Skills />
      <Experience />
      <Projects />
      <Education />
      <Contact />
    </>
  );
};

export default RootPage;
