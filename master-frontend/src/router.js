import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'

Vue.use(Router)

const router = new Router({
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'home',
      meta: {
        title: 'Home'
      },
      component: Home
    },
    {
      path: '/hosts',
      name: 'hosts',
      meta: {
        title: 'Hosts',
        authRequired: true
      },
      component: () => import(/* webpackChunkName: "hosts" */ './views/Hosts.vue')
    },
    {
      path: '/game-servers',
      name: 'game-servers',
      meta: {
        title: 'Game Servers',
        authRequired: true
      },
      component: () => import(/* webpackChunkName: "game-servers" */ './views/GameServers.vue')
    },
    {
      path: '/license',
      name: 'license',
      meta: {
        title: 'License'
      },
      component: () => import(/* webpackChunkName: "license" */ './views/License.vue')
    },
    {
      path: '/help',
      name: 'help',
      meta: {
        title: 'Help'
      },
      // route level code-splitting
      // this generates a separate chunk (about.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import(/* webpackChunkName: "help" */ './views/Help.vue')
    },
    {
      path: '/host/:id',
      name: 'host',
      meta: {
        title: 'Host'
      },
      component: () => import(/* webpackChunkName: "host" */ './views/Host.vue')
    },
    {
      path: '/login',
      name: 'login',
      meta: {
        title: 'Login'
      },
      component: () => import(/* webpackChunkName: "login" */ './views/Login.vue')
    },
    {
      path: '*',
      name: '404',
      meta: {
        title: '404 - Page Not Found'
      },
      component: () => import(/* webpackChunkName: "404" */ './views/404.vue')
    }
  ]
});

router.beforeEach((to, from, next) => {
  if (to.meta.authRequired) {
    let token = localStorage.getItem('token')
    if (token) {
      next();
    } else {
      next('/login');
    }
  }
  next();
});


export default router;