import { VerticalTimelineElement } from "react-vertical-timeline-component";
import { SiYoutube, SiYoutubeHex } from "@icons-pack/react-simple-icons";
import { Link } from "react-router-dom";

import { Badge } from "@/components/shadcn/badge";

interface EducationItemProps {
  source: string;
  name: string;
  logo: string;
  sub?: {
    name: string;
    desc: string;
    logo: string;
  };
  department?: string;
  date: string;
  honor: string;
  desc: string;
  study?: {
    title: string;
    desc: string;
    stack: string[];
    demo: string;
  };
  projects?: {
    title: string;
    desc: string;
    stack: string[];
    demo?: string;
  }[];
  type: string;
}

const EducationItem = ({
  source,
  name,
  logo,
  sub,
  department,
  date,
  honor,
  desc,
  study,
  projects,
}: EducationItemProps) => {
  return (
    <VerticalTimelineElement
      contentStyle={{
        background: "hsl(var(--secondary) / 0.2)",
        border: "1px solid hsl(var(--secondary) / 0.5)",
        color: "hsl(var(--secondary-foreground))",
        display: "flex",
        flexDirection: "column",
        boxShadow: "hsl(var(--primary) / 0.2) 0px 4px 24px",
        borderRadius: "8px",
      }}
      contentArrowStyle={{
        borderRight: "8px solid  hsl(var(--secondary) / 0.8)",
      }}
      date={date}
      iconStyle={{
        backgroundColor: "#FFF",
        boxShadow: "hsl(var(--primary)) 0px 4px 24px",
        outline: "2px solid hsl(var(--primary))",
      }}
      icon={
        <Link to={source} target="_blank">
          <img
            src={logo}
            alt={name}
            className="rounded-full w-full h-full object-cover scale-90"
          />
        </Link>
      }
    >
      <div className="flex items-start gap-x-4">
        <img
          src={logo}
          alt={name}
          className="rounded-sm drop-shadow-primary-glow h-16 w-16 object-cover"
        />
        <div className="leading-tight space-y-0.5">
          <h3 className="font-bold">{name}</h3>
          <h4 className="text-sm text-secondary-foreground/80">{department}</h4>
          <h5 className="text-sm text-secondary-foreground/40">{date}</h5>
          <h6 className="text-xs text-secondary-foreground/40 font-semibold">
            {honor}
          </h6>
        </div>
      </div>

      {sub && (
        <div className="flex gap-x-4 mt-2">
          <img
            src={sub.logo}
            alt={sub.name}
            className="rounded-sm drop-shadow-primary-glow h-16 w-16 object-cover"
          />
          <div className="leading-tight space-y-0.5">
            <h3 className="font-bold">{sub.name}</h3>
            <h6 className="text-xs text-secondary-foreground/40 font-semibold">
              {sub.desc}
            </h6>
          </div>
        </div>
      )}

      <p className="!text-sm text-muted-foreground">{desc}</p>

      {study && (
        <div className="space-y-1.5">
          <p className="!font-semibold !text-sm">Thesis Study</p>
          <h5 className="!text-sm !font-normal text-foreground/80">
            {study.title}
          </h5>
          <p className="!text-xs text-muted-foreground">{study.desc}</p>
          <div className="flex items-start flex-wrap gap-1">
            {study.stack.map((s) => (
              <Badge key={s} className="bg-background/50">
                {s}
              </Badge>
            ))}
          </div>
          <Link to={study.demo} target="_blank">
            <Badge className="bg-secondary/20 whitespace-nowrap gap-x-2 py-1 px-2.5 cursor-pointer w-fit mt-2">
              <SiYoutube color={SiYoutubeHex} className="w-4 h-4" />
              {study.title!.split(" ")[0].replace(":", "")} Demo
            </Badge>
          </Link>
        </div>
      )}

      {projects && (
        <div className="space-y-1.5">
          <p className="!font-semibold !text-sm">Projects</p>
          <div className="space-y-2.5 mt-4">
            {projects.map((p) => (
              <div key={p.title} className="space-y-2">
                <h5 className="!text-sm !font-bold text-foreground/80">
                  {p.title}
                </h5>
                <p className="!text-xs text-muted-foreground !m-0">{p.desc}</p>
                <div className="flex items-start flex-wrap gap-1">
                  {p.stack.map((s) => (
                    <Badge key={s} className="bg-background/50">
                      {s}
                    </Badge>
                  ))}
                </div>
                {p.demo && (
                  <Link to={p.demo} target="_blank">
                    <Badge
                      key={p.title}
                      className="bg-secondary/20 whitespace-nowrap gap-x-2 py-1 px-2.5 cursor-pointer w-fit mt-1.5"
                    >
                      <SiYoutube color={SiYoutubeHex} className="w-4 h-4" />
                      {p.title} Demo
                    </Badge>
                  </Link>
                )}
              </div>
            ))}
          </div>
        </div>
      )}
    </VerticalTimelineElement>
  );
};

export default EducationItem;
