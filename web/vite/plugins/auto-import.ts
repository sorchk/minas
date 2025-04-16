import autoImport from 'unplugin-auto-import/vite';
import { Plugin } from 'vite';

export default function createAutoImport(): Plugin {
  return autoImport({
    imports: ['vue', 'vue-router', 'vuex'],
    dts: false,
  });
}
