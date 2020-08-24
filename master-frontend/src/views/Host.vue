<template>
  <v-container grid-list-md>
    <v-row class="mb-6">
      <v-col>
        <v-card>
          <v-card-text>
            <v-btn @click="update()" class="ma-2" color="green" dark>Update</v-btn>
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
                <v-list-item-title>{{ host.platform }} {{ host.platform_version }}</v-list-item-title>
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
            <v-btn @click="remove()" color="red darken-2 accent-4" text>Remove host</v-btn>
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
import {Host, Metric} from '@/types/api';
import Chart from '@/components/Chart.vue';

export default Vue.extend({
  name: 'host',
  components: {
    Chart,
  },
  data: () => ({
    host: {} as Host,
    metric: {} as Metric,
    removeSnackbar: false,
  }),
  mounted(): void {
    this.$http.get('/v1/host/' + this.$route.params.id).then((res) => {
      this.host = res.data.host;
    }).catch(e => {
      this.$auth.checkResponse(e.response.status);
    });
    //this.getGameServers();
  },
  methods: {
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
