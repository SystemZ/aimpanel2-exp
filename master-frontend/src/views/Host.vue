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
    </v-row>
    <v-row>
      <v-col>
        <!-- make loader bar green or something -->
        <v-card :loading="allMetrics.length < 1">
          <v-card-title>Charts</v-card-title>
          <v-card-text v-if="noChartsYet">
            <p>
              No data to show on charts :(
            </p>
          </v-card-text>
          <v-card-text v-else>
            <v-row>
              <v-col cols="12" xl="1" md="3" xs="12">
                <v-text-field
                  v-model="metricIntervalS"
                  label="Metric interval in sec"
                ></v-text-field>
              </v-col>
              <v-col>
                <v-btn @click.native="getChart" color="green">Update</v-btn>
              </v-col>
            </v-row>
            <host-performance-chart v-if="allMetrics.length > 0" :metrics="allMetrics"/>
          </v-card-text>
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
import HostPerformanceChart from '@/components/HostPerformanceChart.vue';

export default Vue.extend({
  name: 'host',
  components: {
    HostPerformanceChart,
  },
  data: () => ({
    host: {} as Host,
    metric: {} as Metric,
    metricIntervalS: 3600,
    allMetrics: {} as Array<Metric>,
    removeSnackbar: false,
    noMetricsYet: false,
    noChartsYet: false,
    // gameServers: [],
  }),
  mounted(): void {
    this.$http.get('/v1/host/' + this.$route.params.id).then((res) => {
      this.host = res.data.host;
    }).catch(e => {
      this.$auth.checkResponse(e.response.status);
    });
    this.getChart();
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
    getChart(): void {
      this.allMetrics = [];
      this.$http.get('/v1/host/' + this.$route.params.id + '/metric?interval=' + this.metricIntervalS).then((res) => {
        if (res.data.metrics.length < 1) {
          // no data, skip assigning
          this.noMetricsYet = true;
          this.noChartsYet = true;
          return;
        }

        // modify records for chart.js purposes (spanGaps)
        let intervalS = this.metricIntervalS;
        let newMetrics = [];
        // first record is always OK and it's our guideline
        newMetrics[0] = res.data.metrics[0];
        for (let i = 1; i < res.data.metrics.length; i++) {
          let dateAInt = (res.data.metrics[i - 1].t * 1000) + (intervalS * 1000);
          // let dateAStr = new Date(dateAInt);
          let dateBInt = res.data.metrics[i].t * 1000;
          // let dateBStr = new Date(dateBInt);
          let diff = dateBInt - dateAInt;
          // console.log(dateAInt + ' vs ' + dateBInt);
          // console.log(dateAStr + ' vs ' + dateBStr);
          if (dateAInt != dateBInt) {
            // console.log('wrong! diff:');
            // console.log(diff);
            let nullsToAdd = diff / (intervalS * 1000);
            // console.log(nullsToAdd);
            // take 1 from nulls to add to account adding real data as last point
            for (let n = 0; n < nullsToAdd - 1; n++) {
              // console.log(n);
              let emptyDateInt = dateAInt + (intervalS * 1000 * (n + 1));
              // console.log(emptyDateInt);
              // console.log(new Date(emptyDateInt));
              newMetrics.push({
                't': emptyDateInt / 1000,
                'min': NaN,
                'avg': NaN,
                'max': NaN,
              });
            }
          }
          // else {
          //   console.log('ok');
          // }
          // needs NaN or not, we need to add this value
          newMetrics.push(res.data.metrics[i]);

        }
        // use our data from backend enriched with NaNs for chart.js
        this.allMetrics = newMetrics;

        // FIXME pie charts are empty
        this.metric = res.data.metrics[0];

        this.metric.disk_free = +(this.metric.disk_free as number / 1024).toFixed(0);
        this.metric.disk_total = +(this.metric.disk_total / 1024).toFixed(0);
        this.metric.disk_used = this.metric.disk_total - this.metric.disk_free;

        this.metric.ram_total = +(this.metric.ram_total / 1024).toFixed(1);
        this.metric.ram_free = +(this.metric.ram_free / 1024).toFixed(1);
        this.metric.ram_used = +(this.metric.ram_total - this.metric.ram_free).toFixed(1);
      }).catch(e => {
        this.$auth.checkResponse(e.response.status);
      });
    }
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
