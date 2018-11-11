import '@babel/polyfill'
import Vue from 'vue'
import './plugins/vuetify'
import App from './App.vue'
import router from './router'
import './registerServiceWorker'
import auth from './auth'
import axios from 'axios';

axios.defaults.baseURL = 'http://localhost:8000';
Vue.prototype.$http = axios;

Vue.config.productionTip = false

Vue.prototype.$auth = auth;

auth.checkAuthentication();

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
