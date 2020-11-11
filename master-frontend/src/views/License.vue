<template>
  <v-container>
    <v-row>
      <v-col cols="12" lg="4" md="4" sm="12">
        <v-card disabled>
          <v-card-title>
            <span class="headline">Free</span>
          </v-card-title>
          <v-card-subtitle>
            $0 monthly
          </v-card-subtitle>
          <v-card-text>
            <div>
              Features:
              <ul>
                <li>1 host</li>
                <li>1 game server</li>
              </ul>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" lg="4" md="4" sm="12">
        <v-card>
          <v-card-title>
            <span class="headline">
              Basic
            <v-chip
              color="primary"
              small
            >
              Current
            </v-chip>
            </span>
          </v-card-title>
          <v-card-subtitle>
            $5 monthly
          </v-card-subtitle>
          <v-card-text>
            <div>
              Features:
              <ul>
                <li>3 hosts</li>
                <li>10 game servers</li>
              </ul>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" lg="4" md="4" sm="12">
        <v-card
          light
        >
          <v-card-title>
            <span class="headline">Pro</span>
          </v-card-title>
          <v-card-subtitle>
            $10 monthly
          </v-card-subtitle>
          <v-card-text>
            <div>
              Features:
              <ul>
                <li>unlimited hosts</li>
                <li>unlimited game servers</li>
              </ul>
            </div>
          </v-card-text>
          <v-card-actions>
            <v-btn
              color="primary"
            >
              UPGRADE
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>


