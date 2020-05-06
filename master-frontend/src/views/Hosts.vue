<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
    <v-container>
        <v-row>
            <v-col xs12>
                <v-dialog max-width="600px" persistent v-model="createHost.dialog">
                    <template v-slot:activator="{ on }">
                        <v-btn color="info" v-on="on">
                            <v-icon left small>fa-plus</v-icon>
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
            </v-col>
        </v-row>
        <v-row>
            <v-col xs12>
                <v-data-table
                        :headers="headers"
                        :items="hosts"
                        @click:row="goToHost"
                        class="elevation-1"
                        hide-default-footer
                        :loading="hostListLoading"
                >
                    <template v-slot:item.state="{ item }">
                        <span v-if="item.state === 1">
                            <v-icon class="green--text" small>fa-circle</v-icon> Active
                        </span>
                        <span v-else>
                            <v-icon class="red--text" small>fa-circle</v-icon> Locked
                        </span>
                    </template>
                </v-data-table>
            </v-col>
        </v-row>
    </v-container>
</template>

<script lang="ts">
    import Vue from 'vue';
    import {Host} from '@/types/api';

    export default Vue.extend({
        name: 'Hosts',
        data: () => ({
            hostListLoading: true,
            headers: [
                {
                    text: 'Name',
                    align: 'left',
                    sortable: true,
                    value: 'name'
                },
                {
                    text: 'IP',
                    align: 'right',
                    value: 'ip'
                },
                {
                    text: 'Game servers',
                    align: 'right',
                    value: 'gs'
                },
                {
                    text: 'State',
                    align: 'right',
                    value: 'state'
                }
            ],
            hosts: [] as Host[],
            createHost: {
                dialog: false,
                step: 0,
                host: {
                    name: '',
                    ip: '',
                },
                token: ''
            },
            timer: 0,
        }),
        methods: {
            goToHost(row: Host): void {
                this.$router.push('/host/' + row._id);
            },
            getHosts(): void {
                this.hostListLoading = true;
                this.$http.get('/v1/host').then(res => {
                    this.hosts = res.data.hosts;
                    this.hostListLoading = false;
                }).catch(e => {
                    this.$auth.checkResponse(e.response.status);
                });
            },
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
                this.getHosts();
            }
        },
        mounted(): void {
            this.getHosts();
            this.timer = setInterval(() => {
                this.getHosts();
            }, 10 * 1000);
        },
        beforeDestroy(): void {
            clearInterval(this.timer);
        }
    });
</script>

<style>
    .clickable {
        cursor: pointer;
    }
</style>
