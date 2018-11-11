import '@babel/polyfill'
import Vue from 'vue'
import './plugins/vuetify'
import App from './App.vue'
import router from './router'
import './registerServiceWorker'
import auth from './auth'
import axios from 'axios';

auth.checkAuth();

axios.defaults.headers.common['Authorization'] = auth.getAuthorizationHeader();
axios.defaults.baseURL = 'http://localhost:8000';

Vue.config.productionTip = false

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
