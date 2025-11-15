import tailwindcss from '@tailwindcss/vite';
import react from '@vitejs/plugin-react-swc';
import * as path from 'path';
import { defineConfig, loadEnv } from 'vite';

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '');

  return {
    plugins: [react(), tailwindcss()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
    server: {
      proxy: {
        '/github': {
          target: 'https://api.github.com',
          changeOrigin: true,
          rewrite: (requestPath) => requestPath.replace(/^\/github/, ''), // remove /github prefix
          configure: (proxy) => {
            proxy.on('proxyReq', (proxyReq) => {
              if (env.GITHUB_TOKEN) {
                proxyReq.setHeader(
                  'Authorization',
                  `Bearer ${env.GITHUB_TOKEN}`,
                );
              }
              proxyReq.setHeader(
                'Accept',
                'application/vnd.github.mercy-preview+json',
              );
            });
          },
        },
      },
    },
  };
});
