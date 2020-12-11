<template>
  <v-text-field
    v-if="mask.length"
    v-model="val"
    v-mask="mask"
    :name="name"
    :label="lablel"
    :suffix="suffix"
    :rules="rules"
    :type="type"
    :disabled="disabled"
    :placeholder="placeholder"
    :outlined="outlined"
    :prepend-icon="prependicon"
    :prepend-inner-icon="prependinnericon"
    :append-icon="appendicon"
    :append-outer-icon="appendoutericon"
    :dense="dense"
    @keyup.enter="enterkeyPressed"
    @change="changeValue"
  />
  <v-text-field
    v-else
    v-model="val"
    :name="name"
    :label="label"
    :suffix="suffix"
    :rules="rules"
    :type="type"
    :disabled="disabled"
    :placeholder="placeholder"
    :outlined="outlined"
    :prepend-icon="prependicon"
    :prepend-inner-icon="prependinnericon"
    :append-icon="appendicon"
    :append-outer-icon="appendoutericon"
    :dense="dense"
    @keyup.enter="enterkeyPressed"
    @change="changeValue"
  />
</template>

<script lang="ts">
import { Component, Vue, Prop, Emit } from "nuxt-property-decorator";

@Component({
  name: "",
  components: {}
})
export default class CBText extends Vue {
  @Prop({ type: String, default: "" }) readonly name!: string;
  @Prop({ type: String, default: "" }) label!: string;
  @Prop({ type: String, default: "" }) suffix!: string;
  @Prop({ type: [String, Number], default: "" }) readonly value!: any;
  @Prop({ type: [String, Number], default: "" }) readonly orgValue!: any;
  @Prop({ type: String, default: "text" }) readonly type!: string;
  @Prop({ type: Array, default: () => [] }) readonly rules!: any;
  @Prop({ type: Boolean, default: false }) readonly disabled!: boolean;
  @Prop({ type: Boolean, default: false }) readonly dense!: boolean;
  @Prop({ type: String, default: "" }) readonly mask!: string;
  @Prop({ type: Boolean, default: true }) inputLabel!: boolean;
  @Prop({ type: Object, default: () => {} }) nodeData!: object;
  @Prop({ type: String, default: "" }) public prependicon!: string;
  @Prop({ type: String, default: "" }) public prependinnericon!: string;
  @Prop({ type: String, default: "" }) public appendicon!: string;
  @Prop({ type: String, default: "" }) public appendoutericon!: string;
  @Prop(String) readonly placeholder!: string;
  @Prop({ type: Boolean, default: false }) public outlined!: Boolean;
  @Emit("input") public changeValue(_val: [string, number], _nodata?: object) {}
  @Emit("change")
  @Emit("watch")
  private watchValue(val: any, orgVal: any) {
    return val === orgVal.toString();
  }

  // ---------------------------------
  // Fields
  // ---------------------------------
  // ---------------------------------
  // Getters/Setters
  // ---------------------------------

  private get val() {
    return this.value;
  }

  private set val(newVal: [string, number]) {
    this.changeValue(newVal, this.nodeData);
    this.watchValue(newVal, this.orgValue);
  }

  // ---------------------------------
  // Methods
  // ---------------------------------

  private enterkeyPressed() {
    this.$emit("enterkey", this);
  }

  // ---------------------------------
  // Lifecycle Events
  // ---------------------------------
}
</script>

<style lang="scss" scoped>
.v-text-field {
  &::v-deep .v-input__slot {
    border: 1px solid rgba(255, 255, 255, 0.1);
    background: rgba(0, 0, 0, 0.1);
  }

  // &::v-deep .v-text-field__details {
  //   position: absolute;
  //   margin: 14px 100px;
  //   padding: 0;
  // }
}
</style>
