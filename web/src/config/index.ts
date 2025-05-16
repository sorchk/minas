const config = window.config || {};
config.baseUrl = config.baseUrl == 'null' ? '' : (config.baseUrl || import.meta.env.VITE_APP_BASE);
export const baseUrl = config.baseUrl;
config.frontBaseUrl = config.frontBaseUrl == 'null' ? '' : (config.frontBaseUrl || import.meta.env.VITE_APP_BASE);
config.indexPage = config.indexPage || '/';
export const frontBaseUrl = config.frontBaseUrl;
export const apiUrl = config.apiUrl || import.meta.env.VITE_APP_BASE;
export default config;