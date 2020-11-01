<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card>
          <v-card-title>Backup Management</v-card-title>
          <v-card-text>
              <v-btn @click="backup()" color="green">
                Make backup
              </v-btn>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col>
        <v-card>
          <v-card-title>Backup list</v-card-title>
          <v-list>
            <v-list-item
              :key="backup"
              v-for="backup in backups"
            >
              <v-list-item-content>
                <v-list-item-title v-text="backup"></v-list-item-title>
              </v-list-item-content>
              <v-list-item-action>
                <v-btn @click="restore(backup)" color="orange">
                  RESTORE
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
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class GsBackup extends Vue {
  @Prop({
    type: String, required: true, default: () => {
      return '';
    }
  })
  serverId !: String;

  @Prop({
    type: String, required: true, default: () => {
      return '';
    }
  })
  hostId !: String;

  backups = [];

  mounted() {
    this.getBackups();
  }

  restore(backup: string) {
    this.$http.put('/v1/host/' + this.hostId + '/server/' + this.serverId + '/backup/restore', {
      backup_filename: backup
    }).then(res => {
    }).catch(e => {
      this.$auth.checkResponse(e.response.status);
    });
  }

  getBackups() {
    this.$http.get('/v1/host/' + this.hostId + '/server/' + this.serverId + '/backup/list').then(res => {
      this.backups = res.data.backups;
    }).catch(e => {
      this.$auth.checkResponse(e.response.status);
    });
  }

  backup() {
    this.$http.put('/v1/host/' + this.hostId + '/server/' + this.serverId + '/backup').then(res => {
      setTimeout(() => {
        this.getBackups()
      }, 2000)
    }).catch(e => {
      this.$auth.checkResponse(e.response.status);
    });
  }
}
</script>
