import vue from '@vitejs/plugin-vue';
import type { Plugin } from 'vite';

import createAutoImport from './auto-import';
import createCompression from './compression';

export default function createVitePlugins(viteEnv: Record<string, any>, isBuild: boolean = false): Plugin[] {
  const vitePlugins: Plugin[] = [vue()];
  vitePlugins.push(createAutoImport());
  isBuild && vitePlugins.push(...createCompression(viteEnv));
  return vitePlugins;
}
