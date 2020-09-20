<template>
  <v-card>
    <v-toolbar dense flat>
      <v-toolbar-title>
        Console
      </v-toolbar-title>
      <v-spacer/>
      <small>Log limit</small><input v-model="logLimit" type="text"/>
      <v-checkbox v-model="shouldAutoscroll"
                  :label="`Autoscroll: ${shouldAutoscroll.toString()}`"></v-checkbox>
    </v-toolbar>
    <v-card id="log-lines" class="pa-5" height="50vh"
            style="overflow-y: scroll">
      <template v-for="item in logs">{{ item }}<br/></template>
    </v-card>
    <v-card-actions>
      <v-text-field full-width
                    hide-details
                    label="Type some message here"
                    v-model="message"
                    :loading="commandSending"
                    :disabled="commandSending"
                    v-on:keyup.enter="sendMessage()">
        <template v-slot:append-outer>
          <v-btn color="primary"
                 :loading="commandSending"
                 :disabled="commandSending"
                 @click="sendMessage()"
          >
            <v-icon class="mr-2">
              {{ mdiCogs }}
            </v-icon>
            Execute
          </v-btn>
        </template>
      </v-text-field>
    </v-card-actions>
  </v-card>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { mdiCogs } from '@mdi/js';

@Component
export default class GsConsole extends Vue {
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

  game_server = {};
  logs = [] as string[];
  shouldAutoscroll = true;
  logLimit = 100;
  message = '';
  serverUrl = '';
  stream = '' as any;
  commandSending = false;
  //icons
  mdiCogs = mdiCogs;

  mounted() {
    this.serverUrl = '/v1/host/' + this.hostId + '/server/' + this.serverId;

    this.$http.get(this.serverUrl).then(res => {
      this.game_server = res.data.game_server;
    }).catch(e => {
      this.$auth.checkResponse(e.response.status);
    });

    if (this.stream === '' || this.stream === undefined) {
      this.setupStream();
    }
  }

  beforeDestroy() {
    this.stream.close();
  }

  scrollToBottom() {
    if (!this.shouldAutoscroll) {
      return;
    }
    setTimeout(() => {
      const log: HTMLElement | null = document.getElementById('log-lines');
      if (log) {
        log.scrollTo(0, log.scrollHeight);
      }
    }, 10);
  }

  sendMessage() {
    if (this.message.length < 1) {
      return;
    }
    this.commandSending = true;
    this.$http.put(this.serverUrl + '/command', {
      command: this.message
    }).then(res => {
      this.message = '';
      this.commandSending = false;
    }).catch(e => {
      this.$auth.checkResponse(e.response.status);
    });
  }

  setupStream() {
    this.stream = new this.$eventSource(this.$apiUrl + this.serverUrl + '/console', {
      headers: {
        Authorization: this.$auth.getAuthorizationHeader()
      }
    });

    this.stream.onerror = (event: any) => {
      console.error(event);
    };

    this.stream.addEventListener('message', (event: any) => {
      if (event.data === 'heartbeat') {
        return;
      }

      let data = atob(event.data);
      // clear old messages to prevent browser clogging
      while (this.logs.length > this.logLimit - 1) {
        this.logs.shift();
      }
      this.logs.push(data);
      this.scrollToBottom();
    }, false);
  }
}
</script>
