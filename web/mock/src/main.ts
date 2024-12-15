/**
 * main.ts
 *
 * Bootstraps Vuetify and other plugins then mounts the App
 */

// Components
import { createApp } from 'vue';
import App from './App.vue';

// Create and mount the app
const app = createApp(App);
app.mount('#app');