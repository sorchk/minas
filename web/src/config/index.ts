const config = window.config || {};
config.baseUrl = config.baseUrl == 'null' ? '' : (config.baseUrl || import.meta.env.VITE_APP_API);
export const baseUrl = config.baseUrl;
config.frontBaseUrl = config.frontBaseUrl == 'null' ? '' : (config.frontBaseUrl || import.meta.env.VITE_APP_BASE);
config.indexPage = config.indexPage || '/';
export const frontBaseUrl = config.frontBaseUrl;
export default config;