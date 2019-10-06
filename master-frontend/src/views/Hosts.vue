<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
    <v-container fluid>
        <v-row>
            <v-col xs12>
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
                                            <v-text-field label="Name" required
                                                          v-model="createHost.host.name"></v-text-field>
                                        </v-flex>
                                        <v-flex xs12>
                                            <v-text-field label="IP" required
                                                          v-model="createHost.host.ip"></v-text-field>
                                        </v-flex>
                                    </v-layout>
                                </v-container>

                                <v-btn color="primary"
                                       @click="addHost()">
                                    Next
                                </v-btn>

                                <v-btn text @click="createHostCancel()">Cancel</v-btn>
                            </v-stepper-content>
                            <v-stepper-content step="2">
                                <v-container grid-list-md>
                                    <p>Host was successfully created. This is a start token which you need to pass when
                                        you run host app.</p>
                                    <code>{{createHost.token}}</code>
                                </v-container>
                                <v-btn text @click="finish()">Close</v-btn>
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
                        hide-default-footer
                        class="elevation-1"
                >
                    <template v-slot:body="{ items }">
                        <tbody>
                        <tr v-for="item in items" :key="item.id" class="clickable" @click="goToHost(item.id)">
                            <td>{{item.name}}</td>
                            <td class="text-right">{{item.ip}}</td>
                            <td class="text-right">N/A</td>
                            <td class="text-right">
                                    <span v-if="item.state === 1">
                                        <v-icon class="green--text" small>fa-circle</v-icon> Active
                                    </span>
                                <span v-else>
                                        <v-icon class="red--text" small>fa-circle</v-icon> Locked
                                    </span>
                            </td>
                        </tr>
                        </tbody>
                    </template>
                </v-data-table>
            </v-col>
        </v-row>
    </v-container>
</template>

<script lang="ts">
    import Vue from "vue";

    export default Vue.extend({
        name: "Hosts",
        data: () => ({
            headers: [
                {
                    text: "Name",
                    align: "left",
                    sortable: true,
                    value: "name"
                },
                {
                    text: "IP",
                    align: "right",
                    value: "ip"
                },
                {
                    text: "Game servers",
                    align: "right",
                    value: "gs"
                },
                {
                    text: "State",
                    align: "right",
                    value: "state"
                }
            ],
            hosts: [],
            createHost: {
                dialog: false,
                step: 0,
                host: {
                    name: "",
                    ip: "",
                },
                token: ""
            },
        }),
        methods: {
            goToHost(id: any) {
                console.log('tes');
                this.$router.push("/host/" + id);
            },
            getHosts() {
                this.$http.get("/v1/host").then(res => {
                    this.hosts = res.data;
                }).catch(e => {
                    console.error(e);
                });
            },
            createHostCancel() {
                this.createHost.dialog = false;
                this.createHost.step = 1;
            },
            addHost() {
                this.$http.post("/v1/host", this.createHost.host).then(res => {
                    this.createHost.token = res.data.token;
                    this.createHost.step = 2;
                }).catch(e => {
                    console.error(e);
                });
            },
            finish() {
                this.createHost = {
                    dialog: false,
                    step: 0,
                    host: {
                        name: "",
                        ip: "",
                    },
                    token: ""
                };
                this.getHosts();
            }
        },
        mounted() {
            this.getHosts();
        }
    });
</script>

<style>
    .clickable {
        cursor: pointer;
    }
</style>