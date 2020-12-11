<template>
  <div id="app-toolbar">
    <v-app-bar clipped-left app flat height="52" dark>
      <template>
        <nuxt-link to="/">
          <img
            class="logo-box-img"
            src="/images/logo_cb.png"
            height="30"
            width="100"
            max-height="45"
            max-width="190"
            alt="Cloud-Barista API G/W Admin"
          />
        </nuxt-link>
        <v-spacer></v-spacer>
        <h1 class="header">Restful API G/W Admin</h1>
      </template>
      <v-spacer></v-spacer>
      <v-btn icon @click="handleFullScreen()">
        <v-icon>mdi-fullscreen</v-icon>
      </v-btn>
      <v-menu
        offset-y
        origin="center center"
        :nudge-right="140"
        :nudge-bottom="10"
        :close-on-content-click="closeOnContentClick"
      >
        <template v-slot:activator="{ on }">
          <v-btn icon large text v-on="on">
            <v-icon>mdi-account-cog</v-icon>
          </v-btn>
        </template>
        <v-list class="pa-0">
          <v-list-group
            v-for="item in themeItems"
            :key="item.title"
            :prepend-icon="item.action"
            no-action
          >
            <template v-slot:activator>
              <v-list-item-content>
                <v-list-item-title v-text="item.title"></v-list-item-title>
              </v-list-item-content>
            </template>
            <v-list-item v-for="subItem in item.items" :key="subItem.title">
              <v-list-item-content>
                <v-list-item-title>
                  <a @click.prevent="onThemeOptionChagned">{{
                    subItem.title
                  }}</a>
                </v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list-group>
          <v-list-item
            v-for="(item, index) in userItems"
            :key="index"
            :to="item.href"
            :nuxt="true"
            ripple="ripple"
            rel="noopener"
            @click="item.click"
          >
            <v-list-item-action v-if="item.icon">
              <v-icon>{{ item.icon }}</v-icon>
            </v-list-item-action>
            <v-list-item-content>
              <v-list-item-title>{{ item.title }}</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </v-menu>
    </v-app-bar>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "nuxt-property-decorator";

import Util from "@/utils";

@Component({
  name: "app-toolbar",
  components: {}
})
export default class AppToolbar extends Vue {
  // ---------------------------------
  // Fields
  // ---------------------------------

  private closeOnContentClick = false;
  private themeItems = [
    {
      action: "mdi-invert-colors",
      title: "Theme",
      items: [{ title: "Light" }, { title: "Dark" }]
    }
  ];

  private userItems = [
    {
      icon: "mdi-logout",
      href: "#",
      title: "Logout",
      click: this.handleLogout
    }
  ];

  // ---------------------------------
  // Getters/Setters
  // ---------------------------------

  private get theme() {
    return this.$vuetify.theme.dark ? "dark" : "light";
  }

  // ---------------------------------
  // Methods
  // ---------------------------------

  private touchAll() {
    const value = this.$vuetify.theme.themes[this.theme];
    this.$vuetify.theme.themes[this.theme] = {
      primary: "#000",
      secondary: "#000",
      accent: "#000",
      error: "#000",
      info: "#000",
      warning: "#000",
      success: "#000"
    };
    this.$vuetify.theme.themes[this.theme] = value;
  }

  private onThemeOptionChagned(e: any) {
    const value = e.target.textContent;
    this.$vuetify.theme.dark = value === "Dark";
  }

  private async handleLogout() {
    await this.$auth.logout();
  }

  private handleFullScreen() {
    Util.toggleFullScreen();
  }

  // ---------------------------------
  // Lifecycle Events
  // ---------------------------------
}
</script>

<style lang="scss" scoped></style>
