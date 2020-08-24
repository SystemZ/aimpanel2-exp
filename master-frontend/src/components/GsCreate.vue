<template>
  <v-dialog max-width="600px" v-model="createGameServer.dialog">
    <template v-slot:activator="{ on }">
      <v-btn color="info" v-on="on" class="mt-2 mb-2">
        <v-icon class="mr-2">{{ mdiPlus }}</v-icon>
        Create new game server
      </v-btn>
    </template>
    <v-stepper v-model="createGameServer.step">
      <v-stepper-header>
        <v-stepper-step :complete="createGameServer.step > 1" step="1">Details</v-stepper-step>
        <v-divider></v-divider>
        <v-stepper-step :complete="createGameServer.step > 2" step="2">Install</v-stepper-step>
      </v-stepper-header>
      <v-stepper-items>
        <v-stepper-content step="1">
          <v-container grild-list-md>
            <v-layout wrap>
              <v-flex xs12>
                <v-text-field label="Name on list" required
                              v-model="createGameServer.game.name"></v-text-field>
              </v-flex>
              <v-flex xs12>
                <v-select
                  :loading="hostsLoading"
                  :disabled="hostsLoading"
                  :items="hosts"
                  item-text="name"
                  item-value="id"
                  label="Host"
                  v-model="createGameServer.selectedHost">
                </v-select>
              </v-flex>
              <v-flex xs12>
                <v-select
                  :loading="gamesLoading"
                  :disabled="gamesLoading"
                  :items="games"
                  item-text="name"
                  item-value="id"
                  label="Select game"
                  v-model="createGameServer.game.game_id">
                </v-select>
              </v-flex>
              <v-flex xs12>
                <v-select
                  :items="createGameServer.versions"
                  label="Select game version"
                  v-model="createGameServer.game.game_version">
                </v-select>
              </v-flex>
            </v-layout>
          </v-container>

          <v-btn @click="addGameServer()"
                 color="primary">
            Next
          </v-btn>

          <v-btn @click="createGameServerCancel()" text>Cancel</v-btn>
        </v-stepper-content>
        <v-stepper-content step="2">
          <v-container grid-list-md>
            <p>Game server was successfully created. Do you want to install it now?</p>
          </v-container>
          <v-btn @click="installGsEngine()" color="info">Yes, install now</v-btn>
          <v-btn @click="gsAddFinish()" text>Close</v-btn>
        </v-stepper-content>
      </v-stepper-items>
    </v-stepper>
  </v-dialog>
</template>

<script lang="ts">
import {Component, Vue, Watch} from 'vue-property-decorator';

import {mdiPlus} from '@mdi/js';
import {Game, GameServer, Host} from '@/types/api';

@Component
export default class GsCreate extends Vue {
  hostsLoading = true;
  hosts = [] as Host[];
  gamesLoading = true;
  games = [] as Game[];
  gameVersions = [] as string[];
  gameServers = [] as GameServer[];
  createGameServer = {
    dialog: false,
    step: 0,
    selectedHost: [],
    versions: [] as string[],
    game: {
      name: '',
      game_id: 0,
      game_version: 0,
    },
    gameId: 0,
  };

  //icons
  mdiPlus = mdiPlus;

  @Watch('createGameServer.game.game_id')
  onGameIdChanged(value: number) {
    this.createGameServer.versions = this.games.filter((g) => {
      return g.id === value;
    })[0].versions;
  }

  mounted() {
    this.getGames();
    this.getHosts();
    this.getGameServers();
  }

  getGames() {
    this.gamesLoading = true;
    return this.$http.get('/v1/game').then(res => {
      this.games = res.data;
      this.gamesLoading = false;
    }).catch(e => {
      this.$auth.checkResponse(e.response.status);
    });
  }

  getHosts() {
    this.hostsLoading = true;
    return this.$http.get('/v1/host').then(res => {
      this.hosts = res.data.hosts;
      this.hostsLoading = false;
    }).catch(e => {
      this.$auth.checkResponse(e.response.status);
    });
  }

  getGameServers() {
    return this.$http.get('/v1/host/my/server').then(res => {
      this.gameServers = res.data.game_servers;
    }).catch(e => {
      this.$auth.checkResponse(e.response.status);
    });
  }

  addGameServer(): void {
    this.$http.post('/v1/host/' + this.createGameServer.selectedHost + '/server',
      this.createGameServer.game).then(res => {
        this.createGameServer.gameId = res.data.id;
        this.createGameServer.step = 2;
      }
    );
  }

  installGsEngine(): void {
    this.$http.put('/v1/host/' + this.createGameServer.selectedHost +
      '/server/' + this.createGameServer.gameId + '/install').then(res => {
      console.log(res);
    });
  }

  createGameServerCancel(): void {
    this.createGameServer.dialog = false;
    this.createGameServer.step = 1;
    this.gsAddFinish();
  }

  gsAddFinish(): void {
    this.createGameServer = {
      dialog: false,
      step: 0,
      selectedHost: [],
      versions: [],
      game: {
        name: '',
        game_id: 0,
        game_version: 0,
      },
      gameId: 0,
    };
  }
}
</script>
