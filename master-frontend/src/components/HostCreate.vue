<template>
  <v-dialog max-width="800px" v-model="createHost.dialog">
    <template v-slot:activator="{ on }">
      <v-btn color="info" v-on="on" class="mt-2 mb-2">
        <v-icon class="mr-2">{{ mdiPlus }}</v-icon>
        Create new host
      </v-btn>
    </template>
    <v-stepper v-model="createHost.step">
      <v-stepper-header>
        <v-stepper-step :complete="createHost.step > 1" step="1">Details
        </v-stepper-step>
        <v-divider>
        </v-divider>
        <v-stepper-step :complete="createHost.step > 2" step="2">Setup host
        </v-stepper-step>
      </v-stepper-header>
      <v-stepper-items>
        <v-stepper-content step="1">
          <v-container grild-list-md>
            <v-layout wrap>
              <v-flex xs12>
                <v-text-field label="Name" required
                              v-model="createHost.host.name">
                </v-text-field>
              </v-flex>
              <v-flex xs12>
                <v-text-field label="IP" required
                              v-model="createHost.host.ip">
                </v-text-field>
              </v-flex>
            </v-layout>
          </v-container>

          <v-btn @click="addHost()"
                 color="primary">
            Next
          </v-btn>

          <v-btn @click="createHostCancel()" text>Cancel</v-btn>
        </v-stepper-content>
        <v-stepper-content step="2">
          <v-container grid-list-md>
            <p>
              Host was successfully added.<br>
              Please run command below on your host to finish installation
            </p>
            <kbd>wget https://{{ apiHostname }}/i/{{ createHost.token }} -O- |
              bash -</kbd>
          </v-container>
          <v-btn @click="finish()" text>Close</v-btn>
        </v-stepper-content>
      </v-stepper-items>
    </v-stepper>
  </v-dialog>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { mdiPlus } from '@mdi/js';

@Component
export default class HostCreate extends Vue {
  createHost = {
    dialog: false,
    step: 0,
    host: {
      name: '',
      ip: '',
    },
    token: ''
  };

  //icons
  mdiPlus = mdiPlus;

  //computed
  get apiHostname() {
    // FIXME get this from backend
    // console.log(window.location.port)
    // let port = ""
    // if (window.location.port !== 443) {
    //   port = location.port
    // }
    return window.location.hostname;
  }

  createHostCancel() {
    this.createHost.dialog = false;
    this.createHost.step = 1;
  }

  addHost() {
    this.$http.post('/v1/host', this.createHost.host).then(res => {
      this.createHost.token = res.data.token;
      this.createHost.step = 2;
    }).catch(e => {
      this.$auth.checkResponse(e.response.status);
    });
  }

  finish() {
    this.createHost = {
      dialog: false,
      step: 0,
      host: {
        name: '',
        ip: '',
      },
      token: ''
    };
    // TODO send event to refresh host list after adding new
  }
}
</script>
