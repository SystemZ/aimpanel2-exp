<template>
    <v-container fluid>
        <v-layout row wrap>
            <v-flex xs12>
                <v-dialog v-model="createGameServer.dialog" persistent max-width="600px">
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
                                                    item-value="id"
                                                    v-model="createGameServer.selectedHost"
                                                    label="Select host">
                                            </v-select>
                                        </v-flex>
                                        <v-flex xs12>
                                            <v-text-field label="Name" required v-model="createGameServer.game.name"></v-text-field>
                                        </v-flex>
                                        <v-flex xs12>
                                            <v-select
                                                    :items="games"
                                                    item-text="name"
                                                    item-value="id"
                                                    v-model="createGameServer.game.game_id"
                                                    label="Select game">
                                            </v-select>
                                        </v-flex>
                                    </v-layout>
                                </v-container>

                                <v-btn color="primary"
                                       @click="addGameServer()">
                                    Next
                                </v-btn>

                                <v-btn flat @click="createGameServerCancel()">Cancel</v-btn>
                            </v-stepper-content>
                            <v-stepper-content step="2">
                                <v-container grid-list-md>
                                    <p>Game server was successfully created. Do you want to install it now?</p>
                                    <v-btn color="info" @click="install()">Yes, install now</v-btn>
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
                        :items="gameServers"
                        hide-actions
                        class="elevation-1"
                >
                    <template slot="items" slot-scope="props">
                        <td @click="goToGameServer(props.item.id)" class="clickable">{{ props.item.name }}</td>
                        <td @click="goToGameServer(props.item.id)" class="text-xs-right clickable">
                            {{ hosts.find(x => x.id === props.item.host_id).name || '' }}</td>
                        <td @click="goToGameServer(props.item.id)" class="text-xs-right clickable">
                            {{ games.find(x => x.id === props.item.game_id).name || '' }}</td>
                    </template>
                </v-data-table>
            </v-flex>
        </v-layout>
    </v-container>
</template>

<script>
    export default {
        name: 'game_servers',
        data: () => ({
            headers: [
                {
                    text: 'Name',
                    align: 'left',
                    sortable: true,
                    value: 'name'
                },
                {
                    text: 'Host',
                    align: 'right',
                    value: 'host'
                },
                {
                    text: 'Game',
                    align: 'right',
                    value: 'game'
                }
            ],
            games: [],
            gameServers: [],
            hosts: [],
            createGameServer: {
                dialog: false,
                step: 0,
                selectedHost: [],
                game: {
                    name: '',
                    game_id: 0
                },
                gameId: '',
            }
        }),
        methods: {
            goToGameServer(id) {
              //this.$router.push('/host/' + id)
            },
            getGames() {
                this.$http.get('/v1/games').then(res => {
                    this.games = res.data;
                }).catch(e => {
                    console.error(e);
                })
            },
            getGameServers() {
                this.$http.get('/v1/hosts/my/servers').then(res => {
                    this.gameServers = res.data;
                }).catch(e => {
                    console.error(e)
                })
            },
            getHosts() {
                this.$http.get('/v1/hosts').then(res => {
                    this.hosts = res.data;
                }).catch(e => {
                    console.error(e)
                })
            },
            addGameServer() {
                this.$http.post('/v1/hosts/' + this.createGameServer.selectedHost + '/servers',
                    this.createGameServer.game).then(res => {

                    this.createGameServer.gameId = res.data.id;

                    this.createGameServer.step = 2;
                })
            },
            createGameServerCancel() {
                this.createGameServer.dialog = false;
                this.createGameServer.step = 1;
            },
            finish() {
                this.createGameServer = {
                    dialog: false,
                    step: 0,
                    game: {
                        name: '',
                        game_id: 0
                    }
                }
            },
            install() {
                this.$http.get('/v1/hosts/' + this.createGameServer.selectedHost +
                    '/servers/' + this.createGameServer.gameId + '/install').then(res => {
                        console.log(res)
                })
            }
        },
        mounted() {
            this.getGames();
            this.getHosts();
            this.getGameServers();
        }
    }
</script>