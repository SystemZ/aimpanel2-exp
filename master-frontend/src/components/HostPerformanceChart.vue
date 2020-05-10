<template>
    <div>
        <apexchart height="300px" :options="chartOptions" :series="series"></apexchart>
    </div>
</template>
<script lang="ts">
    import Vue from 'vue';

    export default Vue.extend({
        name: 'host-performance-chart',
        props: {
            metrics: {
                // FIXME WTF is this ? how to even type this?
                // type: Array<Metric>, // ? ? ?
                required: true,
            },
        },
        data: () => ({
            chartOptions: {
                theme: {
                    mode: 'dark',
                },
                chart: {
                    id: 'vuechart-example',
                    type: 'area',
                    stacked: false,
                    height: 350,
                    zoom: {
                        type: 'x',
                        enabled: true,
                        autoScaleYaxis: true
                    },
                    toolbar: {
                        autoSelected: 'zoom'
                    },
                },
                stroke: {
                    width: 1,
                    curve: 'smooth',
                },
                grid: {
                    borderColor: '#555',
                    clipMarkers: false,
                    yaxis: {
                        lines: {
                            show: true
                        }
                    }
                },
                dataLabels: {
                    enabled: false
                },
                fill: {
                    gradient: {
                        enabled: true,
                        opacityFrom: 0.55,
                        opacityTo: 0
                    }
                },
                // markers: {
                //     size: 5,
                //     colors: ['#000524'],
                //     strokeColor: '#00BAEC',
                //     strokeWidth: 3
                // },
                title: {
                    text: 'Host performance',
                    align: 'left'
                },
                yaxis: {
                    min: 0,
                    tickAmount: 4
                },
                /*
                yaxis: {
                    labels: {
                        // @ts-ignore
                        // formatter: function(val) {
                        //     return (val / 1000000).toFixed(0);
                        // },
                    },
                    title: {
                        text: 'Price'
                    },
                },
                 */
                xaxis: {
                    type: 'datetime',
                },
                tooltip: {
                    theme: 'dark',
                    // shared: true,
                    x: {
                        format: 'HH:mm:ss dd MMM yyyy',
                    },
                    // y: {
                    // @ts-ignore
                    // formatter: function(val) {
                    //     return (val / 1000000).toFixed(0);
                    // }
                    // }
                }
            },
            series: [],
        }),
        methods: {
            getTsFromId(str: String): Number {
                return parseInt(str.substring(0, 8), 16) * 1000;
            },
            updateChart(): void {
                // @ts-ignore
                //let steal = [];
                let user = [];
                // @ts-ignore
                for (let i = 0; i < this.metrics.length; i++) {
                    // @ts-ignore
                    //let ts = this.getTsFromId(this.metrics[i].id);
                    // @ts-ignore
                    //steal.push([ts, this.metrics[i].steal]);
                    // @ts-ignore
                    user.push([this.metrics[i].t * 1000, this.metrics[i].v]);
                }
                this.series = [
                    // {
                    //     // @ts-ignore
                    //     name: 'Steal',
                    //     // @ts-ignore
                    //     data: steal
                    // },
                    {
                        // @ts-ignore
                        name: 'User',
                        // @ts-ignore
                        data: user
                    }
                ];
            }
        },
        mounted() {
            this.updateChart();
        }
    });
</script>