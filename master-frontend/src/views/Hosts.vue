<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
    <v-container fluid>
        <v-layout row wrap>
            <v-flex xs12>
                <v-dialog v-model="createHost.dialog" persistent max-width="600px">
                    <template v-slot:activator="{ on }">
                        <v-btn color="info" v-on="on">
                            <v-icon left small>fa-plus</v-icon>
                            Create new host
                        </v-btn>
                    </template>
                    <v-stepper v-model="createHost.step">
                        <v-stepper-header>
                            <v-stepper-step :complete="createHost.step > 1" step="1">Details</v-stepper-step>
                            <v-divider></v-divider>
                            <v-stepper-step :complete="createHost.step > 2" step="2">Setup host</v-stepper-step>
                        </v-stepper-header>
                        <v-stepper-items>
                            <v-stepper-content step="1">
                                <v-container grild-list-md>
                                    <v-layout wrap>
                                        <v-flex xs12>
                                            <v-text-field label="Name" required v-model="createHost.host.name"></v-text-field>
                                        </v-flex>
                                        <v-flex xs12>
                                            <v-text-field label="IP" required v-model="createHost.host.ip"></v-text-field>
                                        </v-flex>
                                    </v-layout>
                                </v-container>

                                <v-btn color="primary"
                                       @click="addHost()">
                                    Next
                                </v-btn>

                                <v-btn flat @click="createHostCancel()">Cancel</v-btn>
                            </v-stepper-content>
                            <v-stepper-content step="2">
                                <v-container grid-list-md>
                                    <p>Host was successfully created. This is a start token which you need to pass when you run host app.</p>
                                    <code>Token: {{createHost.token}}</code>
                                </v-container>
                                <v-btn flat @click="finish()">Close</v-btn>
                            </v-stepper-content>
                        </v-stepper-items>
                    </v-stepper>
                </v-dialog>

            </v-flex>
        </v-layout>
        <v-layout row wrap>
            <v-flex xs12>
                <v-data-table
                        :headers="headers"
                        :items="hosts"
                        hide-actions
                        class="elevation-1"
                >
                    <template slot="items" slot-scope="props">
                        <td @click="goToHost(props.item.id)" class="clickable">{{ props.item.name }}</td>
                        <td @click="goToHost(props.item.id)" class="text-xs-right clickable">{{ props.item.ip }}</td>
                        <td @click="goToHost(props.item.id)" class="text-xs-right clickable">
                    <span v-if="props.item.state === 1">
                        <v-icon class="green--text" small>fa-circle</v-icon> Active
                    </span>
                            <span v-else>
                        <v-icon class="red--text" small>fa-circle</v-icon> Locked
                    </span>
                        </td>
                    </template>
                </v-data-table>
            </v-flex>
        </v-layout>
    </v-container>
</template>

<script>
    export default {
        name: 'hosts',
        data: () => ({
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
                    text: 'State',
                    align: 'right',
                    value: 'state'
                }
            ],
            hosts: [],
            createHost: {
                dialog: false,
                step: 0,
                host: {
                    name: '',
                    ip: '',
                },
                token: ''
            },
        }),
        methods: {
            goToHost(id) {
                this.$router.push('/host/' + id)
            },
            getHosts() {
                this.$http.get('/v1/hosts').then(res => {
                    this.hosts = res.data;
                }).catch(e => {
                    console.error(e)
                })
            },
            createHostCancel() {
                this.createHost.dialog = false;
                this.createHost.step = 1;
            },
            addHost() {
                this.$http.post('/v1/hosts', this.createHost.host).then(res => {
                    this.createHost.token = res.data.token;
                    this.createHost.step = 2;
                }).catch(e => {
                    console.error(e)
                })
            },
            finish() {
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
        mounted() {
            this.getHosts();
        }
    }
</script>

<style>
    .clickable {
        cursor: pointer;
    }
</style>