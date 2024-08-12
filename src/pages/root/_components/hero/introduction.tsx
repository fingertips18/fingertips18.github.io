import { useResize } from "@/lib/hooks/use-resize";

const Introduction = () => {
  const { width } = useResize();

  const lg = width > 1024;

  return lg ? (
    <>
      <p className="text-muted-foreground text-sm mt-2">
        I am a dedicated{" "}
        <span className="text-foreground/80">full-stack developer</span> with a
        strong foundation in both{" "}
        <span className="text-foreground/80">web</span> and{" "}
        <span className="text-foreground/80">mobile technologies.</span>{" "}
        Leveraging a{" "}
        <span className="text-foreground/80">
          Bachelor's degree in Computer Science
        </span>{" "}
        and <span className="text-foreground/80">3</span> professional{" "}
        <span className="text-foreground/80">years</span> of{" "}
        <span className="text-foreground/80">experience</span>.
      </p>
      <p className="text-muted-foreground text-sm">
        I specialize in creating robust, scalable solutions using{" "}
        <span className="text-foreground/80">React.js</span>,{" "}
        <span className="text-foreground/80">React Native</span>, and{" "}
        <span className="text-foreground/80">Flutter</span> for front-end
        development, coupled with{" "}
        <span className="text-foreground/80">Express.js</span> for backend
        services. My expertise extends to working with databases and cloud
        platforms, including <span className="text-foreground/80">MongoDB</span>
        , <span className="text-foreground/80">Supabase</span>, and{" "}
        <span className="text-foreground/80">Firebase</span>, as well as
        utilizing <span className="text-foreground/80">Prisma</span> for ORM and{" "}
        <span className="text-foreground/80">TensorFlow</span> for machine
        learning applications.
      </p>
      <p className="text-muted-foreground text-sm">
        In addition to my primary focus on building user-centric web and mobile
        applications, I am also passionate about{" "}
        <span className="text-foreground/80">game development</span>. I create
        games using <span className="text-foreground/80">Unity</span>,{" "}
        <span className="text-foreground/80">Flutter</span> and{" "}
        <span className="text-foreground/80">Vanilla JavaScript</span> as a
        hobby, which enhances my problem-solving skills and creativity in
        software design.
      </p>
    </>
  ) : (
    <p className="text-center text-muted-foreground text-xs sm:text-sm mt-2 max-w-screen-sm">
      I am a <span className="text-foreground/80">full-stack developer</span>{" "}
      with a{" "}
      <span className="text-foreground/80">
        Bachelor's degree in Computer Science
      </span>{" "}
      and <span className="text-foreground/80">3</span> professional{" "}
      <span className="text-foreground/80">years</span> of{" "}
      <span className="text-foreground/80">experience</span>. I excel in
      creating scalable solutions using{" "}
      <span className="text-foreground/80">React.js</span>,{" "}
      <span className="text-foreground/80">React Native</span>,{" "}
      <span className="text-foreground/80">Flutter</span>, and{" "}
      <span className="text-foreground/80">Express.js</span>. My skills include
      working with <span className="text-foreground/80">MongoDB</span>,
      <span className="text-foreground/80">Supabase</span>,{" "}
      <span className="text-foreground/80">Firebase</span>, and{" "}
      <span className="text-foreground/80">Prisma</span>. Additionally, I am
      passionate about game development, creating games with{" "}
      <span className="text-foreground/80">Unity</span>,{" "}
      <span className="text-foreground/80">Flutter</span>, and{" "}
      <span className="text-foreground/80">Vanilla JavaScript</span> to enhance
      my problem-solving and design skills.
    </p>
  );
};

export { Introduction };
