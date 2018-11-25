import Vue from 'vue'
import axios from 'axios';
import router from '../router';
import jwt_decode from 'jwt-decode';

export default new Vue({
  data: {
    email: '',
    username: '',
    logged: false,
  },

  methods: {
    login(context, data, redirect) {
      this.$http.post('/v1/auth/login', data).then(res => {
        if (res.data.token) {
          let decoded = jwt_decode(res.data.token);
          this.username = decoded.username;
          this.email = decoded.email;

          localStorage.setItem('token', res.data.token);
          this.logged = true;

          axios.defaults.headers.common['Authorization'] = this.getAuthorizationHeader();

          if (redirect) {
            router.push(redirect);
          }
        }
      }).catch(e => {
        context.loginError = e.response.data.message
      })
    },
    register(context, data, redirect) {
      this.$http.post('/v1/auth/register', data).then(res => {
        if (res.data.token) {
          localStorage.setItem('token', res.data.token);
          this.logged = true;

          axios.defaults.headers.common['Authorization'] = this.getAuthorizationHeader();

          if (redirect) {
            router.push(redirect);
          }
        }
      }).catch(e => {
        context.registerError = e.response.data.message
      });
    },
    logout() {
      this.logged = false;
      localStorage.removeItem('token');
      router.push('/');
    },
    checkAuthentication() {
      let token = localStorage.getItem('token');
      if (token) {
        this.logged = true;
        let decoded = jwt_decode(token);
        this.username = decoded.username;
        this.email = decoded.email;
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