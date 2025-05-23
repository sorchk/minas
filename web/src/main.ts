import { createApp } from 'vue'
import App from './App.vue'
import { router } from './router/router'
import { store } from './store'
import { directive } from '@/directive';
import i18n from './locales'

import DagDesignerInstaller from '@/components/dagdesigner/installer';

const app = createApp(App).use(router).use(store).use(i18n).use(DagDesignerInstaller);

directive(app);
app.mount('#app');
