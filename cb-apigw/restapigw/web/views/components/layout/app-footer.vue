<template>
  <v-footer app height="42" class="white pa-3 app--footer">
    <v-spacer></v-spacer>
    <span class="caption">&copy; {{ new Date().getFullYear() }}</span>
    &nbsp;&nbsp;&nbsp;
    <span class="caption mr-1">
      Cloud-Barista.
    </span>
  </v-footer>
</template>

<script lang="ts">
import { Component, Vue } from "nuxt-property-decorator";

@Component({
  name: "app-footer",
  components: {}
})
export default class AppFooter extends Vue {
  // ---------------------------------
  // Fields
  // ---------------------------------
  // ---------------------------------
  // Getters/Setters
  // ---------------------------------
  // ---------------------------------
  // Methods
  // ---------------------------------
  // ---------------------------------
  // Lifecycle Events
  // ---------------------------------
}
</script>

<style lang="scss" scoped>
.caption {
  color: black;
}

// .app--footer {
//   position: absolute;
//   bottom: 0;
//   width: 100%;
// }
</style>
