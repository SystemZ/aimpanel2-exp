import Vue from 'vue';
import Router from 'vue-router';
import Home from './views/Home.vue';

Vue.use(Router);

export default new Router({
    routes: [
        {
            path: '/',
            name: 'home',
            meta: {
                title: 'Home',
            },
            component: Home,
        },
        {
            path: '/login',
            name: 'login',
            meta: {
                title: 'Login',
            },
            component: () => import(/* webpackChunkName: "login" */ './views/Login.vue'),
        },
        {
            path: '/license',
            name: 'license',
            meta: {
                title: 'License',
            },
            component: () => import(/* webpackChunkName: "license" */ './views/License.vue'),
        },
        {
            path: '/help',
            name: 'help',
            meta: {
                title: 'Help',
            },
            component: () => import(/* webpackChunkName: "help" */ './views/Help.vue'),
        },
        {
            path: '/profile',
            name: 'profile',
            meta: {
                title: 'Profile',
                authRequired: true,
            },
            component: () => import(/* webpackChunkName: "hosts" */ './views/Profile.vue'),
        },
        {
            path: '/hosts',
            name: 'hosts',
            meta: {
                title: 'Hosts',
                authRequired: true,
            },
            component: () => import(/* webpackChunkName: "hosts" */ './views/Hosts.vue'),
        },
        {
            path: '/host/:id',
            name: 'host',
            meta: {
                title: 'Host',
                authRequired: true,
            },
            component: () => import(/* webpackChunkName: "hosts" */ './views/Host.vue'),
        },
        {
            path: '/game-servers',
            name: 'game-servers',
            meta: {
                title: 'Game Servers',
                authRequired: true,
            },
            component: () => import(/* webpackChunkName: "game-servers" */ './views/GameServers.vue'),
        },
        {
            path: '*',
            name: '404',
            meta: {
                title: '404 - Page Not Found'
            },
            component: () => import(/* webpackChunkName: "404" */ './views/404.vue'),
        },
    ],
});
