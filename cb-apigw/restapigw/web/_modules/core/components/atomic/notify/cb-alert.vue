<template>
  <v-alert
    v-model="isVisible"
    :type="type"
    transition="scale-transition"
    dismissible
    outlined
  >
    {{ value }}
  </v-alert>
</template>

<script lang="ts">
import { Component, Emit, Prop, Vue } from "nuxt-property-decorator";

@Component({
  name: "cb-alert"
})
export default class CBAlert extends Vue {
  @Prop(String) public value!: string;
  @Prop(String) public type!: string;
  @Emit("input") public changeValue(_newVal: string) {}

  private get isVisible() {
    return this.value !== "";
  }

  private set isVisible(newVal: boolean) {
    if (newVal === false) this.changeValue("");
  }
}
</script>

<style scoped></style>
