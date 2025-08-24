import { create } from 'zustand';
import { createJSONStorage, persist } from 'zustand/middleware';

import { ROOTSECTION } from '@/constants/enums';

interface RootSectionStore {
  active: ROOTSECTION;
  onActive: (active: ROOTSECTION) => void;
}

const useRootSectionStore = create(
  persist<RootSectionStore>(
    (set) => ({
      active: ROOTSECTION.about,
      onActive: (active: ROOTSECTION) => set({ active }),
    }),
    { name: 'root-section', storage: createJSONStorage(() => sessionStorage) },
  ),
);

export { useRootSectionStore };
