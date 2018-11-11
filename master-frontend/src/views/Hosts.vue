<template>
    <v-container fluid>
        <v-data-table
                :headers="headers"
                :items="hosts"
                hide-actions
                class="elevation-1"
        >
            <template slot="items" slot-scope="props">
                <td @click="goToHost(props.item.id)" class="clickable">{{ props.item.name }}</td>
                <td @click="goToHost(props.item.id)" class="text-xs-right clickable">{{ props.item.ip }}</td>
                <td @click="goToHost(props.item.id)" class="text-xs-right clickable">
                    <span v-if="props.item.state === 1">
                        <v-icon class="green--text" small>fa-circle</v-icon> Active
                    </span>
                    <span v-else>
                        <v-icon class="red--text" small>fa-circle</v-icon> Locked
                    </span>
                </td>
            </template>
        </v-data-table>
    </v-container>
</template>

<script>
  export default {
    name: 'hosts',
    data: () => ({
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
          text: 'State',
          align: 'right',
          value: 'state'
        }
      ],
      hosts: []
    }),
    methods: {
      goToHost(id) {
        this.$router.push('/host/' + id)
      },
      getHosts() {
        this.$http.get('/v1/hosts').then(res => {
          this.hosts = res.data;
        }).catch(e => {
          console.error(e)
        })
      }
    },
    mounted() {
      this.getHosts();
    }
  }
</script>

<style>
    .clickable {
        cursor: pointer;
    }
</style>