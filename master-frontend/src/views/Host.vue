<template>
    <v-container grid-list-md>
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
                    </v-card-text>
                    <v-card-actions>
                        <v-btn color="red darken-2 accent-4" text>Remove host</v-btn>
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
                                        :size="120"
                                        :width="15"
                                        :rotate="-90"
                                        :value="metric.cpu_usage"
                                        color="primary">
                                    {{metric.cpu_usage}}% CPU
                                </v-progress-circular>
                                <div class="text-xs-center">CPU</div>
                            </v-col>
                            <v-col xs3>
                                <v-progress-circular
                                        :size="120"
                                        :width="15"
                                        :rotate="-90"
                                        :value="(metric.ram_used/metric.ram_total) * 100"
                                        color="primary">
                                     {{this.metric.ram_used }}/{{this.metric.ram_total }} GB
                                </v-progress-circular>
                                <div class="text-xs-center">RAM</div>
                            </v-col>
                            <v-col xs3>
                                <v-progress-circular
                                        :size="120"
                                        :width="15"
                                        :rotate="-90"
                                        :value="(metric.disk_used/metric.disk_total) * 100"
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
        <v-row>
            <v-col xs6>

            </v-col>
            <v-col xs6>
<!--                <br>-->
<!--                <div class="display-1 grey&#45;&#45;text text&#45;&#45;darken-1">File Manager</div>-->
<!--                <br>-->
<!--                <v-row>-->
<!--                    <v-col xs12>-->
<!--                        <v-card>-->
<!--                            <v-list two-line subheader>-->
<!--                                <v-subheader inset>Current directory {{ fileManager.current_dir }}</v-subheader>-->
<!--                                <v-subheader inset>Directories</v-subheader>-->

<!--                                <v-list-item-->
<!--                                        v-for="dir in fileManager.directories"-->
<!--                                        :key="dir.title"-->
<!--                                        avatar-->
<!--                                        @click="">-->
<!--                                    <v-list-item-avatar>-->
<!--                                        <v-icon>{{ dir.icon }}</v-icon>-->
<!--                                    </v-list-item-avatar>-->

<!--                                    <v-list-item-content>-->
<!--                                        <v-list-item-title>{{ dir.title }}</v-list-item-title>-->
<!--                                        <v-list-item-sub-title>{{ dir.last_modification }}</v-list-item-sub-title>-->
<!--                                    </v-list-item-content>-->

<!--                                    <v-list-item-action>-->
<!--                                        <v-btn icon ripple>-->
<!--                                            <v-icon color="red lighten-1">fa-trash-o</v-icon>-->
<!--                                        </v-btn>-->
<!--                                    </v-list-item-action>-->
<!--                                </v-list-item>-->

<!--                                <v-divider inset></v-divider>-->

<!--                                <v-subheader inset>Files</v-subheader>-->

<!--                                <v-list-item-->
<!--                                        v-for="file in fileManager.files"-->
<!--                                        :key="file.title"-->
<!--                                        avatar-->
<!--                                        @click=""-->
<!--                                >-->
<!--                                    <v-list-item-avatar>-->
<!--                                        <v-icon>{{ file.icon }}</v-icon>-->
<!--                                    </v-list-item-avatar>-->

<!--                                    <v-list-item-content>-->
<!--                                        <v-list-item-title>{{ file.title }}</v-list-item-title>-->
<!--                                        <v-list-item-sub-title>{{ file.last_modification }}</v-list-item-sub-title>-->
<!--                                    </v-list-item-content>-->

<!--                                    <v-list-item-action>-->
<!--                                        <v-btn icon ripple>-->
<!--                                            <v-icon color="red lighten-1">fa-trash-o</v-icon>-->
<!--                                        </v-btn>-->
<!--                                    </v-list-item-action>-->
<!--                                </v-list-item>-->
<!--                            </v-list>-->
<!--                        </v-card>-->
<!--                    </v-col>-->
<!--                </v-row>-->
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
    export default {
        name: 'host',
        data: () => ({
            host: {},
            metric: {},
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
            }
        }),
        mounted() {
            this.$http.get('/v1/host/' + this.$route.params.id).then(res => {
                this.host = res.data
                console.log(this.host)
            }).catch(e => {
                console.error(e)
            });

            this.$http.get('/v1/host/' + this.$route.params.id + '/metric').then(res => {
                this.metric = res.data;

                this.metric.disk_free = (this.metric.disk_free / 1024).toFixed(0);
                this.metric.disk_total = (this.metric.disk_total / 1024).toFixed(0);
                this.metric.disk_used = this.metric.disk_total - this.metric.disk_free;

                this.metric.ram_total = (this.metric.ram_total / 1024).toFixed(1);
                this.metric.ram_free = (this.metric.ram_free / 1024).toFixed(1);
                this.metric.ram_used = this.metric.ram_total - this.metric.ram_free

                console.log(this.metric)
            }).catch(e => {
                console.error(e)
            })
        }
    }
</script>
