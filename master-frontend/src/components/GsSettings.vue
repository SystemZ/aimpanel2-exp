<template>
  <v-card>
    <v-dialog v-model="changeGameEngineDialog" persistent max-width="600px">
      <v-card>
        <v-card-title>
          <span class="headline">Change game engine</span>
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-layout wrap>
              <v-flex xs12>
                <v-select
                  :loading="gamesLoading"
                  :disabled="gamesLoading"
                  :items="games"
                  item-text="name"
                  item-value="id"
                  label="Select game"
                  v-model="game.game_id">
                </v-select>
              </v-flex>
              <v-flex xs12>
                <v-select
                  :items="gameVersions"
                  label="Select game version"
                  v-model="game.game_version">
                </v-select>
              </v-flex>
            </v-layout>
          </v-container>
          <small>*required field</small>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="red darken-1" text @click="changeGameEngineDialog = false">Close</v-btn>
          <v-btn color="green darken-1" @click="saveGameEngine">Save</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-card-title>Generic settings</v-card-title>

    <v-form>
      <v-container>
        <v-row>
          <v-col
            cols="12"
          >
            <v-textarea
              v-model="gsInfo.custom_cmd_start"
              color="green"
              filled
              label="Custom CMD"
              auto-grow
              counter
              rows="1"
            ></v-textarea>
          </v-col>
        </v-row>
      </v-container>
    </v-form>
    <v-card-title>Ports</v-card-title>
    <v-card-text>
      Restart GS to apply port forwarding settings
    </v-card-text>

    <v-container>
      <v-row>
        <v-col cols="12" xs="12" md="9" lg="6">
          <div :key="i" v-for="(templateKey,i) in gsInfo.ports">
            <v-row>
              <v-col>
                <v-text-field
                  label="Host IP"
                  color="green"
                  filled
                  v-model="gsInfo.ports[i].host"
                ></v-text-field>
              </v-col>
              <v-col>
                <v-text-field
                  label="Port on host"
                  type="number"
                  color="green"
                  filled
                  v-model="gsInfo.ports[i].port_host"
                ></v-text-field>
              </v-col>
              <v-col>
                <v-text-field
                  label="Port in container"
                  type="number"
                  color="green"
                  filled
                  v-model="gsInfo.ports[i].port_container"
                ></v-text-field>
              </v-col>
              <v-col>
                <v-text-field
                  label="Protocol"
                  color="green"
                  filled
                  v-model="gsInfo.ports[i].protocol"
                ></v-text-field>
              </v-col>
              <v-col>
                <v-btn @click.native="removePort({i : i})">Remove port</v-btn>
              </v-col>
            </v-row>
          </div>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-btn @click.native="addNewPort">Add new port</v-btn>
        </v-col>
      </v-row>
    </v-container>

    <v-card-title>Game engine</v-card-title>
    <v-list-item two-line>
      <v-list-item-content>
        <v-list-item-title>{{ currentGameName }}</v-list-item-title>
        <v-list-item-subtitle>Game</v-list-item-subtitle>
      </v-list-item-content>
      <v-list-item-action>
        <v-btn
          @click="changeGameEngineDialog = true"
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
        <v-list-item-title>{{ gsInfo.game_version }}</v-list-item-title>
        <v-list-item-subtitle>Version</v-list-item-subtitle>
      </v-list-item-content>
    </v-list-item>

    <v-form>
      <v-container>

        <v-row>
          <v-col>
            <v-btn @click="saveSettings()" color="green">Save</v-btn>
          </v-col>
        </v-row>

      </v-container>
    </v-form>


  </v-card>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from 'vue-property-decorator';
import { mdiArrowLeft, mdiPencil } from '@mdi/js';
import { Game } from '@/types/api';

interface Number {
  i: any;
}

@Component
export default class GsSettings extends Vue {
  @Prop({
    type: String, required: true, default: () => {
      return '';
    }
  })
  serverId !: String;

  @Prop({
    type: String, required: true, default: () => {
      return '';
    }
  })
  hostId !: String;

  serverUrl = '';
  gsInfo = {
    game_id: 0,
    game_version: '',
    custom_cmd_start: '',
    ports: [
      {
        'host': '0.0.0.0',
        'protocol': 'tcp',
        'port_container': 7777,
        'port_host': 7777,
      }
    ],
  };
  examplePort = {
    'host': '0.0.0.0',
    'protocol': 'tcp',
    'port_container': 7777,
    'port_host': 7777,
  };

  currentGameName = ''

  game = {
    game_id: 0,
    game_version: 0,
  }
  gamesLoading = true;
  games = [] as Game[];
  gameVersions = [] as string[];
  changeGameEngineDialog = false;
  mdiPencil = mdiPencil;
  //icons
  mdiArrowLeft = mdiArrowLeft

  mounted() {
    this.serverUrl = '/v1/host/' + this.hostId + '/server/' + this.serverId;
    this.getGames();
    this.getSettings();
  }

  getSettings() {
    this.$http.get(this.serverUrl).then(res => {
      this.gsInfo = res.data.game_server;
    }).catch(e => {
      this.$auth.checkResponse(e.response.status);
    });
  }

  addNewPort() {
    this.gsInfo.ports.push(this.examplePort);
  }

  removePort({i}: Number) {
    this.gsInfo.ports.splice(i, 1);
  }

  saveSettings() {
    // TODO disable text field and show loading during save

    let enforcedPorts: never[] | { 'host': string; 'protocol': string; 'port_container': number; 'port_host': number; }[] = [];
    this.gsInfo.ports.forEach(onePortDef => {
      let enforcedPort = onePortDef;
      enforcedPort.port_host = Number(enforcedPort.port_host);
      enforcedPort.port_container = Number(enforcedPort.port_container);
      // @ts-ignore FIXME
      enforcedPorts.push(enforcedPort);
    });

    let data = {
      'custom_cmd_start': this.gsInfo.custom_cmd_start,
      'ports': enforcedPorts,
    };
    this.$http.put(this.serverUrl, data).then(res => {
      this.gsInfo = res.data.game_server;
      // TODO toast with saved info
    }).catch(e => {
      this.$auth.checkResponse(e.response.status);
    });
  }

  @Watch('game.game_id')
  onGameIdChanged(value: number) {
    this.gameVersions = this.games.filter((g) => {
      return g.id === value;
    })[0].versions;
  }

  @Watch('gsInfo.game_id')
  onGsInfoGameIdChanged(value: number) {
    this.games.forEach(g => {
      if(g.id == value) {
        this.currentGameName = g.name
      }
    })
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

  saveGameEngine(): void {
    this.$http.put(this.serverUrl, {
      game_id: this.game.game_id,
      game_version: this.game.game_version
    }).then(res => {
      this.getSettings();
      this.game = {
        game_id: 0,
        game_version: 0,
      }
      this.changeGameEngineDialog = false
    }).catch(e => {
      this.$auth.checkResponse(e.response.status);
    })
  }
}
</script>
