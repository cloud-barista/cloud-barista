<template>
  <v-app id="404">
    <template v-if="error.statusCode === 404">
      <v-container fluid fill-height>
        <v-layout align-center justify-center>
          <div class="mr-3 hidden-sm-and-down">
            <img
              src="@/utils/static/error/404.svg"
              alt="Sorry, page (resources) not found"
            />
          </div>
          <div class="text-md-center">
            <h1>404</h1>
            <h2 class="my-3 headline ">
              Sorry, page (resources) not found
            </h2>
            <div>
              <v-btn color="primary" @click="goHome">
                Go Home
              </v-btn>
            </div>
          </div>
        </v-layout>
      </v-container>
    </template>
    <template v-if="error.statusCode === 403">
      <v-container fluid fill-height>
        <v-layout align-center justify-center>
          <div class="mr-3 hidden-sm-and-down">
            <img
              src="@/utils/static/error/403.svg"
              alt="Sorry, access denied"
            />
          </div>
          <div class="text-md-center">
            <h1>403</h1>
            <h2 class="my-3 headline ">
              Sorry, access denied.
            </h2>
            <div>
              <v-btn color="primary" @click="goHome">
                Go Home
              </v-btn>
            </div>
          </div>
        </v-layout>
      </v-container>
    </template>
    <template v-if="error.statusCode === 500">
      <v-container fluid fill-height>
        <v-layout align-center justify-center row>
          <div class="mr-3 hidden-sm-and-down">
            <img
              src="@/utils/static/error/500.svg"
              alt="Sorry, the server has problems"
            />
          </div>
          <div class="text-md-center">
            <h1>500</h1>
            <h2 class="my-3 headline ">
              Sorry, the server has problems.
            </h2>
            <div>
              <v-btn color="primary" @click="goHome">
                Go Home
              </v-btn>
            </div>
          </div>
        </v-layout>
      </v-container>
    </template>
  </v-app>
</template>

<script lang="ts">
import { Component, Vue } from "nuxt-property-decorator";

@Component({
  props: ["error"],
  layout: "default",
  components: {}
})
export default class ErrorLayout extends Vue {
  // ---------------------------------
  // Fields
  // ---------------------------------
  // ---------------------------------
  // Getters/Setters
  // ---------------------------------
  // ---------------------------------
  // Methods
  // ---------------------------------
  private goHome() {
    this.$router.push({ path: "/" });
  }

  // ---------------------------------
  // Lifecycle Events
  // ---------------------------------
}
</script>

<style lang="scss" scoped>
h1 {
  font-size: 150px;
  line-height: 150px;
  font-weight: 700;
  color: #252932;
  text-shadow: rgba(61, 61, 61, 0.3) 1px 1px, rgba(61, 61, 61, 0.2) 2px 2px,
    rgba(61, 61, 61, 0.3) 3px 3px;
}
</style>
