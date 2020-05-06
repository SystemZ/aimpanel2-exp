<template>
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
                            <v-icon class="green--text" small>{{mdiCheckboxBlankCircle}}</v-icon> Active
                        </span>
            <span v-else>
                            <v-icon class="red--text" small>{{mdiCheckboxBlankCircle}}</v-icon> Locked
                        </span>
        </template>
    </v-data-table>
</template>

<script lang="ts">
    import Vue from 'vue';
    import {Host} from '@/types/api';
    import {mdiCheckboxBlankCircle} from '@mdi/js';

    export default Vue.extend({
        name: 'host-create',
        data: () => ({
            hostListLoading: true,
            refreshInterval: 0,
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
            //icons
            mdiCheckboxBlankCircle: mdiCheckboxBlankCircle,
        }),
        mounted() {
            this.getHosts();
            this.refreshInterval = setInterval(() => {
                this.getHosts();
            }, 10 * 1000);
        },
        beforeDestroy(): void {
            clearInterval(this.refreshInterval);
        },
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
        }
    });
</script>
