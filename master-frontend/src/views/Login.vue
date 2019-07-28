<template>
    <v-container fluid fill-height>
        <v-layout align-center justify-center>
            <v-flex xs12 sm8 md4>
                <v-card class="elevation-12">
                    <v-card-text>
                        <v-form v-model="loginValid" @keyup.native.enter="loginValid && login()">
                            <v-text-field
                                    label="Login"
                                    name="login"
                                    prepend-icon="fa-user"
                                    type="text"
                                    v-model="loginForm.username"
                                    :rules="rules.username"
                            ></v-text-field>

                            <v-text-field
                                    id="password"
                                    label="Password"
                                    name="password"
                                    prepend-icon="fa-key"
                                    type="password"
                                    v-model="loginForm.password"
                                    :rules="rules.password"
                            ></v-text-field>
                        </v-form>
                    </v-card-text>
                    <v-card-actions>
                        <v-spacer></v-spacer>
                        <v-btn color="primary" :disabled="!loginValid" @click="login()">Login</v-btn>
                    </v-card-actions>
                </v-card>
            </v-flex>
        </v-layout>
    </v-container>
</template>

<script lang="ts">
    import Vue from "vue";

    export default Vue.extend({
        name: "Login",
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
                    (v: string) => v.length >= 6 || 'Password must be of minimum 6 characters'
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
        }),
        methods: {
            login() {
                this.loginError = null;
                this.$auth.login(this, this.loginForm, 'hosts')
            }
        }
    });
</script>
