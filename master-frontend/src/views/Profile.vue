<template>
    <v-container grid-list-md>
        <v-layout row>
            <v-flex xs6>
                <div class="display-1 grey--text text--darken-1">Profile</div>
                <br>
                <v-layout row>
                    <v-flex xs12>
                        <v-card>
                            <v-container>
                                <v-text-field label="Customer ID" v-model="user.id" disabled></v-text-field>
                                <v-text-field label="Username" v-model="user.username" disabled></v-text-field>
                                <v-text-field label="E-mail" v-model="user.email" disabled></v-text-field>

                                <v-dialog v-model="changeEmailDialog" persistent max-width="500px">
                                    <v-btn slot="activator" color="primary" dark>Change e-mail</v-btn>
                                    <v-card>
                                        <v-card-title>
                                            <span class="headline">Change e-mail</span>
                                        </v-card-title>
                                        <v-card-text>
                                            <v-alert
                                                    :value="changeEmailError"
                                                    type="error">
                                                {{changeEmailError}}
                                            </v-alert>
                                            <v-form v-model="changeEmailValid"
                                                    @keyup.native.enter="changeEmailValid && changeEmail()">
                                                <v-text-field label="Current e-mail*" required
                                                              :rules="rules.email"
                                                              v-model="changeEmailForm.email"></v-text-field>
                                                <v-text-field label="New e-mail*" required :rules="rules.email"
                                                              v-model="changeEmailForm.new_email"></v-text-field>
                                                <v-text-field label="Repeat new e-mail*" required
                                                              :rules="rules.email"
                                                              :error-messages="confirmation.emailConfirmation"
                                                              v-model="changeEmailForm.new_email_repeat"></v-text-field>
                                            </v-form>
                                        </v-card-text>
                                        <v-card-actions>
                                            <v-spacer></v-spacer>
                                            <v-btn color="blue darken-1" flat @click="changeEmailDialog = false">Close
                                            </v-btn>
                                            <v-btn color="primary" :disabled="!changeEmailValid" flat
                                                   @click="changeEmail()">Save
                                            </v-btn>
                                        </v-card-actions>
                                    </v-card>
                                </v-dialog>
                                <v-dialog v-model="changePasswordDialog" persistent max-width="500px">
                                    <v-btn slot="activator" color="primary" dark>Change password</v-btn>
                                    <v-card>
                                        <v-card-title>
                                            <span class="headline">Change password</span>
                                        </v-card-title>
                                        <v-card-text>
                                            <v-alert
                                                    :value="changePasswordError"
                                                    type="error">
                                                {{changePasswordError}}
                                            </v-alert>
                                            <v-form v-model="changePasswordValid"
                                                    @keyup.native.enter="changePasswordValid && changePassword()">
                                                <v-text-field label="Current password*" required type="password"
                                                              :rules="rules.password"
                                                              v-model="changePasswordForm.password"></v-text-field>
                                                <v-text-field label="New password*" required :rules="rules.password" type="password"
                                                              v-model="changePasswordForm.new_password"></v-text-field>
                                                <v-text-field label="Repeat new password*" required
                                                              :rules="rules.password"
                                                              :error-messages="confirmation.passwordConfirmation" type="password"
                                                              v-model="changePasswordForm.new_password_repeat"></v-text-field>
                                            </v-form>
                                        </v-card-text>
                                        <v-card-actions>
                                            <v-spacer></v-spacer>
                                            <v-btn color="blue darken-1" flat @click="changePasswordDialog = false">Close
                                            </v-btn>
                                            <v-btn color="primary" :disabled="!changePasswordValid" flat
                                                   @click="changePassword()">Save
                                            </v-btn>
                                        </v-card-actions>
                                    </v-card>
                                </v-dialog>
                            </v-container>
                        </v-card>
                    </v-flex>
                </v-layout>
            </v-flex>

        </v-layout>
    </v-container>
</template>

<script>
  export default {
    name: 'profile',
    data: () => ({
      user: {},
      changeEmailDialog: false,
      changeEmailForm: {
        email: '',
        new_email: '',
        new_email_repeat: ''
      },
      changePasswordDialog: false,
      changePasswordForm: {
        password: '',
        new_password: '',
        new_password_repeat: ''
      },
      rules: {
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
        emailConfirmation: '',
        passwordConfirmation: ''
      },
      changeEmailValid: false,
      changeEmailError: null,
      changePasswordValid: false,
      changePasswordError: null
    }),
    mounted() {
      this.getProfile()
    },
    methods: {
      changeEmail() {
        this.changeEmailError = null;
        this.$http.post('/v1/user/change_email', this.changeEmailForm).then(res => {
          this.changeEmailDialog = false;
          this.getProfile();
        }).catch(e => {
          this.changeEmailError = e.response.data.message
        })
      },
      changePassword() {
        this.changePasswordError = null;
        this.$http.post('/v1/user/change_password', this.changePasswordForm).then(res => {
          this.changePasswordDialog = false;
        }).catch(e => {
          this.changePasswordError = e.response.data.message
        })
      },
      getProfile() {
        this.$http.get('/v1/user/profile').then(res => {
          this.user = res.data
        }).catch(e => {
          console.error(e)
        });
      }
    },
    watch: {
      'changeEmailForm.new_email_repeat': function () {
        if (this.changeEmailForm.new_email !== this.changeEmailForm.new_email_repeat) {
          this.confirmation.emailConfirmation = 'Emails do not match'
        } else {
          this.confirmation.emailConfirmation = ''
        }
      },
      'changePasswordForm.new_password_repeat': function () {
        if (this.changePasswordForm.new_password !== this.changePasswordForm.new_password_repeat) {
          this.confirmation.passwordConfirmation = 'Passwords do not match'
        } else {
          this.confirmation.passwordConfirmation = ''
        }
      }
    }
  }
</script>