<template>
    <v-container>
        <v-row>
            <v-col>
                <v-card>
                    <v-card-text>
                        <v-text-field label="Name" required
                                      v-model="createJob.name"></v-text-field>

                        <v-text-field label="Cron expression" required
                                      v-model="createJob.cron_expression"></v-text-field>
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
    </v-container>
</template>

<script lang="ts">
    import Vue from 'vue';

    export default Vue.extend({
        name: 'gs-scheduler',
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
            jobs: [],
            createJob: {
                name: '',
                cron_expression: '* * * * *',
                task_id: {} as any,
                body: '',
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
        }),
        mounted() {
            this.getJobs();
        },
        methods: {
            removeJob(jobId: string): void {
                this.$http.delete('/v1/host/' + this.hostId + '/job/' + jobId).then(res => {
                    console.log(res);
                    this.getJobs();
                }).catch(e => {
                    this.$auth.checkResponse(e.response.status);
                });
            },
            getJobs(): void {
                this.$http.get('/v1/host/' + this.hostId + '/job').then(res => {
                    this.jobs = res.data.jobs;
                }).catch(e => {
                    this.$auth.checkResponse(e.response.status);
                });
            },
            addJob(): void {
                this.$http.post('/v1/host/' + this.hostId + '/job', {
                    name: this.createJob.name,
                    cron_expression: this.createJob.cron_expression,
                    task_message: {
                        task_id: this.createJob.task_id.id,
                        game_server_id: this.serverId,
                        body: this.createJob.body
                    }
                }).then(res => {
                    console.log(res);
                    this.createJob = {
                        name: '',
                        cron_expression: '* * * * *',
                        task_id: 0,
                        body: '',
                    };
                    this.getJobs();
                });
            }
        },
    });
</script>
