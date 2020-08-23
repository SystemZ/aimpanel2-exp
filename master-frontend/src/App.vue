<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
    <v-app :style="{background: $vuetify.theme.themes['dark'].background}">
        <v-navigation-drawer app v-model="drawer">
            <v-list-item>
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Aimpanel v2
                    </v-list-item-title>
                    <v-list-item-subtitle>
                        Pre-alpha
                    </v-list-item-subtitle>
                </v-list-item-content>
            </v-list-item>
            <v-divider></v-divider>

            <v-list dense nav>
                <v-list-item
                        :key="item.title"
                        :to="item.path"
                        active-class="red--text red--darken-1"
                        link
                        v-for="item in menu"
                        v-if="(item.authRequired && loggedIn) || (!item.authRequired && !loggedIn) || !item.authRequired">
                    <v-list-item-icon>
                        <v-icon>{{item.icon}}</v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                        <v-list-item-title>{{item.title}}</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>

                <v-list-item active-class="red--text red--darken-1" to="login" v-if="!loggedIn">
                    <v-list-item-action>
                        <v-icon>{{mdiLogin}}</v-icon>
                    </v-list-item-action>
                    <v-list-item-content>
                        <v-list-item-title>Login</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>
                <v-list-item v-else @click="$auth.logout()">
                    <v-list-item-action>
                        <v-icon>{{mdiLogout}}</v-icon>
                    </v-list-item-action>
                    <v-list-item-content>
                        <v-list-item-title>Logout</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>
            </v-list>
        </v-navigation-drawer>

        <v-app-bar app color="red darken-2">
            <v-toolbar-title>
                {{$route.meta.title}}
            </v-toolbar-title>
            <v-spacer></v-spacer>
            <span v-if="loggedIn">
                <v-menu offset-y>
                    <template v-slot:activator="{ on }">
                        <v-btn v-on="on">
                            {{$auth.username}}
                            <v-icon>
                                {{mdiChevronDown}}
                            </v-icon>
                        </v-btn>
                    </template>
                    <v-list>
                        <v-list-item to="profile">
                            <v-list-item-title>Profile</v-list-item-title>
                        </v-list-item>
                        <v-list-item @click="$auth.logout()">
                            <v-list-item-title>Log out</v-list-item-title>
                        </v-list-item>
                    </v-list>
                </v-menu>
            </span>
        </v-app-bar>

        <v-content>
            <router-view/>
        </v-content>
    </v-app>
</template>

<script lang="ts">
    import Vue from 'vue';
    import {mdiAccount, mdiExitRun, mdiHome, mdiLogin} from '@mdi/js';

    export default Vue.extend({
        name: 'App',
        computed: {
            loggedIn() {
                return this.$auth.logged;
            },
        },
        created() {
            this.$root.$on('sessionExpired', this.$auth.logout);
        },
        destroyed() {
            this.$root.$off('sessionExpired', this.$auth.logout);
        },
        data: () => ({
            drawer: null,
            menu: [
                {
                    title: 'Home',
                    icon: mdiHome,
                    path: '/',
                    authRequired: false
                },
                /*
                {
                    title: 'License',
                    icon: '',
                    path: '/license',
                    authRequired: false
                },
                {
                    title: 'Help',
                    icon: '',
                    path: '/help',
                    authRequired: false
                }
                */
            ],
            // icons
            mdiLogin: mdiLogin,
            mdiHome: mdiHome,
            mdiAccount: mdiAccount,
            mdiLogout: mdiExitRun,
        }),
    });
</script>
