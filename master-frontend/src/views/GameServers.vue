<template>
    <v-container fluid>
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
        }),
        methods: {
            goToGameServer(id) {
              //this.$router.push('/host/' + id)
            },
            getGames() {
                this.$http.get('/v1/games').then(res => {
                    console.log(res)
                    this.games = res.data;
                }).catch(e => {
                    console.error(e);
                })
            },
            getGameServers() {
                this.$http.get('/v1/hosts/my/servers').then(res => {
                    this.gameServers = res.data;
                    console.log(this.gameServers)
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
        },
        mounted() {
            this.getGames();
            this.getHosts();
            this.getGameServers();
        }
    }
</script>