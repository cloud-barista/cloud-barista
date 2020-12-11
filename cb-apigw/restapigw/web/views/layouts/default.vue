<template>
  <div id="defaultLayout">
    <v-app id="app" class="app">
      <!-- Top toolbar -->
      <app-toolbar v-if="$auth.loggedIn" />
      <!-- Content -->
      <v-main fixed>
        <!-- Page Content -->
        <div class="page-wrapper">
          <nuxt />
        </div>
      </v-main>
      <!-- Page Footer -->
      <app-footer />
    </v-app>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "nuxt-property-decorator";

import { AppToolbar, AppFooter } from "@/views/components";

@Component({
  name: "",
  components: {
    "app-toolbar": AppToolbar,
    "app-footer": AppFooter
  },
  auth: true
})
export default class DefaultLayout extends Vue {
  // ---------------------------------
  // Fields
  // ---------------------------------
  // ---------------------------------
  // Getters/Setters
  // ---------------------------------
  // ---------------------------------
  // Methdos
  // ---------------------------------
  // ---------------------------------
  // Lifecycle Events
  // ---------------------------------
}
</script>

<style lang="scss" scoped>
.v-main {
  overflow-y: auto;
  flex: 0 1 auto;
  height: calc(100vh - 35px);
}
</style>
