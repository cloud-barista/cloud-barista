<template>
  <v-text-field
    v-model="password"
    :rules="[requiredRules, passwordRules]"
    :prepend-icon="prependicon"
    :prepend-inner-icon="prependinnericon"
    :append-icon="appendicon || showViewIcon"
    :append-inner-icon="appendinnericon || showViewIcon"
    name="password"
    :label="label || $t('password')"
    :type="show ? 'text' : 'password'"
    :outlined="outlined"
    @keyup.enter="keyEnter"
    @click:append="show = !show"
  />
</template>

<script lang="ts">
import { Component, Emit, Prop, Vue } from "nuxt-property-decorator";
import { required, password } from "@/utils/validation";

@Component({
  name: "cb-password"
})
export default class CBPassword extends Vue {
  @Prop(String) public value!: string;
  @Prop({ type: String, default: "" }) public label!: string;
  @Prop({ type: Boolean, default: false }) public outlined!: Boolean;
  @Prop({ type: Boolean, default: false }) readonly dense!: boolean;
  @Prop({ type: String, default: "" }) public prependicon!: string;
  @Prop({ type: String, default: "" }) public prependinnericon!: string;
  @Prop({ type: String, default: "" }) public appendicon!: string;
  @Prop({ type: String, default: "" }) public appendinnericon!: string;
  @Prop({ type: Boolean, default: false }) readonly showView!: boolean;
  @Emit("input") public changeValue(_val: string) {}
  @Emit("keyup") keyEnter() {}

  // ---------------------------------
  // Fields
  // ---------------------------------

  private requiredRules = required;
  private passwordRules = password;
  private show = false;

  // ---------------------------------
  // Getters/Setters
  // ---------------------------------

  private get password(): string {
    return this.value;
  }

  private set password(newVal: string) {
    this.changeValue(newVal);
  }

  private get showViewIcon(): string {
    return this.showView ? (this.show ? "mdi-eye" : "mdi-eye-off") : "";
  }

  // ---------------------------------
  // Methods
  // ---------------------------------

  // ---------------------------------
  // Lifecycle Events
  // ---------------------------------
}
</script>

<style scoped lang="scss">
.v-text-field {
  &::v-deep .v-input__slot {
    border: 1px solid rgba(255, 255, 255, 0.1);
    background: rgba(0, 0, 0, 0.1);
  }
}
</style>
