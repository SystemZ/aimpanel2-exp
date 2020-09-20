<template>
  <v-card flat>
    <!--
      :loading="isLoading"
      :search-input.sync="search"
            prepend-icon="mdi-database-search"
            return-object

    -->
    <v-row>
      <v-col cols="12" xl="1" md="3" xs="6">
        <v-select
          v-model="metricLastS"
          :items="metricLastSPresets"
          item-text="label"
          item-value="v"
          label="Time"
          @change="getChart"
        ></v-select>
      </v-col>
      <v-col cols="12" xl="1" md="3" xs="6">
        <v-select
          v-model="metricIntervalS"
          :items="metricIntervalPresets"
          item-text="label"
          item-value="v"
          label="Interval"
          @change="getChart"
        ></v-select>
      </v-col>
      <!--
      <v-col cols="12" xl="1" md="3" xs="6">
        <v-text-field
          v-model="metricIntervalS"
          label="Metric interval in sec"
        ></v-text-field>
      </v-col>
      -->
      <v-col cols="12" xl="3" md="3" xs="12">
        <v-autocomplete
          v-model="selectedMetric"
          :items='availableMetrics'
          color="white"
          hide-no-data
          hide-selected
          label="Metric"
          placeholder="Data to track"
          item-text="label"
          item-value="v"
          @change="getChart"
        ></v-autocomplete>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <host-performance-chart
          v-if="allMetrics.length > 0"
          :title="selectedMetricTitle"
          :unit="metricUnit"
          :metrics="allMetrics"
        />
      </v-col>
    </v-row>

  </v-card>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from 'vue-property-decorator';
import HostPerformanceChart from '@/components/HostPerformanceChart.vue';
import { Metric } from '@/types/api';

@Component({
  components: {
    HostPerformanceChart
  }
})
export default class Chart extends Vue {
  @Prop({
    type: String, required: true, default: () => {
      return '';
    }
  })
  hostId !: String;

  serverUrl = '';
  metricLastS = 86400;
  metricLastSPresets = [
    {'label': '30d', 'v': 2592000},
    {'label': '14d', 'v': 1209600},
    {'label': '7d', 'v': 604800},
    {'label': '3d', 'v': 259200},
    {'label': '1d', 'v': 86400},
    /*
    {'label': '???', 'v': 69778},
    {'label': '12h', 'v': 43200},
    {'label': '6h', 'v': 21600},
    {'label': '3h', 'v': 10800},
    {'label': '1h', 'v': 3600},
    {'label': '30m', 'v': 1800},
    {'label': '15m', 'v': 900},
    {'label': '10m', 'v': 600},
    {'label': '5m', 'v': 300},
    */
  ];
  metricIntervalS = 900;
  metricIntervalPresets = [
    {'label': '24h', 'v': 86400},
    {'label': '12h', 'v': 43200},
    {'label': '6h', 'v': 21600},
    {'label': '3h', 'v': 10800},
    {'label': '1h', 'v': 3600},
    {'label': '30m', 'v': 1800},
    {'label': '15m', 'v': 900},
    {'label': '10m', 'v': 600},
    {'label': '5m', 'v': 300},
    {'label': '1m', 'v': 60},
    {'label': '30s', 'v': 30},
    {'label': '15s', 'v': 15},
  ];
  availableMetrics = [
    {'unit': '%', 'label': 'CPU usage', 'v': 'cpu_usage'},
    {'unit': '%', 'label': 'CPU user', 'v': 'cpu_user'},
    {'unit': '%', 'label': 'CPU system', 'v': 'cpu_system'},
    {'unit': '%', 'label': 'CPU idle', 'v': 'cpu_idle'},
    {'unit': '%', 'label': 'CPU nice', 'v': 'cpu_nice'},
    {'unit': '%', 'label': 'CPU guest', 'v': 'cpu_guest'},
    {'unit': '%', 'label': 'CPU guest nice', 'v': 'cpu_guest_nice'},
    {'unit': '%', 'label': 'CPU steal', 'v': 'cpu_steal'},
    {'unit': '%', 'label': 'CPU iowait', 'v': 'cpu_iowait'},
    {'unit': '%', 'label': 'CPU irq', 'v': 'cpu_irq'},
    {'unit': '%', 'label': 'CPU soft irq', 'v': 'cpu_irq_soft'},
    {'unit': 'MB', 'label': 'RAM usage', 'v': 'ram_usage'},
    {'unit': 'MB', 'label': 'RAM available', 'v': 'ram_available'},
    {'unit': 'MB', 'label': 'RAM free', 'v': 'ram_free'},
    {'unit': 'MB', 'label': 'RAM total', 'v': 'ram_total'},
    {'unit': 'MB', 'label': 'RAM buffers', 'v': 'ram_buffers'},
    {'unit': 'MB', 'label': 'RAM cache', 'v': 'ram_cache'},
    {'unit': 'MB', 'label': 'Disk free', 'v': 'disk_free'},
    {'unit': 'MB', 'label': 'Disk used', 'v': 'disk_used'},
    {'unit': 'MB', 'label': 'Disk total', 'v': 'disk_total'}
  ];
  selectedMetric = 'ram_available';
  selectedMetricTitle = '';
  metricUnit = '';
  allMetrics = {} as Array<Metric>;
  noMetricsYet = false;
  noChartsYet = false;

  mounted() {
    // this.serverUrl = '/v1/host/' + this.hostId + '/server/' + this.serverId;
    this.serverUrl = '/v1/host/' + this.hostId;
    this.getChart();
  }

  @Watch('metricIntervalS')
  onMetricIntervalSChanged() {
    console.log('changed!')
  }

  getChart() {
    this.allMetrics = [];
    let metricUrl = '/v1/host/' + this.$route.params.id + '/metric?name=' + this.selectedMetric + '&last=' + this.metricLastS + '&interval=' + this.metricIntervalS;
    this.$http.get(metricUrl).then((res) => {
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
      // select proper tooltip unit for this data
      for (let i = 0; i < this.availableMetrics.length; i++) {
        if (this.availableMetrics[i].v === this.selectedMetric) {
          this.metricUnit = this.availableMetrics[i].unit;
          this.selectedMetricTitle = this.availableMetrics[i].label;
        }
      }

      // use our data from backend enriched with NaNs for chart.js
      this.allMetrics = newMetrics;
    }).catch(e => {
      this.$auth.checkResponse(e.response.status);
    });
  }
}
</script>
