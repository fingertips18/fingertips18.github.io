import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './index.css';

// eslint-disable-next-line @typescript-eslint/no-unnecessary-type-assertion
createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <main>
      <h1>Admin Dashboard</h1>
    </main>
  </StrictMode>,
);
