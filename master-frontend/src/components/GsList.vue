<template>
  <v-data-table
    :headers="headers"
    :items="gameServers"
    @click:row="goToGameServer"
    class="elevation-1"
    hide-default-footer
    :loading="gsListLoading"
  >
    <template v-slot:item.state="{ item }">
      <template v-if="item.state === 1">
        <v-icon class="green--text" small>{{ mdiCheckboxBlankCircle }}</v-icon>
        -
      </template>
      <template v-else>
        <v-icon class="red--text" small>{{ mdiCheckboxBlankCircle }}</v-icon>
        -
      </template>
    </template>

    <template v-slot:item.host="{ item }">
      <template v-text="getHostName(item.host_id)"></template>
    </template>

    <template v-slot:item.game="{ item }">
      <template v-text="getGameName(item.game_id)"></template>
    </template>
  </v-data-table>
</template>

<script lang="ts">
import Vue from 'vue';
import Component from 'vue-class-component';
import {mdiCheckboxBlankCircle} from '@mdi/js';
import {Game, GameServer, Host} from '@/types/api';

@Component
export default class GsList extends Vue {
  gsListLoading = true;
  gameServers = [] as GameServer[];
  hosts = [] as Host[];
  games = [] as Game[];
  headers = [
    {
      text: 'Name',
      align: 'left',
      sortable: true,
      value: 'name'
    },
    {
      text: 'Host',
      align: 'right',
      value: 'host'
    },
    {
      text: 'Game',
      align: 'right',
      value: 'game'
    },
    {
      text: 'State',
      align: 'right',
      value: 'state'
    }
  ];
  refreshInterval = 0;
  mdiCheckboxBlankCircle = mdiCheckboxBlankCircle;


  mounted() {
    this.getGameServers();
    this.refreshInterval = setInterval(() => {
      this.getGameServers();
    }, 10 * 1000);
  }

  beforeDestroy() {
    clearInterval(this.refreshInterval);
  }

  getHostName(hostId: string): string {
    if (this.hosts && this.hosts.length > 0) {
      let host = this.hosts.find(x => {
        const {id} = x;
        return id === hostId;
      });
      if (host) {
        return host.name;
      }
    }
    return '';
  }

  getGameName(gameId: number) {
    if (this.games && this.games.length > 0) {
      let game = this.games.find(x => {
        const {id} = x;
        return id === gameId;
      });
      if (game) {
        return game.name;
      }
    }
    return '';
  }

  goToGameServer(row: GameServer) {
    this.$router.push('/host/' + row.host_id + '/server/' + row.id);
  }

  getGameServers() {
    this.gsListLoading = true;
    return this.$http.get('/v1/host/my/server').then(res => {
      this.gameServers = res.data.game_servers;
      this.gsListLoading = false;
    }).catch(e => {
      this.$auth.checkResponse(e.response.status);
    });
  }

  getHosts() {
    return this.$http.get('/v1/host').then(res => {
      this.hosts = res.data.hosts;
    }).catch(e => {
      this.$auth.checkResponse(e.response.status);
    });
  }
}
</script>
