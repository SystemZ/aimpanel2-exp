<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
    <v-app>
        <v-navigation-drawer v-model="drawer" app>
            <v-list-item>
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Aimpanel
                    </v-list-item-title>
                    <v-list-item-subtitle>
                        v2.0.0
                    </v-list-item-subtitle>
                </v-list-item-content>
            </v-list-item>
            <v-divider></v-divider>

            <v-list dense nav>
                <v-list-item
                        v-for="item in menu"
                        :key="item.title"
                        :to="item.path"
                        v-if="(item.authRequired && loggedIn) || (!item.authRequired && !loggedIn) || !item.authRequired"
                        link
                        active-class="red--text red--darken-1">
                    <v-list-item-icon>
                        <v-icon>{{item.icon}}</v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                        <v-list-item-title>{{item.title}}</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>

                <v-list-item v-if="!loggedIn" to="login" active-class="red--text red--darken-1">
                    <v-list-item-action>
                        <v-icon>fa-sign-in</v-icon>
                    </v-list-item-action>
                    <v-list-item-content>
                        <v-list-item-title>Login</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>
            </v-list>
        </v-navigation-drawer>

        <v-app-bar color="red darken-1" dark app>
            <v-toolbar-title>
                {{$route.meta.title}}
            </v-toolbar-title>
            <v-spacer></v-spacer>
            <span v-if="loggedIn">
                <v-menu offset-y>
                    <template v-slot:activator="{ on }">
                        <v-btn v-on="on">
                            {{$auth.username}} <v-icon size="12">fa-chevron-down fa-small</v-icon>
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
    import Vue from "vue";

    export default Vue.extend({
        name: "App",
        computed: {
            loggedIn() {
                return this.$auth.logged;
            },
        },
        data: () => ({
            drawer: null,
            menu: [
                {
                    title: "Home",
                    icon: "fa-home",
                    path: "/",
                    authRequired: false
                },
                {
                    title: "Hosts",
                    icon: "fa-server",
                    path: "/hosts",
                    authRequired: true
                },
                {
                    title: "Game servers",
                    icon: "fa-gamepad",
                    path: "/game-servers",
                    authRequired: true
                },
                {
                    title: "License",
                    icon: "fa-certificate",
                    path: "/license",
                    authRequired: false
                },
                {
                    title: "Help",
                    icon: "fa-question",
                    path: "/help",
                    authRequired: false
                }
            ]
        }),
    });
</script>
