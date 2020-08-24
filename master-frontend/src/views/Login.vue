<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" md="4" sm="8">
        <v-card class="elevation-12">
          <v-tabs
            color="red darken-1"
            dark
            fixed-tabs
            slider-color="red"
            v-model="active">
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
                  {{ loginError }}
                </v-alert>
                <v-form @keyup.native.enter="loginValid && login()" v-model="loginValid">
                  <v-text-field
                    :rules="rules.username"
                    label="Login"
                    name="login"
                    :prepend-icon="mdiAccount"
                    type="text"
                    v-model="loginForm.username"
                  ></v-text-field>

                  <v-text-field
                    :rules="rules.password"
                    id="password"
                    label="Password"
                    name="password"
                    :prepend-icon="mdiLock"
                    type="password"
                    v-model="loginForm.password"
                  ></v-text-field>
                </v-form>
              </v-card-text>
              <v-card-actions>
                <v-btn @click="active = 1" color="primary">
                  <v-icon class="mr-2">{{ mdiAccountPlus }}</v-icon>
                  Create new account
                </v-btn>
                <v-spacer></v-spacer>
                <v-btn :disabled="!loginValid" @click="login()" color="success">
                  <v-icon class="mr-2">{{ mdiLogin }}</v-icon>
                  Login
                </v-btn>
              </v-card-actions>
            </v-tab-item>
            <v-tab-item>
              <v-card-text>
                <v-alert
                  :value="registerError"
                  type="error"
                >
                  {{ registerError }}
                </v-alert>
                <v-form @keyup.native.enter="registerValid && register()" v-model="registerValid">
                  <v-text-field :rules="rules.username" label="Username"
                                :prepend-icon="mdiAccount"
                                required
                                type="text" v-model="registerForm.username"></v-text-field>

                  <v-text-field :rules="rules.password" label="Password"
                                :prepend-icon="mdiLock"
                                required
                                type="password" v-model="registerForm.password"></v-text-field>

                  <v-text-field :error-messages="confirmation.passwordConfirmation"
                                :rules="rules.password"
                                label="Repeat password"
                                :prepend-icon="mdiLockOpenCheck"
                                required
                                type="password"
                                v-model="registerForm.password_repeat"></v-text-field>

                  <v-text-field :rules="rules.email" label="Email"
                                :prepend-icon="mdiEmail" required
                                type="email" v-model="registerForm.email"></v-text-field>

                  <v-text-field :error-messages="confirmation.emailConfirmation" :rules="rules.email"
                                label="Repeat email"
                                :prepend-icon="mdiEmailCheck"
                                required
                                type="email"
                                v-model="registerForm.email_repeat"></v-text-field>
                </v-form>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn :disabled="!registerValid" @click="register()" color="success">
                  <v-icon class="mr-2">{{ mdiRocket }}</v-icon>
                  Register
                </v-btn>
              </v-card-actions>
            </v-tab-item>
          </v-tabs>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue';
import {
  mdiAccount,
  mdiAccountPlus,
  mdiEmail,
  mdiEmailCheck,
  mdiLock,
  mdiLockOpenCheck,
  mdiLogin,
  mdiRocket
} from '@mdi/js';

export default Vue.extend({
  name: 'Login',
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
        (v: string) => !!v || 'Username is required',
        (v: string) => v.length >= 3 || 'Username must be of minimum 3 characters'
      ],
      password: [
        (v: string) => !!v || 'Password is required',
        (v: string) => v.length >= 5 || 'Password must be of minimum 5 characters'
      ],
      email: [
        (v: string) => !!v || 'Email is required',
        (v: string) => /.+@.+/.test(v) || 'Email must be valid'
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

    active: null,
    // icons
    mdiAccount: mdiAccount,
    mdiLock: mdiLock,
    mdiAccountPlus: mdiAccountPlus,
    mdiLogin: mdiLogin,
    mdiLockOpenCheck: mdiLockOpenCheck,
    mdiEmail: mdiEmail,
    mdiEmailCheck: mdiEmailCheck,
    mdiRocket: mdiRocket,
  }),
  methods: {
    login() {
      this.loginError = null;
      this.$auth.login(this, this.loginForm, '/');
    },
    register() {
      this.registerError = null;
      this.$auth.register(this, this.registerForm, '/');
    }
  },
  watch: {
    'registerForm.password_repeat': function() {
      if (this.registerForm.password !== this.registerForm.password_repeat) {
        this.confirmation.passwordConfirmation = 'Passwords do not match';
      } else {
        this.confirmation.passwordConfirmation = '';
      }
    },
    'registerForm.email_repeat': function() {
      if (this.registerForm.email !== this.registerForm.email_repeat) {
        this.confirmation.emailConfirmation = 'Emails do not match';
      } else {
        this.confirmation.emailConfirmation = '';
      }
    }
  }
});
</script>
