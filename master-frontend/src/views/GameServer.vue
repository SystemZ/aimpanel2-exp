<template>
  <v-container>
    <v-dialog v-model="newGSNameDialog" persistent max-width="600px">
      <v-card>
        <v-card-title>
          <span class="headline">Edit Game Server name</span>
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  v-model="game_server.name"
                  label="Current name"
                  disabled
                ></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-text-field
                  v-model="newGSName"
                  label="New name*"
                  required></v-text-field>
              </v-col>
            </v-row>
          </v-container>
          <small>*required field</small>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="red darken-1" text @click="newGSNameDialog = false">Close</v-btn>
          <v-btn color="green darken-1" @click="saveGSName">Save</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-card>
      <v-tabs
        v-model="tab"
        background-color="transparent"
        color="basil"
        grow
        :vertical="$vuetify.breakpoint.xsOnly"
      >
        <v-tab left>
          State & control
        </v-tab>
        <v-tab>
          File manager
        </v-tab>
        <v-tab>
          Scheduler
        </v-tab>
        <v-tab>
          Settings
        </v-tab>
      </v-tabs>
      <v-tabs-items v-model="tab">
        <v-tab-item>
          <v-container>
            <v-row class="mb-6">
              <v-col cols="3" md="3" sm="12" xs="12">
                <v-card>
                  <v-card-title>Actions</v-card-title>
                  <v-card-text>
                    <v-btn @click="start()" class="ma-2" color="green" dark>Start</v-btn>
                    <v-btn @click="shutdown()" class="ma-2" color="yellow darken-2" dark>Shutdown</v-btn>
                    <v-btn @click="stop()" class="ma-2" color="red" dark>Stop</v-btn>
                    <v-btn @click="restart()" class="ma-2" color="orange darken-2" dark>Restart</v-btn>
                    <v-btn @click="install()" class="ma-2" color="blue" dark>Install</v-btn>
                  </v-card-text>
                </v-card>
                <v-card class="mt-5">
                  <v-card-title>Details</v-card-title>
                  <v-card-text>
                    <v-list-item two-line>
                      <v-list-item-content>
                        <v-list-item-title>{{ game_server.name }}</v-list-item-title>
                        <v-list-item-subtitle>Name</v-list-item-subtitle>
                      </v-list-item-content>
                      <v-list-item-action>
                        <v-btn
                          @click="newGSNameDialog = true"
                          small
                          color="blue"
                          class="my-2"
                        >
                          <v-icon size="20">{{ mdiPencil }}</v-icon>
                        </v-btn>
                      </v-list-item-action>
                    </v-list-item>

                    <v-list-item two-line>
                      <v-list-item-content>
                        <v-list-item-title>{{ game_server.state == 1 ? 'Active' : 'Locked' }}
                        </v-list-item-title>
                        <v-list-item-subtitle>Status</v-list-item-subtitle>
                      </v-list-item-content>
                    </v-list-item>

                    <v-list-item two-line>
                      <v-list-item-content>
                        <v-list-item-title>{{
                            game_server.metric_frequency
                          }}s
                        </v-list-item-title>
                        <v-list-item-subtitle>Metric frequency</v-list-item-subtitle>
                      </v-list-item-content>
                    </v-list-item>

                    <v-list-item two-line>
                      <v-list-item-content>
                        <v-list-item-title>{{ game_server.stop_timeout }}s</v-list-item-title>
                        <v-list-item-subtitle>Stop timeout</v-list-item-subtitle>
                      </v-list-item-content>
                    </v-list-item>
                  </v-card-text>
                  <v-card-actions>
                    <v-btn @click="remove()" color="red darken-3">
                      <v-icon class="mr-2">{{ mdiTrashCan }}</v-icon>
                      Remove game server
                    </v-btn>
                  </v-card-actions>
                </v-card>
              </v-col>
              <v-col cols="9" md="9" sm="12" xs="12">
                <gs-console :host-id="hostId" :server-id="serverId"/>
              </v-col>
            </v-row>
          </v-container>
        </v-tab-item>
        <v-tab-item>
          <!-- v-if="files.selected" -->
          <v-row class="mb-6">
            <v-col cols="12" md="12" sm="12" xs="12">
              <gs-file-manager :host-id="hostId" :server-id="serverId"/>
            </v-col>
          </v-row>
        </v-tab-item>
        <v-tab-item>
          <gs-scheduler :host-id="hostId" :server-id="serverId"/>
        </v-tab-item>
        <v-tab-item>
          <gs-settings :host-id="hostId" :server-id="serverId"/>
        </v-tab-item>
      </v-tabs-items>
    </v-card>

    <v-snackbar
      v-model="installSnackbar"
    >
      Installing game server...
      <v-btn
        @click="installSnackbar = false"
        color="red"
        text
      >
        Close
      </v-btn>
    </v-snackbar>
    <v-snackbar
      v-model="removeSnackbar"
    >
      Removing game server...
      <v-btn
        @click="removeSnackbar = false"
        color="red"
        text
      >
        Close
      </v-btn>
    </v-snackbar>

  </v-container>
</template>

<script lang="ts">
import Vue from 'vue';
import GsConsole from '@/components/GsConsole.vue';
import GsFileManager from '@/components/GsFileManager.vue';
import GsScheduler from '@/components/GsScheduler.vue';
import GsSettings from '@/components/GsSettings.vue';
import {mdiTrashCan} from '@mdi/js';
import {mdiPencil} from '@mdi/js';

export default Vue.extend({
  name: 'game_server',
  components: {
    GsConsole,
    GsFileManager,
    GsScheduler,
    GsSettings,
  },
  data: () => ({
    tab: 0,
    game_server: {},
    message: '',
    timer: '',
    installSnackbar: false,
    removeSnackbar: false,
    newGSNameDialog: false,
    newGSName: '',
    //icons
    mdiTrashCan: mdiTrashCan,
    mdiPencil: mdiPencil,
  }),
  computed: {
    serverId() {
      return this.$route.params.server_id;
    },
    hostId() {
      return this.$route.params.id;
    },
    serverUrl() {
      return '/v1/host/' + this.$route.params.id + '/server/' + this.$route.params.server_id;
    }
  },
  mounted() {
    this.gsInfo()
  },
  methods: {
    gsInfo() {
      this.$http.get(this.serverUrl).then(res => {
        this.game_server = res.data.game_server;
      }).catch(e => {
        this.$auth.checkResponse(e.response.status);
      });
    },
    start() {
      this.$http.put(this.serverUrl + '/start').then(res => {
        //console.log(res);
      }).catch(e => {
        this.$auth.checkResponse(e.response.status);
      });
    },
    stop() {
      this.$http.put(this.serverUrl + '/stop', {
        type: 1
      }).then(res => {
        //console.log(res);
      }).catch(e => {
        this.$auth.checkResponse(e.response.status);
      });
    },
    restart() {
      this.$http.put(this.serverUrl + '/restart', {
        type: 1
      }).then(res => {
        //console.log(res);
      }).catch(e => {
        this.$auth.checkResponse(e.response.status);
      });
    },
    install() {
      this.$http.put(this.serverUrl + '/install').then(res => {
        this.installSnackbar = true;
      }).catch(e => {
        this.$auth.checkResponse(e.response.status);
      });
    },
    remove() {
      this.$http.delete(this.serverUrl).then(res => {
        this.removeSnackbar = true;
        console.log(res);
        this.$router.push('/');
      }).catch(e => {
        this.$auth.checkResponse(e.response.status);
      });
    },
    shutdown() {
      this.$http.put(this.serverUrl + '/shutdown').then(res => {
        //console.log(res);
      }).catch(e => {
        this.$auth.checkResponse(e.response.status);
      });
    },
    saveGSName() {
      this.$http.put(this.serverUrl, {
        name: this.newGSName
      }).then(res => {
        this.gsInfo()
        this.newGSName = ''
        this.newGSNameDialog = false
      }).catch(e => {
        this.$auth.checkResponse(e.response.status);
      })

    },
  },
});
</script>
