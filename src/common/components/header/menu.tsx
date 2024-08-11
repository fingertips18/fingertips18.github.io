import { MENU } from "@/constants/collections";

const Menu = () => {
  return (
    <nav className="flex-center flex-1 px-4">
      <ul className="flex-center gap-x-8">
        {MENU.map((m, i) => (
          <li
            key={`${m.label}-${i}`}
            className="capitalize font-semibold leading-none hover:scale-95 transition-all cursor-pointer hover:drop-shadow-glow"
          >
            {m.label}
          </li>
        ))}
      </ul>
    </nav>
  );
};

export default Menu;
