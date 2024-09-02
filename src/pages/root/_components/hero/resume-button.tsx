import { Link } from "react-router-dom";

const ResumeButton = () => {
  return (
    <Link
      to={
        "https://drive.google.com/file/d/1kS_OEN015F7EtHYZjqu7W2viITWOCquE/view?usp=sharing"
      }
      target="_blank"
    >
      <button
        className="py-4 w-[256px] bg-gradient-to-r from-[#310055] to-[#DC97FF]
        hover:scale-95 transition-all duration-500 ease-in-out rounded-full 
        hover:drop-shadow-purple-glow font-semibold text-lg mt-8 text-white"
      >
        Check Resume
      </button>
    </Link>
  );
};

export { ResumeButton };
