<template>
  <v-card>
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
import Vue from 'vue';
import {mdiArrowLeft} from '@mdi/js';

interface Number {
  i: any;
}

export default Vue.extend({
  name: 'gs-settings',
  props: {
    serverId: {
      type: String,
      required: true,
    },
    hostId: {
      type: String,
      required: true,
    }
  },
  data: () => ({
    serverUrl: '',
    gsInfo: {
      custom_cmd_start: '',
      ports: [
        {
          'host': '0.0.0.0',
          'protocol': 'tcp',
          'port_container': 7777,
          'port_host': 7777,
        }
      ],
    },
    examplePort: {
      'host': '0.0.0.0',
      'protocol': 'tcp',
      'port_container': 7777,
      'port_host': 7777,
    },
    //icons
    mdiArrowLeft: mdiArrowLeft,
  }),
  mounted() {
    this.serverUrl = '/v1/host/' + this.hostId + '/server/' + this.serverId;
    this.getSettings();
  },
  methods: {
    getSettings() {
      this.$http.get(this.serverUrl).then(res => {
        this.gsInfo = res.data.game_server;
      }).catch(e => {
        this.$auth.checkResponse(e.response.status);
      });
    },
    addNewPort() {
      this.gsInfo.ports.push(this.examplePort);
    },
    removePort({i}: Number) {
      this.gsInfo.ports.splice(i, 1);
    },
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
    },
  },
});
</script>
