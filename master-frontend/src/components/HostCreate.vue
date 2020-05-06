<template>
    <v-dialog max-width="600px" v-model="createHost.dialog">
        <template v-slot:activator="{ on }">
            <v-btn color="info" v-on="on" class="mt-2 mb-2">
                <v-icon class="mr-2">{{mdiPlus}}</v-icon>
                Create new host
            </v-btn>
        </template>
        <v-stepper v-model="createHost.step">
            <v-stepper-header>
                <v-stepper-step :complete="createHost.step > 1" step="1">Details</v-stepper-step>
                <v-divider>
                </v-divider>
                <v-stepper-step :complete="createHost.step > 2" step="2">Setup host</v-stepper-step>
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
                        <blockquote class="blockquote">wget
                            https://exp.upp.pl -O install && bash install {{createHost.token}} ; rm install
                        </blockquote>
                    </v-container>
                    <v-btn @click="finish()" text>Close</v-btn>
                </v-stepper-content>
            </v-stepper-items>
        </v-stepper>
    </v-dialog>
</template>

<script lang="ts">
    import Vue from 'vue';
    import {mdiPlus} from '@mdi/js';

    export default Vue.extend({
        name: 'host-new',
        data: () => ({
            createHost: {
                dialog: false,
                step: 0,
                host: {
                    name: '',
                    ip: '',
                },
                token: ''
            },
            //icons
            mdiPlus: mdiPlus,
        }),
        methods: {
            createHostCancel(): void {
                this.createHost.dialog = false;
                this.createHost.step = 1;
            },
            addHost(): void {
                this.$http.post('/v1/host', this.createHost.host).then(res => {
                    this.createHost.token = res.data.token;
                    this.createHost.step = 2;
                }).catch(e => {
                    this.$auth.checkResponse(e.response.status);
                });
            },
            finish(): void {
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
            },
        },
    });
</script>
