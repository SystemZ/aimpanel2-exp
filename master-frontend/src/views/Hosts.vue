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
      hosts: [
        {
          id: 1,
          name: 'VPS #1',
          ip: '127.0.0.1',
          state: 1,
        },
        {
          id: 2,
          name: 'VPS #2',
          ip: '127.0.0.2',
          state: 0,
        }
      ]
    }),
    methods: {
      goToHost(id) {
        this.$router.push('/host/' + id)
      }
    }
  }
</script>

<style>
    .clickable {
        cursor: pointer;
    }
</style>