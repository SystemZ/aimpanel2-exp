export default {
  user: {
    logged: false
  },

  checkAuth() {
    var jwt = localStorage.getItem('token')
    this.user.logged = !!jwt;
  },

  getAuthorizationHeader() {
    return 'Bearer ' + localStorage.getItem('token');
  }
}