<template>
  <v-btn
    ref="button"
    raised
    elevation="0"
    :large="size === 'large'"
    :small="size === 'small'"
    :text="size === 'text'"
    :block="block"
    :color="color"
    :disabled="disabled || loading"
    :ripple="ripple"
    :loading="loading"
    :icon="iconOnly"
    :v-bind="vBind"
    :v-on="vOn"
    @click="onClicked"
  >
    <v-icon v-if="icon.length && !iconOnly && !iconRight" class="mr-2">{{
      icon
    }}</v-icon>
    {{ text }}
    <slot v-if="!text.length && !iconOnly"></slot>
    <v-icon v-if="icon.length && !iconOnly && iconRight" class="ml-2">{{
      icon
    }}</v-icon>
    <v-icon v-if="iconOnly">{{ icon }}</v-icon>
  </v-btn>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "nuxt-property-decorator";

@Component({
  name: "cb-button",
  components: {}
})
export default class CBButton extends Vue {
  @Prop({ type: String }) readonly size!: string;
  @Prop({ type: String }) readonly color!: string;
  @Prop({ type: String, default: "" }) readonly text!: string;
  @Prop({ type: Boolean, default: false }) readonly disabled!: boolean;
  @Prop({ type: Boolean, default: false }) readonly block!: boolean;
  @Prop({ type: [Boolean, Object] }) readonly ripple!: [boolean, object];
  @Prop({ type: String, default: "" }) readonly icon!: string;
  @Prop({ type: Boolean, default: false }) readonly iconRight!: boolean;
  @Prop({ type: Boolean, default: false }) readonly iconOnly!: boolean;
  @Prop({ type: Object, default: () => {} }) readonly vBind!: object;
  @Prop({ type: Object, default: () => {} }) readonly vOn!: object;

  // ---------------------------------
  // Fields
  // ---------------------------------
  private loading: boolean = false;

  // ---------------------------------
  // Getters/Setters
  // ---------------------------------

  // --------------------------------
  // Methods
  // ---------------------------------

  private onClicked() {
    this.$emit("click", this);
  }

  // ---------------------------------
  // Lifecycle Events
  // --------------------------------
}
</script>

<style lang="scss" scoped></style>
