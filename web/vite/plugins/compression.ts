import compression from 'vite-plugin-compression';
import type { Plugin } from 'vite';

export default function createCompression(env: Record<string, string>): Plugin[] {
  const { VITE_BUILD_COMPRESS } = env;
  const plugin: Plugin[] = [];
  if (VITE_BUILD_COMPRESS) {
    const compressList = VITE_BUILD_COMPRESS.split(',');
    if (compressList.includes('gzip')) {
      plugin.push(
        compression({
          ext: '.gz',
          deleteOriginFile: false,
        })
      );
    }
    if (compressList.includes('brotli')) {
      plugin.push(
        compression({
          ext: '.br',
          algorithm: 'brotliCompress',
          deleteOriginFile: false,
        })
      );
    }
  }
  return plugin;
}
