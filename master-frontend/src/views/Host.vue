<template>
    <v-container grid-list-md>
        <v-row class="mb-6">
            <v-col>
                <v-card>
                    <v-card-text>
                        <v-btn @click="update()" class="ma-2" color="green" dark>Update</v-btn>
                    </v-card-text>
                </v-card>
            </v-col>
        </v-row>
        <v-row>
            <v-col xs6>
                <v-card>
                    <v-card-title>Details</v-card-title>
                    <v-card-text>
                        <v-list-item two-line>
                            <v-list-item-content>
                                <v-list-item-title>Active</v-list-item-title>
                                <v-list-item-subtitle>Status</v-list-item-subtitle>
                            </v-list-item-content>
                        </v-list-item>
                        <v-list-item two-line>
                            <v-list-item-content>
                                <v-list-item-title>{{host.os}}</v-list-item-title>
                                <v-list-item-subtitle>OS</v-list-item-subtitle>
                            </v-list-item-content>
                        </v-list-item>

                        <v-list-item two-line>
                            <v-list-item-content>
                                <v-list-item-title>{{host.platform}} {{host.platform_version}}</v-list-item-title>
                                <v-list-item-subtitle>Platform</v-list-item-subtitle>
                            </v-list-item-content>
                        </v-list-item>

                        <v-list-item two-line>
                            <v-list-item-content>
                                <v-list-item-title>{{host.kernel_version}}</v-list-item-title>
                                <v-list-item-subtitle>Kernel</v-list-item-subtitle>
                            </v-list-item-content>
                        </v-list-item>

                        <!--
                        FIXME token for debug purposes only, remove after implementing better slave auth
                        -->
                        <v-list-item two-line>
                            <v-list-item-content>
                                <v-list-item-title>{{host.token}}</v-list-item-title>
                                <v-list-item-subtitle>Token</v-list-item-subtitle>
                            </v-list-item-content>
                        </v-list-item>
                        
                    </v-card-text>
                    <v-card-actions>
                        <v-btn @click="remove()" color="red darken-2 accent-4" text>Remove host</v-btn>
                    </v-card-actions>
                </v-card>
            </v-col>
            <v-col xs6>
                <v-card>
                    <v-card-title>Performance</v-card-title>
                    <v-card-text align="center">
                        <v-row>
                            <v-col xs3>
                                <v-progress-circular
                                        :rotate="-90"
                                        :size="120"
                                        :value="metric.cpu_usage"
                                        :width="15"
                                        color="primary">
                                    {{metric.cpu_usage}}% CPU
                                </v-progress-circular>
                                <div class="text-xs-center">CPU</div>
                            </v-col>
                            <v-col xs3>
                                <v-progress-circular
                                        :rotate="-90"
                                        :size="120"
                                        :value="(metric.ram_used/metric.ram_total) * 100"
                                        :width="15"
                                        color="primary">
                                    {{metric.ram_used }}/{{metric.ram_total }} GB
                                </v-progress-circular>
                                <div class="text-xs-center">RAM</div>
                            </v-col>
                            <v-col xs3>
                                <v-progress-circular
                                        :rotate="-90"
                                        :size="120"
                                        :value="(metric.disk_used/metric.disk_total) * 100"
                                        :width="15"
                                        color="green">
                                    {{metric.disk_used}}/{{metric.disk_total}} GB
                                </v-progress-circular>
                                <div class="text-xs-center">Disk</div>
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-card>
            </v-col>
        </v-row>
        <v-row class="mb-6">
            <v-col>
                <v-card>
                    <v-card-text>
                        <v-text-field label="Name" required
                                      v-model="createJob.name"></v-text-field>

                        <v-text-field label="Cron expression" required
                                      v-model="createJob.cron_expression"></v-text-field>
                        <v-select
                                :items="gameServers"
                                item-text="name"
                                item-value="_id"
                                label="Select game server"
                                v-model="createJob.game_server_id">
                        </v-select>
                        <v-select :items="tasks"
                                  item-text="name"
                                  item-value="id"
                                  label="Select task"
                                  return-object
                                  single-line
                                  v-model="createJob.task_id"
                        ></v-select>
                        <v-text-field label="Body"
                                      v-model="createJob.body"></v-text-field>


                        <v-btn @click="addJob()" class="ma-2" color="green" dark>Add job</v-btn>
                    </v-card-text>
                </v-card>
            </v-col>
            <v-col>

                <v-card>
                    <v-card-title>Jobs</v-card-title>
                    <v-list>
                        <v-list-item
                                :key="job.name"
                                v-for="job in jobs"
                        >
                            <v-list-item-content>
                                <v-list-item-title v-text="job.name"></v-list-item-title>
                            </v-list-item-content>
                            <v-list-item-action>
                                <v-btn @click="removeJob(job._id)">
                                    DELETE
                                </v-btn>
                            </v-list-item-action>
                        </v-list-item>
                    </v-list>
                </v-card>
            </v-col>
        </v-row>
        <v-snackbar
                v-model="removeSnackbar"
        >
            Removing host...
            <v-btn
                    @click="removeSnackbar = false"
                    color="red"
                    text
            >
                Close
            </v-btn>
        </v-snackbar>
    </v-container>
</template>

<script lang="ts">
    import Vue from 'vue';
    import {Host, Metric} from '@/types/api';

    export default Vue.extend({
        name: 'host',
        data: () => ({
            host: {} as Host,
            metric: {} as Metric,
            fileManager: {
                current_dir: '/home/test',
                directories: [
                    {
                        icon: 'fa-folder',
                        title: 'plugins',
                        last_modification: '04.01.2019 20:44:33'
                    },
                    {
                        icon: 'fa-folder',
                        title: 'logs',
                        last_modification: '01.01.2019 19:43:00'
                    },
                    {
                        icon: 'fa-folder',
                        title: 'world',
                        last_modification: '12.12.2018 12:13:41'
                    }
                ],
                files: [
                    {
                        icon: 'fa-file',
                        title: 'server.properties',
                        last_modification: '14.12.2018 12:13:41'
                    },
                    {
                        icon: 'fa-file',
                        title: 'settings.yml',
                        last_modification: '03.12.2018 12:13:41'
                    }
                ]
            },
            removeSnackbar: false,
            createJob: {
                name: '',
                cron_expression: '* * * * *',
                task_id: {} as any,
                body: '',
                game_server_id: ''
            },
            tasks: [
                {
                    name: 'GAME_COMMAND',
                    id: 2
                },
                {
                    name: 'GAME_STOP_SIGKILL',
                    id: 3,
                },
                {
                    name: 'GAME_STOP_SIGTERM',
                    id: 4,
                }
            ],
            gameServers: [],
            jobs: []
        }),
        mounted(): void {
            this.$http.get('/v1/host/' + this.$route.params.id).then((res) => {
                this.host = res.data.host;
            }).catch(e => {
                this.$auth.checkResponse(e.response.status);
            });

            this.getGameServers();
            this.getJobs();

            this.$http.get('/v1/host/' + this.$route.params.id + '/metric').then((res) => {
                this.metric = res.data.metrics[0];

                this.metric.disk_free = +(this.metric.disk_free as number / 1024).toFixed(0);
                this.metric.disk_total = +(this.metric.disk_total / 1024).toFixed(0);
                this.metric.disk_used = this.metric.disk_total - this.metric.disk_free;

                this.metric.ram_total = +(this.metric.ram_total / 1024).toFixed(1);
                this.metric.ram_free = +(this.metric.ram_free / 1024).toFixed(1);
                this.metric.ram_used = +(this.metric.ram_total - this.metric.ram_free).toFixed(1);
            }).catch(e => {
                this.$auth.checkResponse(e.response.status);
            });
        },
        methods: {
            remove(): void {
                this.$http.delete('/v1/host/' + this.$route.params.id).then(res => {
                    this.removeSnackbar = true;
                    console.log(res);
                    this.$router.push('/hosts');
                }).catch(e => {
                    this.$auth.checkResponse(e.response.status);
                });
            },
            update(): void {
                this.$http.get('/v1/host/' + this.$route.params.id + '/update').then(res => {
                    console.log(res);
                }).catch(e => {
                    this.$auth.checkResponse(e.response.status);
                });
            },
            getGameServers(): void {
                this.$http.get('/v1/host/' + this.$route.params.id + '/server').then(res => {
                    this.gameServers = res.data.game_servers;
                }).catch(e => {
                    this.$auth.checkResponse(e.response.status);
                });
            },
            removeJob(jobId: string): void {
                this.$http.delete('/v1/host/' + this.$route.params.id + '/job/' + jobId).then(res => {
                    console.log(res);
                    this.getJobs();
                }).catch(e => {
                    this.$auth.checkResponse(e.response.status);
                });
            },
            getJobs(): void {
                this.$http.get('/v1/host/' + this.$route.params.id + '/job').then(res => {
                    this.jobs = res.data.jobs;
                }).catch(e => {
                    this.$auth.checkResponse(e.response.status);
                });
            },
            addJob(): void {
                this.$http.post('/v1/host/' + this.$route.params.id + '/job', {
                    name: this.createJob.name,
                    cron_expression: this.createJob.cron_expression,
                    task_message: {
                        task_id: this.createJob.task_id.id,
                        game_server_id: this.createJob.game_server_id,
                        body: this.createJob.body
                    }
                }).then(res => {
                    console.log(res);
                    this.createJob = {
                        name: '',
                        cron_expression: '* * * * *',
                        task_id: 0,
                        body: '',
                        game_server_id: ''
                    };
                    this.getJobs();
                });
            }
        },
    });
</script>
