<template>
    <v-card>
        <v-toolbar dense flat>
            <v-toolbar-title>
                Console
            </v-toolbar-title>
            <v-spacer/>
            <small>Log limit</small><input v-model="logLimit" type="text"/>
            <v-checkbox v-model="shouldAutoscroll" :label="`Autoscroll: ${shouldAutoscroll.toString()}`"></v-checkbox>
        </v-toolbar>
        <v-card id="log-lines" class="pa-5" height="50vh" style="overflow-y: scroll">
            <span v-for="item in logs">{{item}}<br/></span>
        </v-card>
        <v-card-actions>
            <v-text-field full-width
                          hide-details
                          label="Type some message here"
                          v-model="message"
                          v-on:keyup.enter="sendMessage()">
                <v-icon color="grey" slot="append">fa-paper-plane</v-icon>
            </v-text-field>
        </v-card-actions>
    </v-card>
</template>

<script lang="ts">
    import Vue from 'vue';

    export default Vue.extend({
        name: 'gs-console',
        props: {
            serverId: {
                type: String,
                required: true,
            },
            hostId: {
                type: String,
                required: true,
            }
        },
        data: () => ({
            game_server: {},
            logs: [] as string[],
            shouldAutoscroll: true,
            logLimit: 100,
            message: '',
            serverUrl: '',
            stream: '' as any,
        }),
        mounted() {
            this.serverUrl = '/v1/host/' + this.hostId + '/server/' + this.serverId;

            this.$http.get(this.serverUrl).then(res => {
                this.game_server = res.data.game_server;
            }).catch(e => {
                this.$auth.checkResponse(e.response.status);
            });

            if (this.stream === '' || this.stream === undefined) {
                this.setupStream();
            }
        },
        methods: {
            scrollToBottom() {
                if (!this.shouldAutoscroll) {
                    return;
                }
                setTimeout(() => {
                    const log: HTMLElement | null = document.getElementById('log-lines');
                    if (log) {
                        log.scrollTo(0, log.scrollHeight);
                    }
                }, 10);
            },
            sendMessage() {
                this.$http.put(this.serverUrl + '/command', {
                    command: this.message
                }).then(res => {
                    console.log(res);
                    this.message = '';
                }).catch(e => {
                    this.$auth.checkResponse(e.response.status);
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
                        return;
                    }

                    let data = atob(event.data);
                    // clear old messages to prevent browser clogging
                    while (this.logs.length > this.logLimit - 1) {
                        this.logs.shift();
                    }
                    this.logs.push(data);
                    this.scrollToBottom();
                }, false);
            },
        },
        beforeDestroy(): void {
            this.stream.close();
        },
    });
</script>
