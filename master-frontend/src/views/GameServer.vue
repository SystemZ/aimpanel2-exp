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
                        <span v-for="item in logs">{{item}}<br/></span>
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
        <v-row class="mb-6" v-if="files.selected">
            <v-col cols="12" md="12" sm="12" xs="12">
                <v-card>
                    <v-card-title>Files</v-card-title>

                    <v-list two-line subheader>
                        <v-subheader inset>Folders</v-subheader>

                        <v-list-item
                            @click="goToParentDirectory()"
                            v-if="files.selected.info.name !== files.root.info.name"
                        >
                            <v-list-item-avatar>
                                <v-icon class="grey lighten-1 white--text">
                                    fa-level-up
                                </v-icon>
                            </v-list-item-avatar>
                            <v-list-item-content>
                                <v-list-item-title>Go to parent directory</v-list-item-title>
                            </v-list-item-content>
                        </v-list-item>

                        <v-list-item
                                v-for="item in files.selected.children"
                                v-if="item.info.is_dir"
                                :key="item.info.name"
                                @click="files.selected = item"
                        >
                            <v-list-item-avatar>
                                <v-icon class="grey lighten-1 white--text">
                                    fa-folder
                                </v-icon>
                            </v-list-item-avatar>

                            <v-list-item-content>
                                <v-list-item-title v-text="item.info.name"></v-list-item-title>
                                <v-list-item-subtitle v-text="item.info.size"></v-list-item-subtitle>
                            </v-list-item-content>

                            <v-list-item-action>
                                <v-btn icon>
                                    <v-icon color="grey lighten-1">fa-info</v-icon>
                                </v-btn>
                            </v-list-item-action>
                        </v-list-item>

                        <v-divider inset></v-divider>

                        <v-subheader inset>Files</v-subheader>

                        <v-list-item
                                v-for="item in files.selected.children"
                                v-if="!item.info.is_dir"
                                :key="item.info.name"
                                @click=""
                        >
                            <v-list-item-avatar>
                                <v-icon class="blue white--text">fa-file</v-icon>
                            </v-list-item-avatar>

                            <v-list-item-content>
                                <v-list-item-title v-text="item.info.name"></v-list-item-title>
                                <v-list-item-subtitle v-text="item.info.size"></v-list-item-subtitle>
                            </v-list-item-content>

                            <v-list-item-action>
                                <v-btn icon>
                                    <v-icon color="grey lighten-1">fa-info</v-icon>
                                </v-btn>
                            </v-list-item-action>
                        </v-list-item>
                    </v-list>
                </v-card>
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

<script lang="ts">
    import Vue from 'vue';
    import {Node} from "@/types/files";

    interface FileRow {
        icon: string,
        iconClass: string,
        title: string,
        subtitle: string
    }

    export default Vue.extend({
        name: 'game_server',
        data: () => ({
            game_server: {},
            serverId: '',
            hostId: '',
            logs: [] as string[],
            message: '',
            serverUrl: '',
            timer: '',
            installSnackbar: false,
            removeSnackbar: false,
            stream: '' as any,

            files: {
                root: {} as Node,
                selected: {} as Node
            }
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

            if (this.stream === '' || this.stream === undefined) {
                this.setupStream()
            }

            this.getFiles()
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
            },
            setupStream() {
                this.stream = new this.$eventSource(this.$apiUrl + this.serverUrl + '/console', {
                    headers: {
                        Authorization: this.$auth.getAuthorizationHeader()
                    }
                });

                this.stream.onerror = (event: any) => {
                    console.error(event);
                };

                this.stream.addEventListener('message', (event: any) => {
                    if (event.data === 'heartbeat') {
                        return
                    }

                    let data = atob(event.data);
                    this.logs.push(data)
                }, false);
            },
            getFiles() {
                this.$http.get(this.serverUrl + '/file/list').then(res => {
                    this.files.root = res.data
                    this.files.selected = this.files.root
                }).catch(e => {
                    this.$auth.checkResponse(e.response.status)
                });
            },
            getParent(root: Node, name: string) {
                let node = null;
                if(name === '') {
                    return root;
                }

                root.children.some(n => {
                    if(n.info.name === name) {
                        return node = n;
                    }

                    if(n.children) {
                        return node = this.getParent(n, name)
                    }
                });

                return node;
            },
            goToParentDirectory() {
                if(this.files.selected.parent_name === this.files.root.info.name) {
                    this.files.selected = this.files.root;
                } else {
                    let node = this.getParent(this.files.root, this.files.selected.parent_name)
                    this.files.selected = node as Node;
                }
            }
        },
        beforeDestroy(): void {
            this.stream.close()
        },
    });
</script>
