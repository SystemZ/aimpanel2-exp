import Vue from 'vue'
import axios from 'axios';
import router from '../router';

export default new Vue({
  data: {
    logged: false,
  },

  methods: {
    login(data, redirect) {
      this.$http.post('/v1/auth/login', data).then(res => {
        if(res.data.token) {
          localStorage.setItem('token', res.data.token);
          this.logged = true;

          axios.defaults.headers.common['Authorization'] = this.getAuthorizationHeader();

          if(redirect) {
            router.push(redirect);
          }
        }
      }).catch(e => {
        console.error(e);
      })
    },
    logout() {
      this.logged = false;
      localStorage.removeItem('token');
      router.push('/');
    },
    checkAuthentication() {
      let token = localStorage.getItem('token');
      if(token) {
        this.logged = true;
        axios.defaults.headers.common['Authorization'] = this.getAuthorizationHeader();
      } else {
        this.logged = false;
      }
    },
    getAuthorizationHeader() {
      return 'Bearer ' + localStorage.getItem('token');
    }
  }
});