import Vue from 'vue';
import App from './App.vue';
import router from './router';
import store from './store';
import './registerServiceWorker';
import vuetify from './plugins/vuetify';
import '@babel/polyfill';
import axios, {AxiosStatic} from 'axios';
import auth, {AuthInterface} from './auth';
// @ts-ignore
import {EventSourcePolyfill} from 'event-source-polyfill';

Vue.config.productionTip = false;

axios.defaults.baseURL = process.env.VUE_APP_API_URL;

Vue.prototype.$http = axios;
Vue.prototype.$auth = auth;
Vue.prototype.$eventSource = EventSourcePolyfill;
Vue.prototype.$apiUrl = process.env.VUE_APP_API_URL;

declare module 'vue/types/vue' {
    interface Vue {
        $http: AxiosStatic;
        $auth: AuthInterface;
        $eventSource: any;
        $apiUrl: string;
    }
}

auth.checkAuthentication();

new Vue({
    router,
    store,
    vuetify,
    render: (h) => h(App)
}).$mount('#app');
