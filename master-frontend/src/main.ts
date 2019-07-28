import Vue from 'vue';
import App from './App.vue';
import router from './router';
import store from './store';
import './registerServiceWorker';
import vuetify from './plugins/vuetify';
import '@babel/polyfill';
import axios, {AxiosStatic} from 'axios';
import auth from './auth';

Vue.config.productionTip = false;

axios.defaults.baseURL = process.env.API_URL;

Vue.prototype.$http = axios;
declare module  'vue/types/vue' {
  interface Vue {
    $http: AxiosStatic;
  }
}

Vue.prototype.$auth = auth;
auth.checkAuthentication();

new Vue({
  router,
  store,
  vuetify,
  render: (h) => h(App)
}).$mount('#app');
