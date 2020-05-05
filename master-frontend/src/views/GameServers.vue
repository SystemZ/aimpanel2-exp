<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
    <v-container>
        <v-row>
            <v-col xs12>
                <v-dialog max-width="600px" persistent v-model="createGameServer.dialog">
                    <template v-slot:activator="{ on }">
                        <v-btn color="info" v-on="on">
                            <v-icon left small>fa-plus</v-icon>
                            Create new game server
                        </v-btn>
                    </template>
                    <v-stepper v-model="createGameServer.step">
                        <v-stepper-header>
                            <v-stepper-step :complete="createGameServer.step > 1" step="1">Details</v-stepper-step>
                            <v-divider></v-divider>
                            <v-stepper-step :complete="createGameServer.step > 2" step="2">Install</v-stepper-step>
                        </v-stepper-header>
                        <v-stepper-items>
                            <v-stepper-content step="1">
                                <v-container grild-list-md>
                                    <v-layout wrap>
                                        <v-flex xs12>
                                            <v-select
                                                    :items="hosts"
                                                    item-text="name"
                                                    item-value="_id"
                                                    label="Select host"
                                                    v-model="createGameServer.selectedHost">
                                            </v-select>
                                        </v-flex>
                                        <v-flex xs12>
                                            <v-text-field label="Name" required
                                                          v-model="createGameServer.game.name"></v-text-field>
                                        </v-flex>
                                        <v-flex xs12>
                                            <v-select
                                                    :items="games"
                                                    item-text="name"
                                                    item-value="id"
                                                    label="Select game"
                                                    v-model="createGameServer.game.game_id">
                                            </v-select>
                                        </v-flex>
                                        <v-flex xs12>
                                            <v-select
                                                    :items="createGameServer.versions"
                                                    label="Select game version"
                                                    v-model="createGameServer.game.game_version">
                                            </v-select>
                                        </v-flex>
                                    </v-layout>
                                </v-container>

                                <v-btn @click="addGameServer()"
                                       color="primary">
                                    Next
                                </v-btn>

                                <v-btn @click="createGameServerCancel()" text>Cancel</v-btn>
                            </v-stepper-content>
                            <v-stepper-content step="2">
                                <v-container grid-list-md>
                                    <p>Game server was successfully created. Do you want to install it now?</p>
                                </v-container>
                                <v-btn @click="install()" color="info">Yes, install now</v-btn>
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
                        :items="gameServers"
                        @click:row="goToGameServer"
                        class="elevation-1"
                        hide-default-footer
                >
                    <template v-slot:item.state="{ item }">
                        <span v-if="item.state === 1">
                            <v-icon class="green--text" small>fa-circle</v-icon> Active
                        </span>
                        <span v-else>
                            <v-icon class="red--text" small>fa-circle</v-icon> Locked
                        </span>
                    </template>

                    <template v-slot:item.host="{ item }">
                        <span v-text="getHostName(item.host_id)"></span>
                    </template>

                    <template v-slot:item.game="{ item }">
                        <span v-text="getGameName(item.game_id)"></span>
                    </template>
                </v-data-table>
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
    import Vue from "vue";

    export default Vue.extend({
        name: "game_servers",
        data: () => ({
            headers: [
                {
                    text: "Name",
                    align: "left",
                    sortable: true,
                    value: "name"
                },
                {
                    text: "Host",
                    align: "right",
                    value: "host"
                },
                {
                    text: "Game",
                    align: "right",
                    value: "game"
                },
                {
                    text: "State",
                    align: "right",
                    value: "state"
                }
            ],
            games: [],
            gameVersions: [],
            gameServers: [],
            hosts: [],
            createGameServer: {
                dialog: false,
                step: 0,
                selectedHost: [],
                versions: [],
                game: {
                    name: "",
                    game_id: 0,
                    game_version: ""
                },
                gameId: "",
            },
            timer: ''
        }),
        methods: {
            goToGameServer(row) {
                this.$router.push('/host/' + row.host_id + '/server/' + row._id)
            },
            getGames() {
                return this.$http.get("/v1/game").then(res => {
                    this.games = res.data;
                }).catch(e => {
                    this.$auth.checkResponse(e.response.status)
                });
            },
            getGameServers() {
                return this.$http.get("/v1/host/my/server").then(res => {
                    this.gameServers = res.data.game_servers;
                }).catch(e => {
                    this.$auth.checkResponse(e.response.status)
                });
            },
            getHosts() {
                return this.$http.get("/v1/host").then(res => {
                    this.hosts = res.data.hosts;
                }).catch(e => {
                    this.$auth.checkResponse(e.response.status)
                });
            },
            addGameServer() {
                this.$http.post("/v1/host/" + this.createGameServer.selectedHost + "/server",
                    this.createGameServer.game).then(res => {

                    this.createGameServer.gameId = res.data.id;

                    this.createGameServer.step = 2;
                });
            },
            createGameServerCancel() {
                this.createGameServer.dialog = false;
                this.createGameServer.step = 1;
                this.finish()
            },
            finish() {
                this.createGameServer = {
                    dialog: false,
                    step: 0,
                    selectedHost: [],
                    game: {
                        name: "",
                        game_id: 0,
                        game_version: 0,
                    },
                    gameId: "",
                };
            },
            install() {
                this.$http.put("/v1/host/" + this.createGameServer.selectedHost +
                    "/server/" + this.createGameServer.gameId + "/install").then(res => {
                    console.log(res);
                });
            },
            getHostName(hostId) {
                if (this.hosts && this.hosts.length > 0) {
                    let host = this.hosts.find(x => x._id === hostId)
                    if (host) {
                        return host.name;
                    }
                }
                return ""
            },
            getGameName(gameId) {
                if (this.games && this.games.length > 0) {
                    let game = this.games.find(x => x.id === gameId)
                    if (game) {
                        return game.name;
                    }
                }
                return ""
            }
        },
        mounted() {
            this.getGames().then(() => {
                this.getHosts().then(() => {
                    this.getGameServers();
                });
            });

            this.timer = setInterval(() => {
                this.getGameServers()
            }, 10 * 1000)
        },
        beforeDestroy() {
            clearInterval(this.timer)
        },
        watch: {
            "createGameServer.game.game_id": function (val) {
                console.log(val)
                this.createGameServer.versions = this.games.filter((g) => {
                    return g.id === val
                })[0].versions
            }
        }
    });
</script>

<style>
    .clickable {
        cursor: pointer;
    }
</style>
