<template>
  <v-container grid-list-md>
    <v-dialog v-model="newHostNameDialog" persistent max-width="600px">
      <v-card>
        <v-card-title>
          <span class="headline">Edit Host name</span>
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  v-model="host.name"
                  label="Current name"
                  disabled
                ></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-text-field
                  v-model="newHostName"
                  label="New name*"
                  required></v-text-field>
              </v-col>
            </v-row>
          </v-container>
          <small>*required field</small>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="red darken-1" text @click="newHostNameDialog = false">
            Close
          </v-btn>
          <v-btn color="green darken-1" @click="saveHostName">Save</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-row class="mb-6">
      <v-col>
        <v-card>
          <v-card-text>
            <v-btn @click="update()" class="ma-2" color="green" dark>Update
            </v-btn>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col xs6>
        <v-card>
          <v-card-title>Details</v-card-title>
          <v-card-text>
            <v-list-item two-line>
              <v-list-item-content>
                <v-list-item-title>{{ host.name }}</v-list-item-title>
                <v-list-item-subtitle>Name</v-list-item-subtitle>
              </v-list-item-content>
              <v-list-item-action>
                <v-btn
                  @click="newHostNameDialog = true"
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
                <v-list-item-title>Active</v-list-item-title>
                <v-list-item-subtitle>Status</v-list-item-subtitle>
              </v-list-item-content>
            </v-list-item>
            <v-list-item two-line>
              <v-list-item-content>
                <v-list-item-title>{{ host.os }}</v-list-item-title>
                <v-list-item-subtitle>OS</v-list-item-subtitle>
              </v-list-item-content>
            </v-list-item>

            <v-list-item two-line>
              <v-list-item-content>
                <v-list-item-title>
                  {{ host.platform }} {{host.platform_version }}
                </v-list-item-title>
                <v-list-item-subtitle>Platform</v-list-item-subtitle>
              </v-list-item-content>
            </v-list-item>

            <v-list-item two-line>
              <v-list-item-content>
                <v-list-item-title>{{ host.kernel_version }}</v-list-item-title>
                <v-list-item-subtitle>Kernel</v-list-item-subtitle>
              </v-list-item-content>
            </v-list-item>

            <!--
            FIXME token for debug purposes only, remove after implementing better slave auth
            -->
            <v-list-item two-line>
              <v-list-item-content>
                <v-list-item-title>{{ host.token }}</v-list-item-title>
                <v-list-item-subtitle>Token</v-list-item-subtitle>
              </v-list-item-content>
            </v-list-item>

          </v-card-text>
          <v-card-actions>
            <v-btn @click="remove()" color="red darken-2 accent-4" text>Remove
              host
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>

      <!--
      <v-col xs6>
        <v-card>
          <v-card-title>Performance</v-card-title>
          <v-card-text align="center">
            <v-row v-if="noMetricsYet">
              <v-col>
                <p>No metrics gathered yet</p>
              </v-col>
            </v-row>
            <v-row v-else>
              <v-col xs3>
                <v-progress-circular
                  :rotate="-90"
                  :size="120"
                  :value="metric.cpu_usage"
                  :width="15"
                  color="primary">
                  {{ metric.cpu_usage }}% CPU
                </v-progress-circular>
                <div class="text-xs-center">CPU</div>
              </v-col>
              <v-col xs3>
                <v-progress-circular
                  :rotate="-90"
                  :size="120"
                  :value="(metric.ram_used/metric.ram_total) * 100"
                  :width="15"
                  color="primary">
                  {{ metric.ram_used }}/{{ metric.ram_total }} GB
                </v-progress-circular>
                <div class="text-xs-center">RAM</div>
              </v-col>
              <v-col xs3>
                <v-progress-circular
                  :rotate="-90"
                  :size="120"
                  :value="(metric.disk_used/metric.disk_total) * 100"
                  :width="15"
                  color="green">
                  {{ metric.disk_used }}/{{ metric.disk_total }} GB
                </v-progress-circular>
                <div class="text-xs-center">Disk</div>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
      -->
    </v-row>
    <v-row>
      <v-col>
        <!-- make loader bar green or something -->
        <!--        <v-card :loading="allMetrics.length < 1">-->
        <v-card>
          <v-card-title>Charts</v-card-title>
          <v-card-text>
            <chart :host-id="this.$route.params.id"></chart>
          </v-card-text>
          <!--          <v-card-text v-if="noChartsYet">-->
          <!--            <p>-->
          <!--              No data to show on charts :(-->
          <!--            </p>-->
          <!--          </v-card-text>-->
          <!--           v-else>-->
          <!--           -->
          <!--          </v-card-text>-->
        </v-card>
      </v-col>
    </v-row>
    <v-snackbar
      v-model="removeSnackbar"
    >
      Removing host...
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
import { Host, Metric } from '@/types/api';
import Chart from '@/components/Chart.vue';
import { mdiPencil } from '@mdi/js';

export default Vue.extend({
  name: 'host',
  components: {
    Chart,
  },
  data: () => ({
    host: {} as Host,
    metric: {} as Metric,
    removeSnackbar: false,
    newHostNameDialog: false,
    newHostName: '',
    mdiPencil: mdiPencil,
  }),
  mounted(): void {
    this.hostInfo()
  },
  methods: {
    hostInfo() {
      this.$http.get('/v1/host/' + this.$route.params.id).then((res) => {
        this.host = res.data.host;
      }).catch(e => {
        this.$auth.checkResponse(e.response.status);
      });
    },
    remove(): void {
      this.$http.delete('/v1/host/' + this.$route.params.id).then(res => {
        this.removeSnackbar = true;
        console.log(res);
        this.$router.push('/');
      }).catch(e => {
        this.$auth.checkResponse(e.response.status);
      });
    },
    update(): void {
      this.$http.get('/v1/host/' + this.$route.params.id + '/update').then(res => {
        console.log(res);
      }).catch(e => {
        this.$auth.checkResponse(e.response.status);
      });
    },
    saveHostName() {
      this.$http.put('/v1/host/' + this.$route.params.id, {
        name: this.newHostName
      }).then(res => {
        this.hostInfo()
        this.newHostName = ''
        this.newHostNameDialog = false
      }).catch(e => {
        this.$auth.checkResponse(e.response.status);
      })
    }
    /*
    // FIXME pie charts are empty
        this.metric = res.data.metrics[0];

        this.metric.disk_free = +(this.metric.disk_free as number / 1024).toFixed(0);
        this.metric.disk_total = +(this.metric.disk_total / 1024).toFixed(0);
        this.metric.disk_used = this.metric.disk_total - this.metric.disk_free;

        this.metric.ram_total = +(this.metric.ram_total / 1024).toFixed(1);
        this.metric.ram_free = +(this.metric.ram_free / 1024).toFixed(1);
        this.metric.ram_used = +(this.metric.ram_total - this.metric.ram_free).toFixed(1);
    */
    /*
    getGameServers(): void {
        this.$http.get('/v1/host/' + this.$route.params.id + '/server').then(res => {
            this.gameServers = res.data.game_servers;
        }).catch(e => {
            this.$auth.checkResponse(e.response.status);
        });
    },
    */
  },
});
</script>
