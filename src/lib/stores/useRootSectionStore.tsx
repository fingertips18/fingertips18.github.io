import { create } from 'zustand';
import { createJSONStorage, persist } from 'zustand/middleware';

import { ROOTSECTION } from '@/constants/enums';

interface RootSectionStore {
  active: string;
  onActive: (active: string) => void;
  onClear: () => void;
}

const useRootSectionStore = create(
  persist<RootSectionStore>(
    (set) => ({
      active: ROOTSECTION.about,
      onActive: (active: string) => set({ active }),
      onClear: () => set({ active: ROOTSECTION.about }),
    }),
    { name: 'root-section', storage: createJSONStorage(() => sessionStorage) },
  ),
);

export { useRootSectionStore };
