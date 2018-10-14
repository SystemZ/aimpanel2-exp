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
                <v-list-tile v-for="item in menu"
                             :to="item.path"
                             active-class="red--text red--darken-1">
                    <v-list-tile-action>
                        <v-icon>{{item.icon}}</v-icon>
                    </v-list-tile-action>
                    <v-list-tile-content>
                        <v-list-tile-title>{{item.title}}</v-list-tile-title>
                    </v-list-tile-content>
                </v-list-tile>
            </v-list>
        </v-navigation-drawer>
        <v-toolbar color="red darken-1" dark fixed app>
            <v-toolbar-side-icon @click.stop="drawer = !drawer"></v-toolbar-side-icon>
            <v-toolbar-title>{{$route.meta.title}}</v-toolbar-title>
        </v-toolbar>
        <v-content>
            <router-view/>
        </v-content>
    </v-app>
</template>

<script>
  export default {
    name: 'Aimpanel',
    data: () => ({
      drawer: null,
      menu: [
        {
          title: 'Home',
          icon: 'fa-home',
          path: '/'
        },
        {
          title: 'Hosts',
          icon: 'fa-server',
          path: '/hosts'
        },
        {
          title: 'Game servers',
          icon: 'fa-gamepad',
          path: '/game-servers'
        },
        {
          title: 'License',
          icon: 'fa-certificate',
          path: '/license'
        },
        {
          title: 'Help',
          icon: 'fa-question',
          path: '/help'
        }
      ]
    }),
    mounted() {
      console.log(this.$route)
    }
  }
</script>
