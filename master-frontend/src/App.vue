<template>
    <v-app>
        <v-navigation-drawer
                v-model="drawer"
                app
                dark
                absolute
        >
            <v-toolbar flat>
                <v-list>
                    <v-list-tile>
                        <v-list-tile-title class="title">
                            <p class="text-xs-center font-weight-regular font-weight-thin">Aimpanel</p>
                        </v-list-tile-title>
                    </v-list-tile>
                </v-list>
            </v-toolbar>
            <v-list class="pt-0" dense>
                <v-divider></v-divider>
                <v-list-tile
                        v-for="item in menu"
                        :key="item.title"
                        :to="item.path"
                        v-if="(item.authRequired && loggedIn) || (!item.authRequired && !loggedIn) || !item.authRequired"
                        active-class="red--text red--darken-1">
                    <v-list-tile-action>
                        <v-icon>{{item.icon}}</v-icon>
                    </v-list-tile-action>
                    <v-list-tile-content>
                        <v-list-tile-title>{{item.title}}</v-list-tile-title>
                    </v-list-tile-content>
                </v-list-tile>

                <v-list-tile v-if="!loggedIn" to="login" active-class="red--text red--darken-1">
                    <v-list-tile-action>
                        <v-icon>fa-sign-in</v-icon>
                    </v-list-tile-action>
                    <v-list-tile-content>
                        <v-list-tile-title>Login</v-list-tile-title>
                    </v-list-tile-content>
                </v-list-tile>
            </v-list>
        </v-navigation-drawer>
        <v-toolbar color="red darken-1" dark fixed app>
            <v-toolbar-side-icon @click.stop="drawer = !drawer"></v-toolbar-side-icon>
            <v-toolbar-title>{{$route.meta.title}}</v-toolbar-title>
            <v-spacer></v-spacer>
            <span v-if="loggedIn">
                <v-menu offset-y>
                <v-btn slot="activator">
                    {{$auth.username}}&nbsp; <v-icon size="12">fa-chevron-down fa-small</v-icon>
                </v-btn>
                <v-list>
                    <v-list-tile to="profile">
                        <v-list-tile-title>Profile</v-list-tile-title>
                    </v-list-tile>
                    <v-list-tile @click="$auth.logout()">
                        <v-list-tile-title>Log out</v-list-tile-title>
                    </v-list-tile>
                </v-list>
            </v-menu>
            </span>
        </v-toolbar>
        <v-content>
            <router-view/>
        </v-content>
    </v-app>
</template>

<script>
  export default {
    name: 'Aimpanel',
    computed: {
      loggedIn() {
        return this.$auth.logged;
      }
    },
    data: () => ({
      drawer: null,
      menu: [
        {
          title: 'Home',
          icon: 'fa-home',
          path: '/',
          authRequired: false
        },
        {
          title: 'Hosts',
          icon: 'fa-server',
          path: '/hosts',
          authRequired: true
        },
        {
          title: 'Game servers',
          icon: 'fa-gamepad',
          path: '/game-servers',
          authRequired: true
        },
        {
          title: 'License',
          icon: 'fa-certificate',
          path: '/license',
          authRequired: false
        },
        {
          title: 'Help',
          icon: 'fa-question',
          path: '/help',
          authRequired: false
        }
      ]
    })
  }
</script>
