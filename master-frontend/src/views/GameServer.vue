<template>
    <v-container>
        <v-row class="mb-6">
            <v-col cols="3" md="3" sm="12" xs="12">
                <v-card>
                    <v-card-title>Actions</v-card-title>
                    <v-card-text>
                        <v-btn class="ma-2" color="green" dark @click="start()">Start</v-btn>
                        <v-btn class="ma-2" color="red" dark @click="stop()">Stop</v-btn>
                        <v-btn class="ma-2" color="blue" dark @click="install()">Install</v-btn>
                    </v-card-text>
                </v-card>
                <v-card class="mt-5">
                    <v-card-title>Details</v-card-title>
                    <v-card-text>
                        <v-list-item two-line>
                            <v-list-item-content>
                                <v-list-item-title>{{game_server.name}}</v-list-item-title>
                                <v-list-item-subtitle>Name</v-list-item-subtitle>
                            </v-list-item-content>
                        </v-list-item>

                        <v-list-item two-line>
                            <v-list-item-content>
                                <v-list-item-title>{{game_server.state == 1 ? 'Active' : 'Locked'}}
                                </v-list-item-title>
                                <v-list-item-subtitle>Status</v-list-item-subtitle>
                            </v-list-item-content>
                        </v-list-item>

                        <v-list-item two-line>
                            <v-list-item-content>
                                <v-list-item-title>{{game_server.metric_frequency}}s</v-list-item-title>
                                <v-list-item-subtitle>Metric frequency</v-list-item-subtitle>
                            </v-list-item-content>
                        </v-list-item>

                        <v-list-item two-line>
                            <v-list-item-content>
                                <v-list-item-title>{{game_server.stop_timeout}}s</v-list-item-title>
                                <v-list-item-subtitle>Stop timeout</v-list-item-subtitle>
                            </v-list-item-content>
                        </v-list-item>
                    </v-card-text>
                    <v-card-actions>
                        <v-btn color="red darken-3" @click="remove()">
                            <v-icon class="mr-2">fa-trash</v-icon>
                            Remove game server
                        </v-btn>
                    </v-card-actions>
                </v-card>
            </v-col>
            <v-col cols="9" md="9" sm="12" xs="12">
                <v-card>
                    <v-card-title>Console</v-card-title>
                    <v-card class="pa-5">
                        <span v-for="item in logs" :key="item.id">{{item.log}}<br/></span>
                    </v-card>
                    <v-card-actions>
                        <v-text-field full-width
                                      label="Type some message here"
                                      hide-details
                                      v-model="message"
                                      v-on:keyup.enter="sendMessage()">
                            <v-icon slot="append" color="grey">fa-paper-plane</v-icon>
                        </v-text-field>
                    </v-card-actions>
                </v-card>
            </v-col>
        </v-row>
        <v-row class="mb-6">
            <v-col cols="3">

            </v-col>
        </v-row>
        <v-snackbar
                v-model="installSnackbar"
        >
            Installing game server...
            <v-btn
                    color="red"
                    text
                    @click="installSnackbar = false"
            >
                Close
            </v-btn>
        </v-snackbar>
        <v-snackbar
                v-model="removeSnackbar"
        >
            Removing game server...
            <v-btn
                    color="red"
                    text
                    @click="removeSnackbar = false"
            >
                Close
            </v-btn>
        </v-snackbar>
    </v-container>
</template>

<script>
    export default {
        name: 'game_server',
        data: () => ({
            game_server: {},
            serverId: '',
            hostId: '',
            logs: [],
            message: '',
            serverUrl: '',
            timer: '',
            installSnackbar: false,
            removeSnackbar: false,
        }),
        mounted() {
            this.serverId = this.$route.params.server_id;
            this.hostId = this.$route.params.id;
            this.serverUrl = '/v1/host/' + this.hostId + '/server/' + this.serverId;

            this.$http.get(this.serverUrl).then(res => {
                this.game_server = res.data.game_server;
            }).catch(e => {
                this.$auth.checkResponse(e.response.status)
            });

            let source = new EventSource(process.env.VUE_APP_API_URL + '/v1/console/' + this.serverId)
            var self = this;
            source.onmessage = function (event) {
                let str = atob(event.data)
                self.logs.push(str)
            }

            this.updateLogs();
            this.timer = setInterval(() => {
                this.updateLogs()
            }, 3 * 1000)
        },
        methods: {
            start() {
                this.$http.put(this.serverUrl + '/start').then(res => {
                    console.log(res);
                }).catch(e => {
                    this.$auth.checkResponse(e.response.status)
                })
            },
            stop() {
                this.$http.put(this.serverUrl + '/stop', {
                    type: 1
                }).then(res => {
                    console.log(res);
                }).catch(e => {
                    this.$auth.checkResponse(e.response.status)
                })
            },
            updateLogs() {
                this.$http.get(this.serverUrl + '/logs').then(res => {
                    this.logs = res.data.reverse();
                }).catch(e => {
                    this.$auth.checkResponse(e.response.status)
                })
            },
            sendMessage() {
                this.$http.put(this.serverUrl + '/command', {
                    command: this.message
                }).then(res => {
                    console.log(res);
                    this.message = '';
                }).catch(e => {
                    this.$auth.checkResponse(e.response.status)
                })
            },
            install() {
                this.$http.put(this.serverUrl + "/install").then(res => {
                    this.installSnackbar = true;
                }).catch(e => {
                    this.$auth.checkResponse(e.response.status)
                });
            },
            remove() {
                this.$http.delete(this.serverUrl).then(res => {
                    this.removeSnackbar = true;
                    console.log(res);
                    this.$router.push("/game-servers");
                }).catch(e => {
                    this.$auth.checkResponse(e.response.status)
                });
            }
        },
        beforeDestroy() {
            clearInterval(this.timer)
        }
    }
</script>
