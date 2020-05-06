<template>
    <v-container>
        <v-row>
            <v-col xs12>
                <h2>Game servers</h2>
                <gs-create/>
                <p>
                    Your game servers. You know... for gaming
                </p>
                <gs-list/>

                <h2 class="mt-5">Hosts</h2>
                <host-create/>
                <p>
                    List of hardware or virtualized hosts that can be used to create game servers
                </p>
                <host-list/>
            </v-col>
        </v-row>
    </v-container>
</template>

<script lang="ts">
    import Vue from 'vue';
    import HostList from '@/components/HostList.vue';
    import GsList from '@/components/GsList.vue';
    import HostCreate from '@/components/HostCreate.vue';
    import GsCreate from '@/components/GsCreate.vue';

    export default Vue.extend({
        name: 'Hosts',
        components: {
            HostList,
            HostCreate,
            GsList,
            GsCreate,
        },
        data: () => ({}),
        methods: {},
    });
</script>
