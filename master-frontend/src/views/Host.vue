<template>
    <v-container grid-list-md>
        <v-layout row>
            <v-flex xs6>
                <div class="display-1 grey--text text--darken-1">Info</div>
                <br>
                <v-layout row>
                    <v-flex xs12>
                        <v-card>
                            <v-container>
                                <v-layout row>
                                    <v-flex xs12>
                                        <v-card color="green" class="white--text">
                                            <v-container grid-list-md>
                                                <div class="headline">
                                                    Active
                                                </div>
                                            </v-container>
                                        </v-card>
                                    </v-flex>
                                </v-layout>
                                <v-layout row>
                                    <v-flex xs12>
                                        <span class="title font-weight-light">Distro: </span> <b>Ubuntu 16.04.4 LTS</b>
                                    </v-flex>
                                </v-layout>
                                <v-layout row>
                                    <v-flex xs12>
                                        <span class="title font-weight-light">Kernel: </span> <b>4.4.0-116-generic</b>
                                    </v-flex>
                                </v-layout>
                            </v-container>
                        </v-card>
                    </v-flex>
                </v-layout>
            </v-flex>
            <v-flex xs6>
                <div class="display-1 grey--text text--darken-1">Performance</div>
                <br>
                <v-layout row>
                    <v-flex xs12>
                        <v-card>
                            <v-container>
                                <v-layout>
                                    <v-flex xs3>
                                        <v-progress-circular
                                                :size="120"
                                                :width="15"
                                                :rotate="-90"
                                                :value="10"
                                                color="primary">
                                            10% CPU
                                        </v-progress-circular>
                                        <div class="text-xs-center">CPU</div>
                                    </v-flex>
                                    <v-flex xs3>
                                        <v-progress-circular
                                                :size="120"
                                                :width="15"
                                                :rotate="-90"
                                                :value="90"
                                                color="primary">
                                            3.46 GB
                                        </v-progress-circular>
                                        <div class="text-xs-center">RAM</div>
                                    </v-flex>
                                    <v-flex xs3>
                                        <v-progress-circular
                                                :size="120"
                                                :width="15"
                                                :rotate="-90"
                                                :value="80"
                                                color="green">
                                            80%
                                        </v-progress-circular>
                                        <div class="text-xs-center">Disk</div>
                                    </v-flex>
                                    <v-flex xs3>
                                        <v-progress-circular
                                                :size="120"
                                                :width="15"
                                                :rotate="-90"
                                                :value="24"
                                                color="red">
                                            1.4
                                        </v-progress-circular>
                                        <div class="text-xs-center">1m load</div>
                                    </v-flex>
                                </v-layout>
                            </v-container>
                        </v-card>
                    </v-flex>
                </v-layout>
            </v-flex>
        </v-layout>
        <v-layout>
            <v-flex xs6>
                <br>
                <div class="display-1 grey--text text--darken-1">Managment</div>
                <br>
                <v-layout row>
                    <v-flex xs12>
                        <v-card>
                            <v-container>
                                <v-layout row>
                                    <v-btn color="red" dark large>
                                        <v-icon>fa-fw fa-plug</v-icon>
                                        Remove Host
                                    </v-btn>
                                </v-layout>
                            </v-container>
                        </v-card>
                    </v-flex>
                </v-layout>
            </v-flex>
            <v-flex xs6>
                <br>
                <div class="display-1 grey--text text--darken-1">File Manager</div>
                <br>
                <v-layout row>
                    <v-flex xs12>
                        <v-card>
                            <v-list two-line subheader>
                                <v-subheader inset>Current directory {{ fileManager.current_dir }}</v-subheader>
                                <v-subheader inset>Directories</v-subheader>

                                <v-list-tile
                                        v-for="dir in fileManager.directories"
                                        :key="dir.title"
                                        avatar
                                        @click="">
                                    <v-list-tile-avatar>
                                        <v-icon>{{ dir.icon }}</v-icon>
                                    </v-list-tile-avatar>

                                    <v-list-tile-content>
                                        <v-list-tile-title>{{ dir.title }}</v-list-tile-title>
                                        <v-list-tile-sub-title>{{ dir.last_modification }}</v-list-tile-sub-title>
                                    </v-list-tile-content>

                                    <v-list-tile-action>
                                        <v-btn icon ripple>
                                            <v-icon color="red lighten-1">fa-trash-o</v-icon>
                                        </v-btn>
                                    </v-list-tile-action>
                                </v-list-tile>

                                <v-divider inset></v-divider>

                                <v-subheader inset>Files</v-subheader>

                                <v-list-tile
                                        v-for="file in fileManager.files"
                                        :key="file.title"
                                        avatar
                                        @click=""
                                >
                                    <v-list-tile-avatar>
                                        <v-icon>{{ file.icon }}</v-icon>
                                    </v-list-tile-avatar>

                                    <v-list-tile-content>
                                        <v-list-tile-title>{{ file.title }}</v-list-tile-title>
                                        <v-list-tile-sub-title>{{ file.last_modification }}</v-list-tile-sub-title>
                                    </v-list-tile-content>

                                    <v-list-tile-action>
                                        <v-btn icon ripple>
                                            <v-icon color="red lighten-1">fa-trash-o</v-icon>
                                        </v-btn>
                                    </v-list-tile-action>
                                </v-list-tile>
                            </v-list>
                        </v-card>
                    </v-flex>
                </v-layout>
            </v-flex>
        </v-layout>
    </v-container>
</template>

<script>
  export default {
    name: 'host',
    data: () => ({
      host: {},
      fileManager: {
        current_dir: '/home/test',
        directories: [
          {
            icon: 'fa-folder',
            title: 'plugins',
            last_modification: '04.01.2019 20:44:33'
          },
          {
            icon: 'fa-folder',
            title: 'logs',
            last_modification: '01.01.2019 19:43:00'
          },
          {
            icon: 'fa-folder',
            title: 'world',
            last_modification: '12.12.2018 12:13:41'
          }
        ],
        files: [
          {
            icon: 'fa-file',
            title: 'server.properties',
            last_modification: '14.12.2018 12:13:41'
          },
          {
            icon: 'fa-file',
            title: 'settings.yml',
            last_modification: '03.12.2018 12:13:41'
          }
        ]
      }
    }),
    mounted() {
      this.$http.get('/v1/hosts/' + this.$route.params.id).then(res => {
        this.host = res.data
      }).catch(e => {
        console.error(e)
      });
    }
  }
</script>