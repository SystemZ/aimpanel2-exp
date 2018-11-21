<template>
    <v-container fluid fill-height>
        <v-layout align-center justify-center>
            <v-flex xs12 sm8 md5>
                <v-card class="elevation-12">
                    <v-tabs
                            v-model="active"
                            color="red darken-1"
                            dark
                            slider-color="red"
                            fixed-tabs>
                        <v-tab ripple>
                            Login
                        </v-tab>
                        <v-tab ripple>
                            Register
                        </v-tab>

                        <v-tab-item>
                            <v-card-text>
                                <v-alert
                                        :value="loginError"
                                        type="error"
                                >
                                    {{loginError}}
                                </v-alert>
                                <v-form v-model="loginValid" @keyup.native.enter="loginValid && login()">
                                    <v-text-field prepend-icon="fa-user" label="Login" :rules="rules.username"
                                                  type="text" v-model="loginForm.username" required></v-text-field>

                                    <v-text-field prepend-icon="fa-key" label="Password" :rules="rules.password"
                                                  type="password" v-model="loginForm.password" required></v-text-field>
                                </v-form>
                            </v-card-text>
                            <v-card-actions>
                                <v-btn color="primary" @click="active = 1">Create new account</v-btn>
                                <v-spacer></v-spacer>
                                <v-btn color="success" :disabled="!loginValid" @click="login()">Login</v-btn>
                            </v-card-actions>
                        </v-tab-item>
                        <v-tab-item>
                            <v-card-text>
                                <v-alert
                                        :value="registerError"
                                        type="error"
                                >
                                    {{registerError}}
                                </v-alert>
                                <v-form v-model="registerValid" @keyup.native.enter="registerValid && register()">
                                    <v-text-field prepend-icon="fa-user" label="Username"
                                                  :rules="rules.username" type="text"
                                                  v-model="registerForm.username" required></v-text-field>

                                    <v-text-field prepend-icon="fa-key" label="Password"
                                                  :rules="rules.password" type="password"
                                                  v-model="registerForm.password" required></v-text-field>

                                    <v-text-field prepend-icon="fa-key" label="Repeat password"
                                                  :rules="rules.password" type="password"
                                                  v-model="registerForm.password_repeat"
                                                  :error-messages="confirmation.passwordConfirmation" required></v-text-field>

                                    <v-text-field prepend-icon="fa-envelope" label="Email"
                                                  :rules="rules.email" type="email"
                                                  v-model="registerForm.email" required></v-text-field>

                                    <v-text-field prepend-icon="fa-envelope" label="Repeat email"
                                                  :rules="rules.email" type="email"
                                                  v-model="registerForm.email_repeat"
                                                  :error-messages="confirmation.emailConfirmation" required></v-text-field>
                                </v-form>
                            </v-card-text>
                            <v-card-actions>
                                <v-spacer></v-spacer>
                                <v-btn color="success" :disabled="!registerValid" @click="register()">Register</v-btn>
                            </v-card-actions>
                        </v-tab-item>
                    </v-tabs>
                </v-card>
            </v-flex>
        </v-layout>
    </v-container>
</template>

<script>
  export default {
    name: 'login',
    data: () => ({
      loginForm: {
        username: '',
        password: ''
      },
      registerForm: {
        username: '',
        password: '',
        password_repeat: '',
        email: '',
        email_repeat: '',
      },
      rules: {
        username: [
          v => !!v || 'Username is required',
          v => v.length >= 3 || 'Username must be of minimum 3 characters'
        ],
        password: [
          v => !!v || 'Password is required',
          v => v.length >= 6 || 'Password must be of minimum 6 characters'
        ],
        email: [
          v => !!v || 'Email is required',
          v => /.+@.+/.test(v) || 'Email must be valid'
        ],
      },

      confirmation: {
        passwordConfirmation: '',
        emailConfirmation: '',
      },

      registerError: null,
      loginError: null,

      registerValid: false,
      loginValid: false,
      active: null
    }),
    methods: {
      login() {
        this.loginError = null;
        this.$auth.login(this, this.loginForm, 'hosts')
      },
      register() {
        this.registerError = null;
        this.$auth.register(this, this.registerForm, 'hosts')
      }
    },
    watch: {
      'registerForm.password_repeat': function () {
        if (this.registerForm.password !== this.registerForm.password_repeat) {
          this.confirmation.passwordConfirmation = 'Passwords do not match'
        } else {
          this.confirmation.passwordConfirmation = ''
        }
      },
      'registerForm.email_repeat': function () {
        if (this.registerForm.email !== this.registerForm.email_repeat) {
          this.confirmation.emailConfirmation = 'Emails do not match'
        } else {
          this.confirmation.emailConfirmation = ''
        }
      }
    }
  }
</script>