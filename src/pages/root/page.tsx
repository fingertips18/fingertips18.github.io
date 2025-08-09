import { Contact } from "./_components/contact";
import { Education } from "./_components/education";
import { Experience } from "./_components/experience";
import { Hero } from "./_components/hero";
import { Projects } from "./_components/projects";
import { Skills } from "./_components/skills";

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
