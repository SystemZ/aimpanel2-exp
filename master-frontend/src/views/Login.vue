<template>
    <v-container fluid fill-height>
        <v-layout align-center justify-center>
            <v-flex xs12 sm8 md4>
                <v-card class="elevation-12">
                    <v-toolbar dark color="red darken-1">
                        <v-toolbar-title>Login</v-toolbar-title>
                    </v-toolbar>
                    <v-card-text>
                        <v-form>
                            <v-text-field prepend-icon="fa-user" name="login" label="Login" type="text"
                                          v-model="username" @keyup.enter.native="login()"></v-text-field>
                            <v-text-field id="password" prepend-icon="fa-key" name="password" label="Password"
                                          type="password" v-model="password" @keyup.enter.native="login()"></v-text-field>
                        </v-form>
                    </v-card-text>
                    <v-card-actions>
                        <v-spacer></v-spacer>
                        <v-btn color="success" @click="login()">Login</v-btn>
                    </v-card-actions>
                </v-card>
            </v-flex>
        </v-layout>
    </v-container>
</template>

<script>
  export default {
    name: 'login',
    data: () => ({
      username: '',
      password: ''
    }),
    methods: {
      login() {
        this.$auth.login({
          username: this.username,
          password: this.password,
        }, 'hosts')
      }
    }
  }
</script>