import vue from '@vitejs/plugin-vue';

import createAutoImport from './auto-import';
import createCompression from './compression';

export default function createVitePlugins(viteEnv, isBuild = false) {
  const vitePlugins = [vue()];
  vitePlugins.push(createAutoImport());
  isBuild && vitePlugins.push(...createCompression(viteEnv));
  return vitePlugins;
}
