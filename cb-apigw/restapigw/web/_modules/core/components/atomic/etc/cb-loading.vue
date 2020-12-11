<template>
  <div>
    <v-overlay v-if="overlay" :value="loading" z-index="1005"></v-overlay>
    <v-progress-circular
      v-if="loading"
      :size="size"
      :class="className"
      :color="color"
      indeterminate
    >
      {{ value }}
    </v-progress-circular>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "nuxt-property-decorator";

@Component({
  name: "cb-loading",
  components: {}
})
export default class CBLoading extends Vue {
  @Prop({ type: Boolean, default: true }) readonly overlay!: boolean;
  @Prop({ type: Boolean, default: false }) readonly loading!: boolean;
  @Prop({ type: String, default: null }) readonly value!: String;
  @Prop({ type: String, default: "30" }) readonly size!: String;
  @Prop({ type: String, default: "progress" }) readonly className!: boolean;
  @Prop({ type: String, default: "secondary" }) readonly color!: boolean;

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
.progress {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 1006;
}
</style>
