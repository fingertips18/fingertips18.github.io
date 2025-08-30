declare module '@tailwindcss/vite' {
  import type { PluginOption } from 'vite';
  function tailwindcss(): PluginOption;
  export default tailwindcss;
}
