<template>
  <v-card>
    <v-card-title>Settings</v-card-title>

    <v-form>
      <v-container>
        <v-row>
          <v-col
              cols="12"
          >
            <v-text-field
                v-model="gsInfo.custom_cmd_start"
                label="custom_cmd_start"
            ></v-text-field>
          </v-col>
        </v-row>
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
      custom_cmd_start: ''
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
        console.log(res.data);
      }).catch(e => {
        this.$auth.checkResponse(e.response.status);
      });
    },
    saveSettings() {
      // TODO disable text field and show loading during save
      let data = {
        'custom_cmd_start': this.gsInfo.custom_cmd_start
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
