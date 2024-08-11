import { Link } from "react-router-dom";

import { Button } from "@/common/components/shadcn/button";

const ResumeButton = () => {
  return (
    <Link
      to={
        "https://drive.google.com/file/d/1zeF5O_iNHSwUhJrxzS5HjV6CdZmfxVm_/view?usp=sharing"
      }
      target="_blank"
      className="hidden lg:block"
    >
      <Button className="hover:scale-95 transition-all rounded-full hover:drop-shadow-glow">
        Check Resume
      </Button>
    </Link>
  );
};

export { ResumeButton };
