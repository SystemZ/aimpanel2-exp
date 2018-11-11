<template>
    <v-container grid-list-md>
        <v-layout row>
            <v-flex xs6>
                <div class="display-1 grey--text text--darken-1">Info</div>
                <br>
                <v-layout row>
                    <v-flex xs12>
                        <v-card>
                            <v-container>
                                <v-layout row>
                                    <v-flex xs12>
                                        <v-card color="green" class="white--text">
                                            <v-container grid-list-md>
                                                <div class="headline">
                                                    Active
                                                </div>
                                            </v-container>
                                        </v-card>
                                    </v-flex>
                                </v-layout>
                                <v-layout row>
                                    <v-flex xs12>
                                        <span class="title font-weight-light">Distro: </span> <b>Ubuntu 16.04.4 LTS</b>
                                    </v-flex>
                                </v-layout>
                                <v-layout row>
                                    <v-flex xs12>
                                        <span class="title font-weight-light">Kernel: </span> <b>4.4.0-116-generic</b>
                                    </v-flex>
                                </v-layout>
                            </v-container>
                        </v-card>
                    </v-flex>
                </v-layout>
            </v-flex>
            <v-flex xs6>
                <div class="display-1 grey--text text--darken-1">Performance</div>
                <br>
                <v-layout row>
                    <v-flex xs12>
                        <v-card>
                            <v-container>
                                <v-layout>
                                    <v-flex xs3>
                                        <v-progress-circular
                                                :size="120"
                                                :width="15"
                                                :rotate="-90"
                                                :value="10"
                                                color="primary">
                                            10% CPU
                                        </v-progress-circular>
                                        <div class="text-xs-center">CPU</div>
                                    </v-flex>
                                    <v-flex xs3>
                                        <v-progress-circular
                                                :size="120"
                                                :width="15"
                                                :rotate="-90"
                                                :value="90"
                                                color="primary">
                                            3.46 GB
                                        </v-progress-circular>
                                        <div class="text-xs-center">RAM</div>
                                    </v-flex>
                                    <v-flex xs3>
                                        <v-progress-circular
                                                :size="120"
                                                :width="15"
                                                :rotate="-90"
                                                :value="80"
                                                color="green">
                                            80%
                                        </v-progress-circular>
                                        <div class="text-xs-center">Disk</div>
                                    </v-flex>
                                    <v-flex xs3>
                                        <v-progress-circular
                                                :size="120"
                                                :width="15"
                                                :rotate="-90"
                                                :value="24"
                                                color="red">
                                            1.4
                                        </v-progress-circular>
                                        <div class="text-xs-center">1m load</div>
                                    </v-flex>
                                </v-layout>
                            </v-container>
                        </v-card>
                    </v-flex>
                </v-layout>
            </v-flex>
        </v-layout>
        <v-layout>
            <v-flex xs6>
                <br>
                <div class="display-1 grey--text text--darken-1">Managment</div>
                <br>
                <v-layout row>
                    <v-flex xs12>
                        <v-card>
                            <v-container>
                                <v-layout row>
                                    <v-btn color="red" dark large>
                                        <v-icon>fa-fw fa-plug</v-icon>
                                        Remove Host
                                    </v-btn>
                                </v-layout>
                            </v-container>
                        </v-card>
                    </v-flex>
                </v-layout>
            </v-flex>
        </v-layout>
    </v-container>
</template>

<script>
  import axios from 'axios';

  export default {
    name: 'host',
    data: () => ({
      host: {}
    }),
    mounted() {
      axios.get('/v1/hosts/' + this.$route.params.id).then(res => {
        this.host = res.data
      }).catch(e => {
        console.log(e)
      });
    }
  }
</script>