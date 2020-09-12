<script lang="ts">
import {Component, Prop, Vue, Watch} from 'vue-property-decorator';
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
  @Prop({
    type: String, required: false, default: () => {
      return '';
    }
  })
  title !: String;
  @Prop({
    type: String, required: false, default: () => {
      return '';
    }
  })
  unit !: String;

  renderChart!: (chartData: any, options: any) => void;

  options = {
    responsive: true,
    responsiveAnimationDuration: 400,
    maintainAspectRatio: false,
    legend: {
      position: 'bottom',
      labels: {
        fontColor: 'white'
      }
    },
    showLines: true,
    spanGaps: false,
    tooltips: {
      mode: 'index',
      intersect: false,
      callbacks: {
        label: (tooltipItem: { datasetIndex: string | number; yLabel: number; }, data: { datasets: { [x: string]: { label: string; }; }; }) => {
          let label = data.datasets[tooltipItem.datasetIndex].label || '';
          if (label) {
            label += ': ';
          }
          label += Math.round(tooltipItem.yLabel * 100) / 100;
          if (this.unit.length > 0) {
            // use prettySize from file manager to auto select units
            label += ' ' + this.unit;
          }
          return label;
        }
        // labelColor: (tooltipItem, chart) => ({
        // borderColor: 'rgb(255, 0, 0)',
        // backgroundColor: 'rgb(255, 0, 0)'
        // }),
        // labelTextColor: (tooltipItem, chart) => '#543453'
      }
    },
    title: {
      display: true,
      position: 'top',
      fontSize: 15,
      fontFamily: '"Roboto", sans-serif',
      fontStyle: 'bold',
      text: this.title
    },
    // elements: {
    //   line: {
    //     tension: 0
    //   }
    // },
    hover: {
      mode: 'nearest',
      intersect: false
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
    }
  };

  chartData = {
    labels: this.metrics.map(m => moment.unix(m.t).toDate()),
    datasets: [
      {
        label: 'Avg',
        backgroundColor: '#96fd9a',
        borderColor: '#96fd9a',
        data: this.metrics.map(m => Math.round(m.avg)),
        type: 'line',
        pointRadius: 2,
        fill: false,
        lineTension: 0,
        borderWidth: 2,
      },
      {
        label: 'Min',
        backgroundColor: '#2a9afc',
        borderColor: '#2a9afc',
        data: this.metrics.map(m => m.min),
        type: 'line',
        pointRadius: 2,
        fill: false,
        lineTension: 0,
        borderWidth: 2,
      },
      {
        label: 'Max',
        backgroundColor: '#ff6969',
        borderColor: '#ff6969',
        data: this.metrics.map(m => m.max),
        type: 'line',
        pointRadius: 2,
        fill: false,
        lineTension: 0,
        borderWidth: 2,
      }
    ]
  };

  mounted() {
    this.renderChart(this.chartData, this.options);
  }

  @Watch('metrics')
  onPropertyChanged(value: string, oldValue: string) {
    this.renderChart(this.chartData, this.options);
    //this.$data._chart.update();
  }

}
</script>
