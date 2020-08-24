<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator';
import {Metric} from '@/types/api';
import {Line} from 'vue-chartjs';
import moment from 'moment';

@Component({
  extends: Line
})
export default class HostPerformanceChart extends Vue {
  @Prop({
    type: Array, required: true, default: () => {
      return [];
    }
  })
  metrics !: Array<Metric>;

  renderChart!: (chartData: any, options: any) => void;

  options = {
    responsive: true,
    maintainAspectRatio: false,
    legend: {
      labels: {
        fontColor: 'white'
      }
    },
    scales: {
      yAxes: [{
        gridLines: {
          drawBorder: false,
        },
        scaleLabel: {
          display: true
        },
        ticks: {
          fontColor: 'white'
        }
      }],
      xAxes: [{
        type: 'time',
        distribution: 'linear',
        offset: true,
        ticks: {
          major: {
            enabled: true,
            fontStyle: 'bold',
          },
          autoSkip: true,
          autoSkipPadding: 75,
          maxRotation: 0,
          sampleSize: 100,
          fontColor: 'white'
        },
      }]
    },
    showLines: false,
    tooltips: {
      mode: 'index',
      intersect: false,
    },
    hover: {
      mode: 'nearest',
      intersect: true
    }
  };

  chartData = {
    labels: this.metrics.map(m => moment.unix(m.t).toDate()),
    datasets: [
      {
        label: 'RAM Available',
        backgroundColor: '#43a047',
        borderColor: '#43a047',
        data: this.metrics.map(m => m.avg),
        type: 'line',
        pointRadius: 1,
        fill: false,
        lineTension: 0,
        borderWidth: 2,
      },
      {
        label: 'RAM Available',
        backgroundColor: '#43a047',
        borderColor: '#43a047',
        data: this.metrics.map(m => m.avg),
        type: 'line',
        pointRadius: 1,
        fill: false,
        lineTension: 0,
        borderWidth: 2,
      }
    ]
  };

  mounted() {
    this.renderChart(this.chartData, this.options);
  }

}
</script>
