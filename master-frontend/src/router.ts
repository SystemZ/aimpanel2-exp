import Vue from 'vue';
import Router from 'vue-router';
import Home from './views/Home.vue';
import Login from './views/Login.vue';
import Host from '@/views/Host.vue';
import Hosts from '@/views/Hosts.vue';
import License from '@/views/License.vue';
import Profile from '@/views/Profile.vue';
import GameServer from '@/views/GameServer.vue';
import GameServers from '@/views/GameServers.vue';
import Help from '@/views/Help.vue';
import NotFound from '@/views/404.vue';

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
            component: Login,
        },
        {
            path: '/license',
            name: 'license',
            meta: {
                title: 'License',
            },
            component: License,
        },
        {
            path: '/help',
            name: 'help',
            meta: {
                title: 'Help',
            },
            component: Help,
        },
        {
            path: '/profile',
            name: 'profile',
            meta: {
                title: 'Profile',
                authRequired: true,
            },
            component: Profile,
        },
        {
            path: '/hosts',
            name: 'hosts',
            meta: {
                title: 'Hosts',
                authRequired: true,
            },
            component: Hosts,
        },
        {
            path: '/host/:id',
            name: 'host',
            meta: {
                title: 'Host',
                authRequired: true,
            },
            component: Host,
        },
        {
            path: '/host/:id/server/:server_id',
            name: 'game_server',
            meta: {
                title: 'Game Server',
                authRequired: true,
            },
            component: GameServer,
        },
        {
            path: '/game-servers',
            name: 'game-servers',
            meta: {
                title: 'Game Servers',
                authRequired: true,
            },
            component: GameServers,
        },
        {
            path: '*',
            name: '404',
            meta: {
                title: '404 - Page Not Found'
            },
            component: NotFound,
        },
    ],
});
