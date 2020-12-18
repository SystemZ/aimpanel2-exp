<template>
  <v-container>
    <v-row>
      <v-col cols="12" lg="4" md="4" sm="12">
        <v-card disabled>
          <v-card-text class="pa-12">
            <div class="pb-1 display-1 text--primary text-center">
              <span>Free</span>
            </div>
            <div class="pb-10 subtitle-1 text-center">
              $0 monthly
            </div>
            <div class="pl-5">
              <div>
                <v-icon class="pr-2">{{ mdiCheck }}</v-icon>
                1 host
              </div>
              <div>
                <v-icon class="pr-2">{{ mdiCheck }}</v-icon>
                1 game server
              </div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" lg="4" md="4" sm="12">
        <v-card>
          <v-card-text class="pa-12">
            <div class="pb-1 display-1 text--primary text-center">
              <span>Basic
                <v-chip
                  color="primary"
                  small
                >
                  Current
                </v-chip>
              </span>
            </div>
            <div class="pb-10 subtitle-1 text-center">
              $5 monthly
            </div>
            <div class="pl-5">
              <div>
                <v-icon class="pr-2">{{ mdiCheck }}</v-icon>
                3 hosts
              </div>
              <div>
                <v-icon class="pr-2">{{ mdiCheck }}</v-icon>
                10 game servers
              </div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" lg="4" md="4" sm="12">
        <v-card light>
          <v-card-text class="pa-12">
            <div class="pb-1 display-1 text-center">
              <span>Pro</span>
            </div>
            <div class="pb-10 subtitle-1 text-center">
              $10 monthly
            </div>
            <div class="pl-5">
              <div>
                <v-icon class="pr-2">{{ mdiCheck }}</v-icon>
                unlimited hosts
              </div>
              <div>
                <v-icon class="pr-2">{{ mdiCheck }}</v-icon>
                unlimited game servers
              </div>
            </div>
          </v-card-text>
          <v-card-actions class="pb-10 justify-center">
            <v-btn
              color="primary"
              x-large
            >
              UPGRADE
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import {Component, Vue} from 'vue-property-decorator';
import {mdiCheck} from '@mdi/js';

@Component
export default class LicensePage extends Vue {
// icons
  mdiCheck = mdiCheck;
}
</script>
